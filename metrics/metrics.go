package metrics

import (
	"bufio"
	"fmt"
	"net"
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

// CWMP / ACS Metrics
var (
	// Inform messages received from CPE devices
	CWMPInformsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cwmp_informs_total",
			Help: "Total CWMP Inform messages received",
		},
		[]string{"tenant_id", "event_code", "manufacturer"},
	)

	// Tasks executed via CWMP
	CWMPTasksExecuted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cwmp_tasks_executed_total",
			Help: "Total CWMP tasks executed",
		},
		[]string{"tenant_id", "action", "status"}, // action: "GetParameterValues", "SetParameterValues", etc.; status: "success", "failure"
	)

	// Active CWMP sessions gauge
	CWMPActiveSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cwmp_active_sessions",
			Help: "Current number of active CWMP sessions",
		},
		[]string{"tenant_id"},
	)

	// CWMP session duration
	CWMPSessionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cwmp_session_duration_seconds",
			Help:    "CWMP session duration from first Inform to session end",
			Buckets: []float64{.1, .5, 1, 2.5, 5, 10, 30, 60, 120, 300},
		},
		[]string{"tenant_id"},
	)

	// SOAP parse errors
	CWMPSOAPParseErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cwmp_soap_parse_errors_total",
			Help: "Total SOAP/XML parse errors in CWMP messages",
		},
	)

	// Connection requests sent to CPE devices
	CWMPConnectionRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cwmp_connection_requests_total",
			Help: "Total connection requests sent to CPE devices",
		},
		[]string{"tenant_id", "auth_method", "result"}, // auth_method: "basic", "digest"; result: "success", "failure"
	)

	// Auto-provisioning events
	CWMPAutoProvisioningTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cwmp_auto_provisioning_total",
			Help: "Total auto-provisioning events triggered",
		},
		[]string{"tenant_id", "status"}, // status: "success", "failure", "no_match"
	)

	// Device registrations
	CWMPDeviceRegistrations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cwmp_device_registrations_total",
			Help: "Total new device registrations via CWMP Inform",
		},
		[]string{"tenant_id", "manufacturer"},
	)
)

