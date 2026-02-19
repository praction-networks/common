package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream          StreamName = "SeedAppStream"
	AuthStream             StreamName = "AuthStream"
	TenantStream           StreamName = "TenantStream"
	TenantUserStream       StreamName = "TenantUserStream"
	InventoryStream        StreamName = "InventoryStream"
	SubscriberStream       StreamName = "SubscriberStream"
	RadiusAccountingStream StreamName = "RadiusAccountingStream"
	CaptivePortalStream    StreamName = "CaptivePortalStream"
	PlanStream             StreamName = "PlanStream"
	LogEngineStream        StreamName = "LogEngineStream"
	TicketStream           StreamName = "TicketStream"
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
	CDNProviderCreatedSubject    Subject = "cdnprovider.created"
	CDNProviderUpdateSubject     Subject = "cdnprovider.updated"
	CDNProviderDeleteSubject     Subject = "cdnprovider.deleted"
	PaymentGatewayCreatedSubject Subject = "paymentgateway.created"
	PaymentGatewayUpdateSubject  Subject = "paymentgateway.updated"
	PaymentGatewayDeleteSubject  Subject = "paymentgateway.deleted"
	ExternalRadiusCreatedSubject Subject = "externalradius.created"
	ExternalRadiusUpdateSubject  Subject = "externalradius.updated"
	ExternalRadiusDeleteSubject  Subject = "externalradius.deleted"
	SMSGatewayCreatedSubject     Subject = "smsgateway.created"
	SMSGatewayUpdateSubject      Subject = "smsgateway.updated"
	SMSGatewayDeleteSubject      Subject = "smsgateway.deleted"
	DeviceCreatedSubject         Subject = "device.created"
	DeviceUpdatedSubject         Subject = "device.updated"
	DeviceDeletedSubject         Subject = "device.deleted"

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
	TenantCDNProviderBindingCreatedSubject  Subject = "tenant.cdnbinding.created"
	TenantCDNProviderBindingUpdatedSubject  Subject = "tenant.cdnbinding.updated"
	TenantCDNProviderBindingDeletedSubject  Subject = "tenant.cdnbinding.deleted"

	// Tenant Branding Events
	TenantBrandingCreatedSubject Subject = "tenant.branding.created"
	TenantBrandingUpdatedSubject Subject = "tenant.branding.updated"
	TenantBrandingDeletedSubject Subject = "tenant.branding.deleted"

	//Domain User Service Event Initialization
	TenantUserCreatedSubject            Subject = "tenantuser.created"
	TenantUserUpdatedSubject            Subject = "tenantuser.updated"
	TenantUserDeletedSubject            Subject = "tenantuser.deleted"
	TenantUserPasswordSetSubject        Subject = "tenantuser.password.set" // Initial password set during onboarding
	TenantUserPreferencesUpdatedSubject Subject = "tenantuser.preferences.updated"

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

	// Voucher Events
	VoucherCreatedSubject Subject = "subscriber.voucher.created"
	VoucherUpdatedSubject Subject = "subscriber.voucher.updated"
	VoucherDeletedSubject Subject = "subscriber.voucher.deleted"

	// RadAcct CDC Events (FreeRADIUS accounting from radius-event-manager-service)
	RadiusAccountingRadAcctSessionStartSubject  Subject = "radiusaccounting.radacct.session.start"
	RadiusAccountingRadAcctSessionUpdateSubject Subject = "radiusaccounting.radacct.session.update"
	RadiusAccountingRadAcctSessionEndSubject    Subject = "radiusaccounting.radacct.session.end"

	// Captive Portal Events
	GuestHotspotSubscriberCreatedSubject          Subject = "subscriber.guest.hotspot.created"
	GuestHotspotSubscriberUpdatedSubject          Subject = "subscriber.guest.hotspot.updated"
	GuestHotspotSubscriberValidityExtendedSubject Subject = "subscriber.guest.hotspot.validity.extended"
	GuestHotspotDeviceAddedSubject                Subject = "subscriber.guest.hotspot.device.added"
	VoucherDetailsSubject                         Subject = "captiveportal.voucher.details"

	// Plan Service Events
	// Plan Tenant Pricing Events (CRITICAL - consumed by billing/subscription services)
	PlanTenantPricingCreatedSubject Subject = "plan.tenant.pricing.created"
	PlanTenantPricingUpdatedSubject Subject = "plan.tenant.pricing.updated"
	PlanTenantPricingDeletedSubject Subject = "plan.tenant.pricing.deleted"
	// Plan Events (consumed by subscriber/billing services)
	PlanCreatedSubject Subject = "plan.created"
	PlanUpdatedSubject Subject = "plan.updated"
	PlanDeletedSubject Subject = "plan.deleted"
	// Price Book Events
	PriceBookCreatedSubject Subject = "pricebook.created"
	PriceBookUpdatedSubject Subject = "pricebook.updated"
	PriceBookDeletedSubject Subject = "pricebook.deleted"
	// Promotion Events
	PromotionCreatedSubject Subject = "promotion.created"
	PromotionUpdatedSubject Subject = "promotion.updated"
	PromotionDeletedSubject Subject = "promotion.deleted"
	// Coupon Events
	CouponCreatedSubject Subject = "coupon.created"
	CouponUpdatedSubject Subject = "coupon.updated"
	CouponDeletedSubject Subject = "coupon.deleted"

	// Ticket Service Events
	TicketCreatedSubject                  Subject = "ticket.created"
	TicketUpdatedSubject                  Subject = "ticket.updated"
	TicketAssignedSubject                 Subject = "ticket.assigned"
	TicketStatusChangedSubject            Subject = "ticket.status.changed"
	TicketResolvedSubject                 Subject = "ticket.resolved"
	TicketClosedSubject                   Subject = "ticket.closed"
	TicketEscalatedSubject                Subject = "ticket.escalated"
	TicketReopenedSubject                 Subject = "ticket.reopened"
	TicketMergedSubject                   Subject = "ticket.merged"
	TicketSplitSubject                    Subject = "ticket.split"
	TicketCommentAddedSubject             Subject = "ticket.comment.added" // DEPRECATED: Use TicketMessageAddedSubject
	TicketMessageAddedSubject             Subject = "ticket.message.added"
	TicketCustomerRepliedSubject          Subject = "ticket.customer.replied"
	TicketAttachmentAddedSubject          Subject = "ticket.attachment.added"
	TicketSLABreachedSubject              Subject = "ticket.sla.breached"
	TicketSLAChangedSubject               Subject = "ticket.sla.changed"
	TicketAppointmentScheduledSubject     Subject = "ticket.appointment.scheduled"
	TicketAppointmentStatusChangedSubject Subject = "ticket.appointment.status.changed"
	TicketChecklistCreatedSubject         Subject = "ticket.checklist.created"
	TicketChecklistItemCompletedSubject   Subject = "ticket.checklist.item.completed"
	TicketAssignmentRequestedSubject      Subject = "ticket.assignment.requested"
	TicketAssignmentApprovedSubject       Subject = "ticket.assignment.approved"
	TicketAssignmentDeniedSubject         Subject = "ticket.assignment.denied"
	TicketTechnicianStatusUpdateSubject   Subject = "ticket.technician.status.update"
	TicketClassifiedSubject               Subject = "ticket.classified"
	TicketAutoLinkedSubject               Subject = "ticket.auto.linked"
	TicketCreatedFromEmailSubject         Subject = "ticket.created.from.email"
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
			DeviceCreatedSubject,
			DeviceUpdatedSubject,
			DeviceDeletedSubject,
			// Tenant Provider Binding Events
			TenantSMSProviderBindingCreatedSubject,
			TenantSMSProviderBindingUpdatedSubject,
			TenantSMSProviderBindingDeletedSubject,
			TenantMailServerBindingCreatedSubject,
			TenantMailServerBindingUpdatedSubject,
			TenantMailServerBindingDeletedSubject,
			CDNProviderCreatedSubject,
			CDNProviderUpdateSubject,
			CDNProviderDeleteSubject,
			TenantCDNProviderBindingCreatedSubject,
			TenantCDNProviderBindingUpdatedSubject,
			TenantCDNProviderBindingDeletedSubject,
			TenantKYCProviderBindingCreatedSubject,
			TenantKYCProviderBindingUpdatedSubject,
			TenantKYCProviderBindingDeletedSubject,
			TenantAppMessagingBindingCreatedSubject,
			TenantAppMessagingBindingUpdatedSubject,
			TenantAppMessagingBindingDeletedSubject,
			// Tenant Branding Events
			TenantBrandingCreatedSubject,
			TenantBrandingUpdatedSubject,
			TenantBrandingDeletedSubject,
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
			TenantUserPreferencesUpdatedSubject,
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
			// Voucher Events
			VoucherCreatedSubject,
			VoucherUpdatedSubject,
			VoucherDeletedSubject,
		},
	},
	RadiusAccountingStream: {
		Name:        RadiusAccountingStream,
		Description: "Stream for radius accounting CDC events from radius-event-manager-service",
		Subjects: []Subject{
			RadiusAccountingRadAcctSessionStartSubject,
			RadiusAccountingRadAcctSessionUpdateSubject,
			RadiusAccountingRadAcctSessionEndSubject,
		},
	},
	CaptivePortalStream: {
		Name:        CaptivePortalStream,
		Description: "Stream for captive portal service events",
		Subjects: []Subject{
			GuestHotspotSubscriberCreatedSubject,
			GuestHotspotSubscriberUpdatedSubject,
			GuestHotspotSubscriberValidityExtendedSubject,
			GuestHotspotDeviceAddedSubject,
			HotspotDeviceAddedSubject,
			VoucherDetailsSubject,
		},
	},
	PlanStream: {
		Name:        PlanStream,
		Description: "Stream for plan service events",
		Subjects: []Subject{
			// Plan Tenant Pricing Events (CRITICAL - consumed by billing/subscription services)
			PlanTenantPricingCreatedSubject,
			PlanTenantPricingUpdatedSubject,
			PlanTenantPricingDeletedSubject,
			// Plan Events (consumed by subscriber/billing services)
			PlanCreatedSubject,
			PlanUpdatedSubject,
			PlanDeletedSubject,
			// Price Book Events
			PriceBookCreatedSubject,
			PriceBookUpdatedSubject,
			PriceBookDeletedSubject,
			// Promotion Events
			PromotionCreatedSubject,
			PromotionUpdatedSubject,
			PromotionDeletedSubject,
			// Coupon Events
			CouponCreatedSubject,
			CouponUpdatedSubject,
			CouponDeletedSubject,
		},
	},
	TicketStream: {
		Name:        TicketStream,
		Description: "Stream for ticket service events",
		Subjects: []Subject{
			// Core Ticket Events
			TicketCreatedSubject,
			TicketUpdatedSubject,
			TicketAssignedSubject,
			TicketStatusChangedSubject,
			TicketResolvedSubject,
			TicketClosedSubject,
			TicketEscalatedSubject,
			TicketReopenedSubject,
			TicketMergedSubject,
			TicketSplitSubject,
			// Ticket Interaction Events
			TicketCommentAddedSubject, // DEPRECATED: kept for backward compatibility
			TicketMessageAddedSubject,
			TicketCustomerRepliedSubject,
			TicketAttachmentAddedSubject,
			// Ticket SLA Events
			TicketSLABreachedSubject,
			TicketSLAChangedSubject,
			// Ticket Appointment Events
			TicketAppointmentScheduledSubject,
			TicketAppointmentStatusChangedSubject,
			// Ticket Checklist Events
			TicketChecklistCreatedSubject,
			TicketChecklistItemCompletedSubject,
			// Ticket Assignment Events
			TicketAssignmentRequestedSubject,
			TicketAssignmentApprovedSubject,
			TicketAssignmentDeniedSubject,
			// Ticket System Events
			TicketTechnicianStatusUpdateSubject,
			TicketClassifiedSubject,
			TicketAutoLinkedSubject,
			TicketCreatedFromEmailSubject,
		},
	},
}
