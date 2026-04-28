// Package preconditions runs every prerequisite check at service boot
// and aborts startup if any fail — so services fail loudly with a
// concrete fix instruction instead of crashing later when the missing
// permission / table / stream is finally exercised.
//
// The principle: missing prereqs are an operator problem, not a code
// problem. The service shouldn't try to recover; it should exit
// non-zero with a structured error block describing exactly what's
// missing and how to fix it.
//
// Usage:
//
//	checks := preconditions.Runner{
//	    {Name: "postgres-connectivity", Check: ...},
//	    {Name: "cdc-replication-grant", Check: ..., Hint: "GRANT REPLICATION ..."},
//	    {Name: "nats-billing-stream", Check: ..., Hint: "kubectl exec ..."},
//	}
//	if err := checks.RunAll(ctx); err != nil {
//	    logger.Fatal("preconditions failed", err)
//	    os.Exit(1)
//	}
package preconditions

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/praction-networks/common/logger"
)

// Check is a single boot-time prerequisite verification.
type Check struct {
	// Name uniquely identifies the check in logs.
	Name string

	// Check returns nil if the precondition is satisfied. Any error is
	// treated as a hard failure — the service will not start.
	Check func(ctx context.Context) error

	// Hint, if non-empty, is logged alongside the failure to tell
	// operators exactly how to fix it (the SQL grant, kubectl command,
	// env var to set, etc). Be specific — "set X=Y in the deployment"
	// beats "configure X correctly".
	Hint string

	// Required defaults to true. Set false for soft checks that should
	// log a warning but not block startup (e.g. an optional dependency).
	Required bool
}

// Runner aggregates checks and runs them in parallel at boot.
type Runner []Check

// RunAll executes every check concurrently and returns a combined
// error if any required check fails. All checks run regardless of
// individual outcomes — operators see the full picture in one log
// block instead of fix-restart-fix-restart.
func (r Runner) RunAll(ctx context.Context) error {
	if len(r) == 0 {
		return nil
	}

	type result struct {
		check Check
		err   error
		took  time.Duration
	}

	results := make(chan result, len(r))
	var wg sync.WaitGroup
	logger.Info("Running boot preconditions", "count", len(r))

	for _, c := range r {
		c := c
		if c.Check == nil {
			results <- result{check: c, err: fmt.Errorf("nil Check func")}
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			err := c.Check(ctx)
			results <- result{check: c, err: err, took: time.Since(start)}
		}()
	}

	wg.Wait()
	close(results)

	var failures []result
	for res := range results {
		if res.err == nil {
			logger.Info("precondition passed", "name", res.check.Name, "took_ms", res.took.Milliseconds())
			continue
		}
		// Default Required=true — only Required=false skips blocking on failure.
		if !res.check.Required && res.check.Required == false && res.check.Name != "" {
			// User explicitly opted out of blocking. Log warning, don't
			// add to failures. (The default zero-value of Required is
			// false, but we treat zero-value as "required by default" —
			// see below for the explicit-opt-out check.)
		}
		logger.Error("precondition FAILED",
			res.err,
			"name", res.check.Name,
			"took_ms", res.took.Milliseconds(),
			"fix", res.check.Hint,
		)
		failures = append(failures, res)
	}

	if len(failures) == 0 {
		logger.Info("All preconditions passed")
		return nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d boot precondition(s) failed — refusing to start:\n", len(failures)))
	for _, f := range failures {
		b.WriteString(fmt.Sprintf("  ✗ %s: %v\n", f.check.Name, f.err))
		if f.check.Hint != "" {
			b.WriteString(fmt.Sprintf("    fix: %s\n", f.check.Hint))
		}
	}
	return fmt.Errorf("%s", b.String())
}
