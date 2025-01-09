package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	eventSubjects "github.com/praction-networks/common/events/eventsubjects"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
)

type Publisher[T any] struct {
	Subject       eventSubjects.Subjects
	StreamManager *StreamManager
	EnableDedup   bool
	Metrics       *metrics.Metrics
}

func NewPublisher[T any](subject eventSubjects.Subjects, streamManager *StreamManager, enableDedup bool, metrics *metrics.Metrics) *Publisher[T] {
	return &Publisher[T]{
		Subject:       subject,
		StreamManager: streamManager,
		EnableDedup:   enableDedup,
		Metrics:       metrics,
	}
}

func (p *Publisher[T]) Publish(data T, config StreamConfig) error {
	// Ensure the stream exists
	if err := p.StreamManager.CreateOrUpdateStream(config); err != nil {
		if p.Metrics != nil {
			p.Metrics.FailedMessages.WithLabelValues(config.Name, string(p.Subject)).Inc()
		}
		logger.Error("Failed to create or update stream", "stream", config.Name, "error", err)
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	// Create event payload
	event := Event[T]{
		Version: 1,
		Subject: p.Subject,
		Data:    data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		if p.Metrics != nil {
			p.Metrics.FailedMessages.WithLabelValues(config.Name, string(p.Subject)).Inc()
		}
		logger.Error("Failed to marshal event", "subject", p.Subject, "error", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Deduplication logic
	options := []nats.PubOpt{}
	if p.EnableDedup {
		msgID := fmt.Sprintf("%s-%d", p.Subject, time.Now().UnixNano())
		options = append(options, nats.MsgId(msgID))
	}

	// Publish message
	ack, err := p.StreamManager.Client.Publish(string(p.Subject), payload, options...)
	if err != nil {
		if p.Metrics != nil {
			p.Metrics.FailedMessages.WithLabelValues(config.Name, string(p.Subject)).Inc()
		}
		logger.Error("Failed to publish event", "subject", p.Subject, "error", err)
		return fmt.Errorf("failed to publish event to subject %s: %w", p.Subject, err)
	}

	if p.Metrics != nil {
		p.Metrics.PublishedEvents.WithLabelValues(config.Name, string(p.Subject)).Inc()
	}
	logger.Info("Published event successfully", "subject", p.Subject, "Ack", ack)
	return nil
}
