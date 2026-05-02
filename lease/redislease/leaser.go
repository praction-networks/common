// Package redislease provides a Redis-backed implementation of
// lease.Leaser. Kept in a sub-package so the bare lease package
// (interface + RunAsLeader) stays free of the go-redis dependency —
// services that don't need a Redis-backed lease can use a different
// backend (Postgres advisory locks, or a test fake) without paying
// the import cost.
//
// The Lua scripts implementing Renew + Release are atomic at the
// Redis side: check current holder + extend/delete must happen as a
// single transactional step so a TTL-expiry + concurrent re-acquire
// can never race the renewal/release.
package redislease

import (
	"context"
	"time"

	"github.com/praction-networks/common/lease"
	"github.com/redis/go-redis/v9"
)

// Leaser implements lease.Leaser against a go-redis client.
type Leaser struct {
	client *redis.Client
}

// New constructs a Redis-backed leaser. A nil client is allowed so DI
// graphs without Redis can still compile; every method on the
// zero-value just returns lease.ErrNoClient.
func New(client *redis.Client) *Leaser {
	return &Leaser{client: client}
}

// renewScript atomically extends the TTL only when the caller still
// owns the lease. Returns 1 on success, 0 on holder-mismatch (lease
// expired or taken by another holder).
const renewScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
    return 0
end
`

// releaseScript atomically deletes the key only when the caller is
// still the holder. Guards against the stale-release race where
// holder A's lease expired, holder B took it, then A's late Release
// would have cleared B's fresh lease.
const releaseScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`

// Acquire uses SET key holder NX EX ttl — the canonical single-RTT
// claim. TTL is enforced by Redis: if holder crashes before Release,
// the lease expires naturally.
func (l *Leaser) Acquire(ctx context.Context, key, holder string, ttl time.Duration) (bool, error) {
	if l == nil || l.client == nil {
		return false, lease.ErrNoClient
	}
	ok, err := l.client.SetNX(ctx, key, holder, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Renew extends the TTL only if the caller still holds the lease.
// Returns (false, nil) when the lease was lost (expired or stolen);
// the caller should treat that as "step down."
func (l *Leaser) Renew(ctx context.Context, key, holder string, ttl time.Duration) (bool, error) {
	if l == nil || l.client == nil {
		return false, lease.ErrNoClient
	}
	res, err := l.client.Eval(ctx, renewScript, []string{key}, holder, ttl.Milliseconds()).Result()
	if err != nil {
		return false, err
	}
	n, _ := res.(int64)
	return n == 1, nil
}

// Release deletes the key only when its value matches the holder.
// Eval-based for atomicity; ErrNotHeld is informational (the net
// effect — we don't hold the lease — is correct).
func (l *Leaser) Release(ctx context.Context, key, holder string) error {
	if l == nil || l.client == nil {
		return lease.ErrNoClient
	}
	res, err := l.client.Eval(ctx, releaseScript, []string{key}, holder).Result()
	if err != nil {
		return err
	}
	n, _ := res.(int64)
	if n == 0 {
		return lease.ErrNotHeld
	}
	return nil
}
