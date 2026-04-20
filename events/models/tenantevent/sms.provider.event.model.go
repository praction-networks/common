package tenantevent

// SMSType represents supported SMS gateway providers
type SMSType string

// Unified SMS Provider events matching the refactored tenant-service architecture.
// NATS JetStream handles at-transport dedup via Nats-Msg-Id; Version is the
// application-level decision key for consumers. No envelope is carried on the
// wire — every other *event.model.go in this package follows the same pattern.

type TenantSMSProviderInsertEventModel struct {
	ID                string                   `bson:"_id" json:"id"`
	OwnerTenantID     string                   `bson:"ownerTenantId" json:"ownerTenantId"`
	Channel           string                   `bson:"channel" json:"channel"`
	Scope             string                   `bson:"scope" json:"scope"`
	ExplicitTenantIDs []string                 `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	TemplateCodes     []string                 `bson:"templateCodes,omitempty" json:"templateCodes,omitempty"`
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
	TemplateCodes     []string                 `bson:"templateCodes,omitempty" json:"templateCodes,omitempty"`
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
	Version       int    `bson:"version,omitempty" json:"version,omitempty"`
}
