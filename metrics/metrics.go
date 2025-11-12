package metrics

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	GoroutinesCount = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Current number of goroutines",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	)

	CPUUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_cpu_usage_manual_percent",
		Help: "Manual CPU usage percent (polled every 10s)",
	})

	MemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_memory_usage_bytes_manual",
		Help: "Manual memory usage in bytes (polled every 10s)",
	})

	DiskUsagePercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_disk_usage_percent",
		Help: "Disk usage percent on root /",
	})

	DiskReadBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_disk_read_bytes",
		Help: "Disk read bytes",
	})

	DiskWriteBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_disk_write_bytes",
		Help: "Disk write bytes",
	})
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

// Query Builder Metrics
var (
	QueryBuildsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "query_builds_total",
			Help: "Total number of queries built",
		},
	)

	QueryBuildDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "query_build_duration_seconds",
			Help:    "Query build duration in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1},
		},
	)

	QueryBuildErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_build_errors_total",
			Help: "Total query build errors",
		},
		[]string{"error_type"},
	)

	QueryCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "query_cache_hits_total",
			Help: "Total query cache hits",
		},
	)

	QueryCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "query_cache_misses_total",
			Help: "Total query cache misses",
		},
	)

	QueryRateLimitHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_rate_limit_hits_total",
			Help: "Total query rate limit hits",
		},
		[]string{"key"},
	)

	QueryComplexity = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "query_complexity_score",
			Help:    "Query complexity score (filters + sorts + search)",
			Buckets: []float64{1, 5, 10, 15, 20, 25, 30, 40, 50},
		},
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
		CPUUsage,
		MemoryUsage,
		DiskUsagePercent,
		DiskReadBytes,
		DiskWriteBytes,
		GoroutinesCount,

		// Business metrics
		UsersRegistered,
		APICalls,

		// Query builder metrics
		QueryBuildsTotal,
		QueryBuildDuration,
		QueryBuildErrors,
		QueryCacheHits,
		QueryCacheMisses,
		QueryRateLimitHits,
		QueryComplexity,

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
