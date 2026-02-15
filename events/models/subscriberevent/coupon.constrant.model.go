package subscriberevent

import "time"

// VoucherTemplateStatus represents the status of a voucher template
type VoucherTemplateStatus string

const (
	VoucherTemplateStatusCreated   VoucherTemplateStatus = "CREATED"
	VoucherTemplateStatusIssued    VoucherTemplateStatus = "ISSUED"
	VoucherTemplateStatusApproved  VoucherTemplateStatus = "APPROVED"
	VoucherTemplateStatusRejected  VoucherTemplateStatus = "REJECTED"
	VoucherTemplateStatusCancelled VoucherTemplateStatus = "CANCELLED"
	VoucherTemplateStatusUnused    VoucherTemplateStatus = "UNUSED"
	VoucherTemplateStatusExhausted VoucherTemplateStatus = "EXHAUSTED"
	VoucherTemplateStatusDisabled  VoucherTemplateStatus = "DISABLED"
	VoucherTemplateStatusExpired   VoucherTemplateStatus = "EXPIRED"
	VoucherTemplateStatusDeleted   VoucherTemplateStatus = "DELETED"
	VoucherTemplateStatusActive    VoucherTemplateStatus = "ACTIVE"
	VoucherTemplateStatusArchived  VoucherTemplateStatus = "ARCHIVED"
)

// VoucherStatus represents the lifecycle state of a voucher instance
type VoucherStatus string

const (
	VoucherStatusCreated   VoucherStatus = "CREATED"   // Generated, not yet issued
	VoucherStatusIssued    VoucherStatus = "ISSUED"    // Issued to user/printed
	VoucherStatusUnused    VoucherStatus = "UNUSED"    // Issued but not used
	VoucherStatusActive    VoucherStatus = "ACTIVE"    // Currently in use
	VoucherStatusExhausted VoucherStatus = "EXHAUSTED" // Limit reached (time/data)
	VoucherStatusExpired   VoucherStatus = "EXPIRED"   // Validity expired
	VoucherStatusRevoked   VoucherStatus = "REVOKED"   // Manually revoked
	VoucherStatusArchived  VoucherStatus = "ARCHIVED"  // Immutable, audit-only
)

// DistributionType represents how the voucher is distributed
type DistributionType string

const (
	DistributionTypeSingle     DistributionType = "SINGLE"      // Single code (manual handout)
	DistributionTypeBulk       DistributionType = "BULK"        // Bulk batch (hotel vouchers)
	DistributionTypePrePrinted DistributionType = "PRE_PRINTED" // Pre-printed (cafe cards)
	DistributionTypeDigital    DistributionType = "DIGITAL"     // Digital (SMS/WhatsApp)
)

// SessionStatus represents the status of a voucher session
type SessionStatus string

const (
	SessionStatusActive     SessionStatus = "ACTIVE"     // Currently active
	SessionStatusEnded      SessionStatus = "ENDED"      // Ended normally
	SessionStatusTerminated SessionStatus = "TERMINATED" // Terminated (limit reached, manual, etc.)
	SessionStatusExpired    SessionStatus = "EXPIRED"    // Expired (timeout)
)

// TerminationReason represents why a session was terminated
type TerminationReason string

const (
	TerminationReasonNormal         TerminationReason = "NORMAL"          // User disconnected
	TerminationReasonTimeLimit      TerminationReason = "TIME_LIMIT"      // Time limit reached
	TerminationReasonDataLimit      TerminationReason = "DATA_LIMIT"      // Data limit reached
	TerminationReasonManual         TerminationReason = "MANUAL"          // Admin terminated
	TerminationReasonNewSession     TerminationReason = "NEW_SESSION"     // New session started (single-use)
	TerminationReasonMACMismatch    TerminationReason = "MAC_MISMATCH"    // MAC address changed
	TerminationReasonLocationDenied TerminationReason = "LOCATION_DENIED" // Location not allowed
	TerminationReasonDeviceLimit    TerminationReason = "DEVICE_LIMIT"    // Device limit exceeded
	TerminationReasonExpired        TerminationReason = "EXPIRED"         // Voucher expired
	TerminationReasonRevoked        TerminationReason = "REVOKED"         // Voucher revoked
)



// VoucherValidity defines the validity window for template and generated vouchers
type VoucherValidity struct {
	// Template validity - when this template can be used to generate vouchers
	Start time.Time `bson:"start" json:"start"` // Template valid from
	End   time.Time `bson:"end" json:"end"`     // Template valid until

	// Default TTL for voucher instances generated from this template
	// Number of days from generation/issue date until voucher expires
	DefaultTTLDays int `bson:"defaultTTLDays" json:"defaultTTLDays"` // e.g., 1 (expires 1 day after issue)
}

// DistributionSettings defines how vouchers are generated and distributed
type DistributionSettings struct {
	// Code format pattern - defines how voucher codes are generated
	// Use X for random alphanumeric, # for random numeric
	// Example: "HOTEL-XXXX-XXXX" generates "HOTEL-A3B2-C4D1"
	CodeFormat string `bson:"codeFormat" json:"codeFormat"` // e.g., "HOTEL-XXXX-XXXX"

	// Code length - total length of generated code (excluding separators)
	CodeLength int `bson:"codeLength" json:"codeLength"` // e.g., 12

	// QR code generation - whether to generate QR codes for vouchers
	QRCodeEnabled bool `bson:"qrCodeEnabled" json:"qrCodeEnabled"` // Default: false

	// QR code URL template - template for QR code data URL
	// Use {code} placeholder for voucher code
	// Example: "https://login.example.com/voucher/{code}"
	QRCodeURLTemplate string `bson:"qrCodeUrlTemplate,omitempty" json:"qrCodeUrlTemplate,omitempty"`

	// Batch size - default batch size for bulk generation
	BatchSize int `bson:"batchSize" json:"batchSize"` // e.g., 1000

	// Prefix - optional prefix for all codes (e.g., "HOTEL", "CAFE")
	Prefix string `bson:"prefix,omitempty" json:"prefix,omitempty"`

	// Suffix - optional suffix for all codes
	Suffix string `bson:"suffix,omitempty" json:"suffix,omitempty"`
}
