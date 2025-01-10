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

// NewListener constructor
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

// Listen method
func (l *Listener[T]) Listen(streamName string) error {
	// Validate stream existence
	if err := l.waitForStream(streamName); err != nil {
		return fmt.Errorf("stream validation failed: %w", err)
	}

	// Subscribe based on listener type
	switch l.Type {
	case OneTime:
		return l.setupOneTimeListener(streamName)
	default:
		return l.setupBufferedListener(streamName)
	}
}

// waitForStream checks if the stream exists and waits until it becomes available
func (l *Listener[T]) waitForStream(streamName string) error {
	for {
		_, err := l.StreamManager.Client.StreamInfo(streamName)
		if err == nil {
			logger.Info("Stream is available", "stream", streamName)
			return nil
		}

		logger.Warn("Stream not available. Waiting...", "stream", streamName, "error", err)
		time.Sleep(5 * time.Second) // Wait before retrying
	}
}

// One-Time Listener Setup
func (l *Listener[T]) setupOneTimeListener(streamName string) error {
	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DurableName, func(msg *nats.Msg) {
		l.processMessage(streamName, msg)
	}, nats.ManualAck(), nats.AckWait(l.AckWait), nats.Bind(streamName, l.DurableName))

	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", l.Subject, err)
	}

	l.Subscription = sub
	logger.Info("Listening to subject (One-Time):", "Subject", l.Subject, "Durable:", l.DurableName)
	return nil
}

// Buffered Listener Setup
func (l *Listener[T]) setupBufferedListener(streamName string) error {
	msgCh := make(chan *nats.Msg, 1024)

	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DurableName, func(msg *nats.Msg) {
		select {
		case msgCh <- msg:
		default:
			logger.Warn(fmt.Sprintf("Message dropped: %s", msg.Subject))
			msg.Term()
		}
	}, nats.ManualAck(), nats.AckWait(l.AckWait), nats.Bind(streamName, l.DurableName))

	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", l.Subject, err)
	}

	l.Subscription = sub
	go l.processMessages(streamName, msgCh)

	logger.Info("Listening to subject:", "Subject", l.Subject, "Durable:", l.DurableName)
	return nil
}

// Process Messages
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

// Process Single Message
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

	if err := l.OnMessageFunc(event.Data, msg); err != nil {
		if l.Metrics != nil {
			l.Metrics.FailedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
		}
		logger.Error("Error processing message:", err)
	} else if l.Metrics != nil {
		l.Metrics.ProcessedMessages.WithLabelValues(streamName, string(l.Subject)).Inc()
	}
	msg.Ack()
	logger.Info("Processed message", "Subject", l.Subject, "Duration", time.Since(start))
}

// Stop Listener
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
