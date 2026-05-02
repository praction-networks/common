// Package lease provides a small leader-election primitive used to keep
// platform background work (Mongo change-stream watchers, periodic
// sweepers, queue drainers) from running on every replica concurrently.
//
// The interface is deliberately minimal — Acquire + Renew + Release —
// so an implementation can be backed by Redis, Postgres advisory locks,
// or anything else with TTL + atomic check-and-set semantics. The
// RunAsLeader helper composes those primitives into a renew loop that
// keeps a long-running fn under continuous ownership and unwinds it
// cleanly when ownership is lost.
//
// Consumers usually want one of two patterns:
//
//   1. Per-tick claim — short fn, ttl > fn duration, no renewal needed
//      (the existing acs-service sweeper pattern). Just call Acquire,
//      run the work, call Release.
//
//   2. Long-running leader — fn runs for hours/days (Mongo CDC watcher,
//      drainer). Use RunAsLeader; renewal keeps the lease alive while
//      fn runs and cancels fn if the lease is lost mid-flight.
package lease

import (
	"context"
	"errors"
	"time"
)

// Leaser is the contract every backend implements.
//
// All three methods are idempotent at the caller boundary:
//   - Acquire returns (false, nil) when another holder owns the lease;
//     not an error condition.
//   - Renew returns (false, nil) when the lease has been lost (TTL
//     expired or stolen by another holder); caller treats this as
//     "step down."
//   - Release returns ErrNotHeld when the current value doesn't match
//     the holder; informational, not a failure (the net effect — we
//     no longer hold the lease — is correct).
type Leaser interface {
	// Acquire attempts to claim key for holder with the given TTL.
	// Returns (true, nil) on success, (false, nil) when held by
	// someone else, (false, err) on backend failure.
	Acquire(ctx context.Context, key, holder string, ttl time.Duration) (bool, error)

	// Renew extends the TTL on a lease the caller already holds.
	// Returns (true, nil) when still held + TTL extended,
	// (false, nil) when the lease has been lost (expired or taken
	// by another holder), (false, err) on backend failure.
	//
	// Implementations must be atomic: check current holder + extend
	// TTL must happen as a single transactional step (Lua script
	// for Redis, single SQL UPDATE for Postgres) so a concurrent
	// holder change can never race the renewal.
	Renew(ctx context.Context, key, holder string, ttl time.Duration) (bool, error)

	// Release drops the lease only when the current holder matches.
	// Atomic check-and-delete prevents the stale-release race where
	// holder A's expired lease is reclaimed by B, then A's late
	// Release accidentally clears B's fresh lease.
	Release(ctx context.Context, key, holder string) error
}

// ErrNotHeld is returned by Release when the current key value doesn't
// match the expected holder. Informational — the net effect (we don't
// hold the lease) is the desired outcome regardless.
var ErrNotHeld = errors.New("lease not held by this holder")

// ErrNoClient is returned when a leaser was constructed without its
// backend client. Callers should treat this as "can't coordinate —
// single-replica semantics only" and decide whether to proceed.
var ErrNoClient = errors.New("lease: no backend client configured")
