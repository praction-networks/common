package apperror

import (
	"fmt"
)

// Predefined error codes
const (
	// Informational (3xx)
	ErrorCodeProcessing           = 3300 // Processing in progress
	ErrorCodeMovedPermanently     = 3301 // Resource moved permanently (301)
	ErrorCodeFoundRedirect        = 3302 // Temporary redirect (302)
	ErrorCodeNotModified          = 3304 // Resource not modified (304)
	ErrorCodeTemporaryRedirect    = 3307 // Temporary redirect (307)
	ErrorCodeRedirectToLogin      = 3310 // Redirect to login page
	ErrorCodeRedirectToHome       = 3311 // Redirect to home page
	ErrorCodeTemporaryUnavailable = 3312 // Temporary unavailability

	// Client Errors (4xx)
	ErrorCodeInvalidInput         = 4400 // Invalid input or bad request
	ErrorCodeUnauthorized         = 4401 // Unauthorized access
	ErrorCodeForbidden            = 4403 // Forbidden action
	ErrorCodeResourceNotFound     = 4404 // Resource not found (404)
	ErrorCodeEmptyBody            = 4405 // Empty request body
	ErrorCodeNotAcceptable        = 4406 // Not acceptable
	ErrorCodePayloadTooLarge      = 4407 // Payload too large
	ErrorCodeUnsupportedMedia     = 4408 // Unsupported media type
	ErrorCodeResourceConflict     = 4409 // Conflict (e.g., duplicate resource)
	ErrorCodeSessionExpired       = 4410 // Session expired
	ErrorCodeInvalidCredentials   = 4411 // Invalid username or password
	ErrorCodeTokenExpired         = 4412 // Token has expired
	ErrorCodeCSRFValidationFailed = 4413 // CSRF validation failed

	// Security Errors (4xx)
	ErrorCodeInactiveUser    = 4414 // Inactive user
	ErrorCodeDeletedUser     = 4415 // Deleted user
	ErrorCodeTooManyRequests = 4429 // Too many requests (rate-limited)

	// Repository Errors (4xx/5xx)
	ErrorCodeDuplicateKey      = 4500 // Duplicate key error (Conflict)
	ErrorCodeRecordNotFound    = 4504 // Record not found (Repository)
	ErrorCodeQueryTimeout      = 5501 // Query execution timeout
	ErrorCodeTransactionFailed = 5502 // Transaction failure

	// Service Errors (5xx)
	ErrorCodeInternalServerError = 5500 // Generic internal server error
	ErrorCodeDatabaseConnection  = 5503 // Database connection error
	ErrorCodeExternalService     = 5504 // External service failure
	ErrorCodeDependencyFailure   = 5505 // Dependency service failure
	ErrorCodeConfigurationError  = 5506 // Misconfiguration on the server
	ErrorCodeResourceExhausted   = 5507 // Server resources exhausted
	ErrorCodeServiceUnavailable  = 5508 // Service unavailable
	ErrorCodeGatewayTimeout      = 5509 // Gateway timeout
	ErrorCodeInsufficientStorage = 5510 // Insufficient storage
	ErrorCodeUnknown             = 5599 // Unknown server error
)

// Predefined error types
const (
	ErrorTypeValidation = "validation"
	ErrorTypeDatabase   = "database"
	ErrorTypeService    = "service"
	ErrorTypeUtils      = "utils"
)

// Predefined error codes (unchanged for brevity)

// AppError defines the structure of a standardized application error.
type AppError interface {
	Message() string
	ErrorCode() int
	ErrorType() string
	Layer() string
	Unwrap() error
	Error() string // For string representation
}

type appErrorImpl struct {
	message   string
	errorCode int
	errorType string
	layer     string
	cause     error
}

func (e *appErrorImpl) Message() string   { return e.message }
func (e *appErrorImpl) ErrorCode() int    { return e.errorCode }
func (e *appErrorImpl) ErrorType() string { return e.errorType }
func (e *appErrorImpl) Layer() string     { return e.layer }
func (e *appErrorImpl) Unwrap() error     { return e.cause }
func (e *appErrorImpl) Error() string {
	return fmt.Sprintf("[%d - %s] %s", e.errorCode, e.errorType, e.message)
}

func newAppError(message string, errorCode int, errorType, layer string, cause error) AppError {
	return &appErrorImpl{
		message:   message,
		errorCode: errorCode,
		errorType: errorType,
		layer:     layer,
		cause:     cause,
	}
}

func NewRepositoryError(message string, errorCode int, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "repository", cause)
}

func NewServiceError(message string, errorCode int, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "service", cause)
}

func NewUtilsError(message string, errorCode int, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "utils", cause)
}

func AsAppError(err error) (AppError, bool) {
	if err == nil {
		return nil, false
	}
	appErr, ok := err.(AppError)
	return appErr, ok
}

func NewUnknownError(cause error) AppError {
	return newAppError("An unknown error occurred", ErrorCodeUnknown, "unknown", "unknown", cause)
}
