package lease

import (
	"context"
	"time"

	"github.com/praction-networks/common/logger"
)

// RunAsLeader runs fn while holding the lease named by key. The call
// blocks until fn returns, the parent context is cancelled, or the
// lease is lost mid-renewal.
//
// Lifecycle:
//
//  1. Acquire(key, holderID, ttl). If another holder owns it, return
//     (false, nil) — caller is expected to retry on its own cadence.
//  2. Spawn fn in a goroutine under a derived context that this
//     function will cancel on lease loss.
//  3. Renewal loop: every renewInterval, call Renew. On Renew=false
//     (lease lost) or Renew=err (backend blip beyond a few attempts),
//     cancel the fn-context so fn observes the loss and unwinds.
//  4. When fn returns or ctx is cancelled, best-effort Release using
//     a fresh background context (don't fail Release because the
//     parent ctx already expired).
//
// Returned values:
//   - (true, fnErr) — we held the lease and ran fn; fnErr is fn's
//     return value (nil on success).
//   - (false, nil) — another holder owned the lease; fn was not run.
//   - (false, err) — backend error during initial acquire; fn was
//     not run.
//
// renewInterval should be ≤ ttl/2; if the caller passes 0 or a value
// >= ttl, this defaults to ttl/3 so renewal lands well before expiry
// even under cluster-consensus jitter.
//
// The fn-context cancellation on lease loss is the primary safety
// rail: any fn that touches shared state (DB writes, NATS publishes)
// must respect ctx.Done() promptly so two replicas can't both think
// they're leader during a TTL-expiry handoff.
func RunAsLeader(
	ctx context.Context,
	leaser Leaser,
	key, holderID string,
	ttl, renewInterval time.Duration,
	fn func(context.Context) error,
) (bool, error) {
	if leaser == nil {
		return false, ErrNoClient
	}
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if renewInterval <= 0 || renewInterval >= ttl {
		renewInterval = ttl / 3
	}

	acquired, err := leaser.Acquire(ctx, key, holderID, ttl)
	if err != nil {
		return false, err
	}
	if !acquired {
		return false, nil
	}

	fnCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	var fnErr error
	fnDone := make(chan struct{})
	go func() {
		defer close(fnDone)
		fnErr = fn(fnCtx)
	}()

	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-fnDone:
			releaseHeld(leaser, key, holderID)
			return true, fnErr

		case <-ctx.Done():
			// Parent ctx cancelled (shutdown). Cancel fn,
			// wait for unwind, release.
			cancelFn()
			<-fnDone
			releaseHeld(leaser, key, holderID)
			return true, fnErr

		case <-ticker.C:
			ok, rerr := leaser.Renew(ctx, key, holderID, ttl)
			if rerr != nil {
				// Backend blip — keep running fn; next tick
				// will retry. The next failed tick (or TTL
				// expiry on the backend side) will surface
				// the loss. Don't pre-emptively cancel fn on
				// a single backend hiccup.
				logger.Warn("lease renewal failed (will retry next tick)",
					"key", key, "err", rerr)
				continue
			}
			if !ok {
				// Lost the lease — cancel fn and unwind.
				logger.Warn("lease lost during renewal — cancelling fn",
					"key", key, "holder", holderID)
				cancelFn()
				<-fnDone
				// Don't try to Release — we no longer hold it.
				return true, fnErr
			}
		}
	}
}

// releaseHeld is a best-effort Release. Uses a fresh background ctx
// with a short timeout so the release still goes through when the
// caller's ctx is already cancelled (the common shutdown case).
// ErrNotHeld is informational — log at Debug.
func releaseHeld(leaser Leaser, key, holderID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := leaser.Release(ctx, key, holderID); err != nil {
		if err == ErrNotHeld {
			logger.Debug("lease release: not held (already expired/taken)",
				"key", key, "holder", holderID)
			return
		}
		logger.Warn("lease release failed", "key", key, "holder", holderID, "err", err)
	}
}
