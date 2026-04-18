package tenantevent

// SMSType represents supported SMS gateway providers
type SMSType string

// Unified SMS Provider events matching the refactored tenant-service architecture

type TenantSMSProviderInsertEventModel struct {
	ID                string                   `bson:"_id" json:"id"`
	OwnerTenantID     string                   `bson:"ownerTenantId" json:"ownerTenantId"`
	Channel           string                   `bson:"channel" json:"channel"`
	Scope             string                   `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                 `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	TenantConfig      *TenantSMSProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority          int                      `bson:"priority" json:"priority"`
	IsActive          bool                     `bson:"isActive" json:"isActive"`
	FailoverOn        bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries        int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            int                      `bson:"weight,omitempty" json:"weight,omitempty"`
	Version           int                      `bson:"version" json:"version"`
}

type TenantSMSProviderUpdateEventModel struct {
	ID                string                   `bson:"_id" json:"id"`
	OwnerTenantID     string                   `bson:"ownerTenantId" json:"ownerTenantId"`
	Channel           string                   `bson:"channel" json:"channel"`
	Scope             string                   `bson:"scope,omitempty" json:"scope,omitempty"`
	ExplicitTenantIDs []string                 `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	TenantConfig      *TenantSMSProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority          *int                     `bson:"priority,omitempty" json:"priority,omitempty"`
	IsActive          *bool                    `bson:"isActive,omitempty" json:"isActive,omitempty"`
	FailoverOn        *bool                    `bson:"failoverOn,omitempty" json:"failoverOn,omitempty"`
	MaxRetries        *int                     `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            *int                     `bson:"weight,omitempty" json:"weight,omitempty"`
	Version           int                      `bson:"version" json:"version"`
}

type TenantSMSProviderDeleteEventModel struct {
	ID            string `bson:"_id" json:"id"`
	OwnerTenantID string `bson:"ownerTenantId" json:"ownerTenantId"`
}
