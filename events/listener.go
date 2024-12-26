package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/praction-networks/common/logger"
)

type Listener[T any] struct {
	Subject       Subjects
	StreamManager *StreamManager
	DurableName   string
	AckWait       time.Duration
	MaxRetries    int
	OnMessageFunc func(data T, msg *nats.Msg) error
	stopCh        chan struct{}
	Subscription  *nats.Subscription
}

func NewListener[T any](
	subject Subjects,
	streamManager *StreamManager,
	durableName string,
	ackWait time.Duration,
	maxRetries int,
	onMessage func(data T, msg *nats.Msg) error,
) *Listener[T] {
	return &Listener[T]{
		Subject:       subject,
		StreamManager: streamManager,
		DurableName:   durableName,
		AckWait:       ackWait,
		MaxRetries:    maxRetries,
		OnMessageFunc: onMessage,
		stopCh:        make(chan struct{}),
	}
}

// Listen sets up a subscription to a NATS stream for the specified subject
func (l *Listener[T]) Listen(config StreamConfig) error {
	// Ensure the stream exists, creating or updating as necessary
	if err := l.StreamManager.CreateOrUpdateStream(config); err != nil {
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	// Use a channel to buffer messages and process them in a separate goroutine
	msgCh := make(chan *nats.Msg, 1024)

	// Subscribe to the stream using a queue and durable consumer
	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DurableName, func(msg *nats.Msg) {
		select {
		case msgCh <- msg:
		default:
			// Handle overflow (e.g., drop messages or log warning)
			logger.Warn(fmt.Sprintf("Message dropped: %s", msg.Subject))
			msg.Term()
		}
	}, nats.ManualAck(), nats.AckWait(l.AckWait), nats.Bind(config.Name, l.DurableName))

	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", l.Subject, err)
	}

	// Store the subscription
	l.Subscription = sub

	// Process messages in a separate goroutine
	go func() {
		defer close(msgCh)
		for {
			select {
			case msg := <-msgCh:
				var event Event[T]
				if err := json.Unmarshal(msg.Data, &event); err != nil {
					logger.Error("Failed to unmarshal message", err)
					msg.Nak() // Negative acknowledgment for malformed messages
					continue
				}

				// Pass the message data to the user-defined function for processing
				if err := l.OnMessageFunc(event.Data, msg); err != nil {
					logger.Error("Error processing message:", err)
				}
			case <-l.stopCh:
				logger.Info(fmt.Sprintf("Stopping listener for subject: %s", l.Subject))
				return
			}
		}
	}()

	logger.Info("Listening to subject:", "Subject", l.Subject, "Durable:", l.DurableName)
	return nil
}

// Stop stops the listener gracefully.
func (l *Listener[T]) Stop(ctx context.Context) error {
	close(l.stopCh)

	// Unsubscribe from the subject
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
