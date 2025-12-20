package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream    StreamName = "SeedAppStream"
	AuthStream       StreamName = "AuthStream"
	TenantStream     StreamName = "TenantStream"
	TenantUserStream StreamName = "TenantUserStream"
	InventoryStream  StreamName = "InventoryStream"
	SubscriberStream StreamName = "SubscriberStream"
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

	// Tenant Provider Binding Events
	TenantSMSProviderBindingCreatedSubject  Subject = "tenant.smsbinding.created"
	TenantSMSProviderBindingUpdatedSubject  Subject = "tenant.smsbinding.updated"
	TenantSMSProviderBindingDeletedSubject  Subject = "tenant.smsbinding.deleted"
	TenantMailServerBindingCreatedSubject   Subject = "tenant.mailbinding.created"
	TenantMailServerBindingUpdatedSubject   Subject = "tenant.mailbinding.updated"
	TenantMailServerBindingDeletedSubject   Subject = "tenant.mailbinding.deleted"
	TenantKYCProviderBindingCreatedSubject  Subject = "tenant.kycbinding.created"
	TenantKYCProviderBindingUpdatedSubject  Subject = "tenant.kycbinding.updated"
	TenantKYCProviderBindingDeletedSubject  Subject = "tenant.kycbinding.deleted"
	TenantAppMessagingBindingCreatedSubject Subject = "tenant.appmessagingbinding.created"
	TenantAppMessagingBindingUpdatedSubject Subject = "tenant.appmessagingbinding.updated"
	TenantAppMessagingBindingDeletedSubject Subject = "tenant.appmessagingbinding.deleted"

	//Domain User Service Event Initialization
	TenantUserCreatedSubject     Subject = "tenantuser.created"
	TenantUserUpdatedSubject     Subject = "tenantuser.updated"
	TenantUserDeletedSubject     Subject = "tenantuser.deleted"
	TenantUserPasswordSetSubject Subject = "tenantuser.password.set" // Initial password set during onboarding

	// Tenant Auth Role Event Initialization
	TenantUserRoleCreatedSubject Subject = "tenantuserrole.created"
	TenantUserRoleUpdatedSubject Subject = "tenantuserrole.updated"
	TenantUserRoleDeletedSubject Subject = "tenantuserrole.deleted"

	// Subscriber Service Event Initialization
	SubscriberCreatedSubject Subject = "subscriber.created"
	SubscriberUpdatedSubject Subject = "subscriber.updated"
	SubscriberDeletedSubject Subject = "subscriber.deleted"

	// Broadband Subscription Events
	BroadbandSubscriptionCreatedSubject Subject = "subscriber.broadband.created"
	BroadbandSubscriptionUpdatedSubject Subject = "subscriber.broadband.updated"
	BroadbandSubscriptionDeletedSubject Subject = "subscriber.broadband.deleted"

	// Hotspot Profile Events
	HotspotProfileCreatedSubject Subject = "subscriber.hotspot.created"
	HotspotProfileUpdatedSubject Subject = "subscriber.hotspot.updated"
	HotspotProfileDeletedSubject Subject = "subscriber.hotspot.deleted"
	HotspotDeviceAddedSubject    Subject = "subscriber.hotspot.device.added"
	HotspotDeviceRemovedSubject  Subject = "subscriber.hotspot.device.removed"

	// Field Configuration Events
	FieldConfigCreatedSubject Subject = "subscriber.fieldconfig.created"
	FieldConfigUpdatedSubject Subject = "subscriber.fieldconfig.updated"
	FieldConfigDeletedSubject Subject = "subscriber.fieldconfig.deleted"

	// Form Configuration Events
	FormConfigCreatedSubject Subject = "subscriber.formconfig.created"
	FormConfigUpdatedSubject Subject = "subscriber.formconfig.updated"
	FormConfigDeletedSubject Subject = "subscriber.formconfig.deleted"

	// Voucher Template Events
	VoucherTemplateCreatedSubject Subject = "subscriber.voucher.template.created"
	VoucherTemplateUpdatedSubject Subject = "subscriber.voucher.template.updated"
	VoucherTemplateDeletedSubject Subject = "subscriber.voucher.template.deleted"

	// Voucher Instance Events
	VoucherInstanceCreatedSubject     Subject = "subscriber.voucher.instance.created"
	VoucherInstanceBulkCreatedSubject Subject = "subscriber.voucher.instance.bulk.created"
	VoucherInstanceUsedSubject        Subject = "subscriber.voucher.instance.used"
	VoucherInstanceExpiredSubject     Subject = "subscriber.voucher.instance.expired"
	VoucherInstanceRevokedSubject     Subject = "subscriber.voucher.instance.revoked"
	VoucherInstanceExtendedSubject    Subject = "subscriber.voucher.instance.extended"

	// Theme Events (Tenant Service) - Unified subjects, portalType in event payload
	ThemeCreatedSubject    Subject = "theme.created"
	ThemeUpdatedSubject    Subject = "theme.updated"
	ThemeDeletedSubject    Subject = "theme.deleted"
	ThemeSetDefaultSubject Subject = "theme.set_default"
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
			// Tenant Provider Binding Events
			TenantSMSProviderBindingCreatedSubject,
			TenantSMSProviderBindingUpdatedSubject,
			TenantSMSProviderBindingDeletedSubject,
			TenantMailServerBindingCreatedSubject,
			TenantMailServerBindingUpdatedSubject,
			TenantMailServerBindingDeletedSubject,
			TenantKYCProviderBindingCreatedSubject,
			TenantKYCProviderBindingUpdatedSubject,
			TenantKYCProviderBindingDeletedSubject,
			TenantAppMessagingBindingCreatedSubject,
			TenantAppMessagingBindingUpdatedSubject,
			TenantAppMessagingBindingDeletedSubject,
			// Theme Events (unified subjects, portalType in event payload)
			ThemeCreatedSubject,
			ThemeUpdatedSubject,
			ThemeDeletedSubject,
			ThemeSetDefaultSubject,
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
			TenantUserPasswordSetSubject,
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
	SubscriberStream: {
		Name:        SubscriberStream,
		Description: "Stream for subscriber service events including broadband, hotspot, and field configurations",
		Subjects: []Subject{
			// Subscriber Events
			SubscriberCreatedSubject,
			SubscriberUpdatedSubject,
			SubscriberDeletedSubject,
			// Broadband Subscription Events
			BroadbandSubscriptionCreatedSubject,
			BroadbandSubscriptionUpdatedSubject,
			BroadbandSubscriptionDeletedSubject,
			// Hotspot Profile Events
			HotspotProfileCreatedSubject,
			HotspotProfileUpdatedSubject,
			HotspotProfileDeletedSubject,
			HotspotDeviceAddedSubject,
			HotspotDeviceRemovedSubject,
			// Field Configuration Events
			FieldConfigCreatedSubject,
			FieldConfigUpdatedSubject,
			FieldConfigDeletedSubject,
			// Form Configuration Events
			FormConfigCreatedSubject,
			FormConfigUpdatedSubject,
			FormConfigDeletedSubject,
			// Voucher Template Events
			VoucherTemplateCreatedSubject,
			VoucherTemplateUpdatedSubject,
			VoucherTemplateDeletedSubject,
			// Voucher Instance Events
			VoucherInstanceCreatedSubject,
			VoucherInstanceBulkCreatedSubject,
			VoucherInstanceUsedSubject,
			VoucherInstanceExpiredSubject,
			VoucherInstanceRevokedSubject,
			VoucherInstanceExtendedSubject,
		},
	},
}
