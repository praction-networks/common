package subscriberevent

import "time"

// VoucherTemplateCreatedEvent represents a voucher template creation event
type VoucherTemplateCreatedEvent struct {
	ID          string `json:"id" bson:"_id"`
	TenantID    string `json:"tenantId" bson:"tenantId"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	PlanID      string `json:"planId,omitempty" bson:"planId,omitempty"`
	Status      string `json:"status" bson:"status"`
	Version     int    `json:"version" bson:"version"`
}

// VoucherBatchStatusChangedEvent represents a batch status change event
type VoucherBatchStatusChangedEvent struct {
	BatchID     string             `json:"batchId" bson:"batchId"`
	TenantID    string             `json:"tenantId" bson:"tenantId"`
	BatchName   string             `json:"batchName" bson:"batchName"`
	OldStatus   VoucherBatchStatus `json:"oldStatus" bson:"oldStatus"`
	NewStatus   VoucherBatchStatus `json:"newStatus" bson:"newStatus"`
	Count       int64              `json:"count" bson:"count"` // Number of vouchers affected
	ChangedAt   time.Time          `json:"changedAt" bson:"changedAt"`
	ChangedBy   string             `json:"changedBy,omitempty" bson:"changedBy,omitempty"`
}

// VoucherTemplateUpdatedEvent represents a voucher template update event
type VoucherTemplateUpdatedEvent struct {
	ID          string `json:"id" bson:"_id"`
	TenantID    string `json:"tenantId" bson:"tenantId"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty"`
	Version     int    `json:"version" bson:"version"`
}

