package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type ListenerType int

const (
	Critical ListenerType = iota
	Retryable
	OneTime
)

type Listener[T any] struct {
	Subject       Subjects
	StreamManager *StreamManager
	DeliverGroup  string
	ConsumerName  string
	AckWait       time.Duration
	MaxRetries    int
	Type          ListenerType
	OnMessageFunc func(data T, msg *nats.Msg) error
	stopCh        chan struct{}
	Subscription  *nats.Subscription
	Metrics       *metrics.Metrics
}

// Constructor for Listener
func NewListener[T any](
	subject Subjects,
	streamManager *StreamManager,
	deliverGroup string,
	ackWait time.Duration,
	maxRetries int,
	listenerType ListenerType,
	onMessage func(data T, msg *nats.Msg) error,
	metrics *metrics.Metrics,
) *Listener[T] {
	return &Listener[T]{
		Subject:       subject,
		StreamManager: streamManager,
		DeliverGroup:  deliverGroup,
		AckWait:       ackWait,
		MaxRetries:    maxRetries,
		Type:          listenerType,
		OnMessageFunc: onMessage,
		stopCh:        make(chan struct{}),
		Metrics:       metrics,
	}
}

// Listen Method
func (l *Listener[T]) Listen(streamName string) error {
	// Ensure stream exists
	if err := l.waitForStream(streamName); err != nil {
		return fmt.Errorf("stream validation failed: %w", err)
	}

	// Validate or create consumer
	if err := l.ensureConsumer(streamName); err != nil {
		return fmt.Errorf("consumer validation failed: %w", err)
	}

	// Subscribe to the subject
	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DeliverGroup, func(msg *nats.Msg) {
		l.processMessage(streamName, msg)
	}, nats.ManualAck(), nats.AckWait(l.AckWait))
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", l.Subject, err)
	}

	l.Subscription = sub
	logger.Info("Listening to subject:", "Subject", l.Subject, "DeliverGroup", l.DeliverGroup)
	return nil
}

// Wait for Stream
func (l *Listener[T]) waitForStream(streamName string) error {
	for {
		_, err := l.StreamManager.Client.StreamInfo(streamName)
		if err == nil {
			logger.Info("Stream is available", "stream", streamName)
			return nil
		}

		logger.Warn("Stream not available. Retrying...", "stream", streamName, "error", err)
		time.Sleep(5 * time.Second)
	}
}

// Ensure Consumer Exists
func (l *Listener[T]) ensureConsumer(streamName string) error {
	consumerName := l.ConsumerName
	if consumerName == "" {
		consumerName = fmt.Sprintf("%s-consumer", l.DeliverGroup)
	}

	// Check if consumer already exists
	consumerInfo, err := l.StreamManager.Client.ConsumerInfo(streamName, consumerName)
	if err != nil {
		if err == nats.ErrConsumerNotFound {
			// Create consumer if it doesn't exist
			consumerConfig := &nats.ConsumerConfig{
				Durable:        consumerName,
				DeliverGroup:   l.DeliverGroup,
				AckPolicy:      nats.AckExplicitPolicy,
				FilterSubject:  string(l.Subject),
				DeliverSubject: fmt.Sprintf("%s.deliver", consumerName),
			}

			_, err = l.StreamManager.Client.AddConsumer(streamName, consumerConfig)
			if err != nil {
				return fmt.Errorf("failed to create consumer: %w", err)
			}

			logger.Info("Consumer created successfully", "consumer", consumerName)
			return nil
		}
		// Return other errors
		return err
	}

	// Validate consumer config if it exists
	if consumerInfo != nil {
		if consumerInfo.Config.DeliverGroup != l.DeliverGroup {
			return fmt.Errorf("existing consumer has a different deliver group: %s", consumerInfo.Config.DeliverGroup)
		}
		logger.Info("Consumer exists and is valid", "consumer", consumerName)
	}

	return nil
}

// Process Single Message
func (l *Listener[T]) processMessage(streamName string, msg *nats.Msg) {
	start := time.Now()
	var event Event[T]
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		l.incrementMetric(l.Metrics.FailedMessages, streamName, err)
		msg.Nak()
		return
	}

	if err := l.OnMessageFunc(event.Data, msg); err != nil {
		l.incrementMetric(l.Metrics.FailedMessages, streamName, err)
	} else {
		l.incrementMetric(l.Metrics.ProcessedMessages, streamName, nil)
		msg.Ack()
	}

	duration := time.Since(start).Seconds()
	l.incrementDurationMetric(streamName, duration)
	logger.Info("Processed message", "Subject", l.Subject, "Duration", duration)
}

// Increment Metrics
func (l *Listener[T]) incrementMetric(metric *prometheus.CounterVec, streamName string, err error) {
	if l.Metrics != nil && metric != nil {
		metric.WithLabelValues(streamName, string(l.Subject)).Inc()
		if err != nil {
			logger.Error("Metric increment due to error", err, "stream", streamName, "subject", l.Subject)
		}
	}
}

// Increment Duration Metrics
func (l *Listener[T]) incrementDurationMetric(streamName string, duration float64) {
	if l.Metrics != nil && l.Metrics.Duration != nil {
		l.Metrics.Duration.WithLabelValues(streamName, string(l.Subject)).Observe(duration)
	}
}

// Stop Listener
func (l *Listener[T]) Stop(ctx context.Context) error {
	logger.Info("Stopping listener", "Subject", l.Subject, "DeliverGroup", l.DeliverGroup)
	close(l.stopCh)

	if l.Subscription != nil {
		if err := l.Subscription.Unsubscribe(); err != nil {
			logger.Error("Failed to unsubscribe from subject:", err)
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
		logger.Info("Unsubscribed from subject", "Subject", l.Subject)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
