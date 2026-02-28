package audit

import (
	"context"
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

// extractResource extracts the resource name from the URL path
// e.g., /api/v1/subscribers/sub-001 → "subscriber"
func extractResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Skip "api" and version prefix
	for i, part := range parts {
		if part == "api" || strings.HasPrefix(part, "v") {
			continue
		}
		if i > 0 {
			// Return singular form of the first meaningful path segment
			resource := parts[i]
			resource = strings.TrimSuffix(resource, "s") // naive singularize
			return resource
		}
	}
	if len(parts) > 0 {
		return parts[len(parts)-1]
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
