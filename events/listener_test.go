package events

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/praction-networks/common/logger"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	var buf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&buf, r)
	}()

	fn()
	logger.Sync()
	_ = w.Close()
	wg.Wait()
	return buf.String()
}

func TestLogHandlerOutcomeVerboseToggle(t *testing.T) {
	out := captureStdout(t, func() {
		if err := logger.InitializeLogger(logger.LoggerConfig{LogLevel: "debug"}); err != nil {
			t.Fatalf("init: %v", err)
		}
		verbose := &Listener{StreamName: "TEST", HandlerName: "h1", VerboseSuccessLog: true}
		quiet := &Listener{StreamName: "TEST", HandlerName: "h2", VerboseSuccessLog: false}
		failing := &Listener{StreamName: "TEST", HandlerName: "h3", VerboseSuccessLog: false}

		verbose.logHandlerOutcome("success", "subj.a", uint64(1), 5*time.Millisecond, nil)
		quiet.logHandlerOutcome("success", "subj.b", uint64(2), 5*time.Millisecond, nil)
		failing.logHandlerOutcome("handle", "subj.c", uint64(3), 5*time.Millisecond, errors.New("boom"))
	})

	// Verbose listener emits its success at INFO with handler=h1.
	if !strings.Contains(out, `"level":"INFO"`) || !strings.Contains(out, `"handler":"h1"`) {
		t.Errorf("expected verbose Info success with handler=h1; got:\n%s", out)
	}

	// Quiet listener emits its success at DEBUG with handler=h2.
	hasQuietDebug := false
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, `"handler":"h2"`) && strings.Contains(line, `"level":"DEBUG"`) {
			hasQuietDebug = true
			break
		}
	}
	if !hasQuietDebug {
		t.Errorf("expected quiet Debug success line for handler=h2; got:\n%s", out)
	}

	// Failing listener emits ERROR with the error string and stage=handle.
	hasErr := false
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, `"handler":"h3"`) && strings.Contains(line, `"level":"ERROR"`) &&
			strings.Contains(line, `"error":"boom"`) && strings.Contains(line, `"stage":"handle"`) {
			hasErr = true
			break
		}
	}
	if !hasErr {
		t.Errorf("expected error line for handler=h3 with stage=handle; got:\n%s", out)
	}
}
