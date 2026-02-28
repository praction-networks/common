package audit

import "time"

// AuditAction defines the type of action being audited
type AuditAction string

const (
	ActionCreate  AuditAction = "CREATE"
	ActionRead    AuditAction = "READ"
	ActionUpdate  AuditAction = "UPDATE"
	ActionDelete  AuditAction = "DELETE"
	ActionLogin   AuditAction = "LOGIN"
	ActionLogout  AuditAction = "LOGOUT"
	ActionExport  AuditAction = "EXPORT"
	ActionImport  AuditAction = "IMPORT"
	ActionApprove AuditAction = "APPROVE"
	ActionDeny    AuditAction = "DENY"
	ActionRevoke  AuditAction = "REVOKE"
	ActionGrant   AuditAction = "GRANT"
)

// AuditStatus represents the outcome of the audited action
type AuditStatus string

const (
	StatusSuccess AuditStatus = "SUCCESS"
	StatusFailure AuditStatus = "FAILURE"
	StatusDenied  AuditStatus = "DENIED"
)

// AuditEvent is the canonical audit event structure published by all services.
// It is serialized to JSON and sent to NATS JetStream on AuditGlobalStream.
type AuditEvent struct {
	// ID is a unique identifier for this audit event (UUID v4)
	ID string `json:"id"`

	// TenantID is the tenant context of the action
	TenantID string `json:"tenantId"`

	// UserID is the user who performed the action
	UserID string `json:"userId"`

	// UserName is the display name of the user (optional, for readability)
	UserName string `json:"userName,omitempty"`

	// Action is the type of action performed (CREATE, UPDATE, DELETE, etc.)
	Action AuditAction `json:"action"`

	// Resource is the entity type being acted upon (e.g., "tenant-user", "subscriber", "plan")
	Resource string `json:"resource"`

	// ResourceID is the unique identifier of the affected entity
	ResourceID string `json:"resourceId"`

	// Service is the originating microservice name (e.g., "tenant-user-service")
	Service string `json:"service"`

	// IPAddress is the client IP address from the request
	IPAddress string `json:"ipAddress,omitempty"`

	// UserAgent is the client user agent string
	UserAgent string `json:"userAgent,omitempty"`

	// Status is the outcome of the action (SUCCESS, FAILURE, DENIED)
	Status AuditStatus `json:"status"`

	// StatusCode is the HTTP status code returned
	StatusCode int `json:"statusCode"`

	// Timestamp is when the action occurred
	Timestamp time.Time `json:"timestamp"`

	// Metadata contains additional context (changes, before/after, reason, etc.)
	Metadata map[string]any `json:"metadata,omitempty"`
}
