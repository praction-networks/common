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
	"go.uber.org/zap"
)

type contextKey struct{ name string }

var (
	RequestIDKey         = &contextKey{"request_id"}
	UserIDKey            = &contextKey{"user_id"}
	TenantIDKey          = &contextKey{"tenant_id"}
	AccessibleTenantsKey = &contextKey{"accessible_tenants"}
)


var (
	chromeRegex         = regexp.MustCompile(`Chrome/([\d\.]+)`)
	firefoxRegex        = regexp.MustCompile(`Firefox/([\d\.]+)`)
	safariRegex         = regexp.MustCompile(`Version/([\d\.]+).*Safari`)
	edgeRegex           = regexp.MustCompile(`Edg(?:e|A)/([\d\.]+)`)
	operaRegex          = regexp.MustCompile(`OPR/([\d\.]+)`)
	ieRegex             = regexp.MustCompile(`MSIE ([\d\.]+)`)
	curlRegex           = regexp.MustCompile(`curl/([\d\.]+)`)
	wgetRegex           = regexp.MustCompile(`Wget/([\d\.]+)`)
	goHTTPRegex         = regexp.MustCompile(`Go-http-client/([\d\.]+)`)
	postmanRegex        = regexp.MustCompile(`PostmanRuntime/([\d\.]+)`)
	insomniaRegex       = regexp.MustCompile(`Insomnia/([\d\.]+)`)
	httpieRegex         = regexp.MustCompile(`HTTPie/([\d\.]+)`)
	pythonRequestsRegex = regexp.MustCompile(`python-requests/([\d\.]+)`)
	axiosRegex          = regexp.MustCompile(`axios/([\d\.]+)`)
	nodeFetchRegex      = regexp.MustCompile(`node-fetch/([\d\.]+)`)
	androidRegex        = regexp.MustCompile(`Android ([\d\.]+)`)
	iosSafariRegex      = regexp.MustCompile(`OS ([\d_]+).*like Mac OS X`)
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

// TenantIDMiddleware extracts tenant ID from X-Tenant-ID header and adds it to request context
func TenantIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID != "" {
			ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// MetricsMiddleware captures metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip wrapping ResponseWriter for WebSocket upgrade requests to preserve http.Hijacker interface
		if isWebSocketUpgrade(r) {
			next.ServeHTTP(w, r)
			return
		}

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
		// Skip wrapping ResponseWriter for WebSocket upgrade requests to preserve http.Hijacker interface
		if isWebSocketUpgrade(r) {
			reqID := GetRequestID(r.Context())
			userID := r.Header.Get("X-User-ID")

			// Set default request logger with context fields (only if not empty)
			requestLoggerFields := []zap.Field{}
			if reqID != "" {
				requestLoggerFields = append(requestLoggerFields, zap.String("reqID", reqID))
			}
			if userID != "" {
				requestLoggerFields = append(requestLoggerFields, zap.String("userId", userID))
			}
			if len(requestLoggerFields) > 0 {
				logger.SetDefaultRequestLogger(requestLoggerFields...)
				defer logger.ClearDefaultRequestLogger()
			}

			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rw := metrics.NewResponseWriter(w)

		reqID := GetRequestID(r.Context())
		userID := r.Header.Get("X-User-ID")

		// Set default request logger with context fields (only if not empty)
		requestLoggerFields := []zap.Field{}
		if reqID != "" {
			requestLoggerFields = append(requestLoggerFields, zap.String("reqID", reqID))
		}
		if userID != "" {
			requestLoggerFields = append(requestLoggerFields, zap.String("userId", userID))
		}
		if len(requestLoggerFields) > 0 {
			logger.SetDefaultRequestLogger(requestLoggerFields...)
			defer logger.ClearDefaultRequestLogger()
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		status := rw.Status()
		size := rw.Size()

		// Skip logging for health check endpoints to reduce noise
		if strings.HasSuffix(r.URL.Path, "/health") {
			return
		}

		fields := []interface{}{
			"reqID", reqID,
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

// isWebSocketUpgrade checks if the request is a WebSocket upgrade request
func isWebSocketUpgrade(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetTenantID retrieves the tenant ID from the context
func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// SetAccessibleTenants sets the list of accessible tenant IDs in the context
// This includes the context tenant itself + all its descendants
func SetAccessibleTenants(ctx context.Context, tenantIDs []string) context.Context {
	return context.WithValue(ctx, AccessibleTenantsKey, tenantIDs)
}

// GetAccessibleTenants retrieves the list of accessible tenant IDs from the context
// This list typically includes the user's context tenant + all its descendant tenants
// Returns nil if not set (e.g., for system users who have access to all tenants)
func GetAccessibleTenants(ctx context.Context) []string {
	if tenantIDs, ok := ctx.Value(AccessibleTenantsKey).([]string); ok {
		return tenantIDs
	}
	return nil
}


func simplifyUserAgent(ua string) string {
	// Check for mobile browsers first (they often contain "Mobile" or "Android" or "iPhone")
	if strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") {
		if match := iosSafariRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "iOS/" + strings.ReplaceAll(match[1], "_", ".")
		}
		return "iOS"
	}
	if strings.Contains(ua, "Android") {
		if match := androidRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Android/" + match[1]
		}
		return "Android"
	}

	// Desktop browsers
	switch {
	case strings.Contains(ua, "Edg") || strings.Contains(ua, "Edge"):
		if match := edgeRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Edge/" + match[1]
		}
		return "Edge"
	case strings.Contains(ua, "OPR") || strings.Contains(ua, "Opera"):
		if match := operaRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Opera/" + match[1]
		}
		return "Opera"
	case strings.Contains(ua, "MSIE") || strings.Contains(ua, "Trident"):
		if match := ieRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "IE/" + match[1]
		}
		return "IE"
	case strings.Contains(ua, "Chrome") && !strings.Contains(ua, "Edg"):
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
	// API testing tools and HTTP clients
	case strings.Contains(ua, "Postman"):
		if match := postmanRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Postman/" + match[1]
		}
		return "Postman"
	case strings.Contains(ua, "Insomnia"):
		if match := insomniaRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Insomnia/" + match[1]
		}
		return "Insomnia"
	case strings.Contains(ua, "HTTPie"):
		if match := httpieRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "HTTPie/" + match[1]
		}
		return "HTTPie"
	// Command-line tools
	case strings.Contains(ua, "curl"):
		if match := curlRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "curl/" + match[1]
		}
		return "curl"
	case strings.Contains(ua, "Wget"):
		if match := wgetRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Wget/" + match[1]
		}
		return "Wget"
	// Programming language HTTP clients
	case strings.Contains(ua, "python-requests"):
		if match := pythonRequestsRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "python-requests/" + match[1]
		}
		return "python-requests"
	case strings.Contains(ua, "axios"):
		if match := axiosRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "axios/" + match[1]
		}
		return "axios"
	case strings.Contains(ua, "node-fetch"):
		if match := nodeFetchRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "node-fetch/" + match[1]
		}
		return "node-fetch"
	case strings.Contains(ua, "Go-http-client"):
		if match := goHTTPRegex.FindStringSubmatch(ua); len(match) == 2 {
			return "Go-http-client/" + match[1]
		}
		return "Go-http-client"
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
