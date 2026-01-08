package captiveportalservice

import (
	"time"

	"github.com/praction-networks/common/events/models/subscriberevent"
)

// GuestHotspotSubscriberCreatedEvent represents a guest hotspot subscriber creation event
// This event is published by captive-portal-service when a guest hotspot subscriber is created through the portal
type GuestHotspotSubscriberCreatedEvent struct {
	// Subscriber information (mapped from form/signup fields)
	SubscriberID string                            `json:"subscriberId,omitempty" bson:"subscriberId,omitempty"`
	TenantID     string                            `json:"tenantId" bson:"tenantId"`
	Type         string                            `json:"type" bson:"type"`           // "GUEST_HOTSPOT" - mapped from form
	Status       string                            `json:"status" bson:"status"`       // "ACTIVE" - mapped from form
	Username     string                            `json:"username" bson:"username"`   // RADIUS username
	Password     string                            `json:"password" bson:"password"`   // RADIUS password
	SessionID    string                            `json:"sessionId" bson:"sessionId"` // Portal session ID
	MAC          string                            `json:"mac" bson:"mac"`             // Device MAC address
	IPAddress    string                            `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	NASId        string                            `json:"nasId" bson:"nasId"`                     // NAS device ID
	SSID         string                            `json:"ssid,omitempty" bson:"ssid,omitempty"`   // WiFi SSID
	APMac        string                            `json:"apMac,omitempty" bson:"apMac,omitempty"` // Access Point MAC
	DeviceInfo   DeviceInfo                        `json:"deviceInfo,omitempty" bson:"deviceInfo,omitempty"`
	AuthMethod   subscriberevent.HotspotAuthMethod `json:"authMethod" bson:"authMethod"` // e.g., "OTP"
	UserData     map[string]interface{}            `json:"userData,omitempty" bson:"userData,omitempty"`
	Consent      ConsentInfo                       `json:"consent,omitempty" bson:"consent,omitempty"`
	CreatedAt    time.Time                         `json:"createdAt" bson:"createdAt"`
	CreatedBy    string                            `json:"createdBy,omitempty" bson:"createdBy,omitempty"` // "captive_portal"
}

// GuestHotspotSubscriberUpdatedEvent represents a guest hotspot subscriber update event
type GuestHotspotSubscriberUpdatedEvent struct {
	// Subscriber information
	SubscriberID     string                            `json:"subscriberId" bson:"subscriberId"`
	TenantID         string                            `json:"tenantId" bson:"tenantId"`
	HotspotProfileID string                            `json:"hotspotProfileId" bson:"hotspotProfileId"`
	MAC              string                            `json:"mac" bson:"mac"` // Device MAC address
	IPAddress        string                            `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	NASId            string                            `json:"nasId" bson:"nasId"`                     // NAS device ID
	SSID             string                            `json:"ssid,omitempty" bson:"ssid,omitempty"`   // WiFi SSID
	APMac            string                            `json:"apMac,omitempty" bson:"apMac,omitempty"` // Access Point MAC
	DeviceInfo       DeviceInfo                        `json:"deviceInfo,omitempty" bson:"deviceInfo,omitempty"`
	AuthMethod       subscriberevent.HotspotAuthMethod `json:"authMethod" bson:"authMethod"` // e.g., "OTP"
	UserData         map[string]interface{}            `json:"userData,omitempty" bson:"userData,omitempty"`
	Consent          ConsentInfo                       `json:"consent,omitempty" bson:"consent,omitempty"`
	UpdatedBy        string                            `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	Version          int                               `json:"version" bson:"version"`
}

type GuestHotspotSubscriberValidityExtendedEvent struct {
	SubscriberID     string    `json:"subscriberId" bson:"subscriberId"`
	TenantID         string    `json:"tenantId" bson:"tenantId"`
	HotspotProfileID string    `json:"hotspotProfileId" bson:"hotspotProfileId"`
	MAC              string    `json:"mac" bson:"mac"` // Device MAC address
	ValidUntil       time.Time `json:"validUntil,omitempty" bson:"validUntil,omitempty"`
	Version          int       `json:"version" bson:"version"`
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

// HotspotDeviceAddedEvent represents a device (MAC address) addition event to an existing hotspot profile
type GuestHotspotDeviceAddedEvent struct {
	HotspotProfileID string     `json:"hotspotProfileId" bson:"hotspotProfileId"`         // The hotspot profile ID from the event
	SubscriberID     string     `json:"subscriberId" bson:"subscriberId"`                 // The subscriber ID
	TenantID         string     `json:"tenantId" bson:"tenantId"`                         // The tenant ID
	MacAddresses     []string   `json:"macAddresses" bson:"macAddresses"`                 // New MAC addresses to add
	UserProfileID    string     `json:"userProfileId" bson:"userProfileId"`               // The existing UserProfile ID in RADIUS
	Username         string     `json:"username,omitempty" bson:"username,omitempty"`     // RADIUS username
	Password         string     `json:"password,omitempty" bson:"password,omitempty"`     // RADIUS password
	PlanCode         string     `json:"planCode,omitempty" bson:"planCode,omitempty"`     // Plan code
	ValidUntil       *time.Time `json:"validUntil,omitempty" bson:"validUntil,omitempty"` // Valid until timestamp
}

// HotspotDeviceRemovedEvent represents a device (MAC address) removal event from a hotspot profile
type GuestHotspotDeviceRemovedEvent struct {
	HotspotProfileID string   `json:"hotspotProfileId" bson:"hotspotProfileId"` // The hotspot profile ID
	SubscriberID     string   `json:"subscriberId" bson:"subscriberId"`         // The subscriber ID
	TenantID         string   `json:"tenantId" bson:"tenantId"`                 // The tenant ID
	MacAddresses     []string `json:"macAddresses" bson:"macAddresses"`         // MAC addresses to remove
	UserProfileID    string   `json:"userProfileId" bson:"userProfileId"`       // The UserProfile ID in RADIUS
}
