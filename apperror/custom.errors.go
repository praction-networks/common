package apperror

// Predefined error codes
const (
	// Repository Errors
	ErrorCodeDatabaseConnection = 1001
	ErrorCodeDuplicateKey       = 1002
	ErrorCodeRecordNotFound     = 1003

	// Service Errors
	ErrorCodeBusinessLogicFailure = 2001
	ErrorCodeInvalidInput         = 2002

	// Handler Errors
	ErrorCodeRequestParsingFailed = 4001
	ErrorCodeUnauthorized         = 4003

	// Generic Errors
	ErrorCodeInternalServerError = 5000
	ErrorCodeUnknown             = 9999
)

// AppError defines the structure of a standardized application error.
type AppError interface {
	Error() string                   // Returns the error message.
	ErrorCode() int                  // Returns the unique error code.
	ErrorType() string               // Returns the error type (e.g., "database", "validation").
	Layer() string                   // Returns the layer where the error originated.
	Details() map[string]interface{} // Returns additional context for the error.
	Unwrap() error                   // Returns the original error, if any.
}

// appErrorImpl is the concrete implementation of the AppError interface.
type appErrorImpl struct {
	message   string
	errorCode int
	errorType string
	layer     string
	details   map[string]interface{}
	cause     error
}

// Error returns the error message.
func (e *appErrorImpl) Error() string {
	return e.message
}

// ErrorCode returns the unique error code.
func (e *appErrorImpl) ErrorCode() int {
	return e.errorCode
}

// ErrorType returns the error type.
func (e *appErrorImpl) ErrorType() string {
	return e.errorType
}

// Layer returns the layer where the error originated.
func (e *appErrorImpl) Layer() string {
	return e.layer
}

// Details returns additional context for the error.
func (e *appErrorImpl) Details() map[string]interface{} {
	return e.details
}

// Unwrap returns the original error, if any.
func (e *appErrorImpl) Unwrap() error {
	return e.cause
}

// NewAppError creates a new AppError.
func NewAppError(message string, errorCode int, errorType, layer string, details map[string]interface{}, cause error) AppError {
	return &appErrorImpl{
		message:   message,
		errorCode: errorCode,
		errorType: errorType,
		layer:     layer,
		details:   details,
		cause:     cause,
	}
}

// Predefined factory methods for creating layer-specific errors.

// Repository Errors
func NewRepositoryError(message string, errorCode int, cause error) AppError {
	return NewAppError(message, errorCode, "repository", "repository", nil, cause)
}

// Service Errors
func NewServiceError(message string, errorCode int, cause error) AppError {
	return NewAppError(message, errorCode, "service", "service", nil, cause)
}

// Handler Errors
func NewHandlerError(message string, errorCode int, cause error) AppError {
	return NewAppError(message, errorCode, "handler", "handler", nil, cause)
}

// Validation Errors
func NewValidationError(message string, errorCode int, details map[string]interface{}) AppError {
	return NewAppError(message, errorCode, "validation", "service", details, nil)
}

// Utility function to convert a generic error into an AppError.
func AsAppError(err error) (AppError, bool) {
	appErr, ok := err.(AppError)
	return appErr, ok
}

// Example utility to handle unknown errors.
func NewUnknownError(cause error) AppError {
	return NewAppError("An unknown error occurred", ErrorCodeUnknown, "unknown", "unknown", nil, cause)
}
