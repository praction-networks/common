package audit

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

// Middleware returns a Chi-compatible HTTP middleware that automatically
// publishes audit events for every request. Services only need to enrich
// the event with ResourceName and Changes via context.
//
// Usage:
//
//	r := chi.NewRouter()
//	r.Use(audit.Middleware(publisher, "my-service"))
func Middleware(publisher *Publisher, serviceName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip health checks and metrics
			if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			// Wrap response writer to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Serve the request
			next.ServeHTTP(rw, r)

			// Build audit event from request context
			tenantID := r.Header.Get("x-tenant-id")
			userID := r.Header.Get("x-user-id")
			userName := r.Header.Get("x-user-name")

			if tenantID == "" {
				return // skip non-tenant requests
			}

			event := AuditEvent{
				TenantID:   tenantID,
				UserID:     userID,
				UserName:   userName,
				Action:     AuditAction(httpMethodToAction(r.Method)),
				Resource:   extractResource(r.URL.Path),
				Service:    serviceName,
				IPAddress:  extractIP(r),
				UserAgent:  r.UserAgent(),
				Status:     statusFromCode(rw.statusCode),
				StatusCode: rw.statusCode,
				Timestamp:  time.Now().UTC(),
			}

			// Publish asynchronously — don't block the response
			go func() {
				_ = publisher.Publish(context.Background(), event)
			}()
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Flush proxies http.Flusher so SSE handlers wrapped by this middleware
// can stream. Without it, type assertions like w.(http.Flusher) fail
// because the embedded interface is hidden behind the wrapper struct.
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack proxies http.Hijacker so WebSocket upgrade handlers continue
// to work when this middleware sits in front of them.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// httpMethodToAction maps HTTP methods to audit action types
func httpMethodToAction(method string) string {
	switch method {
	case http.MethodPost:
		return "CREATE"
	case http.MethodPut, http.MethodPatch:
		return "UPDATE"
	case http.MethodDelete:
		return "DELETE"
	default:
		return "READ"
	}
}

// pathResourceMap maps URL path segments to their canonical singular
// resource name. Explicit map prevents the naive TrimSuffix("s") from
// mangling words like radius → "radiu" or series → "serie".
var pathResourceMap = map[string]string{
	"subscribers":  "subscriber",
	"plans":        "plan",
	"tenants":      "tenant",
	"tenant-users": "tenant-user",
	"users":        "user",
	"invoices":     "invoice",
	"payments":     "payment",
	"olts":         "olt",
	"onts":         "ont",
	"tickets":      "ticket",
	"radius":       "radius",
	"series":       "series",
	"status":       "status",
	"inventory":    "inventory",
	"assets":       "asset",
	"vendors":      "vendor",
	"warehouses":   "warehouse",
	"stocks":       "stock",
	"products":     "product",
	"alarms":       "alarm",
	"sessions":     "session",
	"licenses":     "license",
	"audit-logs":   "audit-log",
}

// extractResource pulls the first meaningful path segment from a request
// URL and returns its canonical resource name. Skips "api" and version
// prefixes. Falls back to the raw segment (never empty string) so
// unrecognised paths are still grouped consistently in audit queries.
func extractResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, part := range parts {
		if part == "" {
			continue
		}
		if part == "api" {
			continue
		}
		// Version prefixes like v1, v2, v10 — at most 3 chars starting with 'v'
		if strings.HasPrefix(part, "v") && len(part) <= 3 {
			continue
		}
		if canonical, ok := pathResourceMap[part]; ok {
			return canonical
		}
		// Unknown segment — return as-is rather than guessing a singular
		return part
	}
	return "unknown"
}

// extractIP gets the client IP, preferring X-Forwarded-For
func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr (host:port)
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// statusFromCode converts HTTP status code to audit status
func statusFromCode(code int) AuditStatus {
	if code >= 200 && code < 400 {
		return StatusSuccess
	}
	return StatusFailure
}
