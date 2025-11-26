package tenantevent

// MFAFactorType - supported MFA methods
type MFAFactorType string

const (
	MFAFactorPassword         MFAFactorType = "password"
	MFAFactorOTPSMS           MFAFactorType = "otp_sms"
	MFAFactorOTPEmail         MFAFactorType = "otp_email"
	MFAFactorTOTPApp          MFAFactorType = "totp_app"
	MFAFactorPushNotification MFAFactorType = "push_app"
	MFAFactorWebAuthnPasskey  MFAFactorType = "webauthn_passkey"
	MFAFactorWebAuthnDevice   MFAFactorType = "webauthn_device"
	MFAFactorHardwareKeyU2F   MFAFactorType = "hardware_key_u2f"
	MFAFactorBiometric        MFAFactorType = "biometric"
)

// MFARiskActionType - high-risk actions that may require fresh MFA
type MFARiskActionType string

const (
	MFARiskActionLogin                MFARiskActionType = "login"
	MFARiskActionPasswordChange       MFARiskActionType = "password_change"
	MFARiskActionProfileUpdate        MFARiskActionType = "profile_update"
	MFARiskActionPaymentOperation     MFARiskActionType = "payment_operation"
	MFARiskActionAPIKeyCreate         MFARiskActionType = "api_key_create"
	MFARiskActionAPIKeyRevoke         MFARiskActionType = "api_key_revoke"
	MFARiskActionKYCOperation         MFARiskActionType = "kyc_operation"
	MFARiskActionAdminPortalAccess    MFARiskActionType = "admin_portal_access"
	MFARiskActionIntegrationChange    MFARiskActionType = "integration_change"
	MFARiskActionNetworkConfigChange  MFARiskActionType = "network_config_change"
	MFARiskActionRadiusConfigChange   MFARiskActionType = "radius_config_change"
	MFARiskActionCriticalConfigChange MFARiskActionType = "critical_config_change"
)

type ConfigScope string

const (
	ScopeSystem ConfigScope = "system"
	ScopeTenant ConfigScope = "tenant"
)

// TenantMFAPolicy is the SINGLE aggregate for all MFA-related settings
// for a tenant. One doc per tenant in collection: tenant_mfa_policies.
type TenantMFAPolicy struct {
	ID       string      `json:"id" bson:"_id"`      // cuid2
	Scope    ConfigScope `json:"scope" bson:"scope"` // system|tenant
	Name     string      `json:"name" bson:"name"`
	TenantID string      `json:"tenantID" bson:"tenantID" validate:"required,isCuid2"`

	// --- Inheritance & control ---

	// If true → this tenant does NOT have its own MFA config,
	// it uses the parent's effective MFA policy.
	UseParentPolicy bool `json:"useParentPolicy" bson:"useParentPolicy"`

	// If false on parent → children are not allowed to weaken / override MFA.
	// You enforce this in service layer when saving child policies.
	OverrideAllowed bool `json:"overrideAllowed" bson:"overrideAllowed"`

	// --- Global requirements ---

	// Is MFA enabled and required for this tenant at all?
	IsMFARequired bool `json:"isMFARequired" bson:"isMFARequired"`

	// Minimum count of factors required:
	// 0 = no MFA,
	// 1 = password + 1 factor,
	// 2 = strong MFA / passwordless + extra factor.
	MinFactors int `json:"minFactors" bson:"minFactors" validate:"min=0,max=3"`

	// Which MFA methods are allowed for users in this tenant.
	AllowedMethods []MFAFactorType `json:"allowedMethods,omitempty" bson:"allowedMethods,omitempty"`

	// Which methods are mandatory for certain roles or risk levels (optional).
	// e.g. admins must have TOTP or WebAuthn.
	RequiredMethods []MFAFactorType `json:"requiredMethods,omitempty" bson:"requiredMethods,omitempty"`

	// --- High-risk action rules ---

	// If true → the platform will force fresh MFA for the listed actions,
	// even if the user already has a valid session.
	HighRiskRequireMFA bool                `json:"highRiskRequireMFA" bson:"highRiskRequireMFA"`
	HighRiskActions    []MFARiskActionType `json:"highRiskActions,omitempty" bson:"highRiskActions,omitempty"`

	// --- Session & re-auth policy ---

	// Total session lifetime before full re-login (minutes)
	SessionTimeoutMinutes int `json:"sessionTimeoutMinutes" bson:"sessionTimeoutMinutes" validate:"min=5,max=1440"`

	// Idle timeout before logout (minutes)
	IdleTimeoutMinutes int `json:"idleTimeoutMinutes" bson:"idleTimeoutMinutes" validate:"min=1,max=120"`

	// How often MFA must be re-verified even if session is active (hours)
	MFAReauthIntervalHours int `json:"mfaReauthIntervalHours" bson:"mfaReauthIntervalHours" validate:"min=1,max=168"`

	// --- Optional network-level controls ---

	IPRestrictionsEnabled bool     `json:"ipRestrictionsEnabled" bson:"ipRestrictionsEnabled"`
	AllowedCIDRs          []string `json:"allowedCIDRs,omitempty" bson:"allowedCIDRs,omitempty"`

	GeoRestrictionsEnabled bool     `json:"geoRestrictionsEnabled" bson:"geoRestrictionsEnabled"`
	AllowedCountries       []string `json:"allowedCountries,omitempty" bson:"allowedCountries,omitempty"`

	Version int `json:"version" bson:"version"`
}
