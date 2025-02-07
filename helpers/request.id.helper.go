package helpers

import (
	"context"
	"net/http"
	"regexp"
	"time"

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
		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"size", wrappedWriter.size,
			"request_body_size", r.ContentLength,
			"duration", time.Since(start),
			"client_ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"referer", r.Referer(),
			"protocol", r.Proto,
			"status_code", wrappedWriter.status,
			"status_description", http.StatusText(wrappedWriter.status),
		)
	})
}

// extractIDFromURL extracts the `{id}` parameter from the URL path if present.
// Assumes that `{id}` is a numeric or alphanumeric segment.
func extractIDFromURL(path string) string {

	// ObjectID regex (MongoDB, 24 hex characters)
	var objectIDRegex = regexp.MustCompile(`\b[0-9a-fA-F]{24}\b`)

	// UUID v4 regex (Valid UUID format)
	// Check for ObjectID in full path
	objectIDMatch := objectIDRegex.FindString(path)
	if objectIDMatch != "" {
		return objectIDMatch
	}

	// Return empty string if no valid ID is found
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
