// Package mongolease provides a Mongo-backed implementation of
// lease.Leaser. Every Go service in the platform already requires
// Mongo, so this is the universal fallback for services that do not
// pull in go-redis just to coordinate change-stream watchers and
// schedulers.
//
// Storage shape — one document per lease key:
//
//	{
//	  _id:       "<lease key>",
//	  holder:    "<holder id>",
//	  expiresAt: <ISO timestamp>,
//	  updatedAt: <ISO timestamp>,
//	}
//
// Atomicity is provided by Mongo's per-document FindOneAndUpdate /
// UpdateOne / DeleteOne — each operation is single-document and its
// filter+update is evaluated as a single step, so concurrent acquires
// against the same key cannot both succeed.
//
// EnsureIndexes installs a TTL index on `expiresAt` so dead leases are
// removed in the background by Mongo's TTL monitor (60-second cycle).
// The TTL index is purely for collection-size hygiene — Acquire's
// `expiresAt < now` filter is what provides the take-over-on-expiry
// correctness, so a lagging TTL monitor does not break leasing.
package mongolease

import (
	"context"
	"time"

	"github.com/praction-networks/common/lease"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Leaser implements lease.Leaser against a Mongo collection.
type Leaser struct {
	coll *mongo.Collection
}

// New constructs a Mongo-backed leaser. A nil collection is allowed
// so DI graphs that haven't wired the lease collection yet still
// compile; every method on the zero-value returns lease.ErrNoClient.
func New(coll *mongo.Collection) *Leaser {
	return &Leaser{coll: coll}
}

// EnsureIndexes installs the TTL index on `expiresAt` plus the
// holder index used by Renew/Release filters. Idempotent — calling
// repeatedly is a no-op once the indexes exist. Run at service
// startup; a missing TTL index will not break correctness but the
// _leases collection will grow unbounded over time.
func (l *Leaser) EnsureIndexes(ctx context.Context) error {
	if l == nil || l.coll == nil {
		return lease.ErrNoClient
	}
	_, err := l.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// TTL — 0 seconds means "expire when expiresAt timestamp is
			// in the past." Mongo's TTL monitor sweeps every 60s.
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("expiresAt_ttl"),
		},
	})
	return err
}

// Acquire claims the lease for holder. Uses FindOneAndUpdate with a
// filter that matches either "we already hold it" (renew-on-acquire)
// or "lease has expired" — the upsert path inserts a fresh doc when
// none exists. A duplicate-key error means the document exists,
// belongs to another holder, and has not expired — return (false, nil).
func (l *Leaser) Acquire(ctx context.Context, key, holder string, ttl time.Duration) (bool, error) {
	if l == nil || l.coll == nil {
		return false, lease.ErrNoClient
	}
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)

	filter := bson.M{
		"_id": key,
		"$or": []bson.M{
			{"holder": holder},
			{"expiresAt": bson.M{"$lt": now}},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"holder":    holder,
			"expiresAt": expiresAt,
			"updatedAt": now,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var doc struct {
		Holder string `bson:"holder"`
	}
	err := l.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Document exists, owned by a different holder, not expired.
			return false, nil
		}
		return false, err
	}
	// Should always equal holder on success, but check defensively in
	// case of replication anomalies under heavy contention.
	return doc.Holder == holder, nil
}

// Renew extends the TTL only when the caller still owns the lease.
// Returns (false, nil) when MatchedCount is zero — either the lease
// expired and was reclaimed, or it was deleted. The caller should
// treat that as "step down."
func (l *Leaser) Renew(ctx context.Context, key, holder string, ttl time.Duration) (bool, error) {
	if l == nil || l.coll == nil {
		return false, lease.ErrNoClient
	}
	now := time.Now().UTC()
	res, err := l.coll.UpdateOne(ctx,
		bson.M{"_id": key, "holder": holder},
		bson.M{"$set": bson.M{
			"expiresAt": now.Add(ttl),
			"updatedAt": now,
		}},
	)
	if err != nil {
		return false, err
	}
	return res.MatchedCount > 0, nil
}

// Release deletes the lease document when the caller is the holder.
// ErrNotHeld is informational — the net effect (we don't hold the
// lease) is correct regardless.
func (l *Leaser) Release(ctx context.Context, key, holder string) error {
	if l == nil || l.coll == nil {
		return lease.ErrNoClient
	}
	res, err := l.coll.DeleteOne(ctx, bson.M{"_id": key, "holder": holder})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return lease.ErrNotHeld
	}
	return nil
}