// DNS Filtering / i9Shield Metrics
var (
	// ── DNS Query Lifecycle ──────────────────────────────────────────────

	// Total DNS queries by type, tenant, and result
	DNSQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_queries_total",
			Help: "Total DNS queries processed",
		},
		[]string{"tenant_id", "qtype", "result"}, // result: "blocked", "allowed", "cached", "error"
	)

	// End-to-end DNS query latency
	DNSQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_query_duration_seconds",
			Help:    "DNS query processing duration in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"result"}, // result: "blocked", "allowed", "cached", "error"
	)

	// Currently processing DNS queries
	DNSQueriesInflight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dns_queries_inflight",
			Help: "Number of DNS queries currently being processed",
		},
	)

	// ── DNS Blocking ────────────────────────────────────────────────────

	// Blocked queries per tenant, category, and reason
	DNSBlockedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_blocked_total",
			Help: "Total DNS queries blocked",
		},
		[]string{"tenant_id", "category", "reason"}, // reason: "category", "custom_blocklist", "custom_domain"
	)

	// Allowed queries per tenant
	DNSAllowedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_allowed_total",
			Help: "Total DNS queries allowed (pass-through)",
		},
		[]string{"tenant_id"},
	)

	// ── DNS Cache Performance ───────────────────────────────────────────

	DNSCacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_hits_total",
			Help: "Total DNS answer cache hits",
		},
	)

	DNSCacheMissesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_misses_total",
			Help: "Total DNS answer cache misses",
		},
	)

	DNSCacheSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dns_cache_size",
			Help: "Current number of entries in the DNS answer cache",
		},
	)

	DNSCacheEvictionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_evictions_total",
			Help: "Total DNS cache evictions",
		},
	)

	// ── DNS Upstream Resolution ────────────────────────────────────────

	// Per-upstream query counts
	DNSUpstreamQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_upstream_queries_total",
			Help: "Total queries sent to upstream DNS resolvers",
		},
		[]string{"upstream", "status"}, // status: "success", "failure"
	)

	// Per-upstream latency
	DNSUpstreamLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_upstream_latency_seconds",
			Help:    "Upstream DNS resolver latency in seconds",
			Buckets: []float64{.0005, .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"upstream"},
	)

	// TCP fallback count (truncated UDP → TCP)
	DNSUpstreamTCPFallbacksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_upstream_tcp_fallbacks_total",
			Help: "Total truncated UDP responses that fell back to TCP",
		},
		[]string{"upstream"},
	)

	// All upstream resolvers failed
	DNSUpstreamAllFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_upstream_all_failed_total",
			Help: "Total times all upstream resolvers failed for a query",
		},
	)

	// ── DNS Policy Engine ──────────────────────────────────────────────

	// Policy resolution counts
	DNSPolicyResolutionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_policy_resolutions_total",
			Help: "Total policy resolution attempts",
		},
		[]string{"tenant_id", "result"}, // result: "hit", "miss", "error"
	)

	// Policy resolution latency
	DNSPolicyResolutionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_policy_resolution_duration_seconds",
			Help:    "Policy resolution and merge duration in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1},
		},
	)

	// Policy cache hits/misses
	DNSPolicyCacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_policy_cache_hits_total",
			Help: "Total in-memory policy cache hits",
		},
	)

	DNSPolicyCacheMissesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_policy_cache_misses_total",
			Help: "Total in-memory policy cache misses (DB fetch required)",
		},
	)

	// Policy cache invalidations
	DNSPolicyCacheInvalidationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_policy_cache_invalidations_total",
			Help: "Total policy cache invalidations",
		},
		[]string{"scope"}, // scope: "tenant", "subscriber"
	)

	// ── DNS Subscriber Resolution ──────────────────────────────────────

	// Redis IP → subscriber lookups
	DNSSubscriberLookupsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_subscriber_lookups_total",
			Help: "Total Redis IP-to-subscriber lookups",
		},
		[]string{"result"}, // result: "found", "not_found", "error"
	)

	// Subscriber lookup latency
	DNSSubscriberLookupDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_subscriber_lookup_duration_seconds",
			Help:    "Redis subscriber lookup duration in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1},
		},
	)

	// ── DNS Bloom Filter & Blocklist ───────────────────────────────────

	// Bloom filter checks
	DNSBloomChecksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_bloom_checks_total",
			Help: "Total bloom filter checks",
		},
		[]string{"result"}, // result: "positive", "negative"
	)

	// Bloom false positives (bloom=yes, hashmap=no)
	DNSBloomFalsePositivesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_bloom_false_positives_total",
			Help: "Total bloom filter false positives (bloom positive, hash map negative)",
		},
	)

	// Allowlist short-circuit counter
	DNSAllowlistShortCircuitTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_allowlist_shortcircuit_total",
			Help: "Total queries short-circuited by global allowlist (skipped bloom checks)",
		},
	)

	// Loaded domains per category
	DNSBlocklistDomainsLoaded = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_blocklist_domains_loaded",
			Help: "Number of domains loaded per blocklist category",
		},
		[]string{"category"},
	)

	// Bloom filter memory per category
	DNSBlocklistMemoryBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_blocklist_memory_bytes",
			Help: "Bloom filter memory usage in bytes per category",
		},
		[]string{"category"},
	)

	// Number of active blocking categories
	DNSBlocklistCategoriesActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dns_blocklist_categories_active",
			Help: "Number of active blocklist categories",
		},
	)

	// ── DNS Blocklist Refresh ──────────────────────────────────────────

	// Refresh operations per category
	DNSBlocklistRefreshTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_blocklist_refresh_total",
			Help: "Total blocklist refresh operations",
		},
		[]string{"category", "status"}, // status: "success", "failure"
	)

	// Refresh duration per category
	DNSBlocklistRefreshDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_blocklist_refresh_duration_seconds",
			Help:    "Blocklist refresh duration (fetch + load) in seconds",
			Buckets: []float64{.5, 1, 2.5, 5, 10, 30, 60, 120, 300},
		},
		[]string{"category"},
	)

	// Domains fetched per refresh
	DNSBlocklistRefreshDomainsFetched = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_blocklist_refresh_domains_fetched_total",
			Help: "Total domains fetched during blocklist refreshes",
		},
		[]string{"category"},
	)

	// Last successful refresh timestamp
	DNSBlocklistLastRefreshTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_blocklist_last_refresh_timestamp",
			Help: "Unix timestamp of last successful blocklist refresh",
		},
		[]string{"category"},
	)

	// ── DNS Resilience / Rate Limiting ─────────────────────────────────

	// Stale cache hits (served expired entry because upstream failed)
	DNSStaleCacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_stale_cache_hits_total",
			Help: "Total DNS queries served from stale (expired) cache entries due to upstream failure",
		},
	)

	// NXDOMAIN cache hits (cached negative response)
	DNSNXDOMAINCacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_nxdomain_cache_hits_total",
			Help: "Total DNS queries served from cached NXDOMAIN responses",
		},
	)

	// Rate limiter: tracked IPs (memory pressure indicator)
	DNSRateLimiterTrackedIPs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dns_rate_limiter_tracked_ips",
			Help: "Number of client IPs currently tracked by the rate limiter",
		},
	)

	// Rate limiter: rejections by tier
	DNSRateLimitRejectionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_rate_limit_rejections_total",
			Help: "Total DNS queries rejected by rate limiting",
		},
		[]string{"tier"}, // "global", "tenant"
	)

	// Circuit breaker: current state per upstream
	DNSCircuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_circuit_breaker_state",
			Help: "Current circuit breaker state per upstream (0=closed, 1=open, 2=half-open)",
		},
		[]string{"upstream"},
	)

	// Circuit breaker: state transitions
	DNSCircuitBreakerTransitionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_circuit_breaker_transitions_total",
			Help: "Total circuit breaker state transitions",
		},
		[]string{"upstream", "from", "to"}, // e.g., "closed" → "open"
	)
)

