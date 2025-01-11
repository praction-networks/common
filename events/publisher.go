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

// Publisher represents a generic publisher for JetStream.
type Publisher[T any] struct {
	Subject       string
	StreamManager *JsStreamManager
	EnableDedup   bool
	Metrics       *metrics.Metrics
}

func NewPublisher[T any](subject string, streamManager *JsStreamManager, enableDedup bool, metrics *metrics.Metrics) *Publisher[T] {
	return &Publisher[T]{
		Subject:       subject,
		StreamManager: streamManager,
		EnableDedup:   enableDedup,
		Metrics:       metrics,
	}
}

func (p *Publisher[T]) Publish(ctx context.Context, data T) error {
	// Ensure the stream exists
	streamInfo, err := p.StreamManager.Stream(ctx, p.Subject)
	if err != nil {
		if p.Metrics != nil {
			p.Metrics.FailedMessages.WithLabelValues("unknown", p.Subject).Inc()
		}
		logger.Error("Stream not found for subject", err, "subject", p.Subject)
		return fmt.Errorf("stream not found for subject %s: %w", p.Subject, err)
	}

	// Create event payload
	event := struct {
		Subject string `json:"subject"`
		Data    T      `json:"data"`
	}{
		Subject: p.Subject,
		Data:    data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		if p.Metrics != nil {
			p.Metrics.FailedMessages.WithLabelValues(streamInfo.Config.Name, p.Subject).Inc()
		}
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
	ack, err := p.StreamManager.JsClient.Publish(ctx, p.Subject, payload, options...)
	if err != nil {
		if p.Metrics != nil {
			p.Metrics.FailedMessages.WithLabelValues(streamInfo.Config.Name, p.Subject).Inc()
		}
		logger.Error("Failed to publish event", "subject", p.Subject, "error", err)
		return fmt.Errorf("failed to publish event to subject %s: %w", p.Subject, err)
	}

	if p.Metrics != nil {
		p.Metrics.PublishedEvents.WithLabelValues(streamInfo.Config.Name, p.Subject).Inc()
	}
	logger.Info("Published event successfully", "subject", p.Subject, "Ack", ack)
	return nil
}
