package apperror

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// AppError interface defines the structure of a custom application error.
type AppError interface {
	Error() string                     // Returns the error message
	Type() string                      // Returns the error type (e.g., "validation", "database", etc.)
	StatusCode() int                   // Returns the corresponding HTTP status code
	Cause() error                      // Returns the original error, if any
	Metadata() map[string]interface{}  // Returns additional context for the error
	LogFields() map[string]interface{} // Returns structured fields for logging
}

// appErrorImpl is the concrete implementation of the AppError interface.
type appErrorImpl struct {
	message    string
	errType    string
	statusCode int
	cause      error
	metadata   map[string]interface{}
}

// Error returns the error message.
func (e *appErrorImpl) Error() string {
	return e.message
}

// Type returns the error type.
func (e *appErrorImpl) Type() string {
	return e.errType
}

// StatusCode returns the HTTP status code associated with the error.
func (e *appErrorImpl) StatusCode() int {
	return e.statusCode
}

// Cause returns the original error, if any.
func (e *appErrorImpl) Cause() error {
	return e.cause
}

// Metadata returns additional context for the error.
func (e *appErrorImpl) Metadata() map[string]interface{} {
	return e.metadata
}

// LogFields returns structured fields for logging.
func (e *appErrorImpl) LogFields() map[string]interface{} {
	return map[string]interface{}{
		"message":    e.message,
		"type":       e.errType,
		"statusCode": e.statusCode,
		"metadata":   e.metadata,
	}
}

// Factory methods to create specific errors

func NewValidationError(message string, metadata map[string]interface{}) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "validation",
		statusCode: http.StatusBadRequest,
		metadata:   metadata,
	}
}

func NewDatabaseError(message string, cause error) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "database",
		statusCode: http.StatusInternalServerError,
		cause:      cause,
	}
}

func NewNotFoundError(message string) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "not_found",
		statusCode: http.StatusNotFound,
	}
}

func NewConflictError(message string) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "conflict",
		statusCode: http.StatusConflict,
	}
}

func NewAuthenticationError(message string) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "authentication",
		statusCode: http.StatusUnauthorized,
	}
}

func NewAuthorizationError(message string) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "authorization",
		statusCode: http.StatusForbidden,
	}
}

func NewRateLimitError(message string) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "rate_limit",
		statusCode: http.StatusTooManyRequests,
	}
}

func NewInternalError(message string, cause error) AppError {
	return &appErrorImpl{
		message:    message,
		errType:    "internal",
		statusCode: http.StatusInternalServerError,
		cause:      cause,
	}
}

// Utility function to convert generic error to AppError
func AsAppError(err error) (AppError, bool) {
	appErr, ok := err.(AppError)
	return appErr, ok
}

// MongoDB-specific error checks

// IsDuplicateKeyError checks if the error is a MongoDB duplicate key error.
func IsDuplicateKeyError(err error) bool {
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		for _, we := range writeErr.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

// IsNetworkError checks if the error is a network-related MongoDB error.
func IsNetworkError(err error) bool {
	_, isNetworkError := err.(driver.Error)
	return isNetworkError
}

// IsTimeoutError checks if the error is a timeout error.
func IsTimeoutError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout")
}

// IsReadPrefError checks if the error is related to MongoDB read preference.
func IsReadPrefError(err error) bool {
	return strings.Contains(err.Error(), "invalid read preference")
}

// IsWriteConcernError checks if the error is a MongoDB write concern error.
func IsWriteConcernError(err error) bool {
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		for _, we := range writeErr.WriteErrors {
			// Check if the error relates to write concern
			if strings.Contains(we.Message, "write concern error") {
				return true
			}
		}
	}

	// Fallback: Check for generic write concern error messages
	return strings.Contains(err.Error(), "write concern error")
}

// Example usage of the MongoDB error checks
func HandleMongoError(err error) AppError {
	switch {
	case IsDuplicateKeyError(err):
		return NewConflictError("Duplicate key error in MongoDB")
	case IsNetworkError(err):
		return NewInternalError("Network error in MongoDB", err)
	case IsTimeoutError(err):
		return NewInternalError("Timeout error in MongoDB", err)
	case IsReadPrefError(err):
		return NewDatabaseError("Invalid read preference error", err)
	case IsWriteConcernError(err):
		return NewDatabaseError("Write concern error", err)
	default:
		return NewInternalError("Unknown MongoDB error", err)
	}
}
