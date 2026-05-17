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

			resource, resourceID := extractResourceAndID(r.URL.Path)
			event := AuditEvent{
				TenantID:   tenantID,
				UserID:     userID,
				UserName:   userName,
				Action:     AuditAction(httpMethodToAction(r.Method)),
				Resource:   resource,
				ResourceID: resourceID,
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
	"subscribers":   "subscriber",
	"plans":         "plan",
	"tenants":       "tenant",
	"tenant-users":  "tenant-user",
	"users":         "user",
	"invoices":      "invoice",
	"payments":      "payment",
	"olts":          "olt",
	"onts":          "ont",
	"tickets":       "ticket",
	"radius":        "radius",
	"series":        "series",
	"status":        "status",
	"inventory":     "inventory",
	"assets":        "asset",
	"vendors":       "vendor",
	"warehouses":    "warehouse",
	"stocks":        "stock",
	"products":      "product",
	"devices":       "device",
	"subscriptions": "subscription",
	"alarms":        "alarm",
	"sessions":      "session",
	"licenses":      "license",
	"audit-logs":    "audit-log",
}

// extractResource pulls the first meaningful path segment from a request
// URL and returns its canonical resource name. Skips "api" and version
// prefixes. Falls back to the raw segment (never empty string) so
// unrecognised paths are still grouped consistently in audit queries.
func extractResource(path string) string {
	r, _ := extractResourceAndID(path)
	return r
}

// extractResourceAndID walks the path and returns the canonical resource
// name plus the resource ID if the next segment looks like an identifier
// (cuid2, UUID, or numeric). For paths like /api/v1/olts/abc123 the
// returned pair is ("olt", "abc123"). For /api/v1/olts (no id) the pair
// is ("olt", ""). Identifier shapes that we recognise:
//   - cuid2: 24-32 lowercase alphanumeric chars
//   - UUID:  8-4-4-4-12 hex with dashes
//   - numeric:  all digits, length 1-12
//
// Action-style suffixes (e.g. /olts/{id}/sync, /tenants/{id}/verify) are
// preserved on Resource via the existing segment match so callers can
// still tell read/sync/verify apart through the Action field; the ID
// itself is the {id} between resource and verb.
func extractResourceAndID(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	resource := ""
	for i, part := range parts {
		if part == "" || part == "api" {
			continue
		}
		if strings.HasPrefix(part, "v") && len(part) <= 3 {
			continue
		}
		if canonical, ok := pathResourceMap[part]; ok {
			resource = canonical
		} else {
			resource = part
		}
		// Look at the next segment for an identifier.
		if i+1 < len(parts) {
			candidate := parts[i+1]
			if looksLikeID(candidate) {
				return resource, candidate
			}
		}
		return resource, ""
	}
	return "unknown", ""
}

// looksLikeID returns true when s appears to be an entity identifier
// rather than another path segment (sub-resource or action verb).
func looksLikeID(s string) bool {
	if s == "" {
		return false
	}
	// UUID: contains dashes at 8-4-4-4-12 positions
	if len(s) == 36 && s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-' {
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c == '-' {
				continue
			}
			if !isHex(c) {
				return false
			}
		}
		return true
	}
	// cuid2: 24-32 lowercase alphanumeric, no dashes
	if len(s) >= 24 && len(s) <= 32 {
		for i := 0; i < len(s); i++ {
			c := s[i]
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
				return false
			}
		}
		return true
	}
	// Numeric ID
	if len(s) >= 1 && len(s) <= 12 {
		for i := 0; i < len(s); i++ {
			if s[i] < '0' || s[i] > '9' {
				return false
			}
		}
		return true
	}
	return false
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
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
