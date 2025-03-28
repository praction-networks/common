package metrics

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/shirou/gopsutil/cpu"
)

var registry = prometheus.NewRegistry()

// HTTP Metrics
var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "status", "handler"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "handler"},
	)

	HTTPResponseSizes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response sizes in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "handler"},
	)

	HTTPErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "HTTP error count by type",
		},
		[]string{"method", "handler", "error_type"},
	)
)

// NATS Metrics
var (
	NATSPublishedEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_published_events_total",
			Help: "Total NATS events published",
		},
		[]string{"stream", "subject"},
	)

	NATSPublishFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_publish_failures_total",
			Help: "NATS publish failures",
		},
		[]string{"stream", "subject", "error"},
	)

	NATSEventProcessingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "nats_event_processing_seconds",
			Help:    "NATS event processing time",
			Buckets: []float64{.001, .005, .01, .05, .1, .5, 1, 5},
		},
		[]string{"stream", "subject"},
	)

	NATSInflightMessages = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nats_inflight_messages",
			Help: "Current inflight NATS messages",
		},
		[]string{"stream", "subject"},
	)

	NATSPublishDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "nats_publish_duration_seconds",
			Help:    "Duration taken to publish NATS events",
			Buckets: []float64{.001, .005, .01, .05, .1, .5, 1, 2.5, 5, 10},
		},
		[]string{"stream", "subject", "success"},
	)
)

// System Metrics
var (
	SystemCPUUsage = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "Current system CPU usage percentage",
		},
		func() float64 {
			percentages, err := cpu.Percent(0, false)
			if err != nil || len(percentages) == 0 {
				return 0.0
			}
			return percentages[0]
		},
	)

	ProcessMemoryUsage = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "process_memory_usage_bytes",
			Help: "Current process memory usage",
		},
		func() float64 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return float64(m.Alloc)
		},
	)

	GoroutinesCount = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Current number of goroutines",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	)
)

// Business Metrics
var (
	UsersRegistered = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "users_registered_total",
			Help: "Total registered users",
		},
	)

	APICalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_calls_total",
			Help: "API call counts",
		},
		[]string{"endpoint"},
	)
)

// ResponseWriter tracks response status and size
type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}

func (rw *ResponseWriter) Size() int {
	return rw.size
}

// RegisterAllMetrics registers all metrics with the registry
func RegisterAllMetrics() {
	registry.MustRegister(
		// HTTP metrics
		HTTPRequests,
		HTTPRequestDuration,
		HTTPResponseSizes,
		HTTPErrors,

		// NATS metrics
		NATSPublishedEvents,
		NATSPublishFailures,
		NATSEventProcessingTime,
		NATSInflightMessages,
		NATSPublishDuration,

		// System metrics
		SystemCPUUsage,
		ProcessMemoryUsage,
		GoroutinesCount,

		// Business metrics
		UsersRegistered,
		APICalls,

		// Default collectors
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Build info
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Build information",
		},
		[]string{"version", "commit", "date"},
	)
	registry.MustRegister(buildInfo)
	buildInfo.WithLabelValues(
		getEnv("VERSION", "unknown"),
		getEnv("GIT_COMMIT", "unknown"),
		getEnv("BUILD_DATE", "unknown"),
	).Set(1)

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Registry returns the metrics registry
func Registry() *prometheus.Registry {
	return registry
}

// NATS Metrics Helpers
func RecordNATSPublished(stream, subject string) {
	NATSPublishedEvents.WithLabelValues(stream, subject).Inc()
}

func RecordNATSFailure(stream, subject string, err error) {
	NATSPublishFailures.WithLabelValues(stream, subject, err.Error()).Inc()
}

func RecordNATSProcessingTime(stream, subject string, duration time.Duration) {
	NATSEventProcessingTime.WithLabelValues(stream, subject).Observe(duration.Seconds())
}

func IncNATSInflight(stream, subject string) {
	NATSInflightMessages.WithLabelValues(stream, subject).Inc()
}

func DecNATSInflight(stream, subject string) {
	NATSInflightMessages.WithLabelValues(stream, subject).Dec()
}
