package events

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// -----------------------------------------------------------------------------
// Delivery modes
// -----------------------------------------------------------------------------

type DeliveryGuarantee int

const (
	// Strict at-least-once: retries + early Mongo fallback + dedupe + expect-stream guard
	DeliveryGuaranteed DeliveryGuarantee = iota + 1
	// A few retries, no fallback. Good default for non-critical events.
	DeliveryPreferred
	// Single try, no fallback, no dedupe. Good for high-volume logs/metrics where drops are OK.
	DeliveryBestEffort
)

// -----------------------------------------------------------------------------
// Back-compat entry point (kept as-is)
// -----------------------------------------------------------------------------

// Publish keeps your original signature but internally uses PublishWithOptions.
// Prefer calling PublishWithOptions (or the helpers below) for fine-grained control.
func (p *Publisher[T]) Publish(ctx context.Context, data T, msgID string, retry RetryConfig) error {
	opts := PublishOptions{
		Guarantee:     DeliveryPreferred,
		MsgID:         msgID,
		RetryEnabled:  retry.Enabled,
		RetryAttempts: retry.Attempts,
		BaseBackoff:   retry.RetryInterval,
		MaxBackoff:    retry.RetryInterval,
		Jitter:        retry.RetryInterval / 10,
		// FallbackAfterAttempts is ignored unless FallbackEnabled is true
		FallbackAfterAttempts: 10,
	}
	_, err := p.PublishWithOptions(ctx, data, opts)
	return err
}

// Convenience helpers
func (p *Publisher[T]) PublishGuaranteed(ctx context.Context, data T, msgID string) (*jetstream.PubAck, error) {
	return p.PublishWithOptions(ctx, data, PublishOptions{
		Guarantee: DeliveryGuaranteed,
		MsgID:     msgID,
	})
}
func (p *Publisher[T]) PublishPreferred(ctx context.Context, data T, msgID string) (*jetstream.PubAck, error) {
	return p.PublishWithOptions(ctx, data, PublishOptions{
		Guarantee: DeliveryPreferred,
		MsgID:     msgID,
	})
}
func (p *Publisher[T]) PublishBestEffort(ctx context.Context, data T) (*jetstream.PubAck, error) {
	return p.PublishWithOptions(ctx, data, PublishOptions{
		Guarantee: DeliveryBestEffort,
		// no MsgID on purpose
	})
}

// -----------------------------------------------------------------------------
// Options-rich publish API
// -----------------------------------------------------------------------------

type PublishOptions struct {
	// Mode
	Guarantee DeliveryGuarantee

	// Dedupe / expectations
	EnableDedup      *bool  // nil => defaults from mode / publisher
	MsgID            string // stable id for dedupe; generated if empty and dedupe enabled
	ExpectStreamName string // if empty, defaults to publisher's Stream
	UseExpectStream  *bool  // nil => defaults from mode

	// Retry & backoff
	RetryEnabled  bool
	RetryAttempts int           // total attempts incl. first; <=0 => practically infinite
	BaseBackoff   time.Duration // e.g., 200ms
	MaxBackoff    time.Duration // e.g., 5s
	Jitter        time.Duration // e.g., 200ms

	// Fallback (Mongo)
	FallbackEnabled       *bool // nil => defaults from mode (true for Guaranteed if fallback configured)
	FallbackAfterAttempts int   // upsert to Mongo after this many tries (when enabled)

	// Optional payload guard if your stream uses MaxMsgSize (defaults to 2MB)
	MaxMsgSize int // 0 = use default (2MB)
}

