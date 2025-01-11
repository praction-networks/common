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
type Listener[T any] struct {
	StreamName     string
	Durable        string
	Description    string
	DeliverPolicy  jetstream.DeliverPolicy
	AckPolicy      jetstream.AckPolicy
	AckWait        time.Duration
	MaxDeliver     int
	BackOff        []time.Duration
	ReplayPolicy   jetstream.ReplayPolicy
	RateLimit      uint64
	HeadersOnly    bool
	FilterSubject  *string
	FilterSubjects []string
	StreamManager  *JsStreamManager
	OnMessageFunc  func(data T, msg jetstream.Msg) error
	Metrics        *metrics.Metrics
	stopCh         chan struct{}
}

func NewListener[T any](
	streamName string,
	durable string,
	description string,
	deliverPolicy jetstream.DeliverPolicy,
	ackPolicy jetstream.AckPolicy,
	ackWait time.Duration,
	maxDeliver int,
	backOff []time.Duration,
	replayPolicy jetstream.ReplayPolicy,
	rateLimit uint64,
	headersOnly bool,
	filterSubject *string,
	filterSubjects []string,
	streamManager *JsStreamManager,
	onMessage func(data T, msg jetstream.Msg) error,
	metrics *metrics.Metrics,
) *Listener[T] {
	return &Listener[T]{
		StreamName:     streamName,
		Durable:        durable,
		Description:    description,
		DeliverPolicy:  deliverPolicy,
		AckPolicy:      ackPolicy,
		AckWait:        ackWait,
		MaxDeliver:     maxDeliver,
		BackOff:        backOff,
		ReplayPolicy:   replayPolicy,
		RateLimit:      rateLimit,
		HeadersOnly:    headersOnly,
		FilterSubject:  filterSubject,
		FilterSubjects: filterSubjects,
		StreamManager:  streamManager,
		OnMessageFunc:  onMessage,
		Metrics:        metrics,
		stopCh:         make(chan struct{}),
	}
}

func (l *Listener[T]) Listen(ctx context.Context) error {
	// Ensure stream exists
	stream, err := l.StreamManager.JsClient.Stream(ctx, l.StreamName)
	if err == jetstream.ErrStreamNotFound {
		logger.Warn("Stream not found, attempting to create it", "streamName", l.StreamName)
		streamConfig := jetstream.StreamConfig{
			Name:      l.StreamName,
			Subjects:  l.FilterSubjects,
			Retention: jetstream.LimitsPolicy, // Default retention policy
			Storage:   jetstream.FileStorage,  // Default storage type
			MaxMsgs:   -1,                     // Unlimited messages
			MaxBytes:  -1,                     // Unlimited size
			MaxAge:    0,                      // Unlimited age
		}
		stream, err = l.StreamManager.JsClient.CreateStream(ctx, streamConfig)
		if err != nil {
			logger.Error("Failed to create stream", err, "streamName", l.StreamName)
			return fmt.Errorf("failed to create stream %s: %w", l.StreamName, err)
		}
		logger.Info("Stream created successfully", "streamName", l.StreamName)
	}
	if err != nil {
		logger.Error("Stream not found for subject", err, "StreamName", l.StreamName)
		return fmt.Errorf("stream not found for stream %s: %w", l.StreamName, err)
	}

	// Create or update the consumer
	consumerConfig := jetstream.ConsumerConfig{
		Durable:        l.Durable,
		DeliverPolicy:  l.DeliverPolicy,
		AckPolicy:      l.AckPolicy,
		AckWait:        l.AckWait,
		MaxDeliver:     l.MaxDeliver,
		BackOff:        l.BackOff,
		ReplayPolicy:   l.ReplayPolicy,
		RateLimit:      l.RateLimit,
		HeadersOnly:    l.HeadersOnly,
		FilterSubjects: l.FilterSubjects,
	}
	if l.FilterSubject != nil {
		consumerConfig.FilterSubject = *l.FilterSubject
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	// Consume messages
	_, err = consumer.Consume(func(msg jetstream.Msg) {
		select {
		case <-l.stopCh:
			logger.Info("Listener stopped, skipping message processing", "StreamName", l.StreamName)
			return
		default:
			l.processMessage(msg)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", *l.FilterSubject, err)
	}

	logger.Info("Listening to subject", "FilterSubject", l.FilterSubject)
	return nil
}

func (l *Listener[T]) processMessage(msg jetstream.Msg) {
	start := time.Now()
	var event struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		logger.Error("Failed to unmarshal message", err, "FilterSubject", l.FilterSubject)
		msg.Nak()
		return
	}

	if err := l.OnMessageFunc(event.Data, msg); err != nil {
		logger.Error("Error processing message", err, "FilterSubject", l.FilterSubject)
	} else {
		msg.Ack()
		logger.Info("Message processed successfully", "FilterSubject", l.FilterSubject)
	}
	duration := time.Since(start).Seconds()
	if l.Metrics != nil && l.Metrics.Duration != nil {
		l.Metrics.Duration.WithLabelValues(l.StreamName).Observe(duration)
	}
}

func (l *Listener[T]) Stop(ctx context.Context) error {
	logger.Info("Stopping listener", "FilterSubject", l.FilterSubject)
	close(l.stopCh)
	return nil
}
