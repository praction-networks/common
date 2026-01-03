package subscriberevent

// Core enums related to Subscriber itself.

import "time"

type SubscriberStatus string

const (
	// SubscriberStatusActive indicates the subscriber is active and can have services
	// Service-level status (e.g., HotspotProfile, BroadbandSubscription) determines actual service availability
	SubscriberStatusActive SubscriberStatus = "ACTIVE"

	// SubscriberStatusInactive indicates the subscriber is disabled/archived (soft delete)
	// Inactive subscribers cannot have new services created, but existing services remain in their current state
	SubscriberStatusInactive SubscriberStatus = "INACTIVE"
)

type SubscriberType string

const (
	// End-customer segments
	SubscriberTypeResidential SubscriberType = "RESIDENTIAL" // home broadband, home Wi-Fi
	SubscriberTypeSMB         SubscriberType = "SMB"         // small offices, retail
	SubscriberTypeEnterprise  SubscriberType = "ENTERPRISE"  // corporates, campuses, hotels, malls

	// Special cases
	SubscriberTypeGuestHotspot SubscriberType = "GUEST_HOTSPOT" // hotspot-only guest users
	SubscriberTypeDemo        SubscriberType = "DEMO"        // lab, test, demo accounts
)

type AddressType string

const (
	AddressTypeInstallation AddressType = "INSTALLATION"
	AddressTypeBilling      AddressType = "BILLING"
	AddressTypeShipping     AddressType = "SHIPPING"
	AddressTypeHotspotSite  AddressType = "HOTSPOT_SITE"
)

type KYCStatus string

const (
	KYCStatusNotRequired KYCStatus = "NOT_REQUIRED"
	KYCStatusPending     KYCStatus = "PENDING"
	KYCStatusVerified    KYCStatus = "VERIFIED"
	KYCStatusRejected    KYCStatus = "REJECTED"
)

// Kept here so KYC is near its enums.
// If you prefer, you can move this struct into subscriber.go instead.
type SubscriberKYC struct {
	Status       KYCStatus  `bson:"status" json:"status"`
	DocumentType string     `bson:"documentType,omitempty" json:"documentType,omitempty"` // AADHAAR, PAN, PASSPORT, etc.
	DocumentID   string     `bson:"documentId,omitempty" json:"documentId,omitempty"`
	VerifiedAt   *time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
	VerifiedBy   string     `bson:"verifiedBy,omitempty" json:"verifiedBy,omitempty"`
	Notes        string     `bson:"notes,omitempty" json:"notes,omitempty"`
}
