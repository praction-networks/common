package tenantevent

type KYCProvideInsertEventModel struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name" json:"name"`
	Provider           string         `bson:"provider" json:"provider"`
	OwnerTenantID      string         `bson:"ownerTenantID" json:"ownerTenantID"`
	OwnerTenantType    string         `bson:"ownerTenantType" json:"ownerTenantType"`
	Scope              string         `bson:"scope" json:"scope"`
	AllowedTenantTypes []string       `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string       `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty"`
	AssignedTo         []string       `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"` // Deprecated: Use ExplicitTenantIDs when Scope=ExplicitTenants
	Metadata           map[string]any `bson:"metadata,omitempty" json:"metadata"`
	IsActive           bool           `bson:"isActive" json:"isActive"`
	Version            int            `bson:"version" json:"version"`
}

type KYCProvideUpdateEventModel struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name" json:"name"`
	Provider           string         `bson:"provider" json:"provider"`
	OwnerTenantID      string         `bson:"ownerTenantID" json:"ownerTenantID"`
	OwnerTenantType    string         `bson:"ownerTenantType" json:"ownerTenantType"`
	Scope              string         `bson:"scope" json:"scope"`
	AllowedTenantTypes []string       `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string       `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty"`
	AssignedTo         []string       `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"` // Deprecated: Use ExplicitTenantIDs when Scope=ExplicitTenants
	Metadata           map[string]any `bson:"metadata,omitempty" json:"metadata"`
	IsActive           bool           `bson:"isActive" json:"isActive"`
	Version            int            `bson:"version" json:"version"`
}

type KYCProvideDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
