package provider

// FieldSchema describes a single metadata field for a KYC provider.
// The frontend reads these to dynamically render form inputs.
type FieldSchema struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Placeholder string `json:"placeholder"`
	Required    bool   `json:"required"`
	Sensitive   bool   `json:"sensitive"`
	MinLength   int    `json:"minLength,omitempty"`
	MaxLength   int    `json:"maxLength,omitempty"`
	IsURL       bool   `json:"isUrl,omitempty"`
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
			{Key: "url", Label: "API URL", Placeholder: "https://api.cashfree.com/verification", Required: true, IsURL: true},
			{Key: "x-client-id", Label: "Client ID", Placeholder: "CF12345678", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "x-client-secret", Label: "Client Secret", Placeholder: "cf-secret-key", Required: true, Sensitive: true, MinLength: 2, MaxLength: 100},
			{Key: "x-cf-signature", Label: "CF Signature", Placeholder: "Webhook signature key (optional)", Sensitive: true, MaxLength: 500},
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
