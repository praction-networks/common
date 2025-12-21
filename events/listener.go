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
	OnMessageFunc  func(ctx context.Context, msg Event[json.RawMessage]) error

	// Hardening knobs (all supported by your ConsumerConfig)
	MaxDeliver     int             // -1 or 0 => unlimited (server default)
	MaxAckPending  int             // avoid backpressure stalls (e.g., 2000)
	Backoff        []time.Duration // used for AckWait-expiry redelivery; we also use it to drive NakWithDelay
	InProgressTick time.Duration   // how often to call msg.InProgress() while handler runs

	// Optional: if we detect unrecoverable payload issues (e.g., JSON syntax), handle and drop
	PoisonHandler func(ctx context.Context, subject string, raw []byte, meta *jetstream.MsgMetadata)

	stopCh chan struct{}
}

// NewListener with sane defaults
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

		// defaults
		MaxDeliver:    -1, // explicit unlimited (server default)
		MaxAckPending: 2000,
		Backoff: []time.Duration{
			1 * time.Second, 2 * time.Second, 5 * time.Second,
			10 * time.Second, 30 * time.Second,
			1 * time.Minute, 2 * time.Minute, 5 * time.Minute,
		},
		InProgressTick: 0, // will default to AckWait/2 below
		stopCh:         make(chan struct{}),
	}
}

// Consumer operation constants
const (
	consumerOperationTimeout = 30 * time.Second // Increased timeout for NATS cluster consensus
	consumerCreateMaxRetries = 3                // Number of retry attempts for consumer creation
	consumerRetryBaseDelay   = 2 * time.Second  // Base delay between retries (exponential backoff)
)

// Listen initializes the consumer and starts message processing.
func (l *Listener) Listen(ctx context.Context) error {
	// Ensure stream exists
	stream, err := waitForStream(ctx, l.StreamManager, string(l.StreamName))
	if err != nil {
		logger.Error("Failed to fetch or wait for stream", err, "streamName", l.StreamName)
		return fmt.Errorf("failed to fetch or wait for stream %s: %w", l.StreamName, err)
	}

	// Default in-progress tick if not set
	if l.InProgressTick <= 0 && l.AckWait > 0 {
		l.InProgressTick = l.AckWait / 2
	}

	// For DeliverAllPolicy, always delete and recreate existing consumers to ensure
	// we replay ALL messages from the beginning, including seed events that may have been
	// published before the service started. Since handlers are idempotent, replaying is safe.
	// This ensures notification-service and tenant-user-service catch up on seed tenant events.
	if err := l.deleteExistingConsumerIfNeeded(ctx, stream); err != nil {
		logger.Warn("Failed to delete existing consumer, will try to update",
			"consumer", l.Durable,
			"stream", string(l.StreamName),
			err)
		// Continue anyway - CreateOrUpdateConsumer might still work
	}

	// Create or update the consumer with retry logic
	consumer, err := l.createConsumerWithRetry(ctx, stream)
	if err != nil {
		logger.Error("Failed to create or update consumer after retries", err, "streamName", l.StreamName)
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	sub, err := consumer.Consume(func(msg jetstream.Msg) {
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
	defer sub.Stop()

	logger.Info("Listening to subject(s)", "Stream", l.StreamName, "FilterSubjects", l.FilterSubjects, "FilterSubject", l.FilterSubject)

	// Wait for context cancellation or stop signal
	select {
	case <-ctx.Done():
		logger.Info("Context cancelled, stopping listener")
	case <-l.stopCh:
		logger.Info("Stop signal received, stopping listener")
	}
	return nil
}

func (l *Listener) processMessage(ctx context.Context, msg jetstream.Msg) {
	subject := msg.Subject()
	streamName := string(l.StreamName)

	metrics.IncNATSInflight(streamName, subject)
	defer metrics.DecNATSInflight(streamName, subject)

	start := time.Now()

	meta, err := msg.Metadata()
	if err != nil {
		logger.Error("Failed to get message metadata", err, "Subject", subject)
		metrics.RecordNATSFailure(streamName, subject, err)
		_ = msg.Nak()
		return
	}

	// keep long handlers alive
	done := make(chan struct{})
	var ticker *time.Ticker
	if l.InProgressTick > 0 {
		ticker = time.NewTicker(l.InProgressTick)
		defer ticker.Stop()
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					_ = msg.InProgress()
				}
			}
		}()
	}

	// decode
	var event Event[json.RawMessage]
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		close(done)
		logger.Error("Failed to unmarshal message", err, "Subject", subject, "Sequence", meta.Sequence)
		metrics.RecordNATSFailure(streamName, subject, err)

		// Poison: let a hook capture & persist, then Ack to stop redelivery loop.
		if l.PoisonHandler != nil {
			func() { // protect against panics in the hook
				defer func() { _ = recover() }()
				l.PoisonHandler(ctx, subject, msg.Data(), meta)
			}()
		}
		_ = msg.Ack() // acknowledge so it won't be redelivered forever
		return
	}

	// handle
	if err := l.OnMessageFunc(ctx, event); err != nil {
		close(done)
		logger.Error("OnMessageFunc returned error", err, "Subject", subject, "Sequence", meta.Sequence)
		metrics.RecordNATSFailure(streamName, subject, err)

		// Use NakWithDelay to control cadence (BackOff is only for AckWait timeouts)
		l.nakWithPolicy(msg, meta)
		return
	}

	// ack
	if err := msg.Ack(); err != nil {
		close(done)
		logger.Error("Failed to acknowledge message", err, "Subject", subject, "Sequence", meta.Sequence)
		metrics.RecordNATSFailure(streamName, subject, err)
		return
	}
	close(done)

	// record success (consider renaming to "consumed" in your metrics)
	metrics.RecordNATSPublished(streamName, subject)

	duration := time.Since(start)
	metrics.RecordNATSProcessingTime(streamName, subject, duration)

	logger.Info("Processed message successfully",
		"Stream", streamName, "Subject", subject, "Sequence", meta.Sequence, "Duration", duration)
}