// KYC Verification Metrics
var (
	// Total KYC verifications by type, provider, and outcome
	KYCVerificationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_verifications_total",
			Help: "Total KYC verification requests",
		},
		[]string{"type", "provider", "status"}, // status: "success", "failed", "not_found", "invalid"
	)

	// End-to-end KYC verification latency
	KYCVerificationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kyc_verification_duration_seconds",
			Help:    "KYC verification duration from request to response",
			Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 15, 30},
		},
		[]string{"type", "provider"},
	)

	// Circuit breaker state per provider (0=closed, 1=open, 2=half-open)
	KYCCircuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kyc_circuit_breaker_state",
			Help: "Current KYC circuit breaker state per provider (0=closed, 1=open, 2=half-open)",
		},
		[]string{"provider"},
	)

	// Provider error breakdown
	KYCProviderErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_provider_errors_total",
			Help: "Total KYC provider errors by error code",
		},
		[]string{"provider", "error_code"},
	)

	// Audit record write counts
	KYCAuditRecordsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_audit_records_total",
			Help: "Total KYC audit records written",
		},
		[]string{"status"}, // "success", "failure"
	)
)

// ResponseWriter tracks response status and size
// It preserves http.Hijacker interface for WebSocket support
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

// Hijack implements http.Hijacker interface for WebSocket support
// It delegates to the underlying ResponseWriter if it implements Hijacker
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
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
		RedirectorNASCacheHits,
		RedirectorNASCacheMisses,
		RedirectorNASCacheInvalidations,
		RedirectorMongoDBPoolSize,
		RedirectorRedisPoolSize,

		// Log Storage Re-Encoding metrics
		ReEncodingOperationsTotal,
		ReEncodingDuration,
		ReEncodingFilesProcessed,
		ReEncodingRecordsProcessed,
		ReEncodingBytesProcessed,
		ReEncodingBytesAfter,
		ReEncodingCompressionRatio,
		ReEncodingErrors,
		ReEncodingScheduledRuns,
		ReEncodingScheduledRunDuration,
		ReEncodingActiveWorkers,
		ReEncodingQueueSize,

		// CWMP / ACS metrics
		CWMPInformsTotal,
		CWMPTasksExecuted,
		CWMPActiveSessions,
		CWMPSessionDuration,
		CWMPSOAPParseErrors,
		CWMPConnectionRequests,
		CWMPAutoProvisioningTotal,
		CWMPDeviceRegistrations,

		// DNS Filtering / i9Shield metrics
		DNSQueriesTotal,
		DNSQueryDuration,
		DNSQueriesInflight,
		DNSBlockedTotal,
		DNSAllowedTotal,
		DNSCacheHitsTotal,
		DNSCacheMissesTotal,
		DNSCacheSize,
		DNSCacheEvictionsTotal,
		DNSUpstreamQueriesTotal,
		DNSUpstreamLatency,
		DNSUpstreamTCPFallbacksTotal,
		DNSUpstreamAllFailedTotal,
		DNSPolicyResolutionsTotal,
		DNSPolicyResolutionDuration,
		DNSPolicyCacheHitsTotal,
		DNSPolicyCacheMissesTotal,
		DNSPolicyCacheInvalidationsTotal,
		DNSSubscriberLookupsTotal,
		DNSSubscriberLookupDuration,
		DNSBloomChecksTotal,
		DNSBloomFalsePositivesTotal,
		DNSAllowlistShortCircuitTotal,
		DNSBlocklistDomainsLoaded,
		DNSBlocklistMemoryBytes,
		DNSBlocklistCategoriesActive,
		DNSBlocklistRefreshTotal,
		DNSBlocklistRefreshDuration,
		DNSBlocklistRefreshDomainsFetched,
		DNSBlocklistLastRefreshTimestamp,

		// DNS Resilience metrics
		DNSStaleCacheHitsTotal,
		DNSNXDOMAINCacheHitsTotal,
		DNSRateLimiterTrackedIPs,
		DNSRateLimitRejectionsTotal,
		DNSCircuitBreakerState,
		DNSCircuitBreakerTransitionsTotal,

		// KYC Verification metrics
		KYCVerificationsTotal,
		KYCVerificationDuration,
		KYCCircuitBreakerState,
		KYCProviderErrors,
		KYCAuditRecordsTotal,

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

	// NAS device cache metrics
	RedirectorNASCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redirector_nas_cache_hits_total",
			Help: "Total NAS device cache hits",
		},
	)

	RedirectorNASCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redirector_nas_cache_misses_total",
			Help: "Total NAS device cache misses",
		},
	)

	RedirectorNASCacheInvalidations = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redirector_nas_cache_invalidations_total",
			Help: "Total NAS device cache invalidations",
		},
	)

	// Connection pool metrics
	RedirectorMongoDBPoolSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redirector_mongodb_pool_size",
			Help: "MongoDB connection pool size metrics",
		},
		[]string{"type"}, // type: "max", "min", "available", "in_use"
	)

	RedirectorRedisPoolSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redirector_redis_pool_size",
			Help: "Redis connection pool size metrics",
		},
		[]string{"type"}, // type: "max", "idle", "active", "waiting"
	)
)

