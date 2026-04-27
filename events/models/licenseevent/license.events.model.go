// Package licenseevent holds the canonical payload models for license-service
// outbound events. license-service is the sole publisher; tenant-service,
// billing-service, audit-log-service, and notification-service are consumers.
//
// Subjects defined in events.LicenseStream (subject.and.stream.go).
package licenseevent

import "time"

// LicenseClaims — public projection of the JWS payload. Consumers
// (jwsverify, auth-service middleware, cloud-sync-agent) read this
// shape after a verified parse. Internal license-service code uses
// services.LicenseClaims which is the same shape but not bound to
// the events module's import boundary.
type LicenseClaims struct {
	Iss               string                `json:"iss"`
	Sub               string                `json:"sub"`
	Aud               []string              `json:"aud"`
	Iat               int64                 `json:"iat"`
	Nbf               int64                 `json:"nbf"`
	Exp               int64                 `json:"exp"`
	JTI               string                `json:"jti"`
	TenantID          string                `json:"tenantId"`
	TenantType        string                `json:"tenantType"`
	InstallationID    string                `json:"installationId,omitempty"`
	MTLSFingerprint   string                `json:"mtlsFingerprint,omitempty"`
	Counter           int64                 `json:"counter"`
	Nonce             string                `json:"nonce"`
	LicenseStatus     string                `json:"licenseStatus"`
	LicenseValidUntil int64                 `json:"licenseValidUntil"`
	Entitlements      []EntitlementSnapshot `json:"entitlements"`
}

// LicenseSnapshot is the wire-shape of a License row, denormalized so
// consumers don't have to call back. Trimmed to fields consumers actually
// need; full record stays in license-service Postgres.
type LicenseSnapshot struct {
	ID                    string    `json:"id"`
	LicenseNumber         string    `json:"licenseNumber"`
	TenantID              string    `json:"tenantId"`
	TenantType            string    `json:"tenantType"`
	ParentLicenseID       *string   `json:"parentLicenseId,omitempty"`
	Status                string    `json:"status"`
	IssuedAt              time.Time `json:"issuedAt"`
	ValidFrom             time.Time `json:"validFrom"`
	ValidUntil            time.Time `json:"validUntil"`
	BillingCycleQuantity  int       `json:"billingCycleQuantity"`
	BillingCycleUnit      string    `json:"billingCycleUnit"`
	AutoRenew             bool      `json:"autoRenew"`
	MaxInstallations      int       `json:"maxInstallations"`
	BillingSubscriptionID string    `json:"billingSubscriptionId,omitempty"`
	PackageCode           string    `json:"packageCode,omitempty"`
	Version               int       `json:"version"`
}

// EntitlementSnapshot — one active SKU entry. Consumers (tenant-service)
// project EnabledFeatures from the union of UnlocksFeatures across the
// active set.
type EntitlementSnapshot struct {
	ID              string                  `json:"id"`
	SKUCode         string                  `json:"skuCode"`
	Family          string                  `json:"family"`
	EnforcementMode string                  `json:"enforcementMode"`
	Quantity        int64                   `json:"quantity"`
	Quotas          map[string]int64        `json:"quotas"`
	OveragePolicy   string                  `json:"overagePolicy"`
	Components      []EntitlementComponent  `json:"components"`
	UnlocksFeatures []string                `json:"unlocksFeatures"`
	EffectiveFrom   time.Time               `json:"effectiveFrom"`
	EffectiveTo     *time.Time              `json:"effectiveTo,omitempty"`
}

type EntitlementComponent struct {
	Component  string `json:"component"`
	Deployment string `json:"deployment"`
}

// IssuedEvent — license.issued. Carries the full initial entitlement set
// so consumers can materialise their projections in one event without
// having to subscribe to license.entitlement.changed separately.
type IssuedEvent struct {
	License      LicenseSnapshot       `json:"license"`
	Entitlements []EntitlementSnapshot `json:"entitlements"`
	IssuedBy     string                `json:"issuedBy"`
	IssuedAt     time.Time             `json:"issuedAt"`
}