// VoucherTemplateDeletedEvent represents a voucher template deletion event
type VoucherTemplateDeletedEvent struct {
	ID        string    `json:"id" bson:"_id"`
	TenantID  string    `json:"tenantId" bson:"tenantId"`
	DeletedAt time.Time `json:"deletedAt" bson:"deletedAt"`
	DeletedBy string    `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
}

// VoucherInstanceCreatedEvent represents a voucher instance creation event
// Captive-portal caches: voucherCode, status, expiresAt, planId for fast validation
type VoucherInstanceCreatedEvent struct {
	ID           string     `json:"id" bson:"_id"`
	TenantID     string     `json:"tenantId" bson:"tenantId"`
	TemplateID   string     `json:"templateId" bson:"templateId"`
	PlanID       string     `json:"planId" bson:"planId"`
	VoucherCode  string     `json:"voucherCode" bson:"voucherCode"`
	Status       string     `json:"status" bson:"status"`
	IssuedAt     *time.Time `json:"issuedAt,omitempty" bson:"issuedAt,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt" bson:"expiresAt"`
	BatchID      *string    `json:"batchId,omitempty" bson:"batchId,omitempty"`
	Distribution string     `json:"distributionType" bson:"distributionType"`
	Version      int        `json:"version" bson:"version"`
}

// VoucherInstanceBulkCreatedEvent represents a bulk voucher instance creation event
type VoucherInstanceBulkCreatedEvent struct {
	TenantID    string   `json:"tenantId" bson:"tenantId"`
	TemplateID  string   `json:"templateId" bson:"templateId"`
	BatchID     string   `json:"batchId" bson:"batchId"`
	Count       int      `json:"count" bson:"count"`
	VoucherCodes []string `json:"voucherCodes,omitempty" bson:"voucherCodes,omitempty"` // Optional: list of generated codes
}

// VoucherInstanceUsedEvent represents a voucher instance usage event
type VoucherInstanceUsedEvent struct {
	ID         string    `json:"id" bson:"_id"`
	TenantID   string    `json:"tenantId" bson:"tenantId"`
	TemplateID string    `json:"templateId" bson:"templateId"`
	VoucherCode string   `json:"voucherCode" bson:"voucherCode"`
	Status     string    `json:"status" bson:"status"`
	MACAddress string    `json:"macAddress" bson:"macAddress"`
	ClientIP   string    `json:"clientIp" bson:"clientIp"`
	CGNATIP    *string   `json:"cgnatIp,omitempty" bson:"cgnatIp,omitempty"`
	NASID      string    `json:"nasId" bson:"nasId"`
	LocationID *string   `json:"locationId,omitempty" bson:"locationId,omitempty"`
	SessionID  string    `json:"sessionId" bson:"sessionId"`
	UsedAt     time.Time `json:"usedAt" bson:"usedAt"`
	Version    int       `json:"version" bson:"version"`

	// Binding information (for cache update)
	BoundMAC    *string `json:"boundMac,omitempty" bson:"boundMac,omitempty"`
	BoundUserIP *string `json:"boundUserIp,omitempty" bson:"boundUserIp,omitempty"`
}

// VoucherInstanceExpiredEvent represents a voucher instance expiration event
type VoucherInstanceExpiredEvent struct {
	ID         string    `json:"id" bson:"_id"`
	TenantID   string    `json:"tenantId" bson:"tenantId"`
	TemplateID string    `json:"templateId" bson:"templateId"`
	VoucherCode string   `json:"voucherCode" bson:"voucherCode"`
	Status     string    `json:"status" bson:"status"`
	ExpiredAt  time.Time `json:"expiredAt" bson:"expiredAt"`
	Reason     string    `json:"reason,omitempty" bson:"reason,omitempty"` // "time_limit", "checkout", "auto_invalidate"
	Version    int       `json:"version" bson:"version"`
}

// VoucherInstanceRevokedEvent represents a voucher instance revocation event
type VoucherInstanceRevokedEvent struct {
	ID         string    `json:"id" bson:"_id"`
	TenantID   string    `json:"tenantId" bson:"tenantId"`
	TemplateID string    `json:"templateId" bson:"templateId"`
	VoucherCode string   `json:"voucherCode" bson:"voucherCode"`
	Status     string    `json:"status" bson:"status"`
	BatchID    *string   `json:"batchId,omitempty" bson:"batchId,omitempty"`
	RevokedAt  time.Time `json:"revokedAt" bson:"revokedAt"`
	RevokedBy  string    `json:"revokedBy,omitempty" bson:"revokedBy,omitempty"`
	Reason     string    `json:"reason,omitempty" bson:"reason,omitempty"`
	Version    int       `json:"version" bson:"version"`
}

// VoucherInstanceExtendedEvent represents a voucher instance validity extension event
type VoucherInstanceExtendedEvent struct {
	ID           string    `json:"id" bson:"_id"`
	TenantID     string    `json:"tenantId" bson:"tenantId"`
	TemplateID   string    `json:"templateId" bson:"templateId"`
	VoucherCode  string    `json:"voucherCode" bson:"voucherCode"`
	BatchID      *string   `json:"batchId,omitempty" bson:"batchId,omitempty"`
	OldExpiresAt time.Time `json:"oldExpiresAt" bson:"oldExpiresAt"`
	NewExpiresAt time.Time `json:"newExpiresAt" bson:"newExpiresAt"`
	ExtendedAt   time.Time `json:"extendedAt" bson:"extendedAt"`
	ExtendedBy   string    `json:"extendedBy,omitempty" bson:"extendedBy,omitempty"`
	Version      int       `json:"version" bson:"version"`
}

// VoucherSessionCreatedEvent represents a voucher session creation event
type VoucherSessionCreatedEvent struct {
	ID                string    `json:"id" bson:"_id"`
	TenantID          string    `json:"tenantId" bson:"tenantId"`
	VoucherInstanceID string    `json:"voucherInstanceId" bson:"voucherInstanceId"`
	VoucherCode       string    `json:"voucherCode" bson:"voucherCode"`
	SessionID         string    `json:"sessionId" bson:"sessionId"`
	MACAddress        string    `json:"macAddress" bson:"macAddress"`
	ClientIP          string    `json:"clientIp" bson:"clientIp"`
	CGNATIP           *string   `json:"cgnatIp,omitempty" bson:"cgnatIp,omitempty"`
	NASID             string    `json:"nasId" bson:"nasId"`
	LocationID        *string   `json:"locationId,omitempty" bson:"locationId,omitempty"`
	StartTime         time.Time `json:"startTime" bson:"startTime"`
	Status            string    `json:"status" bson:"status"`
}

// VoucherSessionEndedEvent represents a voucher session end event
type VoucherSessionEndedEvent struct {
	ID                string    `json:"id" bson:"_id"`
	TenantID          string    `json:"tenantId" bson:"tenantId"`
	VoucherInstanceID string    `json:"voucherInstanceId" bson:"voucherInstanceId"`
	VoucherCode       string    `json:"voucherCode" bson:"voucherCode"`
	SessionID         string    `json:"sessionId" bson:"sessionId"`
	StartTime         time.Time `json:"startTime" bson:"startTime"`
	EndTime           time.Time `json:"endTime" bson:"endTime"`
	Duration          int64     `json:"duration" bson:"duration"` // Duration in seconds
	UsedDataMB        int64     `json:"usedDataMB" bson:"usedDataMB"`
	InputOctets       int64     `json:"inputOctets" bson:"inputOctets"`
	OutputOctets      int64     `json:"outputOctets" bson:"outputOctets"`
	TerminationReason *string   `json:"terminationReason,omitempty" bson:"terminationReason,omitempty"`
	Status            string    `json:"status" bson:"status"`
}
