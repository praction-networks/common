package events

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

// Publisher represents a generic publisher for JetStream.
type Publisher[T any] struct {
	Stream          StreamName
	Subject         Subject
	StreamManager   *JsStreamManager
	EnableDedup     bool
	FallbackStorage *mongo.Collection
}

// NewPublisher creates a new publisher.
func NewPublisher[T any](stream StreamName, subject Subject, streamManager *JsStreamManager, enableDedup bool, fallback *mongo.Collection) *Publisher[T] {
	return &Publisher[T]{
		Stream:          stream,
		Subject:         subject,
		StreamManager:   streamManager,
		EnableDedup:     enableDedup,
		FallbackStorage: fallback,
	}
}

type RetryConfig struct {
	Enabled       bool
	Attempts      int
	RetryInterval time.Duration
}

// Publish keeps your original signature but internally uses PublishWithOptions.
// Prefer calling PublishWithOptions for fine-grained control.
func (p *Publisher[T]) Publish(ctx context.Context, data T, msgID string, retry RetryConfig) error {
	opts := PublishOptions{
		// Dedupe / Expectation
		// EnableDedup: nil => will default to p.EnableDedup below
		MsgID: msgID,

		// Retry/backoff (mapped from the legacy shape)
		RetryEnabled:  retry.Enabled,
		RetryAttempts: retry.Attempts,
		BaseBackoff:   retry.RetryInterval, // if not set, defaults are applied later
		MaxBackoff:    retry.RetryInterval,
		Jitter:        retry.RetryInterval / 10, // small jitter

		// Fallback
		FallbackAfterAttempts: 10,
	}
	_, err := p.PublishWithOptions(ctx, data, opts)
	return err
}

// -----------------------------
// New options-rich publish API
// -----------------------------

type PublishOptions struct {
	// Dedupe / expectations
	EnableDedup      *bool  // nil => default from Publisher.EnableDedup
	MsgID            string // stable id for dedupe; generated if empty and dedupe enabled
	ExpectStreamName string // if empty, defaults to publisher's Stream

	// Retry & backoff
	RetryEnabled  bool
	RetryAttempts int           // total attempts incl. first; <=0 => practically infinite
	BaseBackoff   time.Duration // e.g., 200ms
	MaxBackoff    time.Duration // e.g., 5s
	Jitter        time.Duration // e.g., 200ms

	// Fallback
	FallbackAfterAttempts int // upsert to Mongo after this many tries (default: 10)

	// Optional payload guard if your stream uses MaxMsgSize
	MaxMsgSize int // 0 = disabled
}

func (o *PublishOptions) withDefaults(pub *Publisher[any]) PublishOptions {
	cp := *o
	if cp.EnableDedup == nil {
		v := pub.EnableDedup
		cp.EnableDedup = &v
	}
	if !cp.RetryEnabled {
		cp.RetryAttempts = 1
	}
	if cp.BaseBackoff == 0 {
		cp.BaseBackoff = 200 * time.Millisecond
	}
	if cp.MaxBackoff == 0 {
		cp.MaxBackoff = 5 * time.Second
	}
	if cp.Jitter < 0 {
		cp.Jitter = 0
	}
	if cp.FallbackAfterAttempts <= 0 {
		cp.FallbackAfterAttempts = 10
	}
	return cp
}

func nextBackoff(base, max, jitter time.Duration, attempt int) time.Duration {
	// attempt starts at 1
	backoff := base << (attempt - 1) // exponential
	if max > 0 && backoff > max {
		backoff = max
	}
	if jitter > 0 {
		// #nosec G404 (non-crypto jitter is fine here)
		backoff += time.Duration(rand.Int63n(int64(jitter)))
	}
	return backoff
}

