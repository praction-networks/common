package subscriberevent

import "time"

// VoucherBatchStatus represents the lifecycle state of a voucher batch
type VoucherBatchStatus string

const (
	VoucherBatchStatusCreated  VoucherBatchStatus = "CREATED"  // Batch generated, pending review
	VoucherBatchStatusApproved VoucherBatchStatus = "APPROVED" // Ready for distribution
	VoucherBatchStatusRejected VoucherBatchStatus = "REJECTED" // Batch rejected
	VoucherBatchStatusExpired  VoucherBatchStatus = "EXPIRED"  // All vouchers expired
	VoucherBatchStatusArchived VoucherBatchStatus = "ARCHIVED" // Immutable, audit-only
)

// VoucherStatus represents the lifecycle state of an individual voucher
type VoucherStatus string

const (
	VoucherStatusCreated    VoucherStatus = "CREATED"    // Generated, not yet used
	VoucherStatusUsed       VoucherStatus = "USED"       // Voucher has been redeemed
	VoucherStatusLive       VoucherStatus = "LIVE"       // Currently active session
	VoucherStatusExpired    VoucherStatus = "EXPIRED"    // Validity period ended
	VoucherStatusTerminated VoucherStatus = "TERMINATED" // Manually revoked/disabled
	VoucherStatusExhausted  VoucherStatus = "EXHAUSTED"  // Usage/data/time limit reached
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

// VoucherValidity defines the validity window for voucher batches
type VoucherValidity struct {
	Start time.Time `bson:"start" json:"start"` // Valid from
	End   time.Time `bson:"end" json:"end"`     // Valid until
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

	// Logo URL for QR code center overlay — fetched at generation time
	LogoURL string `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`

	// WiFi SSID - when set, QR codes encode a WiFi auto-connect string (WIFI:T:...;S:<SSID>;;)
	// This allows guests to scan and auto-connect to the hotspot WiFi network
	WifiSSID string `bson:"wifiSSID,omitempty" json:"wifiSSID,omitempty"`

	// WiFi password - optional WPA password for the WiFi network
	// If empty, auth type is "nopass" (open network with captive portal)
	WifiPassword string `bson:"wifiPassword,omitempty" json:"wifiPassword,omitempty"`

	// Batch size - default batch size for bulk generation
	BatchSize int `bson:"batchSize" json:"batchSize"` // e.g., 1000

	// Separator - character between prefix, code, and suffix (e.g., "-", "_", ".")
	Separator string `bson:"separator,omitempty" json:"separator,omitempty"`

	// Prefix - optional prefix for all codes (e.g., "HOTEL", "CAFE")
	Prefix string `bson:"prefix,omitempty" json:"prefix,omitempty"`

	// Suffix - optional suffix for all codes
	Suffix string `bson:"suffix,omitempty" json:"suffix,omitempty"`
}
