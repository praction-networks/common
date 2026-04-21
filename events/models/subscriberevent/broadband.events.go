package subscriberevent

// BroadbandSubscriptionCreatedEvent represents a broadband subscription creation event
type BroadbandSubscriptionCreatedEvent struct {
	ID            string                      `json:"id" bson:"id"`
	SubscriberID  string                      `json:"subscriberId" bson:"subscriberId"`
	TenantID      string                      `json:"tenantId" bson:"tenantId"`
	AccountNumber string                      `json:"accountNumber" bson:"accountNumber"`
	Status        BroadbandSubscriptionStatus `json:"status" bson:"status"`
	// ConnectionType reflects the model's `accessType` field. Kept as
	// `connectionType` on the wire for downstream compatibility.
	ConnectionType   BroadbandConnectionType `json:"connectionType" bson:"accessType"`
	AuthMethod       BroabandAuthMethod      `json:"authMethod" bson:"authMethod"`
	AuthConfig       *BroadbandAuthConfig    `json:"authConfig,omitempty" bson:"authConfig,omitempty"`
	StaticIPAddress  string                  `json:"staticIpAddress,omitempty" bson:"staticIpAddress,omitempty"`
	StaticGateway    string                  `json:"staticGateway,omitempty" bson:"staticGateway,omitempty"`
	VLANID           int                     `json:"vlanId,omitempty" bson:"vlanId,omitempty"`
	RadiusAttributes map[string]interface{}  `json:"radiusAttributes,omitempty" bson:"radiusAttributes,omitempty"`
	// PlanID is the mongo document id of the plan; hydrated straight from bson.
	PlanID string `json:"planId" bson:"planId"`
	// PlanCode is resolved by the publisher via a plan lookup (bson:"-" so it
	// does not accidentally read a non-existent field on the subscription doc).
	PlanCode              string                        `json:"planCode" bson:"-"`
	InstallationAddress   *BroadbandInstallationAddress `json:"installationAddress,omitempty" bson:"installationAddress,omitempty"`
	InstallationAddressID string                        `json:"installationAddressId,omitempty" bson:"installationAddressId,omitempty"` // Deprecated: use InstallationAddress
	BillingProfileID      string                        `json:"billingProfileId,omitempty" bson:"billingProfileId,omitempty"`
	OLTDeviceID           string                        `json:"oltDeviceId,omitempty" bson:"oltDeviceId,omitempty"`
	OLTPort               string                        `json:"oltPort,omitempty" bson:"oltPort,omitempty"`
	CPEMetadata           BroadbandCPEMetadata          `json:"cpeMetadata,omitempty" bson:"cpeMetadata,omitempty"`
	HistoryCPEMetadata    []BroadbandCPEMetadata        `json:"historyCpeMetadata,omitempty" bson:"historyCpeMetadata,omitempty"`
	Notes                 string                        `json:"notes,omitempty" bson:"notes,omitempty"`
	Version               int                           `json:"version" bson:"version"`
}

// BroadbandSubscriptionUpdatedEvent represents a broadband subscription update event
type BroadbandSubscriptionUpdatedEvent struct {
	ID               string                  `json:"id" bson:"id"`
	SubscriberID     string                  `json:"subscriberId" bson:"subscriberId"`
	TenantID         string                  `json:"tenantId" bson:"tenantId"`
	AccountNumber    string                  `json:"accountNumber,omitempty" bson:"accountNumber,omitempty"`
	Status           string                  `json:"status,omitempty" bson:"status,omitempty"`
	ConnectionType   BroadbandConnectionType `json:"connectionType,omitempty" bson:"accessType,omitempty"`
	AuthMethod       BroabandAuthMethod      `json:"authMethod,omitempty" bson:"authMethod,omitempty"`
	AuthConfig       *BroadbandAuthConfig    `json:"authConfig,omitempty" bson:"authConfig,omitempty"`
	StaticIPAddress  string                  `json:"staticIpAddress,omitempty" bson:"staticIpAddress,omitempty"`
	VLANID           int                     `json:"vlanId,omitempty" bson:"vlanId,omitempty"`
	RadiusAttributes map[string]interface{}  `json:"radiusAttributes,omitempty" bson:"radiusAttributes,omitempty"`
	PlanID           string                  `json:"planId,omitempty" bson:"planId,omitempty"`
	PlanCode         string                  `json:"planCode,omitempty" bson:"-"`
	InstallationAddress   *BroadbandInstallationAddress `json:"installationAddress,omitempty" bson:"installationAddress,omitempty"`
	InstallationAddressID string                        `json:"installationAddressId,omitempty" bson:"installationAddressId,omitempty"` // Deprecated: use InstallationAddress
	Version               int                           `json:"version" bson:"version"`
}

// BroadbandSubscriptionDeletedEvent represents a broadband subscription deletion event
type BroadbandSubscriptionDeletedEvent struct {
	ID            string `json:"id" bson:"id"`
	SubscriberID  string `json:"subscriberId" bson:"subscriberId"`
	TenantID      string `json:"tenantId" bson:"tenantId"`
	AccountNumber string `json:"accountNumber" bson:"accountNumber"`
}
