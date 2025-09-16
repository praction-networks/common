package tenantevent

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TenantInsertEventModel struct {
	ID               string                         `json:"id"`
	Name             string                         `json:"name"`
	Code             string                         `json:"code"`
	Type             string                         `json:"type"`
	Fqdn             string                         `json:"fqdn,omitempty"`
	Environment      string                         `json:"environment,omitempty"`
	EntType          string                         `json:"entType,omitempty"`
	ParentTenantID   string                         `json:"parentTenantID,omitempty"`
	DefaultEmail     string                         `json:"defaultEmail"`
	DefaultPhone     string                         `json:"defaultPhone"`
	PermanentAddress AddressModel                   `json:"permanentAddress"`
	CurrentAddress   AddressModel                   `json:"currentAddress"`
	TenantGST        []GSTModel                     `json:"tenantGST"`
	TenantPAN        PANModel                       `json:"tenantPAN"`
	TenantTAN        TANModel                       `json:"tenantTAN"`
	TenantCIN        CINModel                       `json:"tenantCIN"`
	EnabledFeatures  EnabledFeatures                `json:"enabledFeatures,omitempty"`
	OLTs             []string                       `json:"olts,omitempty"`
	KYCProvider      ProvidersModel                 `json:"kycProvider,omitempty"`
	PaymentGateway   ProvidersModel                 `json:"paymentGateway,omitempty"`
	SMSProvider      ProvidersModel                 `json:"smsProvider,omitempty"`
	MailProvider     ProvidersModel                 `json:"mailProvider,omitempty"`
	AppsMessanger    []AppMessagingProvidersModel   `json:"appsMessanger,omitempty"`
	ExternalRadius   []ExternalRadiusProvidersModel `json:"externalRadius,omitempty"`
	IsActive         bool                           `json:"isActive"`
	Version          int                            `json:"version,omitempty"`
}

type TenantUpdateEventModel struct {
	ID               string                          `json:"id"`
	Name             *string                         `json:"name"`
	Code             *string                         `json:"code"`
	Type             *string                         `json:"type"`
	Fqdn             *string                         `json:"fqdn,omitempty"`
	Environment      *string                         `json:"environment,omitempty"`
	EntType          *string                         `json:"entType,omitempty"`
	ParentTenantID   *string                         `json:"parentTenantID,omitempty"`
	ChildIDs         *[]string                       `json:"childIDs,omitempty"`
	DefaultEmail     *string                         `json:"defaultEmail"`
	DefaultPhone     *string                         `json:"defaultPhone"`
	PermanentAddress *AddressModel                   `json:"permanentAddress"`
	CurrentAddress   *AddressModel                   `json:"currentAddress"`
	TenantGST        *[]GSTModel                     `json:"tenantGST"`
	TenantPAN        *PANModel                       `json:"tenantPAN"`
	TenantTAN        *TANModel                       `json:"tenantTAN"`
	TenantCIN        *CINModel                       `json:"tenantCIN"`
	EnabledFeatures  *EnabledFeatures                `json:"enabledFeatures,omitempty"`
	OLTs             *[]string                       `json:"olts,omitempty"`
	KYCProvider      *ProvidersModel                 `json:"kycProvider,omitempty"`
	PaymentGateway   *ProvidersModel                 `json:"paymentGateway,omitempty"`
	SMSProvider      *ProvidersModel                 `json:"smsProvider,omitempty"`
	MailProvider     *ProvidersModel                 `json:"mailProvider,omitempty"`
	AppsMessanger    *[]AppMessagingProvidersModel   `json:"appsMessanger,omitempty"`
	ExternalRadius   *[]ExternalRadiusProvidersModel `json:"externalRadius,omitempty"`
	IsActive         *bool                           `json:"isActive"`
	Version          int                             `json:"version,omitempty"`
}

type GSTModel struct {
	State      string    `json:"state,omitempty"`
	GSTIN      string    `json:"gstin,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type TenantDeleteEventModel struct {
	ID string `json:"id"`
}

type ProvidersModel struct {
	DefaultProviderID string   `json:"defaultProviderID,omitempty"`
	Providers         []string `json:"providers,omitempty"`
}

type AppMessagingProvidersModel struct {
	MessageProviderID string `json:"messageProviderID,omitempty"`
	MessageProvider   string `json:"messageProvider,omitempty"`
}

type ExternalRadiusProvidersModel struct {
	ExternalRadiusProviderID string   `json:"externalRadiusProviderID,omitempty"`
	ExternalRadiusProviders  []string `json:"externalRadiusProviders,omitempty"`
}

type PANModel struct {
	PAN        string    `json:"pan,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type TANModel struct {
	TAN        string    `json:"tan,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type CINModel struct {
	CIN        string    `json:"cin,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type NotificationGateways struct {
	SMS      []string `json:"sms,omitempty"`
	Mail     []string `json:"mail,omitempty"`
	WhatsApp []string `json:"whatsapp,omitempty"`
	Telegram []string `json:"telegram,omitempty"`
}

type ExternalRadiusSettings struct {
	JazeID  string `json:"jazeEnabled,omitempty"`
	IPACTID string `json:"ipactEnabled,omitempty"`
}

type EnabledFeatures struct {
	IsOnlinePaymentEnabled            bool `json:"isOnlinePaymentEnabled"`
	IsUserPortalEnabled               bool `json:"isUserPortalEnabled"`
	IsUserMailNotificationEnabled     bool `json:"isUserMailNotificationEnabled"`
	IsUserWhatsappNotificationEnabled bool `json:"isUserWhatsappNotificationEnabled"`
	IsUserTelegramNotificationEnabled bool `json:"isUserTelegramNotificationEnabled"`
	IsUserSMSNotificationEnabled      bool `json:"isUserSMSNotificationEnabled"`
	IsUserKYCEnabled                  bool `json:"isUserKYCEnabled"`
	IsJazeeraRadiusProviderEnabled    bool `json:"isJazeeraRadiusProviderEnabled"`
	IsIPACTRadiusProviderEnabled      bool `json:"isIPACTRadiusProviderEnabled"`
	IsFreeRadiusProviderEnabled       bool `json:"isFreeRadiusProviderEnabled"`
	IsIPTVEnabled                     bool `json:"isIPTVEnabled"`
	IsOTTEnabled                      bool `json:"isOTTEnabled"`
	IsVoiceServiceEnabled             bool `json:"isVoiceServiceEnabled"`
}

type ISPSettings struct {
	Plans            []string `json:"plans,omitempty" bson:"plans,omitempty"`
	MaxBandwidthMbps int      `json:"maxBandwidthMbps,omitempty" bson:"maxBandwidthMbps,omitempty"`
	IPPoolCIDR       string   `json:"ipPoolCIDR,omitempty" bson:"ipPoolCIDR,omitempty"`
	CoverageArea     bson.M   `json:"coverageArea,omitempty" bson:"coverageArea,omitempty"`
	Latitude         string   `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude        string   `json:"longitude,omitempty" bson:"longitude,omitempty"`
	BillingCycle     string   `json:"billingCycle,omitempty" bson:"billingCycle,omitempty"`
	AutoRenewal      bool     `json:"autoRenewal,omitempty" bson:"autoRenewal,omitempty"`
	SupportContact   string   `json:"supportContact,omitempty" bson:"supportContact,omitempty"`
	AssignedRM       string   `json:"assignedRM,omitempty" bson:"assignedRM,omitempty"`
	DeviceIDs        []string `json:"deviceIDs,omitempty" bson:"deviceIDs,omitempty"`
}

type TenantUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type TenantDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}
