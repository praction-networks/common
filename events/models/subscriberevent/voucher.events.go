package subscriberevent

import "time"

// VoucherTemplateCreatedEvent represents a voucher template creation event
type VoucherTemplateCreatedEvent struct {
	ID          string `json:"id" bson:"id"`
	TenantID    string `json:"tenantId" bson:"tenantId"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	PlanID      string `json:"planId,omitempty" bson:"planId,omitempty"`
	Status      string `json:"status" bson:"status"`
	Version     int    `json:"version" bson:"version"`
}

// VoucherTemplateUpdatedEvent represents a voucher template update event
type VoucherTemplateUpdatedEvent struct {
	ID          string `json:"id" bson:"id"`
	TenantID    string `json:"tenantId" bson:"tenantId"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty"`
	Version     int    `json:"version" bson:"version"`
}

// VoucherTemplateDeletedEvent represents a voucher template deletion event
type VoucherTemplateDeletedEvent struct {
	ID        string    `json:"id" bson:"id"`
	TenantID  string    `json:"tenantId" bson:"tenantId"`
	DeletedAt time.Time `json:"deletedAt" bson:"deletedAt"`
	DeletedBy string    `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
}

// VoucherInstanceCreatedEvent represents a voucher instance creation event
// Includes limits and restrictions from template for local caching in captive-portal-service
type VoucherInstanceCreatedEvent struct {
	ID           string     `json:"id" bson:"id"`
	TenantID     string     `json:"tenantId" bson:"tenantId"`
	TemplateID   string     `json:"templateId" bson:"templateId"`
	CouponCode   string     `json:"couponCode" bson:"couponCode"`
	Status       string     `json:"status" bson:"status"`
	IssuedAt     *time.Time `json:"issuedAt,omitempty" bson:"issuedAt,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt" bson:"expiresAt"`
	BatchID      *string    `json:"batchId,omitempty" bson:"batchId,omitempty"`
	Distribution string     `json:"distributionType" bson:"distributionType"`
	Version      int        `json:"version" bson:"version"`

	// Template data for local caching (limits and restrictions)
	Limits       VoucherLimits       `json:"limits" bson:"limits"`
	Restrictions VoucherRestrictions `json:"restrictions" bson:"restrictions"`
}

// VoucherInstanceBulkCreatedEvent represents a bulk voucher instance creation event
type VoucherInstanceBulkCreatedEvent struct {
	TenantID    string   `json:"tenantId" bson:"tenantId"`
	TemplateID  string   `json:"templateId" bson:"templateId"`
	BatchID     string   `json:"batchId" bson:"batchId"`
	Count       int      `json:"count" bson:"count"`
	CouponCodes []string `json:"couponCodes,omitempty" bson:"couponCodes,omitempty"` // Optional: list of generated codes
}

// VoucherInstanceUsedEvent represents a voucher instance usage event
type VoucherInstanceUsedEvent struct {
	ID         string    `json:"id" bson:"id"`
	TenantID   string    `json:"tenantId" bson:"tenantId"`
	TemplateID string    `json:"templateId" bson:"templateId"`
	CouponCode string    `json:"couponCode" bson:"couponCode"`
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
	ID         string    `json:"id" bson:"id"`
	TenantID   string    `json:"tenantId" bson:"tenantId"`
	TemplateID string    `json:"templateId" bson:"templateId"`
	CouponCode string    `json:"couponCode" bson:"couponCode"`
	Status     string    `json:"status" bson:"status"`
	ExpiredAt  time.Time `json:"expiredAt" bson:"expiredAt"`
	Reason     string    `json:"reason,omitempty" bson:"reason,omitempty"` // "time_limit", "checkout", "auto_invalidate"
	Version    int       `json:"version" bson:"version"`
}

// VoucherInstanceRevokedEvent represents a voucher instance revocation event
type VoucherInstanceRevokedEvent struct {
	ID         string    `json:"id" bson:"id"`
	TenantID   string    `json:"tenantId" bson:"tenantId"`
	TemplateID string    `json:"templateId" bson:"templateId"`
	CouponCode string    `json:"couponCode" bson:"couponCode"`
	Status     string    `json:"status" bson:"status"`
	BatchID    *string   `json:"batchId,omitempty" bson:"batchId,omitempty"`
	RevokedAt  time.Time `json:"revokedAt" bson:"revokedAt"`
	RevokedBy  string    `json:"revokedBy,omitempty" bson:"revokedBy,omitempty"`
	Reason     string    `json:"reason,omitempty" bson:"reason,omitempty"`
	Version    int       `json:"version" bson:"version"`
}

// VoucherInstanceExtendedEvent represents a voucher instance validity extension event
type VoucherInstanceExtendedEvent struct {
	ID           string    `json:"id" bson:"id"`
	TenantID     string    `json:"tenantId" bson:"tenantId"`
	TemplateID   string    `json:"templateId" bson:"templateId"`
	CouponCode   string    `json:"couponCode" bson:"couponCode"`
	BatchID      *string   `json:"batchId,omitempty" bson:"batchId,omitempty"`
	OldExpiresAt time.Time `json:"oldExpiresAt" bson:"oldExpiresAt"`
	NewExpiresAt time.Time `json:"newExpiresAt" bson:"newExpiresAt"`
	ExtendedAt   time.Time `json:"extendedAt" bson:"extendedAt"`
	ExtendedBy   string    `json:"extendedBy,omitempty" bson:"extendedBy,omitempty"`
	Version      int       `json:"version" bson:"version"`
}

// VoucherSessionCreatedEvent represents a voucher session creation event
type VoucherSessionCreatedEvent struct {
	ID                string    `json:"id" bson:"id"`
	TenantID          string    `json:"tenantId" bson:"tenantId"`
	VoucherInstanceID string    `json:"voucherInstanceId" bson:"voucherInstanceId"`
	CouponCode        string    `json:"couponCode" bson:"couponCode"`
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
	ID                string    `json:"id" bson:"id"`
	TenantID          string    `json:"tenantId" bson:"tenantId"`
	VoucherInstanceID string    `json:"voucherInstanceId" bson:"voucherInstanceId"`
	CouponCode        string    `json:"couponCode" bson:"couponCode"`
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
