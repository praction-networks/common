package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/praction-networks/common/logger"
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
	}
}

// Listen method to handle different listener types
func (l *Listener[T]) Listen(config StreamConfig) error {
	if err := l.StreamManager.CreateOrUpdateStream(config); err != nil {
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	if l.Type == OneTime {
		return l.setupOneTimeListener(config)
	}
	return l.setupBufferedListener(config)
}

// Setup for One-Time Listeners
func (l *Listener[T]) setupOneTimeListener(config StreamConfig) error {
	sub, err := l.StreamManager.Client.QueueSubscribe(string(l.Subject), l.DurableName, func(msg *nats.Msg) {
		var event Event[T]
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			logger.Error("Failed to unmarshal message", err)
			msg.Nak()
			return
		}

		if err := l.OnMessageFunc(event.Data, msg); err != nil {
			logger.Error("Error processing one-time message:", err)
		}
		msg.Ack()
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
	go l.processMessages(msgCh)

	logger.Info("Listening to subject:", "Subject", l.Subject, "Durable:", l.DurableName)
	return nil
}

// Message processing logic
func (l *Listener[T]) processMessages(msgCh chan *nats.Msg) {
	defer close(msgCh)

	for {
		select {
		case msg := <-msgCh:
			start := time.Now()
			var event Event[T]
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				logger.Error("Failed to unmarshal message", err)
				msg.Nak()
				continue
			}

			switch l.Type {
			case Critical:
				l.handleCriticalMessage(event.Data, msg)
			case Retryable:
				l.handleRetryableMessage(event.Data, msg)
			}

			logger.Info("Processed message", "Subject", l.Subject, "Duration", time.Since(start))
		case <-l.stopCh:
			logger.Info(fmt.Sprintf("Stopping listener for subject: %s", l.Subject))
			return
		}
	}
}

// Handle Critical messages with retries and DLQ
func (l *Listener[T]) handleCriticalMessage(data T, msg *nats.Msg) {
	retries := 0
	backoff := time.Second

	for {
		if err := l.OnMessageFunc(data, msg); err != nil {
			retries++
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
		msg.Ack()
		break
	}
}

// Handle Retryable messages with limited retries
func (l *Listener[T]) handleRetryableMessage(data T, msg *nats.Msg) {
	retries := 0
	for {
		if err := l.OnMessageFunc(data, msg); err != nil {
			retries++
			if retries > l.MaxRetries {
				logger.Warn("Retryable message dropped after max retries", "Subject", l.Subject, "Data", data)
				msg.Term()
				break
			}
			time.Sleep(time.Second * 2)
			continue
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