// Log Storage Re-Encoding Metrics
var (
	// Re-encoding operations
	ReEncodingOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_operations_total",
			Help: "Total re-encoding operations",
		},
		[]string{"status"}, // status: "success", "failure"
	)

	// Re-encoding performance
	ReEncodingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "log_storage_reencoding_duration_seconds",
			Help:    "Re-encoding operation duration in seconds",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600}, // 1s to 1h
		},
		[]string{"log_type"}, // log_type: "SESSION", "NAT_EVENT", "FLOW"
	)

	// Re-encoding results
	ReEncodingFilesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_files_processed_total",
			Help: "Total files processed by re-encoding",
		},
		[]string{"log_type", "status"}, // status: "created", "replaced", "skipped"
	)

	ReEncodingRecordsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_records_processed_total",
			Help: "Total records processed by re-encoding",
		},
		[]string{"log_type"},
	)

	ReEncodingBytesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_bytes_processed_total",
			Help: "Total bytes processed by re-encoding (before compression)",
		},
		[]string{"log_type"},
	)

	ReEncodingBytesAfter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_bytes_after_total",
			Help: "Total bytes after re-encoding (after compression)",
		},
		[]string{"log_type"},
	)

	ReEncodingCompressionRatio = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "log_storage_reencoding_compression_ratio",
			Help:    "Compression ratio achieved by re-encoding (before/after)",
			Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 20, 25, 30},
		},
		[]string{"log_type"},
	)

	// Re-encoding errors
	ReEncodingErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_errors_total",
			Help: "Total re-encoding errors",
		},
		[]string{"log_type", "error_type"}, // error_type: "read_failed", "write_failed", "sort_failed", "replace_failed"
	)

	// Re-encoding scheduler
	ReEncodingScheduledRuns = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "log_storage_reencoding_scheduled_runs_total",
			Help: "Total scheduled re-encoding runs",
		},
	)

	ReEncodingScheduledRunDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "log_storage_reencoding_scheduled_run_duration_seconds",
			Help:    "Duration of scheduled re-encoding runs",
			Buckets: []float64{60, 300, 600, 1800, 3600, 7200, 10800}, // 1m to 3h
		},
	)

	// Re-encoding worker metrics
	ReEncodingActiveWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "log_storage_reencoding_active_workers",
			Help: "Current number of active re-encoding workers",
		},
	)

	ReEncodingQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "log_storage_reencoding_queue_size",
			Help: "Current number of file groups in re-encoding queue",
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

