package tenantevent

// MFAFactorType - supported MFA methods
type MFAFactorType string

const (
	MFAFactorPassword         MFAFactorType = "password"
	MFAFactorOTPSMS           MFAFactorType = "otpSms"
	MFAFactorOTPEmail         MFAFactorType = "otpEmail"
	MFAFactorTOTPApp          MFAFactorType = "totpApp"
	MFAFactorPushNotification MFAFactorType = "pushApp"
	MFAFactorWebAuthnPasskey  MFAFactorType = "webauthnPasskey"
	MFAFactorWebAuthnDevice   MFAFactorType = "webauthnDevice"
	MFAFactorHardwareKeyU2F   MFAFactorType = "hardwareKeyU2F"
	MFAFactorBiometric        MFAFactorType = "biometric"
)

// MFARiskActionType - high-risk actions that may require fresh MFA
type MFARiskActionType string

const (
	MFARiskActionLogin                MFARiskActionType = "login"
	MFARiskActionPasswordChange       MFARiskActionType = "passwordChange"
	MFARiskActionProfileUpdate        MFARiskActionType = "profile_update"
	MFARiskActionPaymentOperation     MFARiskActionType = "paymentOperation"
	MFARiskActionAPIKeyCreate         MFARiskActionType = "apiKeyCreate"
	MFARiskActionAPIKeyRevoke         MFARiskActionType = "apiKeyRevoke"
	MFARiskActionKYCOperation         MFARiskActionType = "kycOperation"
	MFARiskActionAdminPortalAccess    MFARiskActionType = "adminPortalAccess"
	MFARiskActionIntegrationChange    MFARiskActionType = "integrationChange"
	MFARiskActionNetworkConfigChange  MFARiskActionType = "networkConfigChange"
	MFARiskActionRadiusConfigChange   MFARiskActionType = "radiusConfigChange"
	MFARiskActionCriticalConfigChange MFARiskActionType = "criticalConfigChange"
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
	TenantID string      `json:"tenantId" bson:"tenantId"`

	// --- Inheritance & control ---

	// If true → this tenant does NOT have its own MFA config,
	// it uses the parent's effective MFA policy.
	UseParentPolicy bool `json:"useParentPolicy" bson:"useParentPolicy"`

	// If false on parent → children are not allowed to weaken / override MFA.
	// You enforce this in service layer when saving child policies.
	OverrideAllowed bool `json:"overrideAllowed" bson:"overrideAllowed"`

	// --- Global requirements ---

	// Is MFA enabled and required for this tenant at all?
	IsMFARequired bool `json:"isMfaRequired" bson:"isMfaRequired"`

	// Minimum count of factors required:
	// 0 = no MFA,
	// 1 = password + 1 factor,
	// 2 = strong MFA / passwordless + extra factor.
	MinFactors int `json:"minFactors" bson:"minFactors"`

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
	SessionTimeoutMinutes int `json:"sessionTimeoutMinutes" bson:"sessionTimeoutMinutes"`

	// Idle timeout before logout (minutes)
	IdleTimeoutMinutes int `json:"idleTimeoutMinutes" bson:"idleTimeoutMinutes"`

	// How often MFA must be re-verified even if session is active (hours)
	MFAReauthIntervalHours int `json:"mfaReauthIntervalHours" bson:"mfaReauthIntervalHours"`

	// --- Optional network-level controls ---

	IPRestrictionsEnabled bool     `json:"ipRestrictionsEnabled" bson:"ipRestrictionsEnabled"`
	AllowedCIDRs          []string `json:"allowedCidrs,omitempty" bson:"allowedCidrs,omitempty"`

	GeoRestrictionsEnabled bool     `json:"geoRestrictionsEnabled" bson:"geoRestrictionsEnabled"`
	AllowedCountries       []string `json:"allowedCountries,omitempty" bson:"allowedCountries,omitempty"`

	Version int `json:"version" bson:"version"`
}
