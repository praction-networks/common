package tenantevent

// CDN provider event models. Previously co-located in the (now deleted)
// mail.server.event.model.go file — restored here in a properly named file
// so the tenant-service CDN publisher/listener can continue to publish
// insert/update/delete events for the CDNProviderModel collection.

// CDNProviderInsertEventModel defines the model for CDN provider insert events
type CDNProviderInsertEventModel struct {
	ID                 string   `bson:"_id" json:"id"`
	Name               string   `bson:"name" json:"name"`
	SortCode           string   `bson:"sortCode" json:"sortCode"`
	OwnerTenantID      string   `bson:"ownerTenantId" json:"ownerTenantId"`
	OwnerTenantType    string   `bson:"ownerTenantType" json:"ownerTenantType"`
	Scope              string   `bson:"scope" json:"scope"`
	AllowedTenantTypes []string `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	IsActive           bool     `bson:"isActive" json:"isActive"`
	Version            int      `bson:"version" json:"version"`
}

// CDNProviderUpdateEventModel defines the model for CDN provider update events
type CDNProviderUpdateEventModel struct {
	ID                 string   `bson:"_id" json:"id"`
	Name               string   `bson:"name,omitempty" json:"name,omitempty"`
	SortCode           string   `bson:"sortCode,omitempty" json:"sortCode,omitempty"`
	OwnerTenantID      string   `bson:"ownerTenantId,omitempty" json:"ownerTenantId,omitempty"`
	OwnerTenantType    string   `bson:"ownerTenantType,omitempty" json:"ownerTenantType,omitempty"`
	Scope              string   `bson:"scope,omitempty" json:"scope,omitempty"`
	AllowedTenantTypes []string `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	IsActive           *bool    `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Version            *int     `bson:"version,omitempty" json:"version,omitempty"`
}

// CDNProviderDeleteEventModel defines the model for CDN provider delete events
type CDNProviderDeleteEventModel struct {
	ID      string `bson:"_id" json:"id"`
	Version int    `bson:"version" json:"version"`
}
