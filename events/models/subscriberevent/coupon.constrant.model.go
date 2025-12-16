package subscriberevent

// VoucherTemplateStatus represents the status of a voucher template
type VoucherTemplateStatus string

const (
	VoucherTemplateStatusCreated    VoucherTemplateStatus = "CREATED"
	VoucherTemplateStatusIssued    VoucherTemplateStatus = "ISSUED"
	VoucherTemplateStatusApproved    VoucherTemplateStatus = "APPROVED"
	VoucherTemplateStatusRejected    VoucherTemplateStatus = "REJECTED"
	VoucherTemplateStatusCancelled    VoucherTemplateStatus = "CANCELLED"
	VoucherTemplateStatusUnused    VoucherTemplateStatus = "UNUSED"
	VoucherTemplateStatusExhausted    VoucherTemplateStatus = "EXHAUSTED"
	VoucherTemplateStatusDisabled    VoucherTemplateStatus = "DISABLED"
	VoucherTemplateStatusExpired    VoucherTemplateStatus = "EXPIRED"
	VoucherTemplateStatusDeleted    VoucherTemplateStatus = "DELETED"
	VoucherTemplateStatusActive   VoucherTemplateStatus = "ACTIVE"
	VoucherTemplateStatusArchived VoucherTemplateStatus = "ARCHIVED"
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
	VoucherStatusRevoked   VoucherStatus = "REVOKED"  // Manually revoked
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
	TerminationReasonManual         TerminationReason = "MANUAL"         // Admin terminated
	TerminationReasonNewSession     TerminationReason = "NEW_SESSION"     // New session started (single-use)
	TerminationReasonMACMismatch    TerminationReason = "MAC_MISMATCH"   // MAC address changed
	TerminationReasonLocationDenied TerminationReason = "LOCATION_DENIED" // Location not allowed
	TerminationReasonDeviceLimit    TerminationReason = "DEVICE_LIMIT"   // Device limit exceeded
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