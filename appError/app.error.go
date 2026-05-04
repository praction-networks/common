// Package appError defines the canonical structured error type used across all
// services in this monorepo. AppError carries a stable machine-readable Code,
// a human Message, an HTTP status, and an optional wrapped cause. It is the
// single shape every service layer should return to its handlers; the HTTP
// transport layer (see common/helpers and common/response) renders it.
//
// Two matching idioms are supported and equivalent:
//
//	appError.HasCode(err, appError.EntityNotFound)   // by-code (preferred)
//	errors.Is(err, appError.ErrEntityNotFound)       // by-sentinel
//
// `errors.Is` works for every ErrorCode, not just the explicitly named
// sentinels — see Sentinel and the auto-population in init.
package appError

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

// AppError is a structured error with a stable code, a user-facing message,
// an HTTP status, and an optional wrapped cause.
//
// The Err field is omitted from JSON to avoid leaking internal failure detail
// to API clients; log it on the server side via the Error/Unwrap methods.
type AppError struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	HTTPCode int       `json:"http_code"`
	Err      error     `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap exposes the wrapped cause for errors.Is / errors.As traversal.
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Is reports whether target matches this AppError's code-derived sentinel.
// This makes `errors.Is(appErr, ErrEntityNotFound)` work even when AppError
// was constructed without an explicit cause.
func (e *AppError) Is(target error) bool {
	if e == nil || target == nil {
		return false
	}
	if s := Sentinel(e.Code); s != nil && errors.Is(s, target) {
		return true
	}
	return false
}

// MarshalJSON renders AppError without leaking the wrapped cause.
func (e *AppError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	return json.Marshal(struct {
		Code     ErrorCode `json:"code"`
		Message  string    `json:"message"`
		HTTPCode int       `json:"http_code"`
	}{e.Code, e.Message, e.HTTPCode})
}

// New constructs an AppError. Pass httpCode == 0 to auto-derive the HTTP
// status from `code` (see DefaultHTTPCode). If err is nil, the canonical
// sentinel for `code` is wrapped so errors.Is keeps working.
func New(code ErrorCode, message string, httpCode int, err error) *AppError {
	if httpCode == 0 {
		httpCode = DefaultHTTPCode(code)
	}
	if err == nil {
		err = Sentinel(code)
	}
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
		Err:      err,
	}
}

// Newf is New with printf-style message formatting.
func Newf(code ErrorCode, httpCode int, err error, format string, args ...any) *AppError {
	return New(code, fmt.Sprintf(format, args...), httpCode, err)
}

// Wrap wraps an existing error as an AppError, auto-deriving HTTP status
// from the code. Use this at adapter boundaries (DB, NATS, HTTP clients)
// where you have a cause but no specific HTTP status in mind.
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}
	if ae := As(err); ae != nil {
		return ae
	}
	return New(code, message, DefaultHTTPCode(code), err)
}

// As returns err as *AppError if it is one (anywhere in the unwrap chain),
// or nil otherwise. Convenience wrapper around errors.As.
func As(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return nil
}

// HasCode reports whether err is (or wraps) an AppError with the given code.
func HasCode(err error, code ErrorCode) bool {
	if ae := As(err); ae != nil {
		return ae.Code == code
	}
	return false
}

// HasAnyCode reports whether err is (or wraps) an AppError whose code is in codes.
func HasAnyCode(err error, codes ...ErrorCode) bool {
	ae := As(err)
	if ae == nil {
		return false
	}
	for _, c := range codes {
		if ae.Code == c {
			return true
		}
	}
	return false
}

// CodeOf returns the ErrorCode of err if it is an AppError, or empty otherwise.
func CodeOf(err error) ErrorCode {
	if ae := As(err); ae != nil {
		return ae.Code
	}
	return ""
}

// HTTPStatusOf returns the HTTP status of err if it is an AppError, or
// http.StatusInternalServerError otherwise.
func HTTPStatusOf(err error) int {
	if ae := As(err); ae != nil {
		return ae.HTTPCode
	}
	return http.StatusInternalServerError
}

// ----------------------------------------------------------------------------
// Predicates: high-level classifications used by retry / fallback / fan-out
// logic that doesn't care about the specific code.
// ----------------------------------------------------------------------------

// IsNotFound reports whether err represents a missing entity / file / session.
func IsNotFound(err error) bool {
	return HasAnyCode(err, EntityNotFound, FileNotFound, SessionNotFound, RedisKeyNotFound)
}

// IsValidation reports whether err is a client input / validation failure.
func IsValidation(err error) bool {
	return HasAnyCode(err, ValidationFailed, InvalidInputError, MissingRequiredField, UnsupportedMediaType, PayloadTooLarge)
}

// IsConflict reports whether err is a duplicate / version / idempotency conflict.
func IsConflict(err error) bool {
	return HasAnyCode(err, DuplicateEntityFound, ResourceConflict, VersionMismatch, IdempotencyConflict)
}

// IsAuth reports whether err is an authentication failure (no/invalid identity).
func IsAuth(err error) bool {
	return HasAnyCode(err,
		UnauthorizedAccess, InvalidCredentials, TokenExpired, TokenInvalid,
		TokenRevoked, SessionExpired, SessionNotFound, MfaRequired, InvalidMfaCode,
		AccountLocked, PasswordExpired,
	)
}

// IsAccessDenied reports whether err is an authorization failure (identity OK, not allowed).
func IsAccessDenied(err error) bool {
	return HasAnyCode(err, AccessDenied, Forbidden, InsufficientPermissions)
}

// IsTimeout reports whether err represents a timeout in any layer.
func IsTimeout(err error) bool {
	return HasAnyCode(err,
		TimeoutError, RequestTimeout, GatewayTimeout,
		DBQueryTimeoutError, DBConnectionTimeOut, RedisConnectionTimeout, NATSRequestTimeout,
	)
}

// IsRetryable reports whether err represents a transient failure worth retrying.
// This is a transport-layer hint; respect retry budgets at the call site.
func IsRetryable(err error) bool {
	return HasAnyCode(err,
		DBRetryableError, DBConnectionError, DBConnectionTimeOut, DBQueryTimeoutError,
		RedisConnectionError, RedisConnectionTimeout,
		NATSConnectionError, NATSRequestTimeout,
		ExternalServiceError, ThirdPartyAPIError, TimeoutError,
		ServiceUnavailable, BadGateway, GatewayTimeout, RequestTimeout,
		CircuitBreakerOpen, DependencyFailure,
	)
}

// ----------------------------------------------------------------------------
// Sentinel registry — every ErrorCode resolves to a stable sentinel error,
// so `errors.Is(err, AnySentinel)` works regardless of how AppError was built.
// ----------------------------------------------------------------------------

var (
	sentinelsMu sync.RWMutex
	sentinels   = map[ErrorCode]error{}
)

// Sentinel returns the canonical sentinel for code. The returned error is
// stable for the process lifetime, so it is safe to compare with errors.Is.
func Sentinel(code ErrorCode) error {
	sentinelsMu.RLock()
	s, ok := sentinels[code]
	sentinelsMu.RUnlock()
	if ok {
		return s
	}
	sentinelsMu.Lock()
	defer sentinelsMu.Unlock()
	if s, ok = sentinels[code]; ok {
		return s
	}
	s = errors.New(string(code))
	sentinels[code] = s
	return s
}

// Named sentinels — kept exported for backward compatibility. The init block
// below ensures Sentinel(<code>) returns the same instance, so an AppError
// constructed with code X satisfies `errors.Is(err, ErrX)`.
var (
	ErrEntityNotFound       = errors.New("entity not found")
	ErrDuplicateEntityFound = errors.New("duplicate entity found")
	ErrValidationFailed     = errors.New("validation failed")
	ErrUnauthorizedAccess   = errors.New("unauthorized access")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenExpired         = errors.New("token expired")
	ErrAccessDenied         = errors.New("access denied")
	ErrFileNotFound         = errors.New("file not found")
)

// errorCodeToSentinel is preserved for legacy direct readers, but Sentinel()
// is the source of truth and resolves all codes (not just these 9).
//
// Deprecated: use Sentinel(code) instead.
var errorCodeToSentinel = map[string]error{
	"ENTITY_NOT_FOUND":       ErrEntityNotFound,
	"DUPLICATE_ENTITY_FOUND": ErrDuplicateEntityFound,
	"VALIDATION_FAILED":      ErrValidationFailed,
	"UNAUTHORIZED_ACCESS":    ErrUnauthorizedAccess,
	"RATE_LIMIT_EXCEEDED":    ErrRateLimitExceeded,
	"INVALID_CREDENTIALS":    ErrInvalidCredentials,
	"TOKEN_EXPIRED":          ErrTokenExpired,
	"ACCESS_DENIED":          ErrAccessDenied,
	"FILE_NOT_FOUND":         ErrFileNotFound,
}

func init() {
	// Seed the sentinel map with the named ones so they wire through Sentinel().
	sentinels[EntityNotFound] = ErrEntityNotFound
	sentinels[DuplicateEntityFound] = ErrDuplicateEntityFound
	sentinels[ValidationFailed] = ErrValidationFailed
	sentinels[UnauthorizedAccess] = ErrUnauthorizedAccess
	sentinels[RateLimitExceeded] = ErrRateLimitExceeded
	sentinels[InvalidCredentials] = ErrInvalidCredentials
	sentinels[TokenExpired] = ErrTokenExpired
	sentinels[AccessDenied] = ErrAccessDenied
	sentinels[FileNotFound] = ErrFileNotFound
}

// ----------------------------------------------------------------------------
// ErrorCode catalog
// ----------------------------------------------------------------------------

// ErrorCode is a stable machine-readable identifier for a failure class.
// Codes are SCREAMING_SNAKE_CASE and MUST be append-only across releases —
// existing values are part of the API contract with consumers and clients.
type ErrorCode string

// String implements fmt.Stringer.
func (c ErrorCode) String() string { return string(c) }

const (
	// ----- General -----
	EntityNotFound       ErrorCode = "ENTITY_NOT_FOUND"
	DuplicateEntityFound ErrorCode = "DUPLICATE_ENTITY_FOUND"
	ValidationFailed     ErrorCode = "VALIDATION_FAILED"
	UnauthorizedAccess   ErrorCode = "UNAUTHORIZED_ACCESS"
	Forbidden            ErrorCode = "FORBIDDEN"
	PaymentRequired      ErrorCode = "PAYMENT_REQUIRED"
	ResourceConflict     ErrorCode = "RESOURCE_CONFLICT"
	RateLimitExceeded    ErrorCode = "RATE_LIMIT_EXCEEDED"
	QuotaExceeded        ErrorCode = "QUOTA_EXCEEDED"
	InternalServerError  ErrorCode = "INTERNAL_SERVER_ERROR"
	NotImplemented       ErrorCode = "NOT_IMPLEMENTED"
	ServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"
	BadGateway           ErrorCode = "BAD_GATEWAY"
	GatewayTimeout       ErrorCode = "GATEWAY_TIMEOUT"
	RequestTimeout       ErrorCode = "REQUEST_TIMEOUT"
	RequestCanceled      ErrorCode = "REQUEST_CANCELED"
	MaintenanceMode      ErrorCode = "MAINTENANCE_MODE"
	FeatureDisabled      ErrorCode = "FEATURE_DISABLED"

	// User-permission flavor of Forbidden — kept distinct because the existing
	// codebase uses it to mean "authenticated but lacks role/permission scope".
	InsufficientPermissionsErrorCode ErrorCode = "INSUFFICIENT_PERMISSIONS"

	// ----- Database -----
	DBConnectionError       ErrorCode = "DB_CONNECTION_ERROR"
	DBConnectionTimeOut     ErrorCode = "DB_CONNECTION_TIMEOUT"
	DBFetchError            ErrorCode = "DB_FETCH_ERROR"
	DBInsertError           ErrorCode = "DB_INSERT_ERROR"
	DBUpdateError           ErrorCode = "DB_UPDATE_ERROR"
	DBDeleteError           ErrorCode = "DB_DELETE_ERROR"
	DBTransactionError      ErrorCode = "DB_TRANSACTION_ERROR"
	DBWriteConcernError     ErrorCode = "DB_WRITE_CONCERN_ERROR"
	DBRetryableError        ErrorCode = "DB_RETRYABLE_ERROR"
	DBQueryError            ErrorCode = "DB_QUERY_ERROR"
	DBQueryTimeoutError     ErrorCode = "DB_QUERY_TIMEOUT_ERROR"
	DBQueryCanceledError    ErrorCode = "DB_QUERY_CANCELED_ERROR"
	DBQueryInterruptedError ErrorCode = "DB_QUERY_INTERRUPTED_ERROR"
	DBQueryFailedError      ErrorCode = "DB_QUERY_FAILED_ERROR"
	DBConstraintViolation   ErrorCode = "DB_CONSTRAINT_VIOLATION"
	DBMigrationError        ErrorCode = "DB_MIGRATION_ERROR"
	DBSerializationFailure  ErrorCode = "DB_SERIALIZATION_FAILURE"
	DBDeadlockDetected      ErrorCode = "DB_DEADLOCK_DETECTED"

	// ----- NATS / messaging -----
	NATSConnectionError      ErrorCode = "NATS_CONNECTION_ERROR"
	NATSSubscriptionError    ErrorCode = "NATS_SUBSCRIPTION_ERROR"
	NATSPublishError         ErrorCode = "NATS_PUBLISH_ERROR"
	NATSRequestTimeout       ErrorCode = "NATS_REQUEST_TIMEOUT"
	NATSConsumerError        ErrorCode = "NATS_CONSUMER_ERROR"
	NATSStreamNotFound       ErrorCode = "NATS_STREAM_NOT_FOUND"
	NATSAckError             ErrorCode = "NATS_ACK_ERROR"
	EventDeserializationFail ErrorCode = "EVENT_DESERIALIZATION_FAILED"
	EventHandlerError        ErrorCode = "EVENT_HANDLER_ERROR"

	// ----- Redis -----
	RedisConnectionError            ErrorCode = "REDIS_CONNECTION_ERROR"
	RedisConnectionTimeoutErrorCode ErrorCode = "REDIS_CONNECTION_TIMEOUT"
	RedisFetchError                 ErrorCode = "REDIS_FETCH_ERROR"
	RedisSetError                   ErrorCode = "REDIS_SET_ERROR"
	RedisDeleteError                ErrorCode = "REDIS_DELETE_ERROR"
	RedisKeyNotFound                ErrorCode = "REDIS_KEY_NOT_FOUND"
	RedisTransactionError           ErrorCode = "REDIS_TRANSACTION_ERROR"
	RedisLockError                  ErrorCode = "REDIS_LOCK_ERROR"
	RedisScriptError                ErrorCode = "REDIS_SCRIPT_ERROR"

	// ----- Auth & authorization -----
	InvalidCredentials   ErrorCode = "INVALID_CREDENTIALS"
	TokenExpired         ErrorCode = "TOKEN_EXPIRED"
	TokenInvalid         ErrorCode = "TOKEN_INVALID"
	TokenRevoked         ErrorCode = "TOKEN_REVOKED"
	TokenSignatureFailed ErrorCode = "TOKEN_SIGNATURE_FAILED"
	AccessDenied         ErrorCode = "ACCESS_DENIED"
	SessionExpired       ErrorCode = "SESSION_EXPIRED"
	SessionNotFound      ErrorCode = "SESSION_NOT_FOUND"
	MfaRequired          ErrorCode = "MFA_REQUIRED"
	InvalidMfaCode       ErrorCode = "INVALID_MFA_CODE"
	AccountLocked        ErrorCode = "ACCOUNT_LOCKED"
	AccountDisabled      ErrorCode = "ACCOUNT_DISABLED"
	PasswordExpired      ErrorCode = "PASSWORD_EXPIRED"
	PasswordReuseError   ErrorCode = "PASSWORD_REUSE_FORBIDDEN"
	PasswordPolicyError  ErrorCode = "PASSWORD_POLICY_VIOLATION"
	TenantSuspended      ErrorCode = "TENANT_SUSPENDED"
	APIKeyInvalid        ErrorCode = "API_KEY_INVALID"
	APIKeyExpired        ErrorCode = "API_KEY_EXPIRED"

	// ----- Files / storage -----
	FileNotFound       ErrorCode = "FILE_NOT_FOUND"
	FileUploadError    ErrorCode = "FILE_UPLOAD_ERROR"
	FileDownloadError  ErrorCode = "FILE_DOWNLOAD_ERROR"
	FileDeleteError    ErrorCode = "FILE_DELETE_ERROR"
	FileTooLarge       ErrorCode = "FILE_TOO_LARGE"
	FileTypeNotAllowed ErrorCode = "FILE_TYPE_NOT_ALLOWED"
	StorageQuotaFull   ErrorCode = "STORAGE_QUOTA_FULL"

	// ----- External services -----
	ExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ThirdPartyAPIError   ErrorCode = "THIRD_PARTY_API_ERROR"
	TimeoutError         ErrorCode = "TIMEOUT_ERROR"
	EmailSendError       ErrorCode = "EMAIL_SEND_ERROR"
	SMSDeliveryError     ErrorCode = "SMS_DELIVERY_ERROR"
	WebhookDeliveryError ErrorCode = "WEBHOOK_DELIVERY_ERROR"
	CircuitBreakerOpen   ErrorCode = "CIRCUIT_BREAKER_OPEN"
	RetryExhausted       ErrorCode = "RETRY_EXHAUSTED"
	DependencyFailure    ErrorCode = "DEPENDENCY_FAILURE"

	// ----- Configuration -----
	ConfigLoadError       ErrorCode = "CONFIG_LOAD_ERROR"
	ConfigValidationError ErrorCode = "CONFIG_VALIDATION_ERROR"
	ConfigMissing         ErrorCode = "CONFIG_MISSING"

	// ----- User input / HTTP semantics -----
	BadRequest           ErrorCode = "BAD_REQUEST"
	InvalidInputError    ErrorCode = "INVALID_INPUT_ERROR"
	MissingRequiredField ErrorCode = "MISSING_REQUIRED_FIELD"
	UnsupportedOperation ErrorCode = "UNSUPPORTED_OPERATION"
	InvalidOperation     ErrorCode = "INVALID_OPERATION"
	UnsupportedMediaType ErrorCode = "UNSUPPORTED_MEDIA_TYPE"
	PayloadTooLarge      ErrorCode = "PAYLOAD_TOO_LARGE"
	PreconditionFailed   ErrorCode = "PRECONDITION_FAILED"
	IdempotencyConflict  ErrorCode = "IDEMPOTENCY_KEY_CONFLICT"
	MethodNotAllowed     ErrorCode = "METHOD_NOT_ALLOWED"
	GoneError            ErrorCode = "GONE"

	// ----- Events / concurrency -----
	VersionMismatch        ErrorCode = "VERSION_MISMATCH_ERROR"
	OptimisticLockConflict ErrorCode = "OPTIMISTIC_LOCK_CONFLICT"

	// ----- WebAuthn / passkeys -----
	WebAuthnRegistrationError ErrorCode = "WEBAUTHN_REGISTRATION_ERROR"
	WebAuthnAssertionError    ErrorCode = "WEBAUTHN_ASSERTION_ERROR"
	DeviceNotTrusted          ErrorCode = "DEVICE_NOT_TRUSTED"

	// ----- Crypto / signing -----
	EncryptionFailed        ErrorCode = "ENCRYPTION_FAILED"
	DecryptionFailed        ErrorCode = "DECRYPTION_FAILED"
	SignatureInvalid        ErrorCode = "SIGNATURE_INVALID"
	WebhookSignatureInvalid ErrorCode = "WEBHOOK_SIGNATURE_INVALID"

	// ----- Billing / payments (used by billing-service handlers) -----
	PaymentFailed        ErrorCode = "PAYMENT_FAILED"
	PaymentDeclined      ErrorCode = "PAYMENT_DECLINED"
	RefundFailed         ErrorCode = "REFUND_FAILED"
	SubscriptionInactive ErrorCode = "SUBSCRIPTION_INACTIVE"
	PlanNotFound         ErrorCode = "PLAN_NOT_FOUND"
)

// Convenience aliases for the historically misnamed codes. New code should
// prefer the un-suffixed forms; the suffixed forms are kept for backward
// compatibility with existing call sites across services.
const (
	InsufficientPermissions = InsufficientPermissionsErrorCode
	RedisConnectionTimeout  = RedisConnectionTimeoutErrorCode
	DBQueryFailed           = DBQueryFailedError
	Gone                    = GoneError
)

// ----------------------------------------------------------------------------
// HTTP code defaults — used when New is called with httpCode == 0.
// Callers that pass an explicit code keep doing so; this is purely additive.
// ----------------------------------------------------------------------------

var defaultHTTPCodes = map[ErrorCode]int{
	// 4xx
	EntityNotFound:                   http.StatusNotFound,
	FileNotFound:                     http.StatusNotFound,
	SessionNotFound:                  http.StatusNotFound,
	RedisKeyNotFound:                 http.StatusNotFound,
	PlanNotFound:                     http.StatusNotFound,
	BadRequest:                       http.StatusBadRequest,
	ValidationFailed:                 http.StatusBadRequest,
	InvalidInputError:                http.StatusBadRequest,
	MissingRequiredField:             http.StatusBadRequest,
	InvalidOperation:                 http.StatusBadRequest,
	UnauthorizedAccess:               http.StatusUnauthorized,
	InvalidCredentials:               http.StatusUnauthorized,
	TokenExpired:                     http.StatusUnauthorized,
	TokenInvalid:                     http.StatusUnauthorized,
	TokenRevoked:                     http.StatusUnauthorized,
	TokenSignatureFailed:             http.StatusUnauthorized,
	SessionExpired:                   http.StatusUnauthorized,
	MfaRequired:                      http.StatusUnauthorized,
	InvalidMfaCode:                   http.StatusUnauthorized,
	APIKeyInvalid:                    http.StatusUnauthorized,
	APIKeyExpired:                    http.StatusUnauthorized,
	PaymentRequired:                  http.StatusPaymentRequired,
	Forbidden:                        http.StatusForbidden,
	AccessDenied:                     http.StatusForbidden,
	InsufficientPermissionsErrorCode: http.StatusForbidden,
	AccountLocked:                    http.StatusForbidden,
	AccountDisabled:                  http.StatusForbidden,
	PasswordExpired:                  http.StatusForbidden,
	PasswordReuseError:               http.StatusForbidden,
	PasswordPolicyError:              http.StatusForbidden,
	TenantSuspended:                  http.StatusForbidden,
	DeviceNotTrusted:                 http.StatusForbidden,
	FeatureDisabled:                  http.StatusForbidden,
	MethodNotAllowed:                 http.StatusMethodNotAllowed,
	RequestTimeout:                   http.StatusRequestTimeout,
	DuplicateEntityFound:             http.StatusConflict,
	ResourceConflict:                 http.StatusConflict,
	VersionMismatch:                  http.StatusConflict,
	OptimisticLockConflict:           http.StatusConflict,
	IdempotencyConflict:              http.StatusConflict,
	GoneError:                        http.StatusGone,
	PreconditionFailed:               http.StatusPreconditionFailed,
	PayloadTooLarge:                  http.StatusRequestEntityTooLarge,
	FileTooLarge:                     http.StatusRequestEntityTooLarge,
	UnsupportedMediaType:             http.StatusUnsupportedMediaType,
	UnsupportedOperation:             http.StatusUnprocessableEntity,
	FileTypeNotAllowed:               http.StatusUnsupportedMediaType,
	RateLimitExceeded:                http.StatusTooManyRequests,
	QuotaExceeded:                    http.StatusTooManyRequests,
	RequestCanceled:                  499, // nginx/APISIX convention for client-closed request
	WebhookSignatureInvalid:          http.StatusUnauthorized,
	SignatureInvalid:                 http.StatusUnauthorized,
	WebAuthnRegistrationError:        http.StatusBadRequest,
	WebAuthnAssertionError:           http.StatusUnauthorized,

	// 5xx
	InternalServerError:      http.StatusInternalServerError,
	DBConnectionError:        http.StatusServiceUnavailable,
	DBConnectionTimeOut:      http.StatusGatewayTimeout,
	DBFetchError:             http.StatusInternalServerError,
	DBInsertError:            http.StatusInternalServerError,
	DBUpdateError:            http.StatusInternalServerError,
	DBDeleteError:            http.StatusInternalServerError,
	DBTransactionError:       http.StatusInternalServerError,
	DBWriteConcernError:      http.StatusServiceUnavailable,
	DBRetryableError:         http.StatusServiceUnavailable,
	DBQueryError:             http.StatusInternalServerError,
	DBQueryTimeoutError:      http.StatusGatewayTimeout,
	DBQueryCanceledError:     499,
	DBQueryInterruptedError:  http.StatusInternalServerError,
	DBQueryFailedError:       http.StatusInternalServerError,
	DBConstraintViolation:    http.StatusConflict,
	DBMigrationError:         http.StatusInternalServerError,
	DBSerializationFailure:   http.StatusConflict,
	DBDeadlockDetected:       http.StatusConflict,
	NATSConnectionError:      http.StatusServiceUnavailable,
	NATSSubscriptionError:    http.StatusInternalServerError,
	NATSPublishError:         http.StatusInternalServerError,
	NATSRequestTimeout:       http.StatusGatewayTimeout,
	NATSConsumerError:        http.StatusInternalServerError,
	NATSStreamNotFound:       http.StatusInternalServerError,
	NATSAckError:             http.StatusInternalServerError,
	EventDeserializationFail: http.StatusInternalServerError,
	EventHandlerError:        http.StatusInternalServerError,
	RedisConnectionError:            http.StatusServiceUnavailable,
	RedisConnectionTimeoutErrorCode: http.StatusGatewayTimeout,
	RedisFetchError:                 http.StatusInternalServerError,
	RedisSetError:                   http.StatusInternalServerError,
	RedisDeleteError:                http.StatusInternalServerError,
	RedisTransactionError:           http.StatusInternalServerError,
	RedisLockError:                  http.StatusInternalServerError,
	RedisScriptError:                http.StatusInternalServerError,
	FileUploadError:                 http.StatusInternalServerError,
	FileDownloadError:               http.StatusInternalServerError,
	FileDeleteError:                 http.StatusInternalServerError,
	StorageQuotaFull:                http.StatusInsufficientStorage,
	ExternalServiceError:            http.StatusBadGateway,
	ThirdPartyAPIError:              http.StatusBadGateway,
	TimeoutError:                    http.StatusGatewayTimeout,
	EmailSendError:                  http.StatusBadGateway,
	SMSDeliveryError:                http.StatusBadGateway,
	WebhookDeliveryError:            http.StatusBadGateway,
	CircuitBreakerOpen:              http.StatusServiceUnavailable,
	RetryExhausted:                  http.StatusServiceUnavailable,
	DependencyFailure:               http.StatusBadGateway,
	ConfigLoadError:                 http.StatusInternalServerError,
	ConfigValidationError:           http.StatusInternalServerError,
	ConfigMissing:                   http.StatusInternalServerError,
	NotImplemented:                  http.StatusNotImplemented,
	ServiceUnavailable:              http.StatusServiceUnavailable,
	BadGateway:                      http.StatusBadGateway,
	GatewayTimeout:                  http.StatusGatewayTimeout,
	MaintenanceMode:                 http.StatusServiceUnavailable,
	EncryptionFailed:                http.StatusInternalServerError,
	DecryptionFailed:                http.StatusInternalServerError,
	PaymentFailed:                   http.StatusBadGateway,
	PaymentDeclined:                 http.StatusPaymentRequired,
	RefundFailed:                    http.StatusBadGateway,
	SubscriptionInactive:            http.StatusForbidden,
}

// DefaultHTTPCode returns the canonical HTTP status for code, or 500 if
// unmapped. Use this only as a fallback; explicit status at the call site
// remains preferred when context demands it (e.g. 404 vs 410 for the same
// code).
func DefaultHTTPCode(code ErrorCode) int {
	if c, ok := defaultHTTPCodes[code]; ok {
		return c
	}
	return http.StatusInternalServerError
}
