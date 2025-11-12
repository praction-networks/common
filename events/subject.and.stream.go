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

// Global Stream names as constants
const (
	NotificationStream StreamName = "NotificationStream"
	AuditStream        StreamName = "AuditStream"
	IntegrationStream  StreamName = "IntegrationStream"
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
)

// Global Subjects - Cross-service events that any service can publish
const (
	// OTP Events (authentication/verification)
	OTPSentSubject     Subject = "otp.sent"
	OTPVerifiedSubject Subject = "otp.verified"
	OTPExpiredSubject  Subject = "otp.expired"
	OTPFailedSubject   Subject = "otp.failed"

	// SMS Notification Events
	SMSSentSubject      Subject = "sms.sent"
	SMSDeliveredSubject Subject = "sms.delivered"
	SMSFailedSubject    Subject = "sms.failed"
	SMSReadSubject      Subject = "sms.read" // Optional: for read receipts

	// Email Notification Events
	EmailSentSubject      Subject = "email.sent"
	EmailDeliveredSubject Subject = "email.delivered"
	EmailFailedSubject    Subject = "email.failed"
	EmailBouncedSubject   Subject = "email.bounced"
	EmailOpenedSubject    Subject = "email.opened"  // Optional: for tracking
	EmailClickedSubject   Subject = "email.clicked" // Optional: for tracking

	// WhatsApp Notification Events
	WhatsAppSentSubject      Subject = "whatsapp.sent"
	WhatsAppDeliveredSubject Subject = "whatsapp.delivered"
	WhatsAppFailedSubject    Subject = "whatsapp.failed"
	WhatsAppReadSubject      Subject = "whatsapp.read"

	// Telegram Notification Events
	TelegramSentSubject      Subject = "telegram.sent"
	TelegramDeliveredSubject Subject = "telegram.delivered"
	TelegramFailedSubject    Subject = "telegram.failed"

	// Push Notification Events
	PushSentSubject      Subject = "push.sent"
	PushDeliveredSubject Subject = "push.delivered"
	PushFailedSubject    Subject = "push.failed"
	PushOpenedSubject    Subject = "push.opened" // Optional: for tracking
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
	NotificationStream: {
		Name:        NotificationStream,
		Description: "Global stream for notification events from all services",
		Subjects: []Subject{
			// OTP Events
			OTPSentSubject,
			OTPVerifiedSubject,
			OTPExpiredSubject,
			OTPFailedSubject,

			// SMS Notification Events
			SMSSentSubject,
			SMSDeliveredSubject,
			SMSFailedSubject,
			SMSReadSubject,

			// Email Notification Events
			EmailSentSubject,
			EmailDeliveredSubject,
			EmailFailedSubject,
			EmailBouncedSubject,
			EmailOpenedSubject,
			EmailClickedSubject,

			// WhatsApp Notification Events
			WhatsAppSentSubject,
			WhatsAppDeliveredSubject,
			WhatsAppFailedSubject,
			WhatsAppReadSubject,

			// Telegram Notification Events
			TelegramSentSubject,
			TelegramDeliveredSubject,
			TelegramFailedSubject,

			// Push Notification Events
			PushSentSubject,
			PushDeliveredSubject,
			PushFailedSubject,
			PushOpenedSubject,
		},
	},
}
