package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
)

// Publisher represents a generic publisher for JetStream.
type Publisher[T any] struct {
	Stream        StreamName
	Subject       Subject
	StreamManager *JsStreamManager
	EnableDedup   bool
}

// NewPublisher creates a new publisher.
func NewPublisher[T any](stream StreamName, subject Subject, streamManager *JsStreamManager, enableDedup bool) *Publisher[T] {
	return &Publisher[T]{
		Stream:        stream,
		Subject:       subject,
		StreamManager: streamManager,
		EnableDedup:   enableDedup,
	}
}

type RetryConfig struct {
	Enabled       bool
	Attempts      int
	RetryInterval time.Duration
}

// Publish publishes an event to JetStream.
func (p *Publisher[T]) Publish(ctx context.Context,
	data T,
	msgID string,
	retry RetryConfig,
) error {

	start := time.Now()
	var success bool

	defer func() {
		// Record duration metric for all publish attempts
		metrics.NATSPublishDuration.WithLabelValues(
			string(p.Stream),
			string(p.Subject),
			strconv.FormatBool(success),
		).Observe(time.Since(start).Seconds())
	}()
	// Ensure the stream exists
	streamInfo, err := p.StreamManager.Stream(ctx, string(p.Stream))
	if err != nil {

		metrics.RecordNATSFailure("unknown", string(p.Stream), err)
		logger.Error("Stream not found for subject", err, "Stream", p.Stream, "Subject", p.Subject)
		return fmt.Errorf("stream %s not found for subject %s: %w", p.Stream, p.Subject, err)
	}

	// Create the event payload
	event := Event[T]{
		Subject: p.Subject,
		Data:    data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
		logger.Error("Failed to marshal event", err, "subject", p.Subject)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Deduplication logic
	// Deduplication logic
	opts := []jetstream.PublishOpt{}
	if p.EnableDedup {
		if msgID == "" {
			msgID = fmt.Sprintf("%s-%d", p.Subject, time.Now().UnixNano())
		}
		opts = append(opts, jetstream.WithMsgID(msgID))
	}

	if !retry.Enabled {
		ack, err := p.StreamManager.JsClient.Publish(ctx, string(p.Subject), payload, opts...)
		if err != nil {
			metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
			logger.Error("Failed to publish event (no retry)", "subject", p.Subject, err)
			return fmt.Errorf("failed to publish event: %w", err)
		}
		success = true
		metrics.RecordNATSPublished(streamInfo.Config.Name, string(p.Subject))
		logger.Info("Published event successfully (no retry)", "subject", p.Subject, "ack", ack)
		return nil
	}

	// Retry mode
	for attempt := 1; ; attempt++ {
		ack, err := p.StreamManager.JsClient.Publish(ctx, string(p.Subject), payload, opts...)
		if err == nil {
			success = true
			metrics.RecordNATSPublished(streamInfo.Config.Name, string(p.Subject))
			logger.Info("Published event successfully (with retry)", "subject", p.Subject, "ack", ack)
			return nil
		}

		metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
		logger.Error("Failed to publish event", "subject", p.Subject, "attempt", attempt, "error", err)

		if retry.Attempts > 0 && attempt >= retry.Attempts {
			return fmt.Errorf("failed to publish event after %d attempts: %w", attempt, err)
		}

		select {
		case <-time.After(retry.RetryInterval):
		case <-ctx.Done():
			return fmt.Errorf("publish cancelled: %w", ctx.Err())
		}
	}
}