func (o *PublishOptions) withDefaults(pub *Publisher[any]) PublishOptions {
	cp := *o

	// Defaults by delivery mode
	switch cp.Guarantee {
	case DeliveryBestEffort:
		// fire-and-forget-ish (still JetStream)
		if cp.EnableDedup == nil {
			v := false
			cp.EnableDedup = &v
		}
		if !cp.RetryEnabled {
			cp.RetryAttempts = 1
		}
		if cp.UseExpectStream == nil {
			v := false // permissive; caller can enable
			cp.UseExpectStream = &v
		}
		if cp.FallbackEnabled == nil {
			v := false
			cp.FallbackEnabled = &v
		}
	case DeliveryPreferred:
		// a few retries, no fallback
		if cp.EnableDedup == nil {
			v := pub.EnableDedup // inherit from publisher
			cp.EnableDedup = &v
		}
		if !cp.RetryEnabled {
			cp.RetryEnabled = true
			if cp.RetryAttempts <= 0 {
				cp.RetryAttempts = 5
			}
			if cp.BaseBackoff == 0 {
				cp.BaseBackoff = 200 * time.Millisecond
			}
			if cp.MaxBackoff == 0 {
				cp.MaxBackoff = 5 * time.Second
			}
			if cp.Jitter == 0 {
				cp.Jitter = 200 * time.Millisecond
			}
		}
		if cp.UseExpectStream == nil {
			v := true
			cp.UseExpectStream = &v
		}
		if cp.FallbackEnabled == nil {
			v := false
			cp.FallbackEnabled = &v
		}
	default: // DeliveryGuaranteed
		if cp.EnableDedup == nil {
			v := true
			cp.EnableDedup = &v
		}
		if !cp.RetryEnabled {
			cp.RetryEnabled = true
			if cp.RetryAttempts <= 0 {
				cp.RetryAttempts = 12
			}
			if cp.BaseBackoff == 0 {
				cp.BaseBackoff = 200 * time.Millisecond
			}
			if cp.MaxBackoff == 0 {
				cp.MaxBackoff = 5 * time.Second
			}
			if cp.Jitter == 0 {
				cp.Jitter = 200 * time.Millisecond
			}
		}
		if cp.UseExpectStream == nil {
			v := true
			cp.UseExpectStream = &v
		}
		if cp.FallbackEnabled == nil {
			// Enable fallback by default if wired
			v := pub.FallbackStorage != nil
			cp.FallbackEnabled = &v
		}
		if cp.FallbackAfterAttempts <= 0 {
			cp.FallbackAfterAttempts = 3
		}
	}

	// General defaults
	if cp.MaxBackoff == 0 {
		cp.MaxBackoff = 5 * time.Second
	}
	if cp.BaseBackoff == 0 {
		cp.BaseBackoff = 200 * time.Millisecond
	}
	if cp.Jitter < 0 {
		cp.Jitter = 0
	}
	if cp.RetryEnabled && cp.RetryAttempts <= 0 {
		cp.RetryAttempts = int(^uint(0) >> 1) // int max
	}
	if cp.MaxMsgSize == 0 {
		// Default 2 MB JSON guard (tune to your StreamConfig.MaxMsgSize if set)
		cp.MaxMsgSize = 2 * 1024 * 1024
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
		// Go 1.20+: no need to Seed; default Source is fine
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
	event := Event[T]{Subject: p.Subject, Data: data}
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
	if opts.UseExpectStream != nil && *opts.UseExpectStream {
		jsOpts = append(jsOpts, jetstream.WithExpectStream(expectStream))
	}

	// Dedupe (generate MsgID from payload bytes if empty)
	if opts.EnableDedup != nil && *opts.EnableDedup {
		if opts.MsgID == "" {
			sum := sha256.Sum256(payload)
			trunc := sum[:16] // 16 bytes (32 hex chars) is plenty
			opts.MsgID = fmt.Sprintf("%s:%s", p.Subject, hex.EncodeToString(trunc))
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
	var lastErr error
	actualAttempts := 0

	for attempt := 1; attempt <= attempts; attempt++ {
		actualAttempts = attempt
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

	// Fallback (idempotent upsert) after threshold (only when enabled)
	if p.FallbackStorage != nil && opts.FallbackEnabled != nil && *opts.FallbackEnabled && actualAttempts >= opts.FallbackAfterAttempts {
		docID := fmt.Sprintf("%s|%s", p.Stream, opts.MsgID)
		lastErrStr := "publish failed"
		if lastErr != nil {
			lastErrStr = lastErr.Error()
		}

		filter := bson.M{"_id": docID}
		update := bson.M{
			"$setOnInsert": bson.M{
				"stream_name": string(p.Stream),
				"subject":     string(p.Subject),
				"payload":     payload, // keep original payload on first insert
				"attempts":    0,
			},
			"$set": bson.M{
				"timestamp":  time.Now(), // last attempt time
				"last_error": lastErrStr,
			},
			"$inc": bson.M{"attempts": actualAttempts}, // record how many we already tried here
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
	return nil, fmt.Errorf("failed to publish after %d attempts: %w", actualAttempts, lastErr)
}

// EnsureFallbackIndexes creates helpful indexes for the fallback collection.
// Call once at boot if you use Mongo fallback.
func EnsureFallbackIndexes(ctx context.Context, coll *mongo.Collection) error {
	if coll == nil {
		return nil
	}
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "timestamp", Value: 1}},
			Options: mopt.Index().SetName("ts_idx"),
		},
		// _id is unique by default (we use "<Stream>|<MsgID>")
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