func RecordRedirectorNASCacheHit() {
	RedirectorNASCacheHits.Inc()
}

func RecordRedirectorNASCacheMiss() {
	RedirectorNASCacheMisses.Inc()
}

func RecordRedirectorNASCacheInvalidation() {
	RedirectorNASCacheInvalidations.Inc()
}

// Connection Pool Metrics Helpers
func SetRedirectorMongoDBPoolSize(metricType string, value float64) {
	RedirectorMongoDBPoolSize.WithLabelValues(metricType).Set(value)
}

func SetRedirectorRedisPoolSize(metricType string, value float64) {
	RedirectorRedisPoolSize.WithLabelValues(metricType).Set(value)
}

// Re-Encoding Metrics Helpers
func RecordReEncodingOperation(status string) {
	ReEncodingOperationsTotal.WithLabelValues(status).Inc()
}

func RecordReEncodingDuration(logType string, duration time.Duration) {
	ReEncodingDuration.WithLabelValues(logType).Observe(duration.Seconds())
}

func RecordReEncodingFilesProcessed(logType string, status string, count int) {
	ReEncodingFilesProcessed.WithLabelValues(logType, status).Add(float64(count))
}

func RecordReEncodingRecordsProcessed(logType string, count int64) {
	ReEncodingRecordsProcessed.WithLabelValues(logType).Add(float64(count))
}

func RecordReEncodingBytesProcessed(logType string, bytes int64) {
	ReEncodingBytesProcessed.WithLabelValues(logType).Add(float64(bytes))
}