func (l *Listener) nakWithPolicy(msg jetstream.Msg, meta *jetstream.MsgMetadata) {
	// choose a delay based on delivery count
	var delay time.Duration
	if len(l.Backoff) > 0 && meta != nil {
		idx := int(meta.NumDelivered) - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= len(l.Backoff) {
			idx = len(l.Backoff) - 1
		}
		delay = l.Backoff[idx]
	}
	if delay > 0 {
		_ = msg.NakWithDelay(delay)
	} else {
		_ = msg.Nak()
	}
}

// Stop gracefully stops the listener.
func (l *Listener) Stop(ctx context.Context) error {
	logger.Info("Stopping listener", "StreamName", l.StreamName)
	close(l.stopCh)
	return nil
}

// deleteExistingConsumerIfNeeded deletes the existing consumer if DeliverAllPolicy is set
// Uses an extended timeout to handle NATS cluster consensus delays
func (l *Listener) deleteExistingConsumerIfNeeded(ctx context.Context, stream jetstream.Stream) error {
	// Use extended timeout for consumer operations
	opCtx, cancel := context.WithTimeout(ctx, consumerOperationTimeout)
	defer cancel()

	existingConsumer, _ := stream.Consumer(opCtx, l.Durable)
	if existingConsumer != nil && l.DeliverPolicy == jetstream.DeliverAllPolicy {
		logger.Info("Deleting existing consumer to ensure replay of all messages with DeliverAllPolicy",
			"consumer", l.Durable,
			"stream", string(l.StreamName),
			"reason", "Need to catch up on all messages including seed events")

		deleteCtx, deleteCancel := context.WithTimeout(ctx, consumerOperationTimeout)
		defer deleteCancel()

		if err := stream.DeleteConsumer(deleteCtx, l.Durable); err != nil {
			return err
		}
		logger.Info("Existing consumer deleted successfully, will create new one to replay all messages",
			"consumer", l.Durable,
			"stream", string(l.StreamName))
	}
	return nil
}

