package apperror

import (
	"fmt"
)

type ErrorCode int

// Predefined error codes
const (
	// Informational (3xx)
	ErrorCodeProcessing           ErrorCode = 3300 // Processing in progress
	ErrorCodeMovedPermanently     ErrorCode = 3301 // Resource moved permanently (301)
	ErrorCodeFoundRedirect        ErrorCode = 3302 // Temporary redirect (302)
	ErrorCodeNotModified          ErrorCode = 3304 // Resource not modified (304)
	ErrorCodeTemporaryRedirect    ErrorCode = 3307 // Temporary redirect (307)
	ErrorCodeRedirectToLogin      ErrorCode = 3310 // Redirect to login page
	ErrorCodeRedirectToHome       ErrorCode = 3311 // Redirect to home page
	ErrorCodeTemporaryUnavailable ErrorCode = 3312 // Temporary unavailability

	// Client Errors (4xx)
	ErrorCodeInvalidInput         ErrorCode = 4400 // Invalid input or bad request
	ErrorCodeUnauthorized         ErrorCode = 4401 // Unauthorized access
	ErrorCodeForbidden            ErrorCode = 4403 // Forbidden action
	ErrorCodeResourceNotFound     ErrorCode = 4404 // Resource not found (404)
	ErrorCodeEmptyBody            ErrorCode = 4405 // Empty request body
	ErrorCodeNotAcceptable        ErrorCode = 4406 // Not acceptable
	ErrorCodePayloadTooLarge      ErrorCode = 4407 // Payload too large
	ErrorCodeUnsupportedMedia     ErrorCode = 4408 // Unsupported media type
	ErrorCodeResourceConflict     ErrorCode = 4409 // Conflict (e.g., duplicate resource)
	ErrorCodeSessionExpired       ErrorCode = 4410 // Session expired
	ErrorCodeInvalidCredentials   ErrorCode = 4411 // Invalid username or password
	ErrorCodeTokenExpired         ErrorCode = 4412 // Token has expired
	ErrorCodeCSRFValidationFailed ErrorCode = 4413 // CSRF validation failed

	// Security Errors (4xx)
	ErrorCodeInactiveUser    ErrorCode = 4414 // Inactive user
	ErrorCodeDeletedUser     ErrorCode = 4415 // Deleted user
	ErrorCodeTooManyRequests ErrorCode = 4429 // Too many requests (rate-limited)

	// Repository Errors (4xx/5xx)
	ErrorCodeDuplicateKey      ErrorCode = 4500 // Duplicate key error (Conflict)
	ErrorCodeRecordNotFound    ErrorCode = 4504 // Record not found (Repository)
	ErrorCodeQueryTimeout      ErrorCode = 5501 // Query execution timeout
	ErrorCodeTransactionFailed ErrorCode = 5502 // Transaction failure

	// Service Errors (5xx)
	ErrorCodeInternalServerError ErrorCode = 5500 // Generic internal server error
	ErrorCodeDatabaseConnection  ErrorCode = 5503 // Database connection error
	ErrorCodeExternalService     ErrorCode = 5504 // External service failure
	ErrorCodeDependencyFailure   ErrorCode = 5505 // Dependency service failure
	ErrorCodeConfigurationError  ErrorCode = 5506 // Misconfiguration on the server
	ErrorCodeResourceExhausted   ErrorCode = 5507 // Server resources exhausted
	ErrorCodeServiceUnavailable  ErrorCode = 5508 // Service unavailable
	ErrorCodeGatewayTimeout      ErrorCode = 5509 // Gateway timeout
	ErrorCodeInsufficientStorage ErrorCode = 5510 // Insufficient storage
	ErrorCodeUnknown             ErrorCode = 5599 // Unknown server error

	ErrorTypeValidation = "validation"
	ErrorTypeDatabase   = "database"
	ErrorTypeService    = "service"
	ErrorTypeUtils      = "utils"
	ErrorTypeMiddleware = "middleware"
)

// Predefined error codes (unchanged for brevity)

// AppError defines the structure of a standardized application error.
type AppError interface {
	Message() string
	ErrorCode() ErrorCode
	ErrorType() string
	Layer() string
	Unwrap() error
	Error() string // For string representation
}

type appErrorImpl struct {
	message   string
	errorCode ErrorCode
	errorType string
	layer     string
	cause     error
}

func (e *appErrorImpl) Message() string      { return e.message }
func (e *appErrorImpl) ErrorCode() ErrorCode { return e.errorCode }
func (e *appErrorImpl) ErrorType() string    { return e.errorType }
func (e *appErrorImpl) Layer() string        { return e.layer }
func (e *appErrorImpl) Unwrap() error        { return e.cause }
func (e *appErrorImpl) Error() string {
	return fmt.Sprintf("[%d - %s] %s", e.errorCode, e.errorType, e.message)
}

func newAppError(message string, errorCode ErrorCode, errorType, layer string, cause error) AppError {
	return &appErrorImpl{
		message:   message,
		errorCode: errorCode,
		errorType: errorType,
		layer:     layer,
		cause:     cause,
	}
}

func NewRepositoryError(message string, errorCode ErrorCode, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "repository", cause)
}

func NewServiceError(message string, errorCode ErrorCode, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "service", cause)
}

func NewUtilsError(message string, errorCode ErrorCode, errorType string, cause error) AppError {
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
