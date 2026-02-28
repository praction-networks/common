package tenantevent

// ==================== STORAGE PROVIDER BINDING EVENTS ====================

// TenantStorageProviderBindingInsertEventModel represents a Storage provider binding creation event
type TenantStorageProviderBindingInsertEventModel struct {
	ID                string                                     `bson:"_id" json:"id"`
	TenantID          string                                     `bson:"tenantId" json:"tenantId"`
	Scope             string                                     `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                                   `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	Configs           map[string]*TenantStorageProviderConfigEvent `bson:"configs" json:"configs"`
	Version           int                                        `bson:"version" json:"version"`
}

// TenantStorageProviderBindingUpdateEventModel represents a Storage provider binding update event
type TenantStorageProviderBindingUpdateEventModel struct {
	ID                string                                     `bson:"_id" json:"id"`
	TenantID          string                                     `bson:"tenantId" json:"tenantId"`
	Scope             string                                     `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                                   `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	Configs           map[string]*TenantStorageProviderConfigEvent `bson:"configs" json:"configs"`
	Version           int                                        `bson:"version" json:"version"`
}

// TenantStorageProviderBindingDeleteEventModel represents a Storage provider binding deletion event
type TenantStorageProviderBindingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`
}

// TenantStorageProviderConfigEvent represents a per-purpose storage provider config in events
type TenantStorageProviderConfigEvent struct {
	Provider      string         `bson:"provider" json:"provider"`
	Metadata      map[string]any `bson:"metadata" json:"metadata"`
	IsActive      bool           `bson:"isActive" json:"isActive"`
	MaxFileSizeMB int            `bson:"maxFileSizeMB,omitempty" json:"maxFileSizeMB,omitempty"`
}
