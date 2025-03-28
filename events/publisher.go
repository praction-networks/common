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

// Publish publishes an event to JetStream.
func (p *Publisher[T]) Publish(ctx context.Context, data T) error {

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
	options := []jetstream.PublishOpt{}
	if p.EnableDedup {
		msgID := fmt.Sprintf("%s-%d", p.Subject, time.Now().UnixNano())
		options = append(options, jetstream.WithMsgID(msgID))
	}

	// Publish message
	ack, err := p.StreamManager.JsClient.Publish(ctx, string(p.Subject), payload, options...)
	if err != nil {
		metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
		logger.Error("Failed to publish event", "subject", p.Subject, "error", err)
		return fmt.Errorf("failed to publish event to subject %s: %w", p.Subject, err)
	}

	success = true
	metrics.RecordNATSPublished(streamInfo.Config.Name, string(p.Subject))
	logger.Info("Published event successfully", "subject", p.Subject, "Ack", ack)
	return nil
}
