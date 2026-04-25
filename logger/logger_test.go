package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

// captureStdout redirects os.Stdout for the duration of fn and returns what
// was written. The logger initializes its async writer against os.Stdout once
// per process, so this helper must be called before the first InitializeLogger
// in this test binary.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w

	var buf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&buf, r)
	}()

	defer func() {
		os.Stdout = orig
	}()

	fn()
	Sync()
	_ = w.Close()
	wg.Wait()
	return buf.String()
}

func TestCanonicalKeys(t *testing.T) {
	cases := map[string]string{
		"KeyTenantID":    KeyTenantID,
		"KeyUserID":      KeyUserID,
		"KeyRequestID":   KeyRequestID,
		"KeyTraceID":     KeyTraceID,
		"KeyComponent":   KeyComponent,
		"KeyServiceName": KeyServiceName,
		"KeyVersion":     KeyVersion,
		"KeyDurationMs":  KeyDurationMs,
		"KeyStream":      KeyStream,
		"KeySubject":     KeySubject,
		"KeySequence":    KeySequence,
		"KeyHandler":     KeyHandler,
	}
	want := map[string]string{
		"KeyTenantID":    "tenant_id",
		"KeyUserID":      "user_id",
		"KeyRequestID":   "request_id",
		"KeyTraceID":     "trace_id",
		"KeyComponent":   "component",
		"KeyServiceName": "service",
		"KeyVersion":     "version",
		"KeyDurationMs":  "duration_ms",
		"KeyStream":      "stream",
		"KeySubject":     "subject",
		"KeySequence":    "sequence",
		"KeyHandler":     "handler",
	}
	for name, got := range cases {
		if got != want[name] {
			t.Errorf("%s = %q, want %q", name, got, want[name])
		}
	}
}

func TestLifecycleHelpersEmitExpectedFields(t *testing.T) {
	out := captureStdout(t, func() {
		if err := InitializeLogger(LoggerConfig{LogLevel: "debug"}); err != nil {
			t.Fatalf("init: %v", err)
		}
		ServiceStarted("auth-service", "1.2.3", "components_online", 4)
		ServiceStopping("auth-service")
		ServiceStopped("auth-service", 250)
		GoroutineStarted("kafka-consumer")
		GoroutineStopped("kafka-consumer", nil)
	})

	wantSubstrings := []string{
		`"lifecycle_event":"service_started"`,
		`"lifecycle_event":"service_stopping"`,
		`"lifecycle_event":"service_stopped"`,
		`"lifecycle_event":"goroutine_started"`,
		`"lifecycle_event":"goroutine_stopped"`,
		`"service":"auth-service"`,
		`"version":"1.2.3"`,
		`"duration_ms":250`,
		`"component":"kafka-consumer"`,
	}
	for _, want := range wantSubstrings {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, out)
		}
	}
}
