package tenantevent

// ==================== ESIGN PROVIDER BINDING EVENTS ====================

// TenantESignProviderConfig contains the ESign provider credentials
type TenantESignProviderConfig struct {
	Provider string         `bson:"provider" json:"provider"`
	Metadata map[string]any `bson:"metadata" json:"metadata"`
}

// TenantESignProviderBindingInsertEventModel represents an ESign provider binding creation event
type TenantESignProviderBindingInsertEventModel struct {
	ID                string                     `bson:"_id" json:"id"`
	TenantID          string                     `bson:"tenantId" json:"tenantId"`
	ProviderType      string                     `bson:"providerType" json:"providerType"`
	Scope             string                     `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                   `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	ProviderID        string                     `bson:"providerId,omitempty" json:"providerId,omitempty"`
	UseTemplate       bool                       `bson:"useTemplate" json:"useTemplate"`
	UseParent         bool                       `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID    string                     `bson:"parentTenantId,omitempty" json:"parentTenantId,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig      *TenantESignProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority          int                        `bson:"priority" json:"priority"`
	IsActive          bool                       `bson:"isActive" json:"isActive"`
	FailoverOn        bool                       `bson:"failoverOn" json:"failoverOn"`
	MaxRetries        int                        `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            int                        `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider  string                     `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version           int                        `bson:"version" json:"version"`
}

// TenantESignProviderBindingUpdateEventModel represents an ESign provider binding update event
type TenantESignProviderBindingUpdateEventModel struct {
	ID                string                     `bson:"_id" json:"id"`
	TenantID          string                     `bson:"tenantId" json:"tenantId"`
	ProviderType      string                     `bson:"providerType" json:"providerType"`
	Scope             string                     `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                   `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	ProviderID        string                     `bson:"providerId,omitempty" json:"providerId,omitempty"`
	UseTemplate       bool                       `bson:"useTemplate" json:"useTemplate"`
	UseParent         bool                       `bson:"useParent,omitempty" json:"useParent,omitempty"`
	ParentTenantID    string                     `bson:"parentTenantId,omitempty" json:"parentTenantId,omitempty"` // Explicit parent tenant ID when useParent is true
	TenantConfig      *TenantESignProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority          int                        `bson:"priority" json:"priority"`
	IsActive          bool                       `bson:"isActive" json:"isActive"`
	FailoverOn        bool                       `bson:"failoverOn" json:"failoverOn"`
	MaxRetries        int                        `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            int                        `bson:"weight,omitempty" json:"weight,omitempty"`
	ResolvedProvider  string                     `bson:"resolvedProvider,omitempty" json:"resolvedProvider,omitempty"`
	Version           int                        `bson:"version" json:"version"`
}

// TenantESignProviderBindingDeleteEventModel represents an ESign provider binding deletion event
type TenantESignProviderBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`
}
