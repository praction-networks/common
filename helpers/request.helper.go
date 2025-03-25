package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/common/logger"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id" // Added context key for X-User-ID
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
