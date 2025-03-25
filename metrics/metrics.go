// common/metrics/metrics.go
package metrics

import (
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var registry = prometheus.NewRegistry()

// HTTP Metrics
var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status", "size", "pod", "deployment", "namespace", "node", "content_type", "client_ip", "user_agent", "protocol"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time taken to process HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status", "size", "pod", "deployment", "namespace", "node", "content_type", "client_ip", "user_agent", "protocol"},
	)

	HTTPRequestLatencySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_latency_seconds_summary",
			Help:       "Summary of HTTP request latency in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path", "status", "size", "pod", "deployment", "namespace", "node", "content_type", "client_ip", "user_agent", "protocol"},
	)
)

// System Metrics
var (
	CPUUsage         = prometheus.NewGauge(prometheus.GaugeOpts{Name: "process_cpu_usage_percent", Help: "Current CPU usage %"})
	MemoryUsage      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "process_memory_usage_bytes", Help: "Current memory usage in bytes"})
	DiskReadBytes    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "disk_read_bytes_total", Help: "Disk read bytes"})
	DiskWriteBytes   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "disk_write_bytes_total", Help: "Disk write bytes"})
	DiskUsagePercent = prometheus.NewGauge(prometheus.GaugeOpts{Name: "disk_usage_percent", Help: "Disk usage %"})
)

// Runtime and Build Info
var (
	RuntimeVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "app_go_runtime_info", Help: "Go runtime info"},
		[]string{"go_version"},
	)
	BuildVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "app_build_info", Help: "App build info"},
		[]string{"version", "git_commit", "build_time"},
	)
)

// NATS Events
var (
	PublishedEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "nats_published_events_total", Help: "Published events"},
		[]string{"stream", "event"},
	)
	ProcessedMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "nats_processed_messages_total", Help: "Processed messages"},
		[]string{"stream", "event"},
	)
	FailedMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "nats_failed_messages_total", Help: "Failed messages"},
		[]string{"stream", "event"},
	)
	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "nats_message_processing_duration_seconds", Help: "Processing durations", Buckets: prometheus.DefBuckets},
		[]string{"stream", "event"},
	)
)

func RegisterAllMetrics() {
	registry.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestLatencySummary,
		CPUUsage,
		MemoryUsage,
		DiskReadBytes,
		DiskWriteBytes,
		DiskUsagePercent,
		RuntimeVersion,
		BuildVersion,
		PublishedEvents,
		ProcessedMessages,
		FailedMessages,
		Duration,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	RuntimeVersion.WithLabelValues(runtime.Version()).Set(1)
	BuildVersion.WithLabelValues(
		os.Getenv("VERSION"),
		os.Getenv("GIT_COMMIT"),
		os.Getenv("BUILD_TIME"),
	).Set(1)
}

func Registry() *prometheus.Registry {
	return registry
}
