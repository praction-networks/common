package provider

// FieldOption represents a predefined option for a field (e.g. sandbox/production URLs)
type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// FieldSchema describes a single metadata field for a KYC provider.
// The frontend reads these to dynamically render form inputs.
type FieldSchema struct {
	Key         string        `json:"key"`
	Label       string        `json:"label"`
	Placeholder string        `json:"placeholder"`
	Required    bool          `json:"required"`
	Sensitive   bool          `json:"sensitive"`
	MinLength   int           `json:"minLength,omitempty"`
	MaxLength   int           `json:"maxLength,omitempty"`
	IsURL       bool          `json:"isUrl,omitempty"`
	Options     []FieldOption `json:"options,omitempty"`
}

// KYCProviderInfo holds display metadata and field definitions for one KYC provider.
type KYCProviderInfo struct {
	Value       string        `json:"value"`
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Fields      []FieldSchema `json:"fields"`
}

// KYCFormConfig is the complete form metadata the frontend needs to render the KYC gateway form.
type KYCFormConfig struct {
	Providers         []KYCProviderInfo `json:"providers"`
	Scopes            []string          `json:"scopes"`
	AllowedOwnerTypes []string          `json:"allowedOwnerTypes"`
}

// KYCProviderRegistry is the single source of truth for all KYC providers.
// Adding a new provider = add one entry here; frontend auto-renders fields.
var KYCProviderRegistry = map[string]KYCProviderInfo{
	"CASHFREE": {
		Value:       "CASHFREE",
		Label:       "Cashfree",
		Description: "Cashfree Verification Suite — PAN, Aadhaar, Bank",
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.cashfree.com/verification", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://sandbox.cashfree.com/verification", Label: "Sandbox"},
				{Value: "https://api.cashfree.com/verification", Label: "Production"},
			}},
			{Key: "x-client-id", Label: "Client ID", Placeholder: "CF12345678", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "x-client-secret", Label: "Client Secret", Placeholder: "cf-secret-key", Required: true, Sensitive: true, MinLength: 2, MaxLength: 100},
			{Key: "x-cf-signature", Label: "CF Signature", Placeholder: "Webhook signature key (optional)", Sensitive: true, MaxLength: 500},
		},
	},
	"SETU": {
		Value:       "SETU",
		Label:       "Setu",
		Description: "Setu DigiLocker Gateway — PAN, Aadhaar, Bank Account",
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
