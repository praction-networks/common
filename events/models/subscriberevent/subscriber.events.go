package subscriberevent

import "time"

// SubscriberCreatedEvent represents a subscriber creation event
type SubscriberCreatedEvent struct {
	ID                     string                 `json:"id" bson:"id"`
	TenantID               string                 `json:"tenantId" bson:"tenantId"`
	DomainID               string                 `json:"domainId,omitempty" bson:"domainId,omitempty"`
	Type                   SubscriberType         `json:"type" bson:"type"`
	Status                 SubscriberStatus       `json:"status" bson:"status"`
	ExternalRef            string                 `json:"externalRef,omitempty" bson:"externalRef,omitempty"`
	FullName               string                 `json:"fullName,omitempty" bson:"fullName,omitempty"`
	OrganizationName       string                 `json:"organizationName,omitempty" bson:"organizationName,omitempty"`
	PrimaryMobile          string                 `json:"primaryMobile" bson:"primaryMobile"`
	PrimaryEmail           string                 `json:"primaryEmail,omitempty" bson:"primaryEmail,omitempty"`
	WhatsAppNumber         string                 `json:"whatsAppNumber,omitempty" bson:"whatsAppNumber,omitempty"`
	PrimaryUserID          string                 `json:"primaryUserId,omitempty" bson:"primaryUserId,omitempty"`
	MaxMacAddresses        int                    `json:"maxMacAddresses,omitempty" bson:"maxMacAddresses,omitempty"`
	MaxSimultaneousUse     int                    `json:"maxSimultaneousUse,omitempty" bson:"maxSimultaneousUse,omitempty"`
	CustomerSpecificFields map[string]interface{} `json:"customerSpecificFields,omitempty" bson:"customerSpecificFields,omitempty"`
	Version                int                    `json:"version" bson:"version"`
}

// SubscriberUpdatedEvent represents a subscriber update event
type SubscriberUpdatedEvent struct {
	ID                     string                 `json:"id" bson:"id"`
	TenantID               string                 `json:"tenantId" bson:"tenantId"`
	Type                   SubscriberType         `json:"type,omitempty" bson:"type,omitempty"`
	Status                 SubscriberStatus       `json:"status,omitempty" bson:"status,omitempty"`
	FullName               string                 `json:"fullName,omitempty" bson:"fullName,omitempty"`
	OrganizationName       string                 `json:"organizationName,omitempty" bson:"organizationName,omitempty"`
	PrimaryMobile          string                 `json:"primaryMobile,omitempty" bson:"primaryMobile,omitempty"`
	PrimaryEmail           string                 `json:"primaryEmail,omitempty" bson:"primaryEmail,omitempty"`
	WhatsAppNumber         string                 `json:"whatsAppNumber,omitempty" bson:"whatsAppNumber,omitempty"`
	MaxMacAddresses        int                    `json:"maxMacAddresses,omitempty" bson:"maxMacAddresses,omitempty"`
	MaxSimultaneousUse     int                    `json:"maxSimultaneousUse,omitempty" bson:"maxSimultaneousUse,omitempty"`
	CustomerSpecificFields map[string]interface{} `json:"customerSpecificFields,omitempty" bson:"customerSpecificFields,omitempty"`
	FirstHotspotLoginAt    *time.Time             `json:"firstHotspotLoginAt,omitempty" bson:"firstHotspotLoginAt,omitempty"` // When subscriber first logged into hotspot
	LastHotspotLoginAt     *time.Time             `json:"lastHotspotLoginAt,omitempty" bson:"lastHotspotLoginAt,omitempty"`   // When subscriber last logged into hotspot
	Version                int                    `json:"version" bson:"version"`
}

// SubscriberDeletedEvent represents a subscriber deletion event
type SubscriberDeletedEvent struct {
	ID        string    `json:"id" bson:"id"`
	TenantID  string    `json:"tenantId" bson:"tenantId"`
	DeletedAt time.Time `json:"deletedAt" bson:"deletedAt"`
	DeletedBy string    `json:"deletedBy,omitempty" bson:"deletedBy,omitempty"`
}

// BroadbandSubscriptionCreatedEvent represents a broadband subscription creation event
type BroadbandSubscriptionCreatedEvent struct {
	ID                    string `json:"id" bson:"id"`
	SubscriberID          string `json:"subscriberId" bson:"subscriberId"`
	TenantID              string `json:"tenantId" bson:"tenantId"`
	AccountNumber         string `json:"accountNumber" bson:"accountNumber"`
	Status                string `json:"status" bson:"status"`
	AccessType            string `json:"accessType" bson:"accessType"`
	PlanID                string `json:"planId" bson:"planId"`
	PlanName              string `json:"planName,omitempty" bson:"planName,omitempty"`
	SpeedDownMbps         int    `json:"speedDownMbps" bson:"speedDownMbps"`
	SpeedUpMbps           int    `json:"speedUpMbps" bson:"speedUpMbps"`
	InstallationAddressID string `json:"installationAddressId,omitempty" bson:"installationAddressId,omitempty"`
	BillingProfileID      string `json:"billingProfileId,omitempty" bson:"billingProfileId,omitempty"`
	PPPoEUsername         string `json:"pppoeUsername,omitempty" bson:"pppoeUsername,omitempty"`
	OLTDeviceID           string `json:"oltDeviceId,omitempty" bson:"oltDeviceId,omitempty"`
	OLTPort               string `json:"oltPort,omitempty" bson:"oltPort,omitempty"`
}

