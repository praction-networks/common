package domainevent

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DomainInsertEventModel struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Code             string       `json:"code"`
	Type             string       `json:"type"`
	EntType          string       `json:"entType,omitempty"`
	GSTNumber        string       `json:"gstNumber,omitempty"`
	PANNumber        string       `json:"panNumber,omitempty"`
	IsActive         bool         `json:"isActive"`
	ParentID         string       `json:"parentID,omitempty"`
	ChildsIDs        []string     `json:"childsIDs,omitempty"`
	PermanentAddress AddressModel `json:"permanentAddress"`
	CurrentAddress   AddressModel `json:"currentAddress"`
	BillingAddress   AddressModel `json:"billingAddress,omitempty"`
	Version          int          `json:"version,omitempty"`
}

type DomainUpdateEventModel struct {
	ID                     string                  `json:"id"`
	Name                   string                  `json:"name"`
	Code                   string                  `json:"code"`
	Type                   string                  `json:"type"`
	Users                  []string                `json:"users,omitempty"`
	EntType                string                  `json:"entType,omitempty"`
	GSTNumber              string                  `json:"gstNumber,omitempty"`
	PANNumber              string                  `json:"panNumber,omitempty"`
	IsActive               bool                    `json:"isActive"`
	ParentID               string                  `json:"parentID,omitempty"`
	ChildsIDs              []string                `json:"childsIDs,omitempty"`
	PermanentAddress       AddressModel            `json:"permanentAddress"`
	CurrentAddress         AddressModel            `json:"currentAddress"`
	BillingAddress         AddressModel            `json:"billingAddress,omitempty"`
	OLTs                   []string                `json:"olts,omitempty"`
	PortalSettings         *PortalSettings         `json:"portalSettings,omitempty"`
	ISPSettings            *ISPSettings            `json:"ispSettings,omitempty"`
	NotificationGateway    *NotificationGateways   `json:"notificationGateway,omitempty"`
	OTPGateway             *NotificationGateways   `json:"otpGateway,omitempty"`
	KYCGateway             []string                `json:"kycGateway,omitempty"`
	PaymentGateway         []string                `json:"paymentGateway,omitempty"`
	ExternalRadiusSettings *ExternalRadiusSettings `json:"externalRadiusSettings,omitempty"`
	Version                int                     `json:"version,omitempty"`
}

type DomainDeleteEventModel struct {
	ID string `json:"id"`
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

type PortalSettings struct {
	UsersPortalEnabled         bool `json:"usersPortalEnabled"`
	UsersNotificationEnabled   bool `json:"usersNotificationEnabled"`
	UsersOTPEnabled            bool `json:"usersOTPEnabled"`
	UsersBillingEnabled        bool `json:"usersBillingEnabled"`
	UsersKYCEnabled            bool `json:"usersKYCEnabled"`
	UsersPaymentGatewayEnabled bool `json:"usersPaymentGatewayEnabled"`
	RadiusProviderEnabled      bool `json:"radiusProviderEnabled"`
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

type DomainUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DomainDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

type DepartmentCreate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DepartmentUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DepartmentDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}