// EntitlementChangedEvent — license.entitlement.changed. Canonical event
// the tenant-service projection consumer listens to. Fired on issue,
// add, remove, suspend, expiry — anywhere the active entitlement set
// for a license changes. Carries the FULL active set, not a diff —
// consumers reconcile by overwriting their projection.
type EntitlementChangedEvent struct {
	LicenseID    string                `json:"licenseId"`
	TenantID     string                `json:"tenantId"`
	Status       string                `json:"status"` // current license status
	Entitlements []EntitlementSnapshot `json:"entitlements"`
	Reason       string                `json:"reason,omitempty"`
	OccurredAt   time.Time             `json:"occurredAt"`
	Version      int                   `json:"version"`
}

// SuspendedEvent / ReinstatedEvent / RenewedEvent / ExpiredEvent / TerminatedEvent
// — emitted in addition to EntitlementChangedEvent when status flips, so sales
// + finance + ops alert pipelines can subscribe without filtering the noisier
// entitlement-changed stream.
type SuspendedEvent struct {
	LicenseID  string    `json:"licenseId"`
	TenantID   string    `json:"tenantId"`
	Reason     string    `json:"reason"`
	Actor      string    `json:"actor"`
	OccurredAt time.Time `json:"occurredAt"`
}

type ReinstatedEvent struct {
	LicenseID  string    `json:"licenseId"`
	TenantID   string    `json:"tenantId"`
	Actor      string    `json:"actor"`
	OccurredAt time.Time `json:"occurredAt"`
}

type RenewedEvent struct {
	LicenseID    string    `json:"licenseId"`
	TenantID     string    `json:"tenantId"`
	NewValidFrom time.Time `json:"newValidFrom"`
	NewValidUntil time.Time `json:"newValidUntil"`
	OccurredAt   time.Time `json:"occurredAt"`
}

type ExpiredEvent struct {
	LicenseID  string    `json:"licenseId"`
	TenantID   string    `json:"tenantId"`
	OccurredAt time.Time `json:"occurredAt"`
}

type TerminatedEvent struct {
	LicenseID  string    `json:"licenseId"`
	TenantID   string    `json:"tenantId"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurredAt"`
}

// TokenTopupEvent — wallet topped up (CSM action or paid invoice).
type TokenTopupEvent struct {
	LicenseID    string     `json:"licenseId"`
	TenantID     string     `json:"tenantId"`
	WalletID     string     `json:"walletId"`
	SKUCode      string     `json:"skuCode"`
	CurrencyUnit string     `json:"currencyUnit"`
	Amount       int64      `json:"amount"`
	NewBalance   int64      `json:"newBalance"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	OccurredAt   time.Time  `json:"occurredAt"`
}

// TokenConsumedEvent — wallet decrement Committed. The billing-service's
// usage-aggregation pipeline picks this up to cut overage invoice lines.
//
// Tax context (RateCardID + UnitPricePaise + Currency + HSNSACCode +
// GSTRatePercent + TaxInclusive) is included so billing-service can
// compute the correct GST split (CGST/SGST when tenant-state ==
// system-state; IGST otherwise) without a callback to license-service.
// This is the per-event tax data; tenant-level tax overrides
// (e.g. composition scheme) still apply via billing's tenant.taxation.
type TokenConsumedEvent struct {
	LicenseID     string    `json:"licenseId"`
	TenantID      string    `json:"tenantId"`
	WalletID      string    `json:"walletId"`
	SKUCode       string    `json:"skuCode"`
	Component     string    `json:"component"`  // INFERENCE, SERVICE_PLANE, ...
	Deployment    string    `json:"deployment"` // VENDOR_GPU / CUSTOMER_GPU / ...
	CurrencyUnit  string    `json:"currencyUnit"`
	ActualQty     int64     `json:"actualQty"`
	BalanceAfter  int64     `json:"balanceAfter"`
	ReservationID string    `json:"reservationId"`

	// Tax context — billing-service uses these to compute the invoice
	// line. Empty when the wallet was billed flat-rate via subscription
	// (no per-event tax computation needed).
	RateCardID      string  `json:"rateCardId,omitempty"`
	UnitPricePaise  int64   `json:"unitPricePaise,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	HSNSACCode      string  `json:"hsnSacCode,omitempty"`
	GSTRatePercent  float64 `json:"gstRatePercent,omitempty"`
	TaxInclusive    bool    `json:"taxInclusive,omitempty"`
	IsOverage       bool    `json:"isOverage,omitempty"` // true if this consumption breached the included pool

	OccurredAt time.Time `json:"occurredAt"`
}

