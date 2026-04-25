package logger

// Lifecycle helpers wrap structured logger calls with a fixed shape so every
// service emits service-start, service-stop, and goroutine-start/stop events
// the same way. Operators can reliably grep for `lifecycle_event` to trace
// boots, restarts, and goroutine crashes across the fleet.

// LifecycleEvent values are emitted under the canonical "lifecycle_event" key.
const (
	lifecycleKey = "lifecycle_event"

	LifecycleServiceStarted   = "service_started"
	LifecycleServiceStopping  = "service_stopping"
	LifecycleServiceStopped   = "service_stopped"
	LifecycleGoroutineStarted = "goroutine_started"
	LifecycleGoroutineStopped = "goroutine_stopped"
)

// ServiceStarted logs that a service finished its boot sequence and is ready
// to serve traffic. Call once after all components (DB, NATS, HTTP, etc.) are
// online. Extra fields are appended verbatim.
func ServiceStarted(name, version string, fields ...interface{}) {
	args := []interface{}{
		lifecycleKey, LifecycleServiceStarted,
		KeyServiceName, name,
		KeyVersion, version,
	}
	args = append(args, fields...)
	Info("Service started", args...)
}

// ServiceStopping logs that the service has received a shutdown signal and is
// beginning graceful teardown. Call once at the top of the shutdown path.
func ServiceStopping(name string, fields ...interface{}) {
	args := []interface{}{
		lifecycleKey, LifecycleServiceStopping,
		KeyServiceName, name,
	}
	args = append(args, fields...)
	Info("Service stopping", args...)
}

// ServiceStopped logs that shutdown completed successfully. durationMs is the
// time spent draining; pass 0 if not measured.
func ServiceStopped(name string, durationMs int64, fields ...interface{}) {
	args := []interface{}{
		lifecycleKey, LifecycleServiceStopped,
		KeyServiceName, name,
		KeyDurationMs, durationMs,
	}
	args = append(args, fields...)
	Info("Service stopped", args...)
}

// GoroutineStarted logs entry into a long-running goroutine. Use the goroutine
// name as a stable identifier (e.g. "kafka-consumer", "metrics-server").
// Pair every call with a deferred GoroutineStopped so unexpected exits surface.
func GoroutineStarted(name string, fields ...interface{}) {
	args := []interface{}{
		lifecycleKey, LifecycleGoroutineStarted,
		KeyComponent, name,
	}
	args = append(args, fields...)
	Info("Goroutine started", args...)
}

// GoroutineStopped logs exit from a long-running goroutine. If err is non-nil
// the event is logged at Error so unexpected goroutine deaths page operators;
// a clean exit logs at Info. Pass via `defer logger.GoroutineStopped(name, err)`.
func GoroutineStopped(name string, err error, fields ...interface{}) {
	args := []interface{}{
		lifecycleKey, LifecycleGoroutineStopped,
		KeyComponent, name,
	}
	args = append(args, fields...)
	if err != nil {
		args = append(args, err)
		Error("Goroutine exited with error", args...)
		return
	}
	Info("Goroutine stopped", args...)
}