func RecordReEncodingBytesAfter(logType string, bytes int64) {
	ReEncodingBytesAfter.WithLabelValues(logType).Add(float64(bytes))
}

func RecordReEncodingCompressionRatio(logType string, ratio float64) {
	ReEncodingCompressionRatio.WithLabelValues(logType).Observe(ratio)
}

func RecordReEncodingError(logType string, errorType string) {
	ReEncodingErrors.WithLabelValues(logType, errorType).Inc()
}

func RecordReEncodingScheduledRun() {
	ReEncodingScheduledRuns.Inc()
}

func RecordReEncodingScheduledRunDuration(duration time.Duration) {
	ReEncodingScheduledRunDuration.Observe(duration.Seconds())
}

func SetReEncodingActiveWorkers(count float64) {
	ReEncodingActiveWorkers.Set(count)
}

func SetReEncodingQueueSize(size float64) {
	ReEncodingQueueSize.Set(size)
}

// CWMP / ACS Metrics Helpers
func RecordCWMPInform(tenantID, eventCode, manufacturer string) {
	CWMPInformsTotal.WithLabelValues(tenantID, eventCode, manufacturer).Inc()
}

func RecordCWMPTaskExecuted(tenantID, action, status string) {
	CWMPTasksExecuted.WithLabelValues(tenantID, action, status).Inc()
}

func SetCWMPActiveSessions(tenantID string, count float64) {
	CWMPActiveSessions.WithLabelValues(tenantID).Set(count)
}

func RecordCWMPSessionDuration(tenantID string, duration time.Duration) {
	CWMPSessionDuration.WithLabelValues(tenantID).Observe(duration.Seconds())
}

func RecordCWMPSOAPParseError() {
	CWMPSOAPParseErrors.Inc()
}

func RecordCWMPConnectionRequest(tenantID, authMethod, result string) {
	CWMPConnectionRequests.WithLabelValues(tenantID, authMethod, result).Inc()
}

func RecordCWMPAutoProvisioning(tenantID, status string) {
	CWMPAutoProvisioningTotal.WithLabelValues(tenantID, status).Inc()
}

func RecordCWMPDeviceRegistration(tenantID, manufacturer string) {
	CWMPDeviceRegistrations.WithLabelValues(tenantID, manufacturer).Inc()
}

// DNS Filtering Metrics Helpers

// ── Query Lifecycle ─────────────────────────────────────────────────

func RecordDNSQuery(tenantID, qtype, result string, duration time.Duration) {
	DNSQueriesTotal.WithLabelValues(tenantID, qtype, result).Inc()
	DNSQueryDuration.WithLabelValues(result).Observe(duration.Seconds())
}

func IncDNSInflight() {
	DNSQueriesInflight.Inc()
}

func DecDNSInflight() {
	DNSQueriesInflight.Dec()
}

// ── Blocking ────────────────────────────────────────────────────────

func RecordDNSBlocked(tenantID, category, reason string) {
	DNSBlockedTotal.WithLabelValues(tenantID, category, reason).Inc()
}

func RecordDNSAllowed(tenantID string) {
	DNSAllowedTotal.WithLabelValues(tenantID).Inc()
}

// ── Cache ───────────────────────────────────────────────────────────

func RecordDNSCacheHit() {
	DNSCacheHitsTotal.Inc()
}

func RecordDNSCacheMiss() {
	DNSCacheMissesTotal.Inc()
}

func SetDNSCacheSize(size float64) {
	DNSCacheSize.Set(size)
}

func RecordDNSCacheEviction() {
	DNSCacheEvictionsTotal.Inc()
}

// ── Upstream ────────────────────────────────────────────────────────

func RecordDNSUpstreamQuery(upstream, status string, duration time.Duration) {
	DNSUpstreamQueriesTotal.WithLabelValues(upstream, status).Inc()
	DNSUpstreamLatency.WithLabelValues(upstream).Observe(duration.Seconds())
}