// TokenLowBalanceEvent — included pool fell below low_balance_threshold.
type TokenLowBalanceEvent struct {
	LicenseID    string    `json:"licenseId"`
	TenantID     string    `json:"tenantId"`
	WalletID     string    `json:"walletId"`
	SKUCode      string    `json:"skuCode"`
	CurrencyUnit string    `json:"currencyUnit"`
	Remaining    int64     `json:"remaining"`
	Threshold    int64     `json:"threshold"`
	OccurredAt   time.Time `json:"occurredAt"`
}

// TokenExhaustedEvent — both pools at zero, overage policy denied or unset.
type TokenExhaustedEvent struct {
	LicenseID    string    `json:"licenseId"`
	TenantID     string    `json:"tenantId"`
	WalletID     string    `json:"walletId"`
	SKUCode      string    `json:"skuCode"`
	CurrencyUnit string    `json:"currencyUnit"`
	OccurredAt   time.Time `json:"occurredAt"`
}

// InstallationEnrolledEvent — customer cluster completed first cloud contact.
type InstallationEnrolledEvent struct {
	LicenseID           string    `json:"licenseId"`
	TenantID            string    `json:"tenantId"`
	InstallationID      string    `json:"installationId"`
	Name                string    `json:"name"`
	Kind                string    `json:"kind"`
	Region              string    `json:"region"`
	MTLSCertFingerprint string    `json:"mtlsCertFingerprint"`
	OccurredAt          time.Time `json:"occurredAt"`
}

// InstallationOfflineEvent — heartbeat lapsed past offline_warn threshold.
type InstallationOfflineEvent struct {
	LicenseID      string    `json:"licenseId"`
	TenantID       string    `json:"tenantId"`
	InstallationID string    `json:"installationId"`
	OfflineSince   time.Time `json:"offlineSince"`
	Severity       string    `json:"severity"` // WARN | ESCALATED | EXCEEDED
	OccurredAt     time.Time `json:"occurredAt"`
}

// InstallationRecoveredEvent — heartbeat resumed.
type InstallationRecoveredEvent struct {
	LicenseID       string    `json:"licenseId"`
	TenantID        string    `json:"tenantId"`
	InstallationID  string    `json:"installationId"`
	WasOfflineHours float64   `json:"wasOfflineHours"`
	OccurredAt      time.Time `json:"occurredAt"`
}

// InstallationDecommissionedEvent — CSM removed a cluster.
type InstallationDecommissionedEvent struct {
	LicenseID      string    `json:"licenseId"`
	TenantID       string    `json:"tenantId"`
	InstallationID string    `json:"installationId"`
	Reason         string    `json:"reason"`
	OccurredAt     time.Time `json:"occurredAt"`
}

// JWSRevokedEvent — broadcast that one JWS (by jti) or every JWS for a
// given licenseId must not be honored for the rest of its TTL.
// Verifiers maintain a short-lived deny set keyed by jti when present,
// or by licenseId when JTI is empty (license-level revocation, fired
// on Suspend / Expire / Terminate).
type JWSRevokedEvent struct {
	JTI        string    `json:"jti,omitempty"`
	LicenseID  string    `json:"licenseId"`
	ExpiresAt  time.Time `json:"expiresAt"` // when the deny entry is safe to drop (= max JWS exp)
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurredAt"`
}

// TierExceededEvent — PER_SUBSCRIBER pricing crossed the licensed
// quantity for an entitlement. billing-service consumes to add an
// overage line to the next invoice; CSM consumes to email upsell.
type TierExceededEvent struct {
	LicenseID         string    `json:"licenseId"`
	TenantID          string    `json:"tenantId"`
	SKUCode           string    `json:"skuCode"`
	LicensedQuantity  int64     `json:"licensedQuantity"`
	CurrentQuantity   int64     `json:"currentQuantity"`
	OveragePolicy     string    `json:"overagePolicy"`
	OccurredAt        time.Time `json:"occurredAt"`
}

// TierRecoveredEvent — subscriber count fell back below the licensed
// quantity (e.g. seasonal churn). Lets billing stop accruing overage.
type TierRecoveredEvent struct {
	LicenseID         string    `json:"licenseId"`
	TenantID          string    `json:"tenantId"`
	SKUCode           string    `json:"skuCode"`
	LicensedQuantity  int64     `json:"licensedQuantity"`
	CurrentQuantity   int64     `json:"currentQuantity"`
	OccurredAt        time.Time `json:"occurredAt"`
}
