package appError

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/praction-networks/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MapMongoError converts Mongo driver errors into appError types with appropriate HTTP codes
func MapMongoError(err error, userMsg string) error {
	if err == nil {
		return nil
	}
	var we mongo.WriteException
	var bwe mongo.BulkWriteException
	var ce mongo.CommandError

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return New(TimeoutError, userMsg, 504, err)
	case errors.Is(err, context.Canceled):
		return New(RequestCanceled, userMsg, 499, err)
	case errors.As(err, &we):
		if we.HasErrorCode(11000) {
			return New(DuplicateEntityFound, userMsg, 409, err)
		}
		if we.WriteConcernError != nil {
			return New(DBWriteConcernError, userMsg, 503, err)
		}
		return New(DBUpdateError, userMsg, 503, err)
	case errors.As(err, &bwe):
		return New(DBUpdateError, userMsg, 503, err)
	case errors.As(err, &ce):
		switch ce.Code {
		case 50: // ExceededTimeLimit
			return New(TimeoutError, userMsg, 504, err)
		case 91, 11600, 11602, 13435, 189: // retryable server state changes
			return New(DBRetryableError, userMsg, 503, err)
		default:
			return New(InternalServerError, userMsg, 503, err)
		}
	default:
		if mongo.IsDuplicateKeyError(err) {
			return New(DuplicateEntityFound, userMsg, 409, err)
		}
		return New(InternalServerError, userMsg, 500, err)
	}
}

// WithTxnRetry runs a Mongo transaction with exponential backoff on transient errors
// WithTxnRetry intentionally omitted for now (labels not used by common errors)

// LogMongoError adds consistent structured logging for mongo operations
func LogMongoError(op, coll string, err error, fields ...any) {
	args := append([]any{"op", op, "collection", coll, "error", err}, fields...)
	logger.Error("mongo op failed", args...)
}

// Context keys commonly propagated through middleware
const (
	CtxKeyRequestID = "request_id"
	CtxKeyTenantID  = "tenant_id"
	CtxKeyUserID    = "user_id"
)

// addContextFields appends common context identifiers to logs if present
func addContextFields(ctx context.Context, fields []any) []any {
	if ctx == nil {
		return fields
	}
	if v := ctx.Value(CtxKeyRequestID); v != nil {
		fields = append(fields, "request_id", v)
	}
	if v := ctx.Value(CtxKeyTenantID); v != nil {
		fields = append(fields, "tenant_id", v)
	}
	if v := ctx.Value(CtxKeyUserID); v != nil {
		fields = append(fields, "user_id", v)
	}
	return fields
}

// LogMongoErrorWithCtx includes context identifiers when available
func LogMongoErrorWithCtx(ctx context.Context, op, coll string, err error, fields ...any) {
	args := append([]any{"op", op, "collection", coll, "error", err}, fields...)
	args = addContextFields(ctx, args)
	logger.Error("mongo op failed", args...)
}

// IsRetryableMongoError determines whether an error is transient and worth retrying
func IsRetryableMongoError(err error) bool {
	if err == nil {
		return false
	}
	// Do not retry canceled contexts
	if errors.Is(err, context.Canceled) {
		return false
	}
	// Retry timeouts
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// Retry common network errors
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout() || ne.Temporary()
	}
	// Retry select command errors
	var ce mongo.CommandError
	if errors.As(err, &ce) {
		switch ce.Code { // a small, pragmatic allowlist
		case 50: // ExceededTimeLimit
			return true
		case 91: // ShutdownInProgress
			return true
		case 11600, 11602: // InterruptedAtShutdown / InterruptedDueToReplStateChange
			return true
		case 13435, 189: // NotMasterNoSlaveOk / PrimarySteppedDown (legacy names)
			return true
		default:
			return false
		}
	}
	// Retry write concern issues
	var we mongo.WriteException
	if errors.As(err, &we) {
		if we.WriteConcernError != nil {
			return true
		}
	}
	return false
}

// WithRetry executes fn with bounded exponential backoff for retryable errors
func WithRetry(ctx context.Context, op string, maxAttempts int, baseDelay time.Duration, fn func(context.Context) error) error {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if baseDelay <= 0 {
		baseDelay = 50 * time.Millisecond
	}
	var err error
	delay := baseDelay
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Respect caller context
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn(ctx)
		if err == nil {
			return nil
		}
		if attempt == maxAttempts || !IsRetryableMongoError(err) {
			return err
		}
		// Log retry and sleep before next attempt
		logArgs := addContextFields(ctx, []any{"op", op, "attempt", attempt, "next_delay", delay, "error", err})
		logger.Warn("mongo transient error, will retry", logArgs...)
		time.Sleep(delay)
		// Exponential backoff with cap
		if delay < 2*time.Second {
			delay *= 2
		}
	}
	return err
}

