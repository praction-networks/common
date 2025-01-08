package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
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
	DurableName   string
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
	durableName string,
	ackWait time.Duration,
	maxRetries int,
	listenerType ListenerType,
	onMessage func(data T, msg *nats.Msg) error,
	metrics *metrics.Metrics,
) *Listener[T] {
	return &Listener[T]{
		Subject:       subject,
		StreamManager: streamManager,
		DurableName:   durableName,
		AckWait:       ackWait,
		MaxRetries:    maxRetries,
		Type:          listenerType,
		OnMessageFunc: onMessage,
		stopCh:        make(chan struct{}),
		Metrics:       metrics,
	}
}

// Listen method to handle different listener types
func (l *Listener[T]) Listen(config StreamConfig) error {
	if err := l.StreamManager.CreateOrUpdateStream(config); err != nil {
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	switch l.Type {
	case OneTime:
		return l.setupOneTimeListener(config)
	default:
		return l.setupBufferedListener(config)
	}
}

// Setup for One-Time Listeners
func (l *Listener[T]) setupOneTimeListener(config StreamConfig) error {
	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DurableName, func(msg *nats.Msg) {
		l.processMessage(config.Name, msg)
	}, nats.ManualAck(), nats.AckWait(l.AckWait), nats.Bind(config.Name, l.DurableName))

	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", l.Subject, err)
	}

	l.Subscription = sub
	logger.Info("Listening to subject (One-Time):", "Subject", l.Subject, "Durable:", l.DurableName)
	return nil
}

// Setup for Buffered Listeners (Critical and Retryable)
func (l *Listener[T]) setupBufferedListener(config StreamConfig) error {
	msgCh := make(chan *nats.Msg, 1024)

	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DurableName, func(msg *nats.Msg) {
		select {
		case msgCh <- msg:
		default:
			logger.Warn(fmt.Sprintf("Message dropped: %s", msg.Subject))
			msg.Term()
		}
	}, nats.ManualAck(), nats.AckWait(l.AckWait), nats.Bind(config.Name, l.DurableName))

	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", l.Subject, err)
	}

	l.Subscription = sub
	go l.processMessages(config.Name, msgCh)

	logger.Info("Listening to subject:", "Subject", l.Subject, "Durable:", l.DurableName)
	return nil
}

// Message processing logic
func (l *Listener[T]) processMessages(streamName string, msgCh chan *nats.Msg) {
	defer close(msgCh)

	for {
		select {
		case msg := <-msgCh:
			l.processMessage(streamName, msg)
		case <-l.stopCh:
			logger.Info(fmt.Sprintf("Stopping listener for subject: %s", l.Subject))
			return
		}
	}
}

// Process a single message
func (l *Listener[T]) processMessage(streamName string, msg *nats.Msg) {
	start := time.Now()
	var event Event[T]
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		if l.Metrics != nil {
			l.Metrics.FailedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
		}
		logger.Error("Failed to unmarshal message", "error", err)
		msg.Nak()
		return
	}

	switch l.Type {
	case Critical:
		l.handleCriticalMessage(streamName, event.Data, msg)
	case Retryable:
		l.handleRetryableMessage(streamName, event.Data, msg)
	default:
		if err := l.OnMessageFunc(event.Data, msg); err != nil {
			if l.Metrics != nil {
				l.Metrics.FailedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
			}
			logger.Error("Error processing message:", err)
		} else if l.Metrics != nil {
			l.Metrics.ProcessedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
		}
		msg.Ack()
	}

	logger.Info("Processed message", "Subject", l.Subject, "Duration", time.Since(start))
}

// Handle Critical messages with retries and DLQ
func (l *Listener[T]) handleCriticalMessage(streamName string, data T, msg *nats.Msg) {
	retries := 0
	backoff := time.Second

	for {
		if err := l.OnMessageFunc(data, msg); err != nil {
			retries++
			if l.Metrics != nil {
				l.Metrics.FailedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
			}

			if retries > l.MaxRetries {
				dlqSubject := fmt.Sprintf("%s.dlq", string(l.Subject))
				l.StreamManager.Client.Publish(dlqSubject, msg.Data)
				logger.Error("Critical message moved to DLQ after max retries", "Subject", l.Subject, "Data", data)
				msg.Ack()
				break
			}

			logger.Warn(fmt.Sprintf("Retrying message (retry %d): %s", retries, l.Subject))
			time.Sleep(backoff)
			backoff *= 2
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			continue
		}
		if l.Metrics != nil {
			l.Metrics.ProcessedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
		}
		msg.Ack()
		break
	}
}

// Handle Retryable messages with limited retries
func (l *Listener[T]) handleRetryableMessage(streamName string, data T, msg *nats.Msg) {
	retries := 0
	for {
		if err := l.OnMessageFunc(data, msg); err != nil {
			retries++
			if l.Metrics != nil {
				l.Metrics.FailedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
			}
			if retries > l.MaxRetries {
				logger.Warn("Retryable message dropped after max retries", "Subject", l.Subject, "Data", data)
				msg.Term()
				break
			}
			time.Sleep(time.Second * 2)
			continue
		}
		if l.Metrics != nil {
			l.Metrics.ProcessedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
		}
		msg.Ack()
		break
	}
}

// Stop listener and clean up resources
func (l *Listener[T]) Stop(ctx context.Context) error {
	close(l.stopCh)

	if l.Subscription != nil {
		if err := l.Subscription.Unsubscribe(); err != nil {
			logger.Error("Failed to unsubscribe from subject:", err)
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
		logger.Info(fmt.Sprintf("Unsubscribed from subject: %s", l.Subject))
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
