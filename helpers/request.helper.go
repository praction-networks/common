package helpers

import (
	"context"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id" // Added context key for X-User-ID
)

var (
	chromeRegex  = regexp.MustCompile(`Chrome/([\d\.]+)`)
	firefoxRegex = regexp.MustCompile(`Firefox/([\d\.]+)`)
	safariRegex  = regexp.MustCompile(`Version/([\d\.]+).*Safari`)
	curlRegex    = regexp.MustCompile(`curl/([\d\.]+)`)
	goHTTPRegex  = regexp.MustCompile(`Go-http-client/([\d\.]+)`)
	postmanRegex = regexp.MustCompile(`PostmanRuntime/([\d\.]+)`)
)

var (
	// Common dynamic patterns
	numericIDRegex   = regexp.MustCompile(`\b\d+\b`)
	uuidRegex        = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}`)
	emailRegex       = regexp.MustCompile(`[\w\.-]+@[\w\.-]+`)
	hexObjectIDRegex = regexp.MustCompile(`\b[0-9a-f]{24}\b`) // MongoDB ObjectID
)

// RequestLoggerMiddleware logs details about each HTTP request and response, including a unique request ID and user ID.
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Extract or generate a unique request ID
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
			r.Header.Set("X-Request-ID", reqID) // Set generated request ID in header
		}

		// Extract X-User-ID from request headers
		userID := r.Header.Get("X-User-ID")

		// Set values into the request context
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		ctx = context.WithValue(ctx, UserIDKey, userID)

		r = r.WithContext(ctx)

		// Set X-Request-ID in the response header
		w.Header().Set("X-Request-ID", reqID)

		// Wrap the response writer to capture status and size
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start).Seconds()

		labels := prometheus.Labels{
			"method":       r.Method,
			"path":         sanitizePath(r.URL.Path),
			"status":       strconv.Itoa(wrappedWriter.status),
			"size":         strconv.Itoa(wrappedWriter.size),
			"pod":          os.Getenv("POD_NAME"),
			"deployment":   os.Getenv("DEPLOYMENT_NAME"),
			"namespace":    os.Getenv("POD_NAMESPACE"),
			"node":         os.Getenv("NODE_NAME"),
			"content_type": r.Header.Get("Content-Type"),
			"client_ip":    getClientIP(r),
			"user_agent":   simplifyUserAgent(r.UserAgent()),
			"protocol":     r.Proto,
		}

		// Record Prometheus metrics
		metrics.HTTPRequestsTotal.With(labels).Inc()
		metrics.HTTPRequestDuration.With(labels).Observe(duration)
		metrics.HTTPRequestLatencySummary.With(labels).Observe(duration)

		// Log the request details
		logger.Info("HTTP request",
			"reqID", reqID,
			"userID", userID, // Log user ID
			"method", r.Method,
			"path", r.URL.Path,
			"origin", r.Header.Get("Origin"),
			"referrer", r.Referer(),
			"size", wrappedWriter.size,
			"duration", time.Since(start),
			"client_ip", getClientIP(r),
			"user_agent", simplifyUserAgent(r.UserAgent()),
			"status_code", wrappedWriter.status,
			"protocol", r.Proto,
			"host", r.Host,
			"content_type", r.Header.Get("Content-Type"),
			"content_length", r.ContentLength,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status and size of the response
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

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// simplifyUserAgent extracts a basic UA tag like "Chrome", "curl", etc.
// simplifyUserAgent extracts simplified name + version
func simplifyUserAgent(ua string) string {
	switch {
	case strings.Contains(ua, "Chrome"):
		if match := chromeRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Chrome/" + match[1]
		}
		return "Chrome"
	case strings.Contains(ua, "Firefox"):
		if match := firefoxRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Firefox/" + match[1]
		}
		return "Firefox"
	case strings.Contains(ua, "Safari") && strings.Contains(ua, "Version"):
		if match := safariRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Safari/" + match[1]
		}
		return "Safari"
	case strings.Contains(ua, "curl"):
		if match := curlRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "curl/" + match[1]
		}
		return "curl"
	case strings.Contains(ua, "Go-http-client"):
		if match := goHTTPRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Go-http-client/" + match[1]
		}
		return "Go-http-client"
	case strings.Contains(ua, "Postman"):
		if match := postmanRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Postman/" + match[1]
		}
		return "Postman"
	default:
		return "Other"
	}
}

// getClientIP extracts a shortened client IP to avoid full IP cardinality
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	// Optional: only return first 2 octets to reduce cardinality
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			return parts[0] + "." + parts[1] + ".x.x"
		}
	}
	return ip
}

// sanitizePath replaces dynamic path segments with placeholders to reduce cardinality
func sanitizePath(path string) string {
	path = uuidRegex.ReplaceAllString(path, ":uuid")
	path = numericIDRegex.ReplaceAllString(path, ":id")
	path = emailRegex.ReplaceAllString(path, ":email")
	path = hexObjectIDRegex.ReplaceAllString(path, ":objectid")

	// Optional: collapse repeated slashes, trim trailing slash
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	}
	return path
}
