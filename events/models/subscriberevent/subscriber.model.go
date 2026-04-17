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

	// SubscriberStatusOnboarding indicates the subscriber is recently created and may still be
	// undergoing KYC or plan assignment before an active connection is established.
	SubscriberStatusOnboarding SubscriberStatus = "ONBOARDING"
)

type SubscriberType string

const (
	// End-customer segments
	SubscriberTypeResidential SubscriberType = "RESIDENTIAL" // home broadband, home Wi-Fi
	SubscriberTypeSMB         SubscriberType = "SMB"         // small offices, retail
	SubscriberTypeEnterprise  SubscriberType = "ENTERPRISE"  // corporates, campuses, hotels, malls

	// Special cases
	SubscriberTypeGuestHotspot SubscriberType = "GUEST_HOTSPOT" // hotspot-only guest users
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

// ── KYC Document Type ───────────────────────────────────────────────────────
// Unified list — same document types can be used for BOTH Person KYC and Address KYC.
// e.g. Aadhaar and Voter ID contain the person's address, so they serve as address proof too.

type KYCDocumentType string

const (
	// Identity documents (primarily Person KYC, but Aadhaar/VoterID/DL also valid for Address KYC)
	KYCDocTypePAN      KYCDocumentType = "PAN"
	KYCDocTypeAadhaar  KYCDocumentType = "AADHAAR" // Valid for both Person & Address KYC
	KYCDocTypePassport KYCDocumentType = "PASSPORT"
	KYCDocTypeVoterID  KYCDocumentType = "VOTER_ID"        // Valid for both Person & Address KYC
	KYCDocTypeDL       KYCDocumentType = "DRIVING_LICENSE" // Also contains address

	// Address proof documents (primarily Address KYC)
	KYCDocTypeUtilityBill   KYCDocumentType = "UTILITY_BILL"
	KYCDocTypeBankStatement KYCDocumentType = "BANK_STATEMENT"
	KYCDocTypeRentAgreement KYCDocumentType = "RENT_AGREEMENT"
	KYCDocTypePropertyTax   KYCDocumentType = "PROPERTY_TAX"
	KYCDocTypeGasBill       KYCDocumentType = "GAS_BILL"
)

// ── KYC Document ────────────────────────────────────────────────────────────

// KYCDocument represents a single KYC document (used in both Person and Address KYC).
type KYCDocument struct {
	DocumentType          KYCDocumentType `bson:"documentType" json:"documentType"`
	DocumentID          string          `bson:"documentId" json:"documentId"`
	VerificationSessionID string        `bson:"verificationSessionId,omitempty" json:"verificationSessionId,omitempty"` // For validation during subscriber creation
	Status              KYCStatus       `bson:"status" json:"status"`
	VerifiedAt          *time.Time      `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
	VerifiedBy          string          `bson:"verifiedBy,omitempty" json:"verifiedBy,omitempty"`
	Notes               string          `bson:"notes,omitempty" json:"notes,omitempty"`
}

// ── Person KYC (Subscriber Level) ───────────────────────────────────────────
// SubscriberKYC is the Person KYC — identity verification of the subscriber.
// Lives on the Subscriber model. Applies to all connections of this subscriber.
// Hotspot profiles also use this (no per-connection KYC for hotspot).
type SubscriberKYC struct {
	Status    KYCStatus     `bson:"status" json:"status"`
	Documents []KYCDocument `bson:"documents,omitempty" json:"documents,omitempty"`
	Notes     string        `bson:"notes,omitempty" json:"notes,omitempty"`
	// Deprecated flat fields — kept for backward compatibility with existing MongoDB data.
	// New code should use Documents[]. A migration can backfill old data.
	DocumentType string     `bson:"documentType,omitempty" json:"documentType,omitempty"`
	DocumentID   string     `bson:"documentId,omitempty" json:"documentId,omitempty"`
	VerifiedAt   *time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
	VerifiedBy   string     `bson:"verifiedBy,omitempty" json:"verifiedBy,omitempty"`
}

// ── Address KYC (Broadband Connection Level) ────────────────────────────────
// AddressKYC is the per-connection address verification.
// Proves the subscriber lives/operates at the installation address.
// Lives on BroadbandSubscription model only (not hotspot).
type AddressKYC struct {
	Status    KYCStatus     `bson:"status" json:"status"`
	Documents []KYCDocument `bson:"documents,omitempty" json:"documents,omitempty"`
	Notes     string        `bson:"notes,omitempty" json:"notes,omitempty"`
}
