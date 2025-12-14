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

// Security Metrics
var (
	// Nonce operations
	NonceStored = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nonce_stored_total",
			Help: "Total nonces stored for replay protection",
		},
		[]string{"status"},
	)

	NonceChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nonce_checks_total",
			Help: "Total nonce existence checks",
		},
		[]string{"result"}, // "exists", "not_exists"
	)

	ActiveNonces = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_nonces_count",
			Help: "Current number of active nonces in Redis",
		},
	)

	// Replay attack detection
	ReplayAttacksDetected = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "replay_attacks_detected_total",
			Help: "Total replay attacks detected and blocked",
		},
		[]string{"type", "nas_id"}, // type: "nonce_reuse", "timestamp_expired", etc.
	)

	// Security events
	SecurityIncidents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_incidents_total",
			Help: "Total security incidents",
		},
		[]string{"type", "severity"}, // type: "invalid_signature", "unauthorized_nas", etc.
	)

	// HMAC validation
	HMACValidations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hmac_validations_total",
			Help: "Total HMAC signature validations",
		},
		[]string{"result"}, // "success", "failure"
	)

	// Timestamp validation
	TimestampValidations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timestamp_validations_total",
			Help: "Total timestamp validations",
		},
		[]string{"result"}, // "valid", "expired", "future"
	)

	// NAS authentication
	NASAuthAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nas_auth_attempts_total",
			Help: "Total NAS authentication attempts",
		},
		[]string{"nas_id", "result"}, // result: "success", "failure"
	)
)

// Redis Operation Metrics
var (
	RedisOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_operations_total",
			Help: "Total Redis operations",
		},
		[]string{"operation", "status"}, // operation: "get", "set", "del", "exists"
	)

	RedisOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_operation_duration_seconds",
			Help:    "Redis operation duration",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .05, .1},
		},
		[]string{"operation"},
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

		// Security metrics
		NonceStored,
		NonceChecks,
		ActiveNonces,
		ReplayAttacksDetected,
		SecurityIncidents,
		HMACValidations,
		TimestampValidations,
		NASAuthAttempts,

		// Redis metrics
		RedisOperations,
		RedisOperationDuration,

		// Redirector metrics
		RedirectorRequests,
		RedirectorRateLimitHits,
		RedirectorTokensGenerated,
		RedirectorURLGenerationDuration,
		RedirectorValidationFailures,
		RedirectorCircuitBreakerState,
		RedirectorCircuitBreakerFailures,
		RedirectorMemoryRateLimiterSize,

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

// Security Metrics Helpers
func RecordNonceStored(status string) {
	NonceStored.WithLabelValues(status).Inc()
}

func RecordNonceCheck(result string) {
	NonceChecks.WithLabelValues(result).Inc()
}

func SetActiveNonces(count float64) {
	ActiveNonces.Set(count)
}

func RecordReplayAttack(attackType string, nasId string) {
	ReplayAttacksDetected.WithLabelValues(attackType, nasId).Inc()
}

func RecordSecurityIncident(incidentType string, severity string) {
	SecurityIncidents.WithLabelValues(incidentType, severity).Inc()
}

func RecordHMACValidation(result string) {
	HMACValidations.WithLabelValues(result).Inc()
}

func RecordTimestampValidation(result string) {
	TimestampValidations.WithLabelValues(result).Inc()
}

func RecordNASAuthAttempt(nasId string, result string) {
	NASAuthAttempts.WithLabelValues(nasId, result).Inc()
}

// Redis Metrics Helpers
func RecordRedisOperation(operation string, status string) {
	RedisOperations.WithLabelValues(operation, status).Inc()
}

func RecordRedisOperationDuration(operation string, duration time.Duration) {
	RedisOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// Redirector Service Metrics
var (
	// Redirector requests
	RedirectorRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirector_requests_total",
			Help: "Total redirector requests",
		},
		[]string{"nas_id", "status"}, // status: "success", "rate_limited", "invalid_hostname", etc.
	)

	// Rate limiting
	RedirectorRateLimitHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirector_rate_limit_hits_total",
			Help: "Total redirector rate limit hits",
		},
		[]string{"type", "nas_id"}, // type: "per_nas", "per_ip"
	)

	// Token generation
	RedirectorTokensGenerated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirector_tokens_generated_total",
			Help: "Total redirector tokens generated",
		},
		[]string{"nas_id", "status"}, // status: "success", "failure"
	)

	// URL generation duration
	RedirectorURLGenerationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redirector_url_generation_duration_seconds",
			Help:    "Redirector URL generation duration",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"nas_id"},
	)

	// Request validation failures
	RedirectorValidationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirector_validation_failures_total",
			Help: "Total redirector validation failures",
		},
		[]string{"type"}, // type: "hostname_length", "query_param_length", "query_string_length", "decode_iterations", "invalid_hostname"
	)

	// Circuit breaker state changes
	RedirectorCircuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redirector_circuit_breaker_state",
			Help: "Redirector circuit breaker state (0=Closed, 1=Open, 2=HalfOpen)",
		},
		[]string{"service"}, // service: "redis", "mongodb"
	)

	// Circuit breaker failures
	RedirectorCircuitBreakerFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirector_circuit_breaker_failures_total",
			Help: "Total circuit breaker failures",
		},
		[]string{"service", "state"}, // service: "redis", "mongodb"; state: "closed", "open", "halfopen"
	)

	// Memory rate limiter usage
	RedirectorMemoryRateLimiterSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "redirector_memory_rate_limiter_size",
			Help: "Current number of entries in memory rate limiter",
		},
	)
)

// Redirector Metrics Helpers
func RecordRedirectorRequest(nasId string, status string) {
	RedirectorRequests.WithLabelValues(nasId, status).Inc()
}

func RecordRedirectorRateLimitHit(limitType string, nasId string) {
	RedirectorRateLimitHits.WithLabelValues(limitType, nasId).Inc()
}

func RecordRedirectorTokenGenerated(nasId string, status string) {
	RedirectorTokensGenerated.WithLabelValues(nasId, status).Inc()
}

func RecordRedirectorURLGenerationDuration(nasId string, duration time.Duration) {
	RedirectorURLGenerationDuration.WithLabelValues(nasId).Observe(duration.Seconds())
}

func RecordRedirectorValidationFailure(validationType string) {
	RedirectorValidationFailures.WithLabelValues(validationType).Inc()
}

func SetRedirectorCircuitBreakerState(service string, state float64) {
	RedirectorCircuitBreakerState.WithLabelValues(service).Set(state)
}

func RecordRedirectorCircuitBreakerFailure(service string, state string) {
	RedirectorCircuitBreakerFailures.WithLabelValues(service, state).Inc()
}

func SetRedirectorMemoryRateLimiterSize(size float64) {
	RedirectorMemoryRateLimiterSize.Set(size)
}
