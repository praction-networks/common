package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
)

// Listener represents a JetStream consumer with custom message handling logic.
type Listener[T any] struct {
	StreamName     string
	Durable        string
	DeliverPolicy  jetstream.DeliverPolicy
	AckPolicy      jetstream.AckPolicy
	AckWait        time.Duration
	FilterSubject  *string
	FilterSubjects []string
	StreamManager  *JsStreamManager
	OnMessageFunc  func(data T, msg jetstream.Msg) error
	Metrics        *metrics.Metrics
	stopCh         chan struct{}
}

// Constructor for Listener
func NewListener[T any](
	streamName string,
	durable string,
	deliverPolicy jetstream.DeliverPolicy,
	ackPolicy jetstream.AckPolicy,
	ackWait time.Duration,
	filterSubject *string,
	filterSubjects []string,
	streamManager *JsStreamManager,
	onMessage func(data T, msg jetstream.Msg) error,
	metrics *metrics.Metrics,
) *Listener[T] {
	return &Listener[T]{
		StreamName:     streamName,
		Durable:        durable,
		DeliverPolicy:  deliverPolicy,
		AckPolicy:      ackPolicy,
		AckWait:        ackWait,
		FilterSubject:  filterSubject,
		FilterSubjects: filterSubjects,
		StreamManager:  streamManager,
		OnMessageFunc:  onMessage,
		Metrics:        metrics,
		stopCh:         make(chan struct{}),
	}
}

// Listen initializes the consumer and starts message processing.
func (l *Listener[T]) Listen(ctx context.Context) error {
	// Ensure stream exists
	stream, err := l.StreamManager.JsClient.Stream(ctx, l.StreamName)
	if err == jetstream.ErrStreamNotFound {
		logger.Error("Stream not found", "streamName", l.StreamName)
		return fmt.Errorf("stream %s not found: %w", l.StreamName, err)
	}
	if err != nil {
		logger.Error("Error fetching stream", err, "streamName", l.StreamName)
		return fmt.Errorf("error fetching stream %s: %w", l.StreamName, err)
	}

	// Create or update the consumer
	consumerConfig := jetstream.ConsumerConfig{
		Name:          l.Durable,
		Durable:       l.Durable,
		DeliverPolicy: l.DeliverPolicy,
		AckPolicy:     l.AckPolicy,
		AckWait:       l.AckWait,
	}
	if l.FilterSubject != nil {
		consumerConfig.FilterSubject = *l.FilterSubject
	} else if len(l.FilterSubjects) > 0 {
		consumerConfig.FilterSubjects = l.FilterSubjects
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		logger.Error("Failed to create or update consumer", err, "streamName", l.StreamName)
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	// Consume messages
	_, err = consumer.Consume(func(msg jetstream.Msg) {
		select {
		case <-l.stopCh:
			logger.Info("Listener stopped, skipping message processing", "StreamName", l.StreamName)
			return
		default:
			l.processMessage(msg)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject: %w", err)
	}

	logger.Info("Listening to subject(s)", "FilterSubjects", l.FilterSubjects, "FilterSubject", l.FilterSubject)
	return nil
}

// processMessage processes a single message.
func (l *Listener[T]) processMessage(msg jetstream.Msg) {
	start := time.Now()
	var event struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		logger.Error("Failed to unmarshal message", err, "FilterSubject", l.FilterSubject)
		msg.Nak()
		return
	}

	if err := l.OnMessageFunc(event.Data, msg); err != nil {
		logger.Error("Error processing message", err, "FilterSubject", l.FilterSubject)
	} else {
		msg.Ack()
		logger.Info("Message processed successfully", "FilterSubject", l.FilterSubject)
	}
	duration := time.Since(start).Seconds()
	if l.Metrics != nil && l.Metrics.Duration != nil {
		l.Metrics.Duration.WithLabelValues(l.StreamName).Observe(duration)
	}
}

// Stop gracefully stops the listener.
func (l *Listener[T]) Stop(ctx context.Context) error {
	logger.Info("Stopping listener", "FilterSubject", l.FilterSubject)
	close(l.stopCh)
	return nil
}
