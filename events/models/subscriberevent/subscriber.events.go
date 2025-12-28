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
	SubscriberFields       map[string]interface{} `json:"subscriberFields,omitempty" bson:"subscriberFields,omitempty"` // All dynamic fields from subscriber model
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
	SubscriberFields       map[string]interface{} `json:"subscriberFields,omitempty" bson:"subscriberFields,omitempty"` // All dynamic fields from subscriber model
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
