package helpers

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/common/logger"
	"go.uber.org/zap"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	URLIDKey     contextKey = "object_id"
)

// RequestLoggerMiddleware logs details about each HTTP request and response, including a unique request ID.
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Extract or generate a unique request ID
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" || !IsUUIDv4(reqID) {
			reqID = uuid.New().String()
			// Add the generated request ID to the request header
			r.Header.Set("X-Request-ID", reqID)
		}
		// Extract `{id}` from the URL if present
		id := extractIDFromURL(r.URL.Path)

		// Set the request ID into the context
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		ctx = context.WithValue(ctx, URLIDKey, id)

		r = r.WithContext(ctx)

		// Create a logger with RequestID
		// Set the request logger

		if id != "" {
			logger.SetDefaultRequestLogger(
				zap.String("Request-ID", reqID),
				zap.String("Object-ID", id),
			)
		} else {
			logger.SetDefaultRequestLogger(zap.String("Request-ID", reqID))
		}

		// Clean up the logger at the end of the request
		defer logger.ClearDefaultRequestLogger()

		w.Header().Set("X-Request-ID", reqID) // Set `X-Request-ID` in the response header

		// Wrap the response writer to capture status and size
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Log the request details
		logger.Info("HTTP request", "method", r.Method, "path", r.URL.Path, "size", wrappedWriter.size, "duration", time.Since(start), "client_ip", r.RemoteAddr, "user_agent", r.UserAgent(), "status_Code", wrappedWriter.status)
	})
}

// extractIDFromURL extracts the `{id}` parameter from the URL path if present.
// Assumes that `{id}` is a numeric or alphanumeric segment.
func extractIDFromURL(path string) string {
	// Example regex to capture numeric IDs: `/resource/{id}`
	re := regexp.MustCompile(`/([^/]+)/?`)
	matches := re.FindAllStringSubmatch(path, -1)
	if len(matches) > 0 {
		return matches[len(matches)-1][1] // Get the last segment, assuming it's the ID
	}
	return ""
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

func IsUUIDv4(s string) bool {
	// Regular expression for UUID v4
	uuidV4Regex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
	return uuidV4Regex.MatchString(s)
}

func GetRequestID(ctx context.Context) string {
	reqID, ok := ctx.Value(RequestIDKey).(string)

	if !ok {
		logger.Warn("Request ID is not a string or missing in context")
		return ""
	}
	return reqID
}

func GetURLID(ctx context.Context) string {
	id, ok := ctx.Value(URLIDKey).(string)
	if !ok {
		logger.Warn("ID is not found in context")
		return ""
	}
	return id
}
