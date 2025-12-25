package subscriberevent

// BroadbandSubscriptionCreatedEvent represents a broadband subscription creation event
type BroadbandSubscriptionCreatedEvent struct {
	ID                    string                      `json:"id" bson:"id"`
	SubscriberID          string                      `json:"subscriberId" bson:"subscriberId"`
	TenantID              string                      `json:"tenantId" bson:"tenantId"`
	AccountNumber         string                      `json:"accountNumber" bson:"accountNumber"`
	Status                BroadbandSubscriptionStatus `json:"status" bson:"status"`
	ConnectionType        BroadbandConnectionType     `json:"connectionType" bson:"connectionType"`
	AuthMethod            BroabandAuthMethod          `json:"authMethod" bson:"authMethod"`
	AuthConfig            *BroadbandAuthConfig        `json:"authConfig,omitempty" bson:"authConfig,omitempty"` // Authentication configuration based on AuthMethod
	StaticIPAddress       string                      `json:"staticIpAddress,omitempty" bson:"staticIpAddress,omitempty"`
	StaticGateway         string                      `json:"staticGateway,omitempty" bson:"staticGateway,omitempty"`
	VLANID                int                         `json:"vlanId,omitempty" bson:"vlanId,omitempty"`
	PlanCode              string                      `json:"planCode" bson:"planCode"`
	InstallationAddressID string                      `json:"installationAddressId,omitempty" bson:"installationAddressId,omitempty"`
	BillingProfileID      string                      `json:"billingProfileId,omitempty" bson:"billingProfileId,omitempty"`
	OLTDeviceID           string                      `json:"oltDeviceId,omitempty" bson:"oltDeviceId,omitempty"`
	OLTPort               string                      `json:"oltPort,omitempty" bson:"oltPort,omitempty"`
	CPEMetadata           BroadbandCPEMetadata        `json:"cpeMetadata,omitempty" bson:"cpeMetadata,omitempty"`
	HistoryCPEMetadata    []BroadbandCPEMetadata      `json:"historyCpeMetadata,omitempty" bson:"historyCpeMetadata,omitempty"`
	Notes                 string                      `json:"notes,omitempty" bson:"notes,omitempty"`
	Version               int                         `json:"version" bson:"version"`
}

// BroadbandSubscriptionUpdatedEvent represents a broadband subscription update event
type BroadbandSubscriptionUpdatedEvent struct {
	ID                    string `json:"id" bson:"id"`
	SubscriberID          string `json:"subscriberId" bson:"subscriberId"`
	TenantID              string `json:"tenantId" bson:"tenantId"`
	AccountNumber         string `json:"accountNumber,omitempty" bson:"accountNumber,omitempty"`
	Status                string `json:"status,omitempty" bson:"status,omitempty"`
	PlanCode              string `json:"planCode,omitempty" bson:"planCode,omitempty"`
	InstallationAddressID string `json:"installationAddressId,omitempty" bson:"installationAddressId,omitempty"`
	Version               int    `json:"version" bson:"version"`
}

// BroadbandSubscriptionDeletedEvent represents a broadband subscription deletion event
type BroadbandSubscriptionDeletedEvent struct {
	ID            string `json:"id" bson:"id"`
	SubscriberID  string `json:"subscriberId" bson:"subscriberId"`
	TenantID      string `json:"tenantId" bson:"tenantId"`
	AccountNumber string `json:"accountNumber" bson:"accountNumber"`
}
