package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream    StreamName = "SeedAppStream"
	AuthStream       StreamName = "AuthStream"
	TenantStream     StreamName = "TenantStream"
	TenantUserStream StreamName = "TenantUserStream"
	InventoryStream  StreamName = "InventoryStream"
)

// Subjects defines the NATS Subjects for different events
type Subject string

const (

	//Domain Service Event Initialization
	TenantCreatedSubject         Subject = "tenant.created"
	TenantUpdatedSubject         Subject = "tenant.updated"
	TenantDeletedSubject         Subject = "tenant.deleted"
	AppMessengerCreateSubject    Subject = "appmessenger.created"
	AppMessengerUpdateSubject    Subject = "appmessenger.updated"
	AppMessengerDeleteSubject    Subject = "appmessenger.deleted"
	KYCGatewayCreatedSubject     Subject = "kycgateway.created"
	KYCGatewayUpdateSubject      Subject = "kycgateway.updated"
	KYCGatewayDeleteSubject      Subject = "kycgateway.deleted"
	MailServerCreatedSubject     Subject = "mailserver.created"
	MailServerUpdateSubject      Subject = "mailserver.updated"
	MailServerDeleteSubject      Subject = "mailserver.deleted"
	PaymentGatewayCreatedSubject Subject = "paymentgateway.created"
	PaymentGatewayUpdateSubject  Subject = "paymentgateway.updated"
	PaymentGatewayDeleteSubject  Subject = "paymentgateway.deleted"
	ExternalRadiusCreatedSubject Subject = "externalradius.created"
	ExternalRadiusUpdateSubject  Subject = "externalradius.updated"
	ExternalRadiusDeleteSubject  Subject = "externalradius.deleted"
	SMSGatewayCreatedSubject     Subject = "smsgateway.created"
	SMSGatewayUpdateSubject      Subject = "smsgateway.updated"
	SMSGatewayDeleteSubject      Subject = "smsgateway.deleted"

	//Domain User Service Event Initialization
	TenantUserCreatedSubject Subject = "tenantuser.created"
	TenantUserUpdatedSubject Subject = "tenantuser.updated"
	TenantUserDeletedSubject Subject = "tenantuser.deleted"

	// Tenant Auth Role Event Initialization
	TenantUserRoleCreatedSubject Subject = "tenantuserrole.created"
	TenantUserRoleUpdatedSubject Subject = "tenantuserrole.updated"
	TenantUserRoleDeletedSubject Subject = "tenantuserrole.deleted"

	// Auth Service Event Initialization
	AuthUserCreatedSubject Subject = "auth.user.created"
	AuthUserUpdatedSubject Subject = "auth.user.updated"
	AuthUserDeletedSubject Subject = "auth.user.deleted"
)

// StreamMetadata defines metadata for streams
type StreamMetadata struct {
	Name        StreamName
	Description string
	Subjects    []Subject
}

// Predefined stream configurations
var Streams = map[StreamName]StreamMetadata{

	TenantStream: {
		Name:        TenantStream,
		Description: "Stream for domain-related events",
		Subjects: []Subject{TenantCreatedSubject,
			TenantUpdatedSubject,
			TenantDeletedSubject,
			AppMessengerCreateSubject,
			AppMessengerUpdateSubject,
			AppMessengerDeleteSubject,
			KYCGatewayCreatedSubject,
			KYCGatewayUpdateSubject,
			KYCGatewayDeleteSubject,
			MailServerCreatedSubject,
			MailServerUpdateSubject,
			MailServerDeleteSubject,
			PaymentGatewayCreatedSubject,
			PaymentGatewayUpdateSubject,
			PaymentGatewayDeleteSubject,
			ExternalRadiusCreatedSubject,
			ExternalRadiusUpdateSubject,
			ExternalRadiusDeleteSubject,
			SMSGatewayCreatedSubject,
			SMSGatewayUpdateSubject,
			SMSGatewayDeleteSubject,
		},
	},

	TenantUserStream: {
		Name:        TenantUserStream,
		Description: "Stream for tenant user service events",
		Subjects: []Subject{
			// Tenant User Events (Extended Profiles)
			TenantUserCreatedSubject,
			TenantUserUpdatedSubject,
			TenantUserDeletedSubject,
		},
	},
	AuthStream: {
		Name:        AuthStream,
		Description: "Stream for auth service events",
		Subjects: []Subject{
			AuthUserCreatedSubject,
			AuthUserUpdatedSubject,
			AuthUserDeletedSubject,

			// Tenant User Role Events
			TenantUserRoleCreatedSubject,
			TenantUserRoleUpdatedSubject,
			TenantUserRoleDeletedSubject,
		},
	},
}
