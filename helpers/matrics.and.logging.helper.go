package helpers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/common/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/process"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed",
		},
		[]string{"method", "path", "status", "pod", "deployment"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status", "pod", "deployment"},
	)

	httpRequestLatencySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_latency_seconds_summary",
			Help:       "Summary of HTTP request latency in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path", "status", "pod", "deployment"},
	)

	goMemStatsAlloc = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_alloc_bytes",
			Help: "Memory allocated and still in use (bytes)",
		},
	)

	goMemStatsSys = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_sys_bytes",
			Help: "Total memory obtained from the OS (bytes)",
		},
	)

	goGoroutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_goroutines",
			Help: "Number of goroutines",
		},
	)

	cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_cpu_usage_percent",
			Help: "CPU usage percentage of the process",
		},
	)

	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_memory_usage_bytes",
			Help: "Memory usage of the process",
		},
	)
	diskUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_usage_percent",
			Help: "Disk usage percentage on root path",
		},
	)
	diskReadBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_read_bytes_total",
			Help: "Total number of bytes read from disk",
		},
	)
	diskWriteBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_write_bytes_total",
			Help: "Total number of bytes written to disk",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestLatencySummary)
	prometheus.MustRegister(goMemStatsAlloc)
	prometheus.MustRegister(goMemStatsSys)
	prometheus.MustRegister(goGoroutines)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)
	prometheus.MustRegister(diskUsage)
	prometheus.MustRegister(diskReadBytes)
	prometheus.MustRegister(diskWriteBytes)

	// Use modern collectors
	prometheus.MustRegister(collectors.NewGoCollector())
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	go collectGoRuntimeMetrics()
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
			r.Header.Set("X-Request-ID", reqID)
		}

		userID := r.Header.Get("X-User-ID")
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", reqID)

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrappedWriter.status)
		pod := os.Getenv("POD_NAME")
		deployment := os.Getenv("DEPLOYMENT_NAME")

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status, pod, deployment).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status, pod, deployment).Observe(duration)
		httpRequestLatencySummary.WithLabelValues(r.Method, r.URL.Path, status, pod, deployment).Observe(duration)

		logger.Info("HTTP request",
			"reqID", reqID,
			"userID", userID,
			"method", r.Method,
			"path", r.URL.Path,
			"origin", r.Header.Get("Origin"),
			"referrer", r.Referer(),
			"size", wrappedWriter.size,
			"duration", time.Since(start),
			"client_ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status_code", wrappedWriter.status,
			"protocol", r.Proto,
			"host", r.Host,
			"content_type", r.Header.Get("Content-Type"),
			"content_length", r.ContentLength,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

func collectGoRuntimeMetrics() {
	pid := int32(os.Getpid())
	proc, _ := process.NewProcess(pid)

	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		goMemStatsAlloc.Set(float64(m.Alloc))
		goMemStatsSys.Set(float64(m.Sys))
		goGoroutines.Set(float64(runtime.NumGoroutine()))

		if cpuPercent, err := proc.CPUPercent(); err == nil {
			cpuUsage.Set(cpuPercent)
		}

		if memInfo, err := proc.MemoryInfo(); err == nil {
			memoryUsage.Set(float64(memInfo.RSS))
		}

		if usage, err := disk.Usage("/"); err == nil {
			diskUsage.Set(usage.UsedPercent)
		}

		if ioStats, err := disk.IOCounters(); err == nil {
			for _, stat := range ioStats {
				diskReadBytes.Set(float64(stat.ReadBytes))
				diskWriteBytes.Set(float64(stat.WriteBytes))
				break
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func ExposePrometheus(port string) {
	go func() {
		mux := http.NewServeMux()

		mux.Handle("/api/v1/domina/metrics", promhttp.Handler())

		mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
				logger.Warn("Failed to write /live response", err)
			}
		})

		mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"status":"ready"}`)); err != nil {
				logger.Warn("Failed to write /readiness response", err)
			}
		})

		logger.Info("Prometheus metrics & probes running",
			"metrics", fmt.Sprintf("http://localhost:%s/api/v1/domina/metrics", port),
			"live", fmt.Sprintf("http://localhost:%s/live", port),
			"readiness", fmt.Sprintf("http://localhost:%s/readiness", port),
		)

		if err := http.ListenAndServe(":"+port, mux); err != nil {
			logger.Error("Failed to start Prometheus metrics server", err)
		}
	}()
}
