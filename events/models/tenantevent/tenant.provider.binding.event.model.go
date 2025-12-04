package tenantevent

// ==================== SMS PROVIDER BINDING EVENTS ====================

// TenantSMSProviderBindingInsertEventModel represents an SMS provider binding creation event
type TenantSMSProviderBindingInsertEventModel struct {
	ID               string                   `bson:"_id" json:"id"`
	TenantID         string                   `bson:"tenantID" json:"tenantID"`
	Channel          string                   `bson:"channel" json:"channel"`
	ProviderID       string                   `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate      bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantSMSProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                      `bson:"priority" json:"priority"`
	IsActive         bool                     `bson:"isActive" json:"isActive"`
	FailoverOn       bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider string                   `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version          int                      `bson:"version" json:"version"`
}

// TenantSMSProviderBindingUpdateEventModel represents an SMS provider binding update event
type TenantSMSProviderBindingUpdateEventModel struct {
	ID               string                   `bson:"_id" json:"id"`
	TenantID         string                   `bson:"tenantID" json:"tenantID"`
	Channel          string                   `bson:"channel" json:"channel"`
	ProviderID       string                   `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate      bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantSMSProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                      `bson:"priority" json:"priority"`
	IsActive         bool                     `bson:"isActive" json:"isActive"`
	FailoverOn       bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider string                   `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version          int                      `bson:"version" json:"version"`
}

// TenantSMSProviderBindingDeleteEventModel represents an SMS provider binding deletion event
type TenantSMSProviderBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantID" json:"tenantID"`
}

// ==================== MAIL SERVER BINDING EVENTS ====================

// TenantMailServerBindingInsertEventModel represents a Mail server binding creation event
type TenantMailServerBindingInsertEventModel struct {
	ID               string                  `bson:"_id" json:"id"`
	TenantID         string                  `bson:"tenantID" json:"tenantID"`
	ProviderID       string                  `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate      bool                    `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                    `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                  `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantMailServerConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                     `bson:"priority" json:"priority"`
	IsActive         bool                    `bson:"isActive" json:"isActive"`
	FailoverOn       bool                    `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                     `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                     `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedSortCode string                  `bson:"resolvedSortCode,omitempty" json:"resolvedSortCode,omitempty"`
	Version          int                     `bson:"version" json:"version"`
}

// TenantMailServerBindingUpdateEventModel represents a Mail server binding update event
type TenantMailServerBindingUpdateEventModel struct {
	ID               string                  `bson:"_id" json:"id"`
	TenantID         string                  `bson:"tenantID" json:"tenantID"`
	ProviderID       string                  `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate      bool                    `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                    `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                  `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantMailServerConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                     `bson:"priority" json:"priority"`
	IsActive         bool                    `bson:"isActive" json:"isActive"`
	FailoverOn       bool                    `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                     `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                     `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedSortCode string                  `bson:"resolvedSortCode,omitempty" json:"resolvedSortCode,omitempty"`
	Version          int                     `bson:"version" json:"version"`
}

// TenantMailServerBindingDeleteEventModel represents a Mail server binding deletion event
type TenantMailServerBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantID" json:"tenantID"`
}

// ==================== KYC PROVIDER BINDING EVENTS ====================

// TenantKYCProviderBindingInsertEventModel represents a KYC provider binding creation event
type TenantKYCProviderBindingInsertEventModel struct {
	ID               string                   `bson:"_id" json:"id"`
	TenantID         string                   `bson:"tenantID" json:"tenantID"`
	ProviderType     string                   `bson:"providerType" json:"providerType"`
	ProviderID       string                   `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate      bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantKYCProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                      `bson:"priority" json:"priority"`
	IsActive         bool                     `bson:"isActive" json:"isActive"`
	FailoverOn       bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider string                   `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version          int                      `bson:"version" json:"version"`
}

// TenantKYCProviderBindingUpdateEventModel represents a KYC provider binding update event
type TenantKYCProviderBindingUpdateEventModel struct {
	ID               string                   `bson:"_id" json:"id"`
	TenantID         string                   `bson:"tenantID" json:"tenantID"`
	ProviderType     string                   `bson:"providerType" json:"providerType"`
	ProviderID       string                   `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate      bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent        bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID   string                   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig     *TenantKYCProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority         int                      `bson:"priority" json:"priority"`
	IsActive         bool                     `bson:"isActive" json:"isActive"`
	FailoverOn       bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries       int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight           int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider string                   `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version          int                      `bson:"version" json:"version"`
}

// TenantKYCProviderBindingDeleteEventModel represents a KYC provider binding deletion event
type TenantKYCProviderBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantID" json:"tenantID"`
}

// ==================== APP MESSAGING BINDING EVENTS ====================

// TenantAppMessagingBindingInsertEventModel represents an App Messaging provider binding creation event
// Contains only operational data needed by consumers (e.g., notification-service)
type TenantAppMessagingBindingInsertEventModel struct {
	ID                       string                    `bson:"_id" json:"id"`
	TenantID                 string                    `bson:"tenantID" json:"tenantID"`
	Channel                  string                    `bson:"channel" json:"channel"`
	ProviderID               string                    `bson:"providerID" json:"providerID"`
	CachedProviderURL        string                    `bson:"cachedProviderURL" json:"cachedProviderURL"`
	CachedProviderAPIVersion string                    `bson:"cachedProviderAPIVersion" json:"cachedProviderAPIVersion"`
	CachedProviderName       string                    `bson:"cachedProviderName" json:"cachedProviderName"`
	TenantConfig             *TenantAppMessagingConfig `bson:"tenantConfig" json:"tenantConfig"`
	TemplateBindings         []TemplateBinding         `bson:"templateBindings" json:"templateBindings"`
	Version                  int                       `bson:"version" json:"version"`
}

// TenantAppMessagingBindingUpdateEventModel represents an App Messaging provider binding update event
// Contains only operational data needed by consumers (e.g., notification-service)
type TenantAppMessagingBindingUpdateEventModel struct {
	ID                       string                    `bson:"_id" json:"id"`
	TenantID                 string                    `bson:"tenantID" json:"tenantID"`
	Channel                  string                    `bson:"channel" json:"channel"`
	ProviderID               string                    `bson:"providerID,omitempty" json:"providerID,omitempty"`
	CachedProviderURL        string                    `bson:"cachedProviderURL,omitempty" json:"cachedProviderURL,omitempty"`
	CachedProviderAPIVersion string                    `bson:"cachedProviderAPIVersion,omitempty" json:"cachedProviderAPIVersion,omitempty"`
	CachedProviderName       string                    `bson:"cachedProviderName,omitempty" json:"cachedProviderName,omitempty"`
	TenantConfig             *TenantAppMessagingConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	TemplateBindings         []TemplateBinding         `bson:"templateBindings,omitempty" json:"templateBindings,omitempty"`
	Version                  int                       `bson:"version" json:"version"`
}

// TenantAppMessagingBindingDeleteEventModel represents an App Messaging provider binding deletion event
type TenantAppMessagingBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantID" json:"tenantID"`
}