// BroadbandSubscriptionUpdatedEvent represents a broadband subscription update event
type BroadbandSubscriptionUpdatedEvent struct {
	ID                    string `json:"id" bson:"id"`
	SubscriberID          string `json:"subscriberId" bson:"subscriberId"`
	TenantID              string `json:"tenantId" bson:"tenantId"`
	AccountNumber         string `json:"accountNumber,omitempty" bson:"accountNumber,omitempty"`
	Status                string `json:"status,omitempty" bson:"status,omitempty"`
	PlanID                string `json:"planId,omitempty" bson:"planId,omitempty"`
	SpeedDownMbps         int    `json:"speedDownMbps,omitempty" bson:"speedDownMbps,omitempty"`
	SpeedUpMbps           int    `json:"speedUpMbps,omitempty" bson:"speedUpMbps,omitempty"`
	InstallationAddressID string `json:"installationAddressId,omitempty" bson:"installationAddressId,omitempty"`
	Version               int    `json:"version" bson:"version"`
}

// BroadbandSubscriptionDeletedEvent represents a broadband subscription deletion event
type BroadbandSubscriptionDeletedEvent struct {
	ID            string    `json:"id" bson:"id"`
	SubscriberID  string    `json:"subscriberId" bson:"subscriberId"`
	TenantID      string    `json:"tenantId" bson:"tenantId"`
	AccountNumber string    `json:"accountNumber" bson:"accountNumber"`
}

// HotspotProfileCreatedEvent represents a hotspot profile creation event
type HotspotProfileCreatedEvent struct {
	ID                string     `json:"id" bson:"id"`
	SubscriberID      string     `json:"subscriberId" bson:"subscriberId"`
	MACAddress        string     `json:"macAddress,omitempty" bson:"macAddress,omitempty"`
	TenantID          string     `json:"tenantId" bson:"tenantId"`
	Status            HotspotProfileStatus     `json:"status" bson:"status"`
	UserID            string     `json:"userId,omitempty" bson:"userId,omitempty"`
	Username          string     `json:"username,omitempty" bson:"username,omitempty"` // RADIUS username
	Password          string     `json:"password,omitempty" bson:"password,omitempty"` // RADIUS password
	DefaultAuthMethod HotspotAuthMethod     `json:"defaultAuthMethod" bson:"defaultAuthMethod"`
	MaxDevices        int        `json:"maxDevices" bson:"maxDevices"`
	AutoLoginEnabled  bool       `json:"autoLoginEnabled" bson:"autoLoginEnabled"`
	CreationSource    string     `json:"creationSource,omitempty" bson:"creationSource,omitempty"` // "CAPTIVE_PORTAL", "ADMIN", "API", "IMPORT"
	ValidFrom         time.Time  `json:"validFrom,omitempty" bson:"validFrom,omitempty"`           // When this profile becomes valid
	ValidUntil        *time.Time `json:"validUntil,omitempty" bson:"validUntil,omitempty"`         // When this profile expires (nil = no expiration)
	PlanCode          string     `json:"planCode,omitempty" bson:"planCode,omitempty"`             // Plan code used for this hotspot profile
	Version           int        `json:"version" bson:"version"`
}

// HotspotProfileUpdatedEvent represents a hotspot profile update event
type HotspotProfileUpdatedEvent struct {
	ID                string     `json:"id" bson:"id"`
	SubscriberID      string     `json:"subscriberId" bson:"subscriberId"`
	MACAddress        string     `json:"macAddress,omitempty" bson:"macAddress,omitempty"`
	TenantID          string     `json:"tenantId" bson:"tenantId"`
	Status            HotspotProfileStatus     `json:"status,omitempty" bson:"status,omitempty"`
	UserID            string     `json:"userId,omitempty" bson:"userId,omitempty"`
	Username          string     `json:"username,omitempty" bson:"username,omitempty"` // RADIUS username
	Password          string     `json:"password,omitempty" bson:"password,omitempty"` // RADIUS password
	DefaultAuthMethod HotspotAuthMethod     `json:"defaultAuthMethod,omitempty" bson:"defaultAuthMethod,omitempty"`
	MaxDevices        int        `json:"maxDevices,omitempty" bson:"maxDevices,omitempty"`
	FirstLoginAt      *time.Time `json:"firstLoginAt,omitempty" bson:"firstLoginAt,omitempty"` // When was the first successful login
	ValidFrom         *time.Time `json:"validFrom,omitempty" bson:"validFrom,omitempty"`       // When this profile becomes valid (optional update)
	ValidUntil        *time.Time `json:"validUntil,omitempty" bson:"validUntil,omitempty"`     // When this profile expires (optional update)
	PlanCode          string     `json:"planCode,omitempty" bson:"planCode,omitempty"`         // Plan code used for this hotspot profile (optional update)
	Version           int        `json:"version" bson:"version"`
}

// HotspotProfileDeletedEvent represents a hotspot profile deletion event
type HotspotProfileDeletedEvent struct {
	ID           string `json:"id" bson:"id"`
	SubscriberID string `json:"subscriberId" bson:"subscriberId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`
	Version      int    `json:"version" bson:"version"`
}