package tenantevent

// ── KYC Policy Configuration ─────────────────────────────────────────────────
// Controls how KYC is enforced per subscriber type for both Person KYC and Address KYC.
// The master toggle is Core.IsUserKYCEnabled — if false, all KYC is skipped entirely.
// When enabled, these policies control granular enforcement.
//
// Architecture:
//   - Master toggle: EnabledFeatures.Core.IsUserKYCEnabled (boolean, managed via feature toggles)
//   - Policy config: TenantModel.KycPolicy (this struct, managed via dedicated KYC policy API)
//   - Gateway config: TenantKYCProviderBinding (separate collection, managed via integrations)
//
// Enforcement modes:
//
//	REQUIRED     = KYC is mandatory. New subscribers/connections cannot be activated without it.
//	OPTIONAL     = KYC form is shown but can be skipped. Useful during initial rollout or bulk onboarding.
//	NOT_REQUIRED = KYC section is hidden entirely. No KYC fields shown.
//
// Bulk onboarding flow:
//  1. Tenant sets KYC to OPTIONAL → imports subscribers in bulk without KYC
//  2. Later, tenant switches to REQUIRED + sets EnforceOnExisting=true
//  3. All existing subscribers without KYC get status set to PENDING
//  4. Subscribers must complete KYC within GracePeriodDays or face suspension

// KYCEnforcement controls KYC behavior per subscriber type
type KYCEnforcement string

const (
	KYCEnforcementRequired    KYCEnforcement = "REQUIRED"     // FORCED — must complete KYC to activate
	KYCEnforcementOptional    KYCEnforcement = "OPTIONAL"     // Shown in form but can be skipped
	KYCEnforcementNotRequired KYCEnforcement = "NOT_REQUIRED" // NONE — hidden entirely
)

// PersonKYCPolicy defines Person KYC rules per subscriber type
// Person KYC = identity verification (PAN, Aadhaar, Passport, etc.)
// Lives on the Subscriber model, applies to all connections.
type PersonKYCPolicy struct {
	SubscriberType string         `json:"subscriberType" bson:"subscriberType"` // RESIDENTIAL, SMB, ENTERPRISE, GUEST_HOTSPOT
	Enforcement    KYCEnforcement `json:"enforcement" bson:"enforcement"`
	MinDocuments   int            `json:"minDocuments" bson:"minDocuments"` // Minimum identity documents required
	AutoVerify     bool           `json:"autoVerify" bson:"autoVerify"`     // Use KYC provider for auto-verification
}

// AddressKYCPolicy defines Address KYC rules for broadband connections per subscriber type
// Address KYC = proof of installation address (utility bill, Aadhaar, rent agreement, etc.)
// Lives on BroadbandSubscription model only (not hotspot).
type AddressKYCPolicy struct {
	SubscriberType string         `json:"subscriberType" bson:"subscriberType"` // RESIDENTIAL, SMB, ENTERPRISE
	Enforcement    KYCEnforcement `json:"enforcement" bson:"enforcement"`
	MinDocuments   int            `json:"minDocuments" bson:"minDocuments"` // Minimum address proof documents required
	AutoVerify     bool           `json:"autoVerify" bson:"autoVerify"`     // Use KYC provider for auto-verification
}

// KYCFeatures — tenant-level KYC policy configuration
// Master toggle: Core.IsUserKYCEnabled must be true for any of this to apply.
// Lives on TenantModel.KycPolicy (NOT inside EnabledFeatures).
type KYCFeatures struct {
	// Per-subscriber-type enforcement policies
	PersonKYCPolicies  []PersonKYCPolicy  `json:"personKycPolicies" bson:"personKycPolicies"`   // Person KYC rules per subscriber type
	AddressKYCPolicies []AddressKYCPolicy `json:"addressKycPolicies" bson:"addressKycPolicies"` // Address KYC rules per subscriber type (broadband only)

	// Accepted document types (tenant can restrict which doc types they accept)
	AcceptedPersonDocTypes []string `json:"acceptedPersonDocTypes" bson:"acceptedPersonDocTypes"` // e.g. ["PAN","AADHAAR","PASSPORT","VOTER_ID","DRIVING_LICENSE"]
	AcceptedAddrDocTypes   []string `json:"acceptedAddrDocTypes" bson:"acceptedAddrDocTypes"`     // e.g. ["AADHAAR","VOTER_ID","UTILITY_BILL","BANK_STATEMENT",...]

	// Grace period — days to complete KYC after creation (0 = must provide immediately)
	// Used for both new subscribers and when EnforceOnExisting is flipped to true.
	GracePeriodDays int `json:"gracePeriodDays" bson:"gracePeriodDays"`

	// Enforce on existing subscribers — when tenant switches enforcement from OPTIONAL to REQUIRED,
	// set this to true to retroactively mark all existing subscribers without KYC as PENDING.
	// They get GracePeriodDays to complete KYC. After that, services can be suspended.
	// This is essential for the bulk onboarding → later enforcement flow.
	EnforceOnExisting bool `json:"enforceOnExisting" bson:"enforceOnExisting"`

	// Block RADIUS authentication after grace period expires for subscribers without completed KYC.
	// When true, subscribers who haven't completed KYC within GracePeriodDays will be denied RADIUS auth.
	// Only effective when EnforceOnExisting is true and GracePeriodDays > 0.
	BlockRadiusOnKYCFailure bool `json:"blockRadiusOnKYCFailure" bson:"blockRadiusOnKYCFailure"`
}
