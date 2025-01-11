package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/logger"
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
	stopCh         chan struct{}
	OnMessageFunc  func(ctx context.Context, data T, msg jetstream.Msg) error // Application-specific handler
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
	onMessageFunc func(ctx context.Context, data T, msg jetstream.Msg) error,
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
		OnMessageFunc:  onMessageFunc,
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
	subscription, err := consumer.Consume(func(msg jetstream.Msg) {
		select {
		case <-l.stopCh:
			logger.Info("Listener stopped, skipping message processing", "StreamName", l.StreamName)
			return
		default:
			l.processMessage(ctx, msg)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject: %w", err)
	}
	defer subscription.Stop()

	logger.Info("Listening to subject(s)", "FilterSubjects", l.FilterSubjects, "FilterSubject", l.FilterSubject)

	// Wait for context cancellation or stop signal
	select {
	case <-ctx.Done():
		logger.Info("Context cancelled, stopping listener")
	case <-l.stopCh:
		logger.Info("Stop signal received, stopping listener")
	}

	return nil
}

// processMessage processes a single message.
func (l *Listener[T]) processMessage(ctx context.Context, msg jetstream.Msg) {
	logger.Info("Processing message", "StreamName", l.StreamName, "Subject", msg.Subject())

	// Parse message data into the expected type
	var data T
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		logger.Error("Failed to unmarshal message", "error", err, "FilterSubject", l.FilterSubject)
		msg.Nak()
		return
	}

	// Call the application-defined message handler
	if l.OnMessageFunc != nil {
		err := l.OnMessageFunc(ctx, data, msg)
		if err != nil {
			logger.Error("Error in OnMessageFunc", "error", err, "FilterSubject", l.FilterSubject)
			msg.Nak()
			return
		}
	}

	// Acknowledge the message
	if err := msg.Ack(); err != nil {
		logger.Error("Failed to acknowledge message", "error", err, "StreamName", l.StreamName)
	}
}

// Stop gracefully stops the listener.
func (l *Listener[T]) Stop(ctx context.Context) error {
	logger.Info("Stopping listener", "FilterSubject", l.FilterSubject)
	close(l.stopCh)
	return nil
}
