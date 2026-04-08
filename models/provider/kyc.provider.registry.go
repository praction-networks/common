package provider

// KYCBindingScope represents who can use a KYC provider binding within a tenant tree.
// Used by tenant-service (binding owner) and consumer services (subscriber-service).
type KYCBindingScope string

const (
	// KYCBindingScopeOwnerOnly means only the owner tenant can use this binding
	KYCBindingScopeOwnerOnly KYCBindingScope = "OwnerOnly"
	// KYCBindingScopeOwnerAndDescendants means owner and all descendant tenants can use this binding
	KYCBindingScopeOwnerAndDescendants KYCBindingScope = "OwnerAndDescendants"
	// KYCBindingScopeExplicitTenants means only explicitly listed tenants can use this binding.
	// Explicit tenants must be descendants of the binding owner.
	KYCBindingScopeExplicitTenants KYCBindingScope = "ExplicitTenants"
)

// WebhookAuthType represents the type of webhook authentication.
// Configured per-tenant in the KYC binding and validated by the webhook handler.
type WebhookAuthType string

const (
	// WebhookAuthNone means no additional auth — rely on signature verification only
	WebhookAuthNone WebhookAuthType = "none"
	// WebhookAuthBearer validates Authorization: Bearer <token>
	WebhookAuthBearer WebhookAuthType = "bearer"
	// WebhookAuthBasic validates Authorization: Basic <base64(user:pass)>
	WebhookAuthBasic WebhookAuthType = "basic"
	// WebhookAuthCustomHeader validates a custom header name against the configured token
	WebhookAuthCustomHeader WebhookAuthType = "custom_header"
)

// WebhookAuthConfig holds the webhook authentication configuration for a KYC provider binding.
// Stored in the binding's Metadata map and validated by subscriber-service before signature verification.
type WebhookAuthConfig struct {
	AuthType   WebhookAuthType `json:"authType" bson:"authType"`
	Token      string          `json:"token,omitempty" bson:"token,omitempty"`           // Bearer token, Basic auth base64, or custom header value
	HeaderName string          `json:"headerName,omitempty" bson:"headerName,omitempty"` // Custom header name (for custom_header type)
}

// FieldOption represents a predefined option for a field (e.g. sandbox/production URLs)
type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// FieldSchema describes a single metadata field for a KYC provider.
// The frontend reads these to dynamically render form inputs.
type FieldSchema struct {
	Key            string        `json:"key"`
	Label          string        `json:"label"`
	Placeholder    string        `json:"placeholder"`
	Required       bool          `json:"required"`
	Sensitive      bool          `json:"sensitive"`
	MinLength      int           `json:"minLength,omitempty"`
	MaxLength      int           `json:"maxLength,omitempty"`
	IsURL          bool          `json:"isUrl,omitempty"`
	Options        []FieldOption `json:"options,omitempty"`
	Computed       bool          `json:"computed,omitempty"`       // If true, field is read-only and auto-built from ComputePattern
	ComputePattern string        `json:"computePattern,omitempty"` // Template using {fieldKey} substitution, e.g. "https://{accountId}.r2.cloudflarestorage.com"
}

// KYCProviderInfo holds display metadata, field definitions, and supported verification types for one KYC provider.
type KYCProviderInfo struct {
	Value          string             `json:"value"`
	Label          string             `json:"label"`
	Description    string             `json:"description"`
	Fields         []FieldSchema      `json:"fields"`
	SupportedTypes []VerificationType `json:"supportedTypes"`
}

// KYCFormConfig is the complete form metadata the frontend needs to render the KYC gateway form.
type KYCFormConfig struct {
	Providers         []KYCProviderInfo `json:"providers"`
	Scopes            []string          `json:"scopes"`
	AllowedOwnerTypes []string          `json:"allowedOwnerTypes"`
}

// allTypes is a shorthand for providers that support all 7 verification types.
var allTypes = []VerificationType{
	VerifyPAN, VerifyDigiLocker, VerifyPennydrop,
	VerifyEStamp, VerifyESign, VerifyGST, VerifyPassport,
}

