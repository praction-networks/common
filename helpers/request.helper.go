package helpers

import (
	"context"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
)

type contextKey struct{ name string }

var (
	RequestIDKey = &contextKey{"request_id"}
	UserIDKey    = &contextKey{"user_id"}
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
	numericIDRegex   = regexp.MustCompile(`\b\d+\b`)
	uuidRegex        = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}`)
	emailRegex       = regexp.MustCompile(`[\w\.-]+@[\w\.-]+`)
	hexObjectIDRegex = regexp.MustCompile(`\b[0-9a-f]{24}\b`)
)

// RequestLoggerMiddleware creates a middleware that handles both logging and metrics
// RequestIDMiddleware generates and injects a request ID
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MetricsMiddleware captures metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := metrics.NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		status := rw.Status()
		size := rw.Size()
		handler := sanitizePath(getRoutePattern(r))

		metrics.HTTPRequests.WithLabelValues(r.Method, http.StatusText(status), handler).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, handler).Observe(duration.Seconds())
		metrics.HTTPResponseSizes.WithLabelValues(r.Method, handler).Observe(float64(size))

		if status >= 400 {
			errorType := http.StatusText(status)
			if errorType == "" {
				errorType = "unknown"
			}
			metrics.HTTPErrors.WithLabelValues(r.Method, handler, errorType).Inc()
		}
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := metrics.NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		status := rw.Status()
		size := rw.Size()
		reqID := GetRequestID(r.Context())
		userID := r.Header.Get("X-User-ID")

		fields := []interface{}{
			"reqID", reqID,
			"userID", userID,
			"method", r.Method,
			"path", r.URL.Path,
			"size", size,
			"duration", duration,
			"client_ip", getClientIP(r),
			"user_agent", simplifyUserAgent(r.UserAgent()),
			"status_code", status,
			"protocol", r.Proto,
			"host", r.Host,
		}

		if status >= 500 {
			logger.Error("HTTP server error", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP client error", fields...)
		} else {
			logger.Info("HTTP request", fields...)
		}
	})
}

// getRoutePattern extracts the Chi route pattern
func getRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return r.URL.Path
	}

	pattern := rctx.RoutePattern()
	pattern = strings.ReplaceAll(pattern, "/*/", "/")

	if pattern == "" {
		return r.URL.Path
	}

	return pattern
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

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ipList := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
		if len(ipList) > 0 && ipList[0] != "" {
			ip = strings.TrimSpace(ipList[0])
		}
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	// Anonymize IPv4 addresses
	if parsedIP := net.ParseIP(ip); parsedIP != nil && parsedIP.To4() != nil {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			return parts[0] + "." + parts[1] + ".x.x"
		}
	}
	return ip
}

func sanitizePath(path string) string {
	path = uuidRegex.ReplaceAllString(path, ":uuid")
	path = numericIDRegex.ReplaceAllString(path, ":id")
	path = emailRegex.ReplaceAllString(path, ":email")
	path = hexObjectIDRegex.ReplaceAllString(path, ":objectid")
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	}
	return path
}
