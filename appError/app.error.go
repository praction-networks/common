package appError

import (
	"fmt"
)

// AppError represents a structured error with an error code and HTTP status.
type AppError struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	HTTPCode int       `json:"http_code"`
	Err      error     `json:"-"`
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows unwrapping nested errors.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with a code, message, HTTP status, and optional wrapped error.
func New(code ErrorCode, message string, httpCode int, err error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
		Err:      err,
	}
}

// Predefined error codes for common scenarios

type ErrorCode string

const (
	// General Errors
	EntityNotFound                   ErrorCode = "ENTITY_NOT_FOUND"         // Entity was not found in the database or system
	DuplicateEntityFound             ErrorCode = "DUPLICATE_ENTITY_FOUND"   // Duplicate entity exists in the database
	ValidationFailed                 ErrorCode = "VALIDATION_FAILED"        // Validation error occurred
	UnauthorizedAccess               ErrorCode = "UNAUTHORIZED_ACCESS"      // Unauthorized access to a resource
	PaymentRequired                  ErrorCode = "PAYMENT_REQUIRED"         // Payment is required to proceed
	InsufficientPermissionsErrorCode ErrorCode = "INSUFFICIENT_PERMISSIONS" // User lacks the required permissions
	ResourceConflict                 ErrorCode = "RESOURCE_CONFLICT"        // Resource conflict occurred (e.g., version mismatch)
	RateLimitExceeded                ErrorCode = "RATE_LIMIT_EXCEEDED"      // Rate limit exceeded for API or service
	InternalServerError              ErrorCode = "INTERNAL_SERVER_ERROR"    // Generic server error

	// Database Errors
	DBConnectionError   ErrorCode = "DB_CONNECTION_ERROR"   // Failed to connect to the database
	DBConnectionTimeOut ErrorCode = "DB_CONNECTION_TIMEOUT" // Database connection timeout
	DBFetchError        ErrorCode = "DB_FETCH_ERROR"        // Error while fetching data from the database
	DBInsertError       ErrorCode = "DB_INSERT_ERROR"       // Error while inserting data into the database
	DBUpdateError       ErrorCode = "DB_UPDATE_ERROR"       // Error while updating data in the database
	DBDeleteError       ErrorCode = "DB_DELETE_ERROR"       // Error while deleting data from the database
	DBTransactionError  ErrorCode = "DB_TRANSACTION_ERROR"  // Error during a database transaction

	// NATS (Messaging System) Errors
	NATSConnectionError   ErrorCode = "NATS_CONNECTION_ERROR"   // Failed to connect to NATS
	NATSSubscriptionError ErrorCode = "NATS_SUBSCRIPTION_ERROR" // Failed to subscribe to a NATS subject
	NATSPublishError      ErrorCode = "NATS_PUBLISH_ERROR"      // Failed to publish a message to NATS
	NATSRequestTimeout    ErrorCode = "NATS_REQUEST_TIMEOUT"    // NATS request timed out

	// Redis Errors
	RedisConnectionError            ErrorCode = "REDIS_CONNECTION_ERROR"   // Failed to connect to Redis
	RedisConnectionTimeoutErrorCode ErrorCode = "REDIS_CONNECTION_TIMEOUT" // Redis connection timeout
	RedisFetchError                 ErrorCode = "REDIS_FETCH_ERROR"        // Error while fetching data from Redis
	RedisSetError                   ErrorCode = "REDIS_SET_ERROR"          // Error while setting data in Redis
	RedisDeleteError                ErrorCode = "REDIS_DELETE_ERROR"       // Error while deleting data from Redis
	RedisKeyNotFound                ErrorCode = "REDIS_KEY_NOT_FOUND"      // Key not found in Redis
	RedisTransactionError           ErrorCode = "REDIS_TRANSACTION_ERROR"  // Error during a Redis transaction

	// Authentication & Authorization Errors
	InvalidCredentials ErrorCode = "INVALID_CREDENTIALS" // Provided credentials are invalid
	TokenExpired       ErrorCode = "TOKEN_EXPIRED"       // Token has expired
	TokenInvalid       ErrorCode = "TOKEN_INVALID"       // Token is invalid
	AccessDenied       ErrorCode = "ACCESS_DENIED"       // Access is denied to the requested resource
	SessionExpired     ErrorCode = "SESSION_EXPIRED"     // User session has expired

	// File Handling Errors
	FileNotFound      ErrorCode = "FILE_NOT_FOUND"      // File not found
	FileUploadError   ErrorCode = "FILE_UPLOAD_ERROR"   // Error during file upload
	FileDownloadError ErrorCode = "FILE_DOWNLOAD_ERROR" // Error during file download
	FileDeleteError   ErrorCode = "FILE_DELETE_ERROR"   // Error during file deletion

	// External Service Errors
	ExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR" // Error occurred while interacting with an external service
	ThirdPartyAPIError   ErrorCode = "THIRD_PARTY_API_ERROR"  // Third-party API returned an error
	TimeoutError         ErrorCode = "TIMEOUT_ERROR"          // Request timed out while communicating with an external service
	EmailSendError       ErrorCode = "EMAIL_SEND_ERROR"       // Error while sending an email
	SMSDeliveryError     ErrorCode = "SMS_DELIVERY_ERROR"     // Error while delivering an SMS

	// Configuration Errors
	ConfigLoadError       ErrorCode = "CONFIG_LOAD_ERROR"       // Error loading configuration
	ConfigValidationError ErrorCode = "CONFIG_VALIDATION_ERROR" // Configuration validation failed

	// User Input Errors
	InvalidInputError    ErrorCode = "INVALID_INPUT_ERROR"    // User input is invalid
	MissingRequiredField ErrorCode = "MISSING_REQUIRED_FIELD" // A required field is missing in the input
	UnsupportedOperation ErrorCode = "UNSUPPORTED_OPERATION"  // Operation is not supported
	InvalidOperation     ErrorCode = "INVALID_OPERATION"      // Operation is invalid
)