// KYCProviderRegistry is the single source of truth for all KYC providers.
// Adding a new provider = add one entry here; frontend auto-renders fields.
var KYCProviderRegistry = map[string]KYCProviderInfo{
	"CASHFREE": {
		Value:          "CASHFREE",
		Label:          "Cashfree",
		Description:    "Cashfree Verification Suite — PAN, Aadhaar, Bank, eSign, Passport",
		SupportedTypes: allTypes,
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.cashfree.com/verification", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://sandbox.cashfree.com/verification", Label: "Sandbox"},
				{Value: "https://api.cashfree.com/verification", Label: "Production"},
			}},
			{Key: "x-client-id", Label: "Client ID", Placeholder: "CF12345678", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "x-client-secret", Label: "Client Secret", Placeholder: "cf-secret-key", Required: true, Sensitive: true, MinLength: 2, MaxLength: 100},
			// Webhook authentication fields (Cashfree Security Checklist — Authentication Validation)
			// These are optional and provide an additional layer on top of HMAC signature verification.
			{Key: "webhook-auth-type", Label: "Webhook Auth Type", Placeholder: "none", Options: []FieldOption{
				{Value: "none", Label: "None (Signature Only)"},
				{Value: "bearer", Label: "Bearer Token"},
				{Value: "basic", Label: "Basic Auth"},
				{Value: "custom_header", Label: "Custom Header"},
			}},
			{Key: "webhook-auth-token", Label: "Webhook Auth Token", Placeholder: "Token value shared with Cashfree", Sensitive: true, MaxLength: 500},
			{Key: "webhook-auth-header", Label: "Custom Auth Header Name", Placeholder: "X-Custom-Auth", MaxLength: 100},
		},
	},
	"SETU": {
		Value:       "SETU",
		Label:       "Setu",
		Description: "Setu DigiLocker Gateway — PAN, Aadhaar, Bank Account, eSign, GST",
		SupportedTypes: []VerificationType{
			VerifyPAN, VerifyDigiLocker, VerifyPennydrop, VerifyESign, VerifyGST,
		},
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://dg.setu.co", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://dg-sandbox.setu.co", Label: "Sandbox"},
				{Value: "https://dg.setu.co", Label: "Production"},
			}},
			{Key: "x-client-id", Label: "Client ID", Placeholder: "Your Setu client ID", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "x-client-secret", Label: "Client Secret", Placeholder: "Your Setu client secret", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
			{Key: "x-product-instance-id", Label: "Product Instance ID", Placeholder: "Product instance identifier", Required: true, MinLength: 8, MaxLength: 100},
		},
	},
	"DIGIO": {
		Value:          "DIGIO",
		Label:          "Digio",
		Description:    "Digio — eKYC, eSign, eStamp, Aadhaar, DigiLocker, Passport",
		SupportedTypes: allTypes,
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.digio.in", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://ext.digio.in", Label: "Sandbox"},
				{Value: "https://api.digio.in", Label: "Production"},
			}},
			{Key: "client-id", Label: "Client ID", Placeholder: "Your Digio client ID", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "client-secret", Label: "Client Secret", Placeholder: "Your Digio client secret", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
		},
	},
	"SIGNZY": {
		Value:          "SIGNZY",
		Label:          "Signzy",
		Description:    "Signzy — Video KYC, AI Document Verification, eSign, eStamp",
		SupportedTypes: allTypes,
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.signzy.app", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://preproduction.signzy.tech", Label: "Sandbox"},
				{Value: "https://api.signzy.app", Label: "Production"},
			}},
			{Key: "authorization", Label: "Authorization Token", Placeholder: "Bearer token from Signzy", Required: true, Sensitive: true, MinLength: 10, MaxLength: 500},
			{Key: "callback-url", Label: "Callback URL", Placeholder: "https://your-domain.com/webhook/signzy", IsURL: true, MaxLength: 500},
		},
	},
	"IDFY": {
		Value:          "IDFY",
		Label:          "IDfy",
		Description:    "IDfy — Face Match, Liveness Detection, PAN/Aadhaar/DL/Passport",
		SupportedTypes: allTypes,
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://eve.idfy.com", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://eve-sandbox.idfy.com", Label: "Sandbox"},
				{Value: "https://eve.idfy.com", Label: "Production"},
			}},
			{Key: "api-key", Label: "API Key", Placeholder: "Your IDfy API key", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
			{Key: "account-id", Label: "Account ID", Placeholder: "Your IDfy account ID", Required: true, MinLength: 5, MaxLength: 100},
		},
	},
	"KARZA": {
		Value:       "KARZA",
		Label:       "Karza (Perfios)",
		Description: "Karza by Perfios — PAN, GST, Bank, Aadhaar, Passport, eSign",
		SupportedTypes: []VerificationType{
			VerifyPAN, VerifyDigiLocker, VerifyPennydrop, VerifyESign, VerifyGST, VerifyPassport,
		},
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.karza.in", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://testapi.karza.in", Label: "Sandbox"},
				{Value: "https://api.karza.in", Label: "Production"},
			}},
			{Key: "x-karza-key", Label: "Karza API Key", Placeholder: "Your Karza API key", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
		},
	},
}

// GetKYCFormConfig returns the complete form configuration for the frontend.
func GetKYCFormConfig() KYCFormConfig {
	providers := make([]KYCProviderInfo, 0, len(KYCProviderRegistry))
	for _, info := range KYCProviderRegistry {
		providers = append(providers, info)
	}
	return KYCFormConfig{
		Providers:         providers,
		Scopes:            []string{"OwnerOnly", "OwnerAndDescendants", "ExplicitTenants"},
		AllowedOwnerTypes: []string{"ISP", "Enterprise"},
	}
}
