package provider

// ESignProviderType defines supported eSign provider types.
type ESignProviderType string

const (
	ESignProviderCashfree ESignProviderType = "CASHFREE"
	ESignProviderSetu     ESignProviderType = "SETU"
	ESignProviderDigio    ESignProviderType = "DIGIO"
	ESignProviderSignzy   ESignProviderType = "SIGNZY"
	ESignProviderLeegality ESignProviderType = "LEEGALITY"
	ESignProviderSignDesk ESignProviderType = "SIGNDESK"
)

// ESignBindingScope represents who can use an eSign binding within a tenant tree.
type ESignBindingScope string

const (
	ESignBindingScopeOwnerOnly            ESignBindingScope = "OwnerOnly"
	ESignBindingScopeOwnerAndDescendants  ESignBindingScope = "OwnerAndDescendants"
	ESignBindingScopeExplicitTenants      ESignBindingScope = "ExplicitTenants"
)

// ESignProviderInfo holds display metadata and field definitions for one eSign provider.
type ESignProviderInfo struct {
	Value       string              `json:"value"`
	Label       string              `json:"label"`
	Description string              `json:"description"`
	Fields      []FieldSchema       `json:"fields"`
}

// ESignFormConfig is the complete form metadata the frontend needs to render the eSign binding form.
type ESignFormConfig struct {
	Providers         []ESignProviderInfo `json:"providers"`
	Scopes            []string            `json:"scopes"`
	AllowedOwnerTypes []string            `json:"allowedOwnerTypes"`
}

// ESignProviderRegistry is the single source of truth for all eSign providers.
// Adding a new provider = add one entry here; frontend auto-renders fields.
var ESignProviderRegistry = map[string]ESignProviderInfo{
	"CASHFREE": {
		Value:       "CASHFREE",
		Label:       "Cashfree",
		Description: "Cashfree Verification Suite — Aadhaar-based eSign for legally binding documents",
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.cashfree.com/verification", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://sandbox.cashfree.com/verification", Label: "Sandbox"},
				{Value: "https://api.cashfree.com/verification", Label: "Production"},
			}},
			{Key: "x-client-id", Label: "Client ID", Placeholder: "CF12345678", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "x-client-secret", Label: "Client Secret", Placeholder: "cf-secret-key", Required: true, Sensitive: true, MinLength: 2, MaxLength: 100},
		},
	},
	"LEEGALITY": {
		Value:       "LEEGALITY",
		Label:       "Leegality",
		Description: "Leegality — Aadhaar eSign, Digital Signature, eStamp integration",
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.leegality.com", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://sandbox.leegality.com", Label: "Sandbox"},
				{Value: "https://api.leegality.com", Label: "Production"},
			}},
			{Key: "api-key", Label: "API Key", Placeholder: "Your Leegality API key", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
		},
	},
	"DIGIO": {
		Value:       "DIGIO",
		Label:       "Digio",
		Description: "Digio — eSign, Digital Signature, eStamp, Smart Agreements",
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.digio.in", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://ext.digio.in", Label: "Sandbox"},
				{Value: "https://api.digio.in", Label: "Production"},
			}},
			{Key: "client-id", Label: "Client ID", Placeholder: "Your Digio client ID", Required: true, MinLength: 8, MaxLength: 100},
			{Key: "client-secret", Label: "Client Secret", Placeholder: "Your Digio client secret", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
		},
	},
	"SIGNDESK": {
		Value:       "SIGNDESK",
		Label:       "SignDesk",
		Description: "SignDesk — eSign, eStamp, document automation for Indian enterprises",
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.signdesk.com", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://sandbox.signdesk.com", Label: "Sandbox"},
				{Value: "https://api.signdesk.com", Label: "Production"},
			}},
			{Key: "api-key", Label: "API Key", Placeholder: "Your SignDesk API key", Required: true, Sensitive: true, MinLength: 8, MaxLength: 200},
			{Key: "account-id", Label: "Account ID", Placeholder: "Your SignDesk account ID", Required: true, MinLength: 5, MaxLength: 100},
		},
	},
	"SIGNZY": {
		Value:       "SIGNZY",
		Label:       "Signzy",
		Description: "Signzy — AI-powered eSign, Video KYC, and document verification",
		Fields: []FieldSchema{
			{Key: "url", Label: "API URL", Placeholder: "https://api.signzy.app", Required: true, IsURL: true, Options: []FieldOption{
				{Value: "https://preproduction.signzy.tech", Label: "Sandbox"},
				{Value: "https://api.signzy.app", Label: "Production"},
			}},
			{Key: "authorization", Label: "Authorization Token", Placeholder: "Bearer token from Signzy", Required: true, Sensitive: true, MinLength: 10, MaxLength: 500},
		},
	},
}

// GetESignFormConfig returns the complete form configuration for the frontend.
func GetESignFormConfig() ESignFormConfig {
	providers := make([]ESignProviderInfo, 0, len(ESignProviderRegistry))
	for _, info := range ESignProviderRegistry {
		providers = append(providers, info)
	}
	return ESignFormConfig{
		Providers:         providers,
		Scopes:            []string{"OwnerOnly", "OwnerAndDescendants", "ExplicitTenants"},
		AllowedOwnerTypes: []string{"ISP"},
	}
}
