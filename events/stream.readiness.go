package events

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/common/logger"
)

// WaitForStreams waits for all required streams to be ready before service startup
// This prevents event loss by ensuring streams exist before services start publishing
// Returns error if any stream is not ready within the timeout
func WaitForStreams(ctx context.Context, streamManager *JsStreamManager, streams []StreamName, timeout time.Duration) error {
	if len(streams) == 0 {
		return nil
	}

	logger.Info("Waiting for required streams to be ready", "streams", streams, "timeout", timeout)

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	ready := make(map[StreamName]bool)

	for {
		allReady := true

		for _, streamName := range streams {
			if ready[streamName] {
				continue
			}

			// Check if stream exists and is ready
			streamInfo, err := streamManager.Stream(ctx, string(streamName))
			if err != nil {
				allReady = false
				continue
			}

			// Stream exists, mark as ready
			ready[streamName] = true
			logger.Info("Stream is ready", "streamName", streamName, "state", streamInfo.Config.Name)
		}

		if allReady {
			logger.Info("All required streams are ready", "streams", streams)
			return nil
		}

		// Check timeout
		if time.Now().After(deadline) {
			missing := []StreamName{}
			for _, streamName := range streams {
				if !ready[streamName] {
					missing = append(missing, streamName)
				}
			}
			return fmt.Errorf("timeout waiting for streams to be ready: %v", missing)
		}

		// Wait before next check
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for streams: %w", ctx.Err())
		}
	}
}