// --- Combined repository execution helpers ---

const (
	defaultReadTimeout   = 5 * time.Second
	defaultWriteTimeout  = 5 * time.Second
	defaultRetryAttempts = 4
	defaultBaseDelay     = 75 * time.Millisecond
)

// CountDocuments wraps CountDocuments with timeout, retry, logging, and error mapping.
func CountDocuments(ctx context.Context, coll *mongo.Collection, filter interface{}) (int64, error) {
	cctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var total int64
	err := WithRetry(cctx, "CountDocuments", defaultRetryAttempts, defaultBaseDelay, func(rc context.Context) error {
		var e error
		total, e = coll.CountDocuments(rc, filter)
		return e
	})
	if err != nil {
		LogMongoErrorWithCtx(ctx, "CountDocuments", coll.Name(), err, "filter", filter)
		return 0, MapMongoError(err, "Failed to count documents")
	}
	return total, nil
}

// FindAll executes a Find and streams decode into a typed slice with retries and error mapping.
func FindAll[T any](ctx context.Context, coll *mongo.Collection, filter interface{}, findOpts *options.FindOptions) ([]T, error) {
	cctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var cur *mongo.Cursor
	err := WithRetry(cctx, "Find", defaultRetryAttempts, defaultBaseDelay, func(rc context.Context) error {
		var e error
		if findOpts != nil {
			cur, e = coll.Find(rc, filter, findOpts)
		} else {
			cur, e = coll.Find(rc, filter)
		}
		return e
	})
	if err != nil {
		LogMongoErrorWithCtx(ctx, "Find", coll.Name(), err, "filter", filter)
		return nil, MapMongoError(err, "Failed to query documents")
	}
	defer cur.Close(cctx)

	var out []T
	for cur.Next(cctx) {
		var item T
		if e := cur.Decode(&item); e != nil {
			LogMongoErrorWithCtx(ctx, "Cursor.Decode", coll.Name(), e)
			return nil, MapMongoError(e, "Failed to decode document")
		}
		out = append(out, item)
	}
	if e := cur.Err(); e != nil {
		LogMongoErrorWithCtx(ctx, "Cursor.Err", coll.Name(), e)
		return nil, MapMongoError(e, "Cursor iteration error")
	}
	return out, nil
}

// InsertOne wraps InsertOne with timeout, retry, logging, and error mapping.
func InsertOne(ctx context.Context, coll *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	cctx, cancel := context.WithTimeout(ctx, defaultWriteTimeout)
	defer cancel()

	var res *mongo.InsertOneResult
	err := WithRetry(cctx, "InsertOne", defaultRetryAttempts, defaultBaseDelay, func(rc context.Context) error {
		var e error
		res, e = coll.InsertOne(rc, document)
		return e
	})
	if err != nil {
		LogMongoErrorWithCtx(ctx, "InsertOne", coll.Name(), err)
		return nil, MapMongoError(err, "Failed to insert document")
	}
	return res, nil
}

// UpdateOne wraps UpdateOne with timeout, retry, logging, and error mapping.
func UpdateOne(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	cctx, cancel := context.WithTimeout(ctx, defaultWriteTimeout)
	defer cancel()

	var res *mongo.UpdateResult
	err := WithRetry(cctx, "UpdateOne", defaultRetryAttempts, defaultBaseDelay, func(rc context.Context) error {
		var e error
		if len(opts) > 0 && opts[0] != nil {
			res, e = coll.UpdateOne(rc, filter, update, opts[0])
		} else {
			res, e = coll.UpdateOne(rc, filter, update)
		}
		return e
	})
	if err != nil {
		LogMongoErrorWithCtx(ctx, "UpdateOne", coll.Name(), err, "filter", filter)
		return nil, MapMongoError(err, "Failed to update document")
	}
	return res, nil
}

// DeleteOne wraps DeleteOne with timeout, retry, logging, and error mapping.
func DeleteOne(ctx context.Context, coll *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	cctx, cancel := context.WithTimeout(ctx, defaultWriteTimeout)
	defer cancel()

	var res *mongo.DeleteResult
	err := WithRetry(cctx, "DeleteOne", defaultRetryAttempts, defaultBaseDelay, func(rc context.Context) error {
		var e error
		res, e = coll.DeleteOne(rc, filter)
		return e
	})
	if err != nil {
		LogMongoErrorWithCtx(ctx, "DeleteOne", coll.Name(), err, "filter", filter)
		return nil, MapMongoError(err, "Failed to delete document")
	}
	return res, nil
}
