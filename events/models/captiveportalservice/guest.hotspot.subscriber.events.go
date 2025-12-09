package captiveportalservice

import "time"

// GuestHotspotSubscriberCreatedEvent represents a guest hotspot subscriber creation event
// This event is published by captive-portal-service when a guest hotspot subscriber is created through the portal
type GuestHotspotSubscriberCreatedEvent struct {
	// Subscriber information (created by subscriber-service)
	SubscriberID     string `json:"subscriberId" bson:"subscriberId"`
	TenantID         string `json:"tenantId" bson:"tenantId"`
	HotspotProfileID string `json:"hotspotProfileId" bson:"hotspotProfileId"`

	// Subscriber details
	Type          string `json:"type" bson:"type"`                   // "GUEST_HOTSPOT"
	Status        string `json:"status" bson:"status"`               // "ACTIVE"
	PrimaryMobile string `json:"primaryMobile" bson:"primaryMobile"` // Required
	FullName      string `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PrimaryEmail  string `json:"primaryEmail,omitempty" bson:"primaryEmail,omitempty"`

	// RADIUS authentication credentials (for NAS/RADIUS authentication)
	Username string `json:"username" bson:"username"` // RADIUS username
	Password string `json:"password" bson:"password"` // RADIUS password

	// Portal session information (from captive portal session)
	SessionID string `json:"sessionId" bson:"sessionId"` // Portal session ID
	MAC       string `json:"mac" bson:"mac"`             // Device MAC address
	IPAddress string `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	NASId     string `json:"nasId" bson:"nasId"`                     // NAS device ID
	SSID      string `json:"ssid,omitempty" bson:"ssid,omitempty"`   // WiFi SSID
	APMac     string `json:"apMac,omitempty" bson:"apMac,omitempty"` // Access Point MAC

	// Device information (from portal session)
	DeviceInfo DeviceInfo `json:"deviceInfo,omitempty" bson:"deviceInfo,omitempty"`

	// Authentication method used (OTP, PASSWORD, VOUCHER, etc.)
	AuthMethod string `json:"authMethod" bson:"authMethod"` // e.g., "OTP"

	// User data submitted during authentication (form data)
	UserData map[string]interface{} `json:"userData,omitempty" bson:"userData,omitempty"`

	// Consent information (if provided during authentication)
	Consent ConsentInfo `json:"consent,omitempty" bson:"consent,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	CreatedBy string    `json:"createdBy,omitempty" bson:"createdBy,omitempty"` // "captive_portal"
}

// GuestHotspotSubscriberUpdatedEvent represents a guest hotspot subscriber update event
type GuestHotspotSubscriberUpdatedEvent struct {
	// Subscriber information
	SubscriberID string `json:"subscriberId" bson:"subscriberId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`

	// Updated user information
	FullName     string `json:"fullName,omitempty" bson:"fullName,omitempty"`
	PrimaryEmail string `json:"primaryEmail,omitempty" bson:"primaryEmail,omitempty"`

	// Updated user data
	UserData map[string]interface{} `json:"userData,omitempty" bson:"userData,omitempty"`

	// Timestamps
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy string    `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	Version   int       `json:"version" bson:"version"`
}

// GuestHotspotSubscriberDeletedEvent represents a guest hotspot subscriber deletion event
type GuestHotspotSubscriberDeletedEvent struct {
	// Subscriber information
	SubscriberID string `json:"subscriberId" bson:"subscriberId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`

	// Timestamps
	DeletedAt time.Time `json:"deletedAt" bson:"deletedAt"`
	DeletedBy string    `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
}

// DeviceInfo contains information about the user's device
type DeviceInfo struct {
	UserAgent      string `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	DeviceType     string `json:"deviceType,omitempty" bson:"deviceType,omitempty"` // mobile, tablet, desktop
	OS             string `json:"os,omitempty" bson:"os,omitempty"`                 // Android, iOS, Windows, macOS
	OSVersion      string `json:"osVersion,omitempty" bson:"osVersion,omitempty"`
	Browser        string `json:"browser,omitempty" bson:"browser,omitempty"` // Chrome, Safari, Firefox
	BrowserVersion string `json:"browserVersion,omitempty" bson:"browserVersion,omitempty"`
	Language       string `json:"language,omitempty" bson:"language,omitempty"`
	ScreenWidth    int    `json:"screenWidth,omitempty" bson:"screenWidth,omitempty"`
	ScreenHeight   int    `json:"screenHeight,omitempty" bson:"screenHeight,omitempty"`
	Timezone       string `json:"timezone,omitempty" bson:"timezone,omitempty"`
}

// ConsentInfo tracks user consent for terms, privacy, marketing
type ConsentInfo struct {
	TermsAccepted     bool      `json:"termsAccepted" bson:"termsAccepted"`
	TermsAcceptedAt   time.Time `json:"termsAcceptedAt,omitempty" bson:"termsAcceptedAt,omitempty"`
	PrivacyAccepted   bool      `json:"privacyAccepted" bson:"privacyAccepted"`
	PrivacyAcceptedAt time.Time `json:"privacyAcceptedAt,omitempty" bson:"privacyAcceptedAt,omitempty"`
	MarketingOptIn    bool      `json:"marketingOptIn" bson:"marketingOptIn"`
	IPAddress         string    `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"` // IP where consent was given (for audit)
}
