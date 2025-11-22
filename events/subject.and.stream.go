package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream      StreamName = "SeedAppStream"
	AuthStream         StreamName = "AuthStream"
	TenantStream       StreamName = "TenantStream"
	TenantUserStream   StreamName = "TenantUserStream"
	InventoryStream    StreamName = "InventoryStream"
	NotificationStream StreamName = "NotificationStream"
)

// Global Stream names as constants
const (
	NotificationGlobalStream StreamName = "NotificationGlobalStream"
	AuditGlobalStream        StreamName = "AuditGlobalStream"
	IntegrationGlobalStream  StreamName = "IntegrationGlobalStream"
)

// Subjects defines the NATS Subjects for different events
type Subject string

const (

	//Domain Service Event Initialization
	TenantCreatedSubject            Subject = "tenant.created"
	TenantUpdatedSubject            Subject = "tenant.updated"
	TenantDeletedSubject            Subject = "tenant.deleted"
	AppMessengerCreateSubject       Subject = "appmessenger.created"
	AppMessengerUpdateSubject       Subject = "appmessenger.updated"
	AppMessengerDeleteSubject       Subject = "appmessenger.deleted"
	KYCGatewayCreatedSubject        Subject = "kycgateway.created"
	KYCGatewayUpdateSubject         Subject = "kycgateway.updated"
	KYCGatewayDeleteSubject         Subject = "kycgateway.deleted"
	MailServerCreatedSubject        Subject = "mailserver.created"
	MailServerUpdateSubject         Subject = "mailserver.updated"
	MailServerDeleteSubject         Subject = "mailserver.deleted"
	PaymentGatewayCreatedSubject    Subject = "paymentgateway.created"
	PaymentGatewayUpdateSubject     Subject = "paymentgateway.updated"
	PaymentGatewayDeleteSubject     Subject = "paymentgateway.deleted"
	ExternalRadiusCreatedSubject    Subject = "externalradius.created"
	ExternalRadiusUpdateSubject     Subject = "externalradius.updated"
	ExternalRadiusDeleteSubject     Subject = "externalradius.deleted"
	SMSGatewayCreatedSubject        Subject = "smsgateway.created"
	SMSGatewayUpdateSubject         Subject = "smsgateway.updated"
	SMSGatewayDeleteSubject         Subject = "smsgateway.deleted"
	MessagingTemplateCreatedSubject Subject = "messagingtemplate.created"
	MessagingTemplateUpdateSubject  Subject = "messagingtemplate.updated"
	MessagingTemplateDeleteSubject  Subject = "messagingtemplate.deleted"

	//Domain User Service Event Initialization
	TenantUserCreatedSubject Subject = "tenantuser.created"
	TenantUserUpdatedSubject Subject = "tenantuser.updated"
	TenantUserDeletedSubject Subject = "tenantuser.deleted"

	// Tenant Auth Role Event Initialization
	TenantUserRoleCreatedSubject Subject = "tenantuserrole.created"
	TenantUserRoleUpdatedSubject Subject = "tenantuserrole.updated"
	TenantUserRoleDeletedSubject Subject = "tenantuserrole.deleted"

	// Template Action Code Event Initialization
	TemplateActionCodeCreatedSubject Subject = "templateactioncode.created"
	TemplateActionCodeUpdatedSubject Subject = "templateactioncode.updated"
	TemplateActionCodeDeletedSubject Subject = "templateactioncode.deleted"
)

// Global Subjects - Cross-service events that any service can publish
const (
	// OTP Events (authentication/verification)
	UserNotifcationSentSubject     Subject = "user.notification.sent"
	UserNotifcationVerifiedSubject Subject = "user.notification.verified"
	UserNotifcationExpiredSubject  Subject = "user.notification.expired"
	UserNotifcationFailedSubject   Subject = "user.notification.failed"

	// Push Notification Events
	UserPushNotificationSentSubject      Subject = "user.push.notification.sent"
	UserPushNotificationDeliveredSubject Subject = "user.push.notification.delivered"
	UserPushNotificationFailedSubject    Subject = "user.push.notification.failed"
	UserPushNotificationOpenedSubject    Subject = "user.push.notification.opened" // Optional: for tracking
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
			MessagingTemplateCreatedSubject,
			MessagingTemplateUpdateSubject,
			MessagingTemplateDeleteSubject,
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

			// Tenant User Role Events
			TenantUserRoleCreatedSubject,
			TenantUserRoleUpdatedSubject,
			TenantUserRoleDeletedSubject,
		},
	},

	// Global Streams - Cross-service events
	NotificationGlobalStream: {
		Name:        NotificationGlobalStream,
		Description: "Global stream for notification events from all services",
		Subjects: []Subject{
			UserNotifcationSentSubject,
			UserNotifcationVerifiedSubject,
			UserNotifcationExpiredSubject,
			UserNotifcationFailedSubject,
			UserPushNotificationSentSubject,
			UserPushNotificationDeliveredSubject,
			UserPushNotificationFailedSubject,
			UserPushNotificationOpenedSubject,
		},
	},

	NotificationStream: {
		Name:        NotificationStream,
		Description: "Stream for notification events from notification service",
		Subjects: []Subject{
			TemplateActionCodeCreatedSubject,
			TemplateActionCodeUpdatedSubject,
			TemplateActionCodeDeletedSubject,
		},
	},
}
