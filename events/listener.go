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
type Listener struct {
	StreamName     StreamName
	Durable        string
	DeliverPolicy  jetstream.DeliverPolicy
	AckPolicy      jetstream.AckPolicy
	AckWait        time.Duration
	FilterSubject  *Subject
	FilterSubjects []Subject
	StreamManager  *JsStreamManager
	OnMessageFunc  func(ctx context.Context, msg Event[json.RawMessage]) error // Application-specific handler
	stopCh         chan struct{}
}

// Constructor for Listener
func NewListener(
	streamName StreamName,
	durable string,
	deliverPolicy jetstream.DeliverPolicy,
	ackPolicy jetstream.AckPolicy,
	ackWait time.Duration,
	filterSubject *Subject,
	filterSubjects []Subject,
	streamManager *JsStreamManager,
	onMessageFunc func(ctx context.Context, msg Event[json.RawMessage]) error,
) *Listener {
	return &Listener{
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
func (l *Listener) Listen(ctx context.Context) error {
	// Ensure stream exists
	stream, err := waitForStream(ctx, l.StreamManager, string(l.StreamName))
	if err != nil {
		logger.Error("Failed to fetch or wait for stream", err, "streamName", l.StreamName)
		return fmt.Errorf("failed to fetch or wait for stream %s: %w", l.StreamName, err)
	}

	// Create or update the consumer
	consumerConfig := jetstream.ConsumerConfig{
		Name:          l.Durable,
		Durable:       l.Durable,
		DeliverPolicy: l.DeliverPolicy,

		AckPolicy: l.AckPolicy,
		AckWait:   l.AckWait,
	}
	if l.FilterSubject != nil {
		consumerConfig.FilterSubject = string(*l.FilterSubject)
	} else if len(l.FilterSubjects) > 0 {
		subjects := make([]string, len(l.FilterSubjects))
		for i, s := range l.FilterSubjects {
			subjects[i] = string(s)
		}
		consumerConfig.FilterSubjects = subjects
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
func (l *Listener) processMessage(ctx context.Context, msg jetstream.Msg) {

	msgMetaData, err := msg.Metadata()
	if err != nil {
		logger.Error("Failed to get message metadata", err, "StreamName", l.StreamName)
		metrics.FailedMessages.WithLabelValues(string(l.StreamName), msg.Subject()).Inc()
		msg.Nak()
		return
	}
	logger.Info("Processing message", "StreamName", l.StreamName, "Subject", msg.Subject(), "Sequence", msgMetaData.Sequence, "Timestamp", msgMetaData.Timestamp, "Consumer", msgMetaData.Consumer, "Data", string(msg.Data()))

	// Parse the message into a generic event wrapper
	var event Event[json.RawMessage]
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		logger.Error("Failed to unmarshal message", err, "Subject", msg.Subject())
		metrics.FailedMessages.WithLabelValues(string(l.StreamName), msg.Subject()).Inc()
		msg.Nak()
		return
	}

	// Call the application-defined message handler
	if l.OnMessageFunc != nil {
		err := l.OnMessageFunc(ctx, event)
		if err != nil {
			logger.Error("Error in OnMessageFunc", err, "Subject", event.Subject)
			metrics.FailedMessages.WithLabelValues(string(l.StreamName), msg.Subject()).Inc()
			msg.Nak()
			return
		}
	}

	// Acknowledge the message
	if err := msg.Ack(); err != nil {
		logger.Error("Failed to acknowledge message", err, "StreamName", l.StreamName)
	}

	metrics.ProcessedMessages.WithLabelValues(string(l.StreamName), msg.Subject()).Inc()

	logger.Info("Message processed", "StreamName", l.StreamName, "Subject", msg.Subject(), "Sequence", msgMetaData.Sequence, "Timestamp", msgMetaData.Timestamp, "Consumer", msgMetaData.Consumer)
}

// Stop gracefully stops the listener.
func (l *Listener) Stop(ctx context.Context) error {
	logger.Info("Stopping listener", "StreamName", l.StreamName)
	close(l.stopCh)
	return nil
}

func waitForStream(ctx context.Context, streamManager *JsStreamManager, streamName string) (jetstream.Stream, error) {
	// Define a maximum wait time
	const maxRetries = 20
	const retryInterval = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		stream, err := streamManager.JsClient.Stream(ctx, streamName)
		if err == nil {
			// Stream found
			return stream, nil
		}
		if err == jetstream.ErrStreamNotFound {
			logger.Warn("Stream not found, retrying...", "streamName", streamName, "attempt", i+1)
			select {
			case <-time.After(retryInterval): // Wait before retrying
				continue
			case <-ctx.Done(): // Context cancelled
				return nil, fmt.Errorf("context cancelled while waiting for stream %s", streamName)
			}
		} else {
			// Other error
			return nil, fmt.Errorf("error fetching stream %s: %w", streamName, err)
		}
	}

	return nil, fmt.Errorf("stream %s not found after %d retries", streamName, maxRetries)
}
