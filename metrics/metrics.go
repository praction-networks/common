package metrics

import (
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Custom registry to avoid default global registry
var registry = prometheus.NewRegistry()

// ----------------------
// HTTP Metrics
// ----------------------
var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status", "pod", "deployment", "namespace", "instance", "service", "content_type", "user_agent"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time taken to process HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status", "pod", "deployment", "namespace", "instance", "service", "content_type", "user_agent"},
	)

	HTTPRequestLatencySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_latency_seconds_summary",
			Help:       "Summary of HTTP request latency in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path", "status", "pod", "deployment", "namespace", "instance", "service", "content_type", "user_agent"},
	)
)

// ----------------------
// System Metrics
// ----------------------
var (
	CPUUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_cpu_usage_percent",
			Help: "Current CPU usage percentage of the process",
		},
	)

	MemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_memory_usage_bytes",
			Help: "Current memory usage of the process in bytes",
		},
	)

	DiskReadBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_read_bytes_total",
			Help: "Total number of bytes read from disk",
		},
	)

	DiskWriteBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_write_bytes_total",
			Help: "Total number of bytes written to disk",
		},
	)

	DiskUsagePercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_usage_percent",
			Help: "Percentage of disk space used on root partition",
		},
	)
)

// ----------------------
// Runtime & Build Info
// ----------------------
var (
	RuntimeVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_go_runtime_info",
			Help: "Go runtime version info",
		},
		[]string{"go_version"},
	)

	BuildVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_build_info",
			Help: "App build info (label-only metric)",
		},
		[]string{"version", "git_commit", "build_time"},
	)
)

// ----------------------
// NATS / Events
// ----------------------
var (
	PublishedEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_published_events_total",
			Help: "Total number of published events, labeled by stream and event",
		},
		[]string{"stream", "event"},
	)

	ProcessedMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_processed_messages_total",
			Help: "Total number of processed messages, labeled by stream and event",
		},
		[]string{"stream", "event"},
	)

	FailedMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_failed_messages_total",
			Help: "Total number of failed messages, labeled by stream and event",
		},
		[]string{"stream", "event"},
	)

	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "nats_message_processing_duration_seconds",
			Help:    "Histogram of message processing durations in seconds, labeled by stream and event",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"stream", "event"},
	)
)

// ----------------------
// Register Metrics (to be called from app)
// ----------------------
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

// ----------------------
// Expose Registry
// ----------------------
func Registry() *prometheus.Registry {
	return registry
}

// func MetricsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

// 		next.ServeHTTP(rw, r)

// 		labels := prometheus.Labels{
// 			"method":       r.Method,
// 			"path":         r.URL.Path,
// 			"status":       strconv.Itoa(rw.status),
// 			"pod":          os.Getenv("POD_NAME"),
// 			"deployment":   os.Getenv("DEPLOYMENT_NAME"),
// 			"namespace":    os.Getenv("POD_NAMESPACE"),
// 			"instance":     os.Getenv("INSTANCE_NAME"),
// 			"service":      os.Getenv("SERVICE_NAME"),
// 			"content_type": r.Header.Get("Content-Type"),
// 			"user_agent":   r.UserAgent(),
// 		}

// 		httpRequestsTotal.With(labels).Inc()
// 		httpRequestDuration.With(labels).Observe(time.Since(start).Seconds())
// 		httpRequestLatencySummary.With(labels).Observe(time.Since(start).Seconds())
// 	})
// }

// func StartMetricsServer(port, metricsPath string) *http.Server {

// 	if port == "" {
// 		port = "9001"
// 	}

// 	metricsStopChan = make(chan struct{})
// 	metricsDoneChan = make(chan struct{})

// 	go collectSystemMetrics()

// 	mux := http.NewServeMux()
// 	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
// 		registry,
// 		promhttp.HandlerFor(registry, promhttp.HandlerOpts{
// 			EnableOpenMetrics: true,
// 		}),
// 	))

// 	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(`{"status":"ok"}`))
// 	})

// 	mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(`{"status":"ready"}`))
// 	})

// 	metricsServer = &http.Server{
// 		Addr:    fmt.Sprintf(":" + port),
// 		Handler: mux,
// 	}

// 	return metricsServer
// }

// func StopMetricsServer() error {
// 	if metricsStopChan != nil {
// 		close(metricsStopChan)
// 	}

// 	if metricsServer != nil {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		if err := metricsServer.Shutdown(ctx); err != nil {
// 			return fmt.Errorf("metrics server shutdown failed: %w", err)
// 		}

// 		<-metricsDoneChan
// 	}

// 	return nil
// }

// func collectSystemMetrics() {
// 	ticker := time.NewTicker(10 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-metricsStopChan:
// 			return
// 		case <-ticker.C:
// 			updateCPUMetrics()
// 			updateMemoryMetrics()
// 			updateDiskIOMetrics()
// 		}
// 	}
// }

// func updateCPUMetrics() {
// 	if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
// 		cpuUsage.Set(percent[0])
// 	}
// }

// func updateMemoryMetrics() {
// 	if memStat, err := mem.VirtualMemory(); err == nil {
// 		memoryUsage.Set(float64(memStat.Used))
// 	}
// }

// func updateDiskIOMetrics() {
// 	if usage, err := disk.Usage("/"); err == nil {
// 		diskUsagePercent.Set(usage.UsedPercent)
// 	}
// 	if ioStats, err := disk.IOCounters(); err == nil {
// 		for _, stat := range ioStats {
// 			diskReadBytes.Set(float64(stat.ReadBytes))
// 			diskWriteBytes.Set(float64(stat.WriteBytes))
// 			break // only the first device
// 		}
// 	} else {
// 		logger.Warn("Failed to fetch disk I/O stats", err)
// 	}
// }

// // responseWriter wraps http.ResponseWriter to capture status code
// type responseWriter struct {
// 	http.ResponseWriter
// 	status int
// }

// func (rw *responseWriter) WriteHeader(statusCode int) {
// 	rw.status = statusCode
// 	rw.ResponseWriter.WriteHeader(statusCode)
// }

// func Handler() http.Handler {
// 	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
// }
