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

// VoucherType represents the type of access limit for a voucher
type VoucherType string

const (
	VoucherTypeTime  VoucherType = "TIME"  // Time-based (e.g., 2 hours)
	VoucherTypeData  VoucherType = "DATA"  // Data-based (e.g., 1 GB)
	VoucherTypeSpeed VoucherType = "SPEED" // Speed-based (e.g., 10 Mbps)
	VoucherTypeCombo VoucherType = "COMBO" // Combo (time OR data, whichever comes first)
)

// VoucherScope represents the scope/restriction of voucher usage
type VoucherScope string

const (
	VoucherScopeSingleDevice  VoucherScope = "SINGLE_DEVICE"  // One MAC only
	VoucherScopeMultiDevice   VoucherScope = "MULTI_DEVICE"   // Up to N devices
	VoucherScopeLocationBound VoucherScope = "LOCATION_BOUND" // Only one AP/NAS
	VoucherScopeGlobal        VoucherScope = "GLOBAL"         // Any AP under tenant
)

// VoucherLimits defines the access limits for vouchers generated from this template
type VoucherLimits struct {
	// Time limit in minutes (optional - nil means no time limit)
	// If set, session will be terminated after this duration
	TimeMinutes *int64 `bson:"timeMinutes,omitempty" json:"timeMinutes,omitempty"` // e.g., 1440 (24 hours)

	// Data limit in megabytes (optional - nil means no data limit)
	// If set, session will be terminated after this amount of data is consumed
	DataMB *int64 `bson:"dataMB,omitempty" json:"dataMB,omitempty"` // e.g., 1024 (1 GB)

	// Speed limit in Mbps (optional - nil means no speed limit)
	// If set, bandwidth will be throttled to these limits
	SpeedMbps *SpeedLimit `bson:"speedMbps,omitempty" json:"speedMbps,omitempty"` // e.g., { download: 10, upload: 5 }

	// Combo limit - if both time and data are set, whichever is reached first terminates the session
	// This is implicit - no separate field needed
}

// SpeedLimit defines download and upload speed limits in Mbps
type SpeedLimit struct {
	Download int `bson:"download" json:"download"` // Download speed in Mbps
	Upload   int `bson:"upload" json:"upload"`     // Upload speed in Mbps
}

// VoucherRestrictions defines usage constraints for vouchers
type VoucherRestrictions struct {
	// Maximum number of devices that can use this voucher simultaneously
	// nil = unlimited devices
	// 1 = single device only
	MaxDevices *int `bson:"maxDevices,omitempty" json:"maxDevices,omitempty"` // e.g., 1, 3, 5

	// Bind to first MAC address that uses the voucher
	// If true, only the first MAC address can use the voucher
	// Prevents MAC spoofing and sharing
	BindMAC bool `bson:"bindMac" json:"bindMac"` // Default: false

	// Location binding - restrict voucher to specific NAS/AP locations
	// If true, voucher can only be used at allowed NAS locations
	BindLocation bool `bson:"bindLocation" json:"bindLocation"` // Default: false

	// Allowed NAS IDs - if BindLocation is true, only these NAS IDs can accept this voucher
	// Empty array with BindLocation=true means no locations allowed (invalid voucher)
	AllowedNASIDs []string `bson:"allowedNasIds,omitempty" json:"allowedNasIds,omitempty"` // e.g., ["ap_12", "ap_13"]

	// Allowed SSIDs - restrict voucher to specific SSIDs
	// Empty array means all SSIDs allowed
	AllowedSSIDs []string `bson:"allowedSsids,omitempty" json:"allowedSsids,omitempty"` // e.g., ["HotelWiFi_Guest"]
}

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
