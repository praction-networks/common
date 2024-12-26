package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher[T any] struct {
	Subject       Subjects
	StreamManager *StreamManager
	EnableDedup   bool
}

func NewPublisher[T any](subject Subjects, streamManager *StreamManager, enableDedup bool) *Publisher[T] {
	return &Publisher[T]{
		Subject:       subject,
		StreamManager: streamManager,
		EnableDedup:   enableDedup,
	}
}

func (p *Publisher[T]) Publish(data T, config StreamConfig) error {
	// Ensure the stream exists
	if err := p.StreamManager.CreateOrUpdateStream(config); err != nil {
		return fmt.Errorf("failed to ensure stream: %w", err)
	}

	// Create event payload
	event := Event[T]{
		Subject: p.Subject,
		Data:    data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
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
		return fmt.Errorf("failed to publish event to subject %s: %w", p.Subject, err)
	}

	log.Printf("Published event to subject: %s, Ack: %+v\n", p.Subject, ack)
	return nil
}
