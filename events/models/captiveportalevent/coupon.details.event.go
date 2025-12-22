package captiveportalevent

// CouponDetailsEvent represents voucher validation details for captive portal
// Published after coupon code validation during authentication flow
type CouponDetailsEvent struct {
	// Coupon information
	CouponCode        string `json:"couponCode" bson:"couponCode"`
	VoucherInstanceID string `json:"voucherInstanceId,omitempty" bson:"voucherInstanceId,omitempty"`
	TemplateID        string `json:"templateId,omitempty" bson:"templateId,omitempty"`

	// Tenant and location
	TenantID   string  `json:"tenantId" bson:"tenantId"`
	NASID      string  `json:"nasId" bson:"nasId"`
	LocationID *string `json:"locationId,omitempty" bson:"locationId,omitempty"`

	// Validation result
	Valid   bool   `json:"valid" bson:"valid"`
	Status  string `json:"status,omitempty" bson:"status,omitempty"`   // UNUSED, ACTIVE, EXPIRED, etc.
	Message string `json:"message,omitempty" bson:"message,omitempty"` // Error message if invalid

	// Session information
	SessionID    string  `json:"sessionId" bson:"sessionId"`
	SubscriberID *string `json:"subscriberId,omitempty" bson:"subscriberId,omitempty"` // Created subscriber ID

	// Device information
	MACAddress string  `json:"macAddress" bson:"macAddress"`
	ClientIP   string  `json:"clientIp" bson:"clientIp"`
	CGNATIP    *string `json:"cgnatIp,omitempty" bson:"cgnatIp,omitempty"`

	// Voucher policy (limits and restrictions)
	Policy *VoucherPolicy `json:"policy,omitempty" bson:"policy,omitempty"`
}

// VoucherPolicy represents the policy applied to a voucher
type VoucherPolicy struct {
	// Limits
	TimeMinutes *int64      `json:"timeMinutes,omitempty" bson:"timeMinutes,omitempty"`
	DataMB      *int64      `json:"dataMB,omitempty" bson:"dataMB,omitempty"`
	SpeedMbps   *SpeedLimit `json:"speedMbps,omitempty" bson:"speedMbps,omitempty"`

	// Restrictions
	MaxDevices    *int     `json:"maxDevices,omitempty" bson:"maxDevices,omitempty"`
	BindMAC       bool     `json:"bindMac" bson:"bindMac"`
	BindLocation  bool     `json:"bindLocation" bson:"bindLocation"`
	AllowedNASIDs []string `json:"allowedNasIds,omitempty" bson:"allowedNasIds,omitempty"`
	AllowedSSIDs  []string `json:"allowedSsids,omitempty" bson:"allowedSsids,omitempty"`
}

// SpeedLimit defines download and upload speed limits
type SpeedLimit struct {
	Download int `json:"download" bson:"download"` // Mbps
	Upload   int `json:"upload" bson:"upload"`     // Mbps
}
