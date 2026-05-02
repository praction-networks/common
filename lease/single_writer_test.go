package lease_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/praction-networks/common/lease"
	"github.com/praction-networks/common/logger"
)

// TestMain initializes the package-level logger that RunAsLeader uses
// for renew-failure / lease-loss diagnostics. The common logger panics
// on first call when uninitialized, which matches the production
// contract (services init it at startup) — so tests follow suit.
func TestMain(m *testing.M) {
	if err := logger.InitializeLogger(logger.LoggerConfig{LogLevel: "debug"}); err != nil {
		panic("test logger init: " + err.Error())
	}
	os.Exit(m.Run())
}

// fakeLeaser is a minimal in-memory Leaser for testing RunAsLeader.
// Concurrent-safe; supports forcing a renewal failure or a renewal
// "stolen by another holder" event so we can drive the lifecycle
// branches without running a real Redis.
type fakeLeaser struct {
	mu             sync.Mutex
	holder         string
	renewFailNext  atomic.Bool // make next Renew return (false, nil) — lease lost
	renewErrNext   atomic.Pointer[error]
	acquireErrNext atomic.Pointer[error]
	renewCalls     atomic.Int32
	releaseCalls   atomic.Int32
}

func (f *fakeLeaser) Acquire(_ context.Context, _ /*key*/, holder string, _ time.Duration) (bool, error) {
	if errp := f.acquireErrNext.Swap(nil); errp != nil {
		return false, *errp
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.holder != "" && f.holder != holder {
		return false, nil
	}
	f.holder = holder
	return true, nil
}

func (f *fakeLeaser) Renew(_ context.Context, _ /*key*/, holder string, _ time.Duration) (bool, error) {
	f.renewCalls.Add(1)
	if errp := f.renewErrNext.Swap(nil); errp != nil {
		return false, *errp
	}
	if f.renewFailNext.Swap(false) {
		f.mu.Lock()
		f.holder = "" // simulate "stolen by another holder"
		f.mu.Unlock()
		return false, nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.holder == holder, nil
}

func (f *fakeLeaser) Release(_ context.Context, _ /*key*/, holder string) error {
	f.releaseCalls.Add(1)
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.holder != holder {
		return lease.ErrNotHeld
	}
	f.holder = ""
	return nil
}

func TestRunAsLeader_HappyPath(t *testing.T) {
	leaser := &fakeLeaser{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var ran atomic.Bool
	held, err := lease.RunAsLeader(
		ctx, leaser, "key", "holder-A",
		200*time.Millisecond, 50*time.Millisecond,
		func(ctx context.Context) error {
			ran.Store(true)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !held {
		t.Fatal("expected held=true on successful acquire")
	}
	if !ran.Load() {
		t.Fatal("fn was not invoked")
	}
	if leaser.releaseCalls.Load() != 1 {
		t.Fatalf("expected 1 release call, got %d", leaser.releaseCalls.Load())
	}
}

func TestRunAsLeader_NotAcquired(t *testing.T) {
	leaser := &fakeLeaser{holder: "someone-else"}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var ran atomic.Bool
	held, err := lease.RunAsLeader(
		ctx, leaser, "key", "holder-A",
		200*time.Millisecond, 50*time.Millisecond,
		func(ctx context.Context) error {
			ran.Store(true)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if held {
		t.Fatal("expected held=false when another holder owns the lease")
	}
	if ran.Load() {
		t.Fatal("fn must not run when lease was not acquired")
	}
	if leaser.releaseCalls.Load() != 0 {
		t.Fatal("must not Release a lease we never held")
	}
}

func TestRunAsLeader_LeaseLostMidRun_CancelsFn(t *testing.T) {
	leaser := &fakeLeaser{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Force the first renewal call to report "lost" — we expect the
	// fn-context to be cancelled and fn to observe the cancellation.
	leaser.renewFailNext.Store(true)

	fnCancelled := make(chan struct{})
	held, err := lease.RunAsLeader(
		ctx, leaser, "key", "holder-A",
		100*time.Millisecond, 20*time.Millisecond,
		func(fnCtx context.Context) error {
			<-fnCtx.Done()
			close(fnCancelled)
			return fnCtx.Err()
		},
	)
	if !held {
		t.Fatal("expected held=true (we did acquire initially)")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected fn to return context.Canceled, got %v", err)
	}
	select {
	case <-fnCancelled:
		// good — fn observed cancellation.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("fn was not cancelled within 500ms of lease loss")
	}
	// We don't Release after losing the lease — Release would be a
	// no-op (ErrNotHeld) and could mask a real bug if added later.
	if leaser.releaseCalls.Load() != 0 {
		t.Fatalf("expected 0 release calls after lease loss, got %d", leaser.releaseCalls.Load())
	}
}

func TestRunAsLeader_ParentCtxCancelled_Unwinds(t *testing.T) {
	leaser := &fakeLeaser{}
	ctx, cancel := context.WithCancel(context.Background())

	fnStarted := make(chan struct{})
	doneCh := make(chan error, 1)
	go func() {
		_, err := lease.RunAsLeader(
			ctx, leaser, "key", "holder-A",
			500*time.Millisecond, 100*time.Millisecond,
			func(fnCtx context.Context) error {
				close(fnStarted)
				<-fnCtx.Done()
				return fnCtx.Err()
			},
		)
		doneCh <- err
	}()

	<-fnStarted
	cancel()

	select {
	case err := <-doneCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled propagated from fn, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunAsLeader did not unwind within 2s of parent cancel")
	}
	if leaser.releaseCalls.Load() != 1 {
		t.Fatalf("expected 1 release call on graceful unwind, got %d", leaser.releaseCalls.Load())
	}
}

func TestRunAsLeader_RenewBackendBlipDoesNotCancel(t *testing.T) {
	// A single Renew error (backend hiccup) should NOT cancel fn —
	// the lease may still be valid on the Redis side; next tick will
	// retry. Cancelling on every transient error would cause spurious
	// failovers under network jitter.
	leaser := &fakeLeaser{}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	transientErr := errors.New("simulated backend blip")
	leaser.renewErrNext.Store(&transientErr)

	fnCtxObserved := make(chan struct{}, 1)
	held, err := lease.RunAsLeader(
		ctx, leaser, "key", "holder-A",
		100*time.Millisecond, 25*time.Millisecond,
		func(fnCtx context.Context) error {
			// Run a few ticks past the blip, then exit cleanly.
			t := time.NewTimer(150 * time.Millisecond)
			defer t.Stop()
			select {
			case <-t.C:
				fnCtxObserved <- struct{}{}
				return nil
			case <-fnCtx.Done():
				return fnCtx.Err()
			}
		},
	)
	if !held {
		t.Fatal("expected held=true")
	}
	if err != nil {
		t.Fatalf("expected nil fn error after transient blip, got %v", err)
	}
	select {
	case <-fnCtxObserved:
		// good — fn ran past the blip and exited cleanly.
	default:
		t.Fatal("fn did not complete naturally — was probably cancelled spuriously")
	}
}

func TestRunAsLeader_NilLeaser(t *testing.T) {
	held, err := lease.RunAsLeader(
		context.Background(), nil, "key", "holder-A",
		100*time.Millisecond, 25*time.Millisecond,
		func(ctx context.Context) error { return nil },
	)
	if held {
		t.Fatal("expected held=false with nil leaser")
	}
	if !errors.Is(err, lease.ErrNoClient) {
		t.Fatalf("expected ErrNoClient, got %v", err)
	}
}

func TestRunAsLeader_AcquireError(t *testing.T) {
	backendErr := errors.New("backend unavailable")
	leaser := &fakeLeaser{}
	leaser.acquireErrNext.Store(&backendErr)

	held, err := lease.RunAsLeader(
		context.Background(), leaser, "key", "holder-A",
		100*time.Millisecond, 25*time.Millisecond,
		func(ctx context.Context) error { return nil },
	)
	if held {
		t.Fatal("expected held=false on acquire error")
	}
	if !errors.Is(err, backendErr) {
		t.Fatalf("expected backendErr, got %v", err)
	}
}
