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
	ErrorCodeEmptyBody           = 9000
	ErrorCodeInactiveUser        = 9001
	ErrorCodeDeletedUser         = 9002
)

// AppError defines the structure of a standardized application error.
type AppError interface {
	Error() string     // Returns the error message.
	ErrorCode() int    // Returns the unique error code.
	ErrorType() string // Returns the error type (e.g., "database", "validation").
	Layer() string     // Returns the layer where the error originated.
	Unwrap() error     // Returns the original error, if any.
}

type appErrorImpl struct {
	message   string
	errorCode int
	errorType string
	layer     string
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

// Unwrap returns the original error, if any.
func (e *appErrorImpl) Unwrap() error {
	return e.cause
}

// newAppError is an internal helper to create a new AppError.
func newAppError(
	message string,
	errorCode int,
	errorType string,
	layer string,
	cause error,
) AppError {
	return &appErrorImpl{
		message:   message,
		errorCode: errorCode,
		errorType: errorType,
		layer:     layer,
		cause:     cause,
	}
}

// Public factory methods for creating layer-specific errors.

// NewRepositoryError creates an AppError for the repository layer.
func NewRepositoryError(message string, errorCode int, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "repository", cause)
}

// NewServiceError creates an AppError for the service layer.
func NewServiceError(message string, errorCode int, errorType string, cause error) AppError {
	return newAppError(message, errorCode, errorType, "service", cause)
}

// AsAppError converts a generic error into an AppError if possible.
func AsAppError(err error) (AppError, bool) {
	appErr, ok := err.(AppError)
	return appErr, ok
}

// NewUnknownError creates an AppError for unknown errors.
func NewUnknownError(cause error) AppError {
	return newAppError("An unknown error occurred", ErrorCodeUnknown, "unknown", "unknown", cause)
}