// createConsumerWithRetry creates or updates the consumer with retry logic and exponential backoff
// If all retries fail, the process exits to trigger container restart by Kubernetes
func (l *Listener) createConsumerWithRetry(ctx context.Context, stream jetstream.Stream) (jetstream.Consumer, error) {
	// Build consumer config
	cc := jetstream.ConsumerConfig{
		Name:          l.Durable,
		Durable:       l.Durable,
		DeliverPolicy: l.DeliverPolicy,
		AckPolicy:     l.AckPolicy,
		AckWait:       l.AckWait,
		MaxDeliver:    l.MaxDeliver,
		MaxAckPending: l.MaxAckPending,
		BackOff:       l.Backoff,
	}
	if l.FilterSubject != nil {
		cc.FilterSubject = string(*l.FilterSubject)
	} else if len(l.FilterSubjects) > 0 {
		subjects := make([]string, len(l.FilterSubjects))
		for i, s := range l.FilterSubjects {
			subjects[i] = string(s)
		}
		cc.FilterSubjects = subjects
	}

	var lastErr error
	for attempt := 1; attempt <= consumerCreateMaxRetries; attempt++ {
		// Use extended timeout for each attempt
		opCtx, cancel := context.WithTimeout(ctx, consumerOperationTimeout)

		consumer, err := stream.CreateOrUpdateConsumer(opCtx, cc)
		cancel() // Always cancel the context after operation

		if err == nil {
			if attempt > 1 {
				logger.Info("Consumer created successfully after retry",
					"consumer", l.Durable,
					"stream", string(l.StreamName),
					"attempt", attempt)
			}
			return consumer, nil
		}

		lastErr = err
		logger.Warn("Failed to create consumer, will retry",
			"consumer", l.Durable,
			"stream", string(l.StreamName),
			"attempt", attempt,
			"maxRetries", consumerCreateMaxRetries,
			err)

		// Don't retry if context is already cancelled
		if ctx.Err() != nil {
			logger.Fatal("Context cancelled during consumer creation, exiting to trigger container restart",
				fmt.Errorf("context cancelled: %w", ctx.Err()),
				"consumer", l.Durable,
				"stream", string(l.StreamName))
		}

		// Don't sleep after the last attempt
		if attempt < consumerCreateMaxRetries {
			// Exponential backoff: 2s, 4s, 8s...
			backoffDelay := consumerRetryBaseDelay * time.Duration(1<<(attempt-1))
			logger.Info("Waiting before retry",
				"consumer", l.Durable,
				"delay", backoffDelay,
				"nextAttempt", attempt+1)

			select {
			case <-time.After(backoffDelay):
				// Continue to next attempt
			case <-ctx.Done():
				logger.Fatal("Context cancelled while waiting for retry, exiting to trigger container restart",
					fmt.Errorf("context cancelled: %w", ctx.Err()),
					"consumer", l.Durable,
					"stream", string(l.StreamName))
			}
		}
	}

	// All retries exhausted - exit process to trigger Kubernetes container restart
	logger.Fatal("Failed to create consumer after all retries, exiting to trigger container restart",
		lastErr,
		"consumer", l.Durable,
		"stream", string(l.StreamName),
		"attempts", consumerCreateMaxRetries)

	// This line won't be reached as logger.Fatal calls os.Exit(1)
	return nil, fmt.Errorf("failed to create consumer after %d attempts: %w", consumerCreateMaxRetries, lastErr)
}

func waitForStream(ctx context.Context, streamManager *JsStreamManager, streamName string) (jetstream.Stream, error) {
	const maxRetries = 20
	const retryInterval = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		stream, err := streamManager.JsClient.Stream(ctx, streamName)
		if err == nil {
			return stream, nil
		}
		if err == jetstream.ErrStreamNotFound {
			logger.Warn("Stream not found, retrying...", "streamName", streamName, "attempt", i+1)
			select {
			case <-time.After(retryInterval):
				continue
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled while waiting for stream %s", streamName)
			}
		} else {
			return nil, fmt.Errorf("error fetching stream %s: %w", streamName, err)
		}
	}
	return nil, fmt.Errorf("stream %s not found after %d retries", streamName, maxRetries)
}