// PublishWithOptions publishes an event with robust options.
// Returns a PubAck on success (stream name, sequence, duplicate flag).
func (p *Publisher[T]) PublishWithOptions(ctx context.Context, data T, userOpts PublishOptions) (*jetstream.PubAck, error) {
	start := time.Now()
	var success bool
	defer func() {
		metrics.NATSPublishDuration.
			WithLabelValues(string(p.Stream), string(p.Subject), strconv.FormatBool(success)).
			Observe(time.Since(start).Seconds())
	}()

	// Ensure the stream exists (fail fast on config errors)
	streamInfo, err := p.StreamManager.Stream(ctx, string(p.Stream))
	if err != nil {
		metrics.RecordNATSFailure("unknown", string(p.Subject), err)
		logger.Error("Stream not found for subject", err, "Stream", p.Stream, "Subject", p.Subject)
		return nil, fmt.Errorf("stream %s not found for subject %s: %w", p.Stream, p.Subject, err)
	}

	// Prepare payload (carry your generic Event[T])
	event := Event[T]{
		Subject: p.Subject,
		Data:    data,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
		logger.Error("Failed to marshal event", err, "subject", p.Subject)
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	opts := userOpts.withDefaults(&Publisher[any]{
		Stream:          p.Stream,
		Subject:         p.Subject,
		StreamManager:   p.StreamManager,
		EnableDedup:     p.EnableDedup,
		FallbackStorage: p.FallbackStorage,
	})

	// Optional MaxMsgSize guard
	if opts.MaxMsgSize > 0 && len(payload) > opts.MaxMsgSize {
		return nil, fmt.Errorf("payload too large: %d > %d", len(payload), opts.MaxMsgSize)
	}

	// Build JetStream publish options
	jsOpts := []jetstream.PublishOpt{}
	expectStream := opts.ExpectStreamName
	if expectStream == "" {
		expectStream = string(p.Stream)
	}
	jsOpts = append(jsOpts, jetstream.WithExpectStream(expectStream))

	// Dedupe (also force dedupe if fallback is configured)
	if *opts.EnableDedup || p.FallbackStorage != nil {
		if opts.MsgID == "" {
			opts.MsgID = fmt.Sprintf("%s-%d", p.Subject, time.Now().UnixNano())
		}
		jsOpts = append(jsOpts, jetstream.WithMsgID(opts.MsgID))
	}

	// One-shot publish (no retry loop)
	if !opts.RetryEnabled {
		ack, err := p.StreamManager.JsClient.Publish(ctx, string(p.Subject), payload, jsOpts...)
		if err != nil {
			metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
			logger.Error("Failed to publish event (no retry)", "subject", p.Subject, err)
			return nil, fmt.Errorf("failed to publish event: %w", err)
		}
		success = true
		metrics.RecordNATSPublished(streamInfo.Config.Name, string(p.Subject))
		logger.Info("Published (no retry)", "subject", p.Subject, "stream", ack.Stream, "seq", ack.Sequence, "duplicate", ack.Duplicate)
		return ack, nil
	}

	// Retry path (exponential backoff + jitter; honor ctx)
	attempts := opts.RetryAttempts
	if attempts <= 0 {
		attempts = int(^uint(0) >> 1) // effectively unbounded
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		ack, err := p.StreamManager.JsClient.Publish(ctx, string(p.Subject), payload, jsOpts...)
		if err == nil {
			success = true
			metrics.RecordNATSPublished(streamInfo.Config.Name, string(p.Subject))
			logger.Info("Published (retry ok)", "subject", p.Subject, "attempt", attempt, "stream", ack.Stream, "seq", ack.Sequence, "duplicate", ack.Duplicate)
			return ack, nil
		}

		lastErr = err
		metrics.RecordNATSFailure(streamInfo.Config.Name, string(p.Subject), err)
		logger.Error("Publish failed", "subject", p.Subject, "attempt", attempt, "error", err)

		// Stop if context is done
		if ctx.Err() != nil {
			break
		}

		// Backoff (except after final attempt)
		if attempt < attempts {
			sleep := nextBackoff(opts.BaseBackoff, opts.MaxBackoff, opts.Jitter, attempt)
			select {
			case <-time.After(sleep):
			case <-ctx.Done():
				return nil, fmt.Errorf("publish cancelled: %w", ctx.Err())
			}
		}
	}

	// Fallback (idempotent upsert) after threshold
	if p.FallbackStorage != nil && attempts >= opts.FallbackAfterAttempts {
		filter := bson.M{"msg_id": opts.MsgID}
		update := bson.M{
			"$set": bson.M{
				"stream_name": string(p.Stream),
				"subject":     string(p.Subject),
				"payload":     payload,
				"timestamp":   time.Now(),
			},
			"$inc": bson.M{"attempts": attempts},
		}
		_, ferr := p.FallbackStorage.UpdateOne(ctx, filter, update, mopt.Update().SetUpsert(true))
		if ferr != nil {
			logger.Error("Fallback upsert failed", ferr, "subject", p.Subject, "msgID", opts.MsgID)
		} else {
			logger.Warn("Fallback upserted", "subject", p.Subject, "msgID", opts.MsgID)
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("publish failed without specific error")
	}
	return nil, fmt.Errorf("failed to publish after %d attempts: %w", attempts, lastErr)
}

// EnsureFallbackIndexes creates helpful indexes for the fallback collection.
// Call once at boot if you use Mongo fallback.
func EnsureFallbackIndexes(ctx context.Context, coll *mongo.Collection) error {
	if coll == nil {
		return nil
	}
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "msg_id", Value: 1}},
			Options: mopt.Index().SetUnique(true).SetName("uniq_msg_id"),
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: 1}},
			Options: mopt.Index().SetName("ts_idx"),
		},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
