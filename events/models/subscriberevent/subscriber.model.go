package subscriberevent

// Core enums related to Subscriber itself.

import "time"



type SubscriberStatus string

const (
	// Funnel / early stages
	SubscriberStatusTrial SubscriberStatus = "TRIAL" // free/demo hotspot, trial broadband, etc.

	// Normal lifecycle
	SubscriberStatusPendingActivation SubscriberStatus = "PENDING_ACTIVATION" // KYC / install / payment pending
	SubscriberStatusActive            SubscriberStatus = "ACTIVE"             // fully active subscriber

	// Temporarily blocked
	SubscriberStatusGrace     SubscriberStatus = "GRACE"     // grace period after due date
	SubscriberStatusSuspended SubscriberStatus = "SUSPENDED" // non-payment, abuse, admin hold

	// End of relationship
	SubscriberStatusTerminated SubscriberStatus = "TERMINATED" // normal closure, not necessarily churn
	SubscriberStatusChurned    SubscriberStatus = "CHURNED"    // confirmed lost to competitor / no comeback expected

	// Risk / fraud
	SubscriberStatusBlacklisted SubscriberStatus = "BLACKLISTED" // fraud / abuse, do not re-activate
)

type SubscriberType string

const (
	// End-customer segments
	SubscriberTypeResidential SubscriberType = "RESIDENTIAL" // home broadband, home Wi-Fi
	SubscriberTypeSMB         SubscriberType = "SMB"         // small offices, retail
	SubscriberTypeEnterprise  SubscriberType = "ENTERPRISE"  // corporates, campuses, hotels, malls

	// Special cases
	SubscriberTypeGuestHotspot SubscriberType = "GUEST_HOTSPOT" // hotspot-only guest users
	SubscriberTypePartner      SubscriberType = "PARTNER"       // LCO, reseller, channel partner
	SubscriberTypeWholesale    SubscriberType = "WHOLESALE"     // carrier, NNI, interconnect
	SubscriberTypeInternal     SubscriberType = "INTERNAL"      // lab, test, demo accounts
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
