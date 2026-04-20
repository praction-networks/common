package tenantevent

// Note: SMS Provider Events have been moved to sms.provider.event.model.go
// as part of the Unified Provider Architecture refactor.
// Note: Mail Provider Events live in mail.provider.event.model.go — the old
// MailServer binding event models were retired when the unified
// TenantMailProvider shape replaced the two-table template+binding stack.

// ==================== KYC PROVIDER BINDING EVENTS ====================

// TenantKYCProviderBindingInsertEventModel represents a KYC provider binding creation event
type TenantKYCProviderBindingInsertEventModel struct {
	ID                string                   `bson:"_id" json:"id"`
	TenantID          string                   `bson:"tenantId" json:"tenantId"`
	ProviderType      string                   `bson:"providerType" json:"providerType"`
	Scope             string                   `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                 `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	ProviderID        string                   `bson:"providerId,omitempty" json:"providerId,omitempty"`
	UseTemplate       bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent         bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID    string                   `bson:"parentTenantId,omitempty" json:"parentTenantId,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig      *TenantKYCProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority          int                      `bson:"priority" json:"priority"`
	IsActive          bool                     `bson:"isActive" json:"isActive"`
	FailoverOn        bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries        int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider  string                   `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version           int                      `bson:"version" json:"version"`
}

// TenantKYCProviderBindingUpdateEventModel represents a KYC provider binding update event
type TenantKYCProviderBindingUpdateEventModel struct {
	ID                string                   `bson:"_id" json:"id"`
	TenantID          string                   `bson:"tenantId" json:"tenantId"`
	ProviderType      string                   `bson:"providerType" json:"providerType"`
	Scope             string                   `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                 `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	ProviderID        string                   `bson:"providerId,omitempty" json:"providerId,omitempty"`
	UseTemplate       bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent         bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID    string                   `bson:"parentTenantId,omitempty" json:"parentTenantId,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig      *TenantKYCProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority          int                      `bson:"priority" json:"priority"`
	IsActive          bool                     `bson:"isActive" json:"isActive"`
	FailoverOn        bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries        int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider  string                   `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version           int                      `bson:"version" json:"version"`
}

// TenantKYCProviderBindingDeleteEventModel represents a KYC provider binding deletion event
type TenantKYCProviderBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`
}

// ==================== APP MESSAGING BINDING EVENTS ====================

// TenantAppMessagingBindingInsertEventModel represents an App Messaging provider binding creation event
// Contains only operational data needed by consumers (e.g., notification-service)
type TenantAppMessagingBindingInsertEventModel struct {
	ID                       string                    `bson:"_id" json:"id"`
	TenantID                 string                    `bson:"tenantId" json:"tenantId"`
	Channel                  string                    `bson:"channel" json:"channel"`
	ProviderID               string                    `bson:"providerId" json:"providerId"`
	CachedProviderURL        string                    `bson:"cachedProviderUrl" json:"cachedProviderUrl"`
	CachedProviderAPIVersion string                    `bson:"cachedProviderApiVersion" json:"cachedProviderApiVersion"`
	CachedProviderName       string                    `bson:"cachedProviderName" json:"cachedProviderName"`
	TenantConfig             *TenantAppMessagingConfig `bson:"tenantConfig" json:"tenantConfig"`
	TemplateBindings         []TemplateBinding         `bson:"templateBindings,omitempty" json:"templateBindings,omitempty"`
	Version                  int                       `bson:"version" json:"version"`
}

// TenantAppMessagingBindingUpdateEventModel represents an App Messaging provider binding update event
// Contains only operational data needed by consumers (e.g., notification-service)
type TenantAppMessagingBindingUpdateEventModel struct {
	ID                       string                    `bson:"_id" json:"id"`
	TenantID                 string                    `bson:"tenantId" json:"tenantId"`
	Channel                  string                    `bson:"channel" json:"channel"`
	ProviderID               string                    `bson:"providerId,omitempty" json:"providerId,omitempty"`
	CachedProviderURL        string                    `bson:"cachedProviderUrl,omitempty" json:"cachedProviderUrl,omitempty"`
	CachedProviderAPIVersion string                    `bson:"cachedProviderApiVersion,omitempty" json:"cachedProviderApiVersion,omitempty"`
	CachedProviderName       string                    `bson:"cachedProviderName,omitempty" json:"cachedProviderName,omitempty"`
	TenantConfig             *TenantAppMessagingConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	TemplateBindings         []TemplateBinding         `bson:"templateBindings,omitempty" json:"templateBindings,omitempty"`
	Version                  int                       `bson:"version" json:"version"`
}

// TenantAppMessagingBindingDeleteEventModel represents an App Messaging provider binding deletion event
type TenantAppMessagingBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`
}

// ==================== CDN PROVIDER BINDING EVENTS ====================

// TenantCDNProviderBindingInsertEventModel represents a CDN provider binding creation event
type TenantCDNProviderBindingInsertEventModel struct {
	ID               string                   `bson:"_id" json:"id"`
	TenantID         string                   `bson:"tenantId" json:"tenantId"`
	ProviderID       string                   `bson:"providerId,omitempty" json:"providerId,omitempty"`
	UseTemplate      bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                   `bson:"parentTenantId,omitempty" json:"parentTenantId,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantCDNProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                      `bson:"priority" json:"priority"`
	IsActive         bool                     `bson:"isActive" json:"isActive"`
	FailoverOn       bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedSortCode string                   `bson:"resolvedSortCode,omitempty" json:"resolvedSortCode,omitempty"`
	Version          int                      `bson:"version" json:"version"`
}

// TenantCDNProviderBindingUpdateEventModel represents a CDN provider binding update event
type TenantCDNProviderBindingUpdateEventModel struct {
	ID               string                   `bson:"_id" json:"id"`
	TenantID         string                   `bson:"tenantId" json:"tenantId"`
	ProviderID       string                   `bson:"providerId,omitempty" json:"providerId,omitempty"`
	UseTemplate      bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                   `bson:"parentTenantId,omitempty" json:"parentTenantId,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantCDNProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                      `bson:"priority" json:"priority"`
	IsActive         bool                     `bson:"isActive" json:"isActive"`
	FailoverOn       bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedSortCode string                   `bson:"resolvedSortCode,omitempty" json:"resolvedSortCode,omitempty"`
	Version          int                      `bson:"version" json:"version"`
}

// TenantCDNProviderBindingDeleteEventModel represents a CDN provider binding deletion event
type TenantCDNProviderBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`
}