func RecordDNSUpstreamTCPFallback(upstream string) {
	DNSUpstreamTCPFallbacksTotal.WithLabelValues(upstream).Inc()
}

func RecordDNSUpstreamAllFailed() {
	DNSUpstreamAllFailedTotal.Inc()
}

// ── Policy Engine ───────────────────────────────────────────────────

func RecordDNSPolicyResolution(tenantID, result string) {
	DNSPolicyResolutionsTotal.WithLabelValues(tenantID, result).Inc()
}

func RecordDNSPolicyResolutionDuration(duration time.Duration) {
	DNSPolicyResolutionDuration.Observe(duration.Seconds())
}

func RecordDNSPolicyCacheHit() {
	DNSPolicyCacheHitsTotal.Inc()
}

func RecordDNSPolicyCacheMiss() {
	DNSPolicyCacheMissesTotal.Inc()
}

func RecordDNSPolicyCacheInvalidation(scope string) {
	DNSPolicyCacheInvalidationsTotal.WithLabelValues(scope).Inc()
}

// ── Subscriber Resolution ───────────────────────────────────────────

func RecordDNSSubscriberLookup(result string, duration time.Duration) {
	DNSSubscriberLookupsTotal.WithLabelValues(result).Inc()
	DNSSubscriberLookupDuration.Observe(duration.Seconds())
}

// ── Bloom Filter & Blocklist ────────────────────────────────────────

func RecordDNSBloomCheck(result string) {
	DNSBloomChecksTotal.WithLabelValues(result).Inc()
}

func RecordDNSBloomFalsePositive() {
	DNSBloomFalsePositivesTotal.Inc()
}

func RecordDNSAllowlistShortCircuit() {
	DNSAllowlistShortCircuitTotal.Inc()
}

func SetDNSBlocklistDomainsLoaded(category string, count float64) {
	DNSBlocklistDomainsLoaded.WithLabelValues(category).Set(count)
}

func SetDNSBlocklistMemoryBytes(category string, bytes float64) {
	DNSBlocklistMemoryBytes.WithLabelValues(category).Set(bytes)
}

func SetDNSBlocklistCategoriesActive(count float64) {
	DNSBlocklistCategoriesActive.Set(count)
}

// ── Blocklist Refresh ───────────────────────────────────────────────

func RecordDNSBlocklistRefresh(category, status string, duration time.Duration) {
	DNSBlocklistRefreshTotal.WithLabelValues(category, status).Inc()
	DNSBlocklistRefreshDuration.WithLabelValues(category).Observe(duration.Seconds())
}

func RecordDNSBlocklistRefreshDomainsFetched(category string, count int) {
	DNSBlocklistRefreshDomainsFetched.WithLabelValues(category).Add(float64(count))
}

func SetDNSBlocklistLastRefreshTimestamp(category string) {
	DNSBlocklistLastRefreshTimestamp.WithLabelValues(category).SetToCurrentTime()
}

// ── DNS Resilience / Rate Limiting ──────────────────────────────────

func RecordDNSStaleCacheHit() {
	DNSStaleCacheHitsTotal.Inc()
}

func RecordDNSNXDOMAINCacheHit() {
	DNSNXDOMAINCacheHitsTotal.Inc()
}

func SetDNSRateLimiterTrackedIPs(count float64) {
	DNSRateLimiterTrackedIPs.Set(count)
}

func RecordDNSRateLimitRejection(tier string) {
	DNSRateLimitRejectionsTotal.WithLabelValues(tier).Inc()
}

func SetDNSCircuitBreakerState(upstream string, state float64) {
	DNSCircuitBreakerState.WithLabelValues(upstream).Set(state)
}

func RecordDNSCircuitBreakerTransition(upstream, from, to string) {
	DNSCircuitBreakerTransitionsTotal.WithLabelValues(upstream, from, to).Inc()
}
