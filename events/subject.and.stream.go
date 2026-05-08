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
	VenueStream            StreamName = "VenueStream"
	BillingStream          StreamName = "BillingStream"
	LicenseStream          StreamName = "LicenseStream"
	OLTEventStream         StreamName = "OLTEventStream"
	OLTManagerStream       StreamName = "OLTManagerStream"
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
	CDNProviderCreatedSubject    Subject = "cdnprovider.created"
	CDNProviderUpdateSubject     Subject = "cdnprovider.updated"
	CDNProviderDeleteSubject     Subject = "cdnprovider.deleted"
	PaymentGatewayCreatedSubject Subject = "paymentgateway.created"
	PaymentGatewayUpdateSubject  Subject = "paymentgateway.updated"
	PaymentGatewayDeleteSubject  Subject = "paymentgateway.deleted"
	ExternalRadiusCreatedSubject Subject = "externalradius.created"
	ExternalRadiusUpdateSubject  Subject = "externalradius.updated"
	ExternalRadiusDeleteSubject  Subject = "externalradius.deleted"
	DeviceCreatedSubject         Subject = "device.created"
	DeviceUpdatedSubject         Subject = "device.updated"
	DeviceDeletedSubject         Subject = "device.deleted"
	OLTCreatedSubject            Subject = "olt.created"
	OLTUpdatedSubject            Subject = "olt.updated"
	OLTDeletedSubject            Subject = "olt.deleted"

	// OLT runtime events — pushed by olt-manager from SNMP traps.
	// Distinct stream (OLTEventStream) so trap storms can't crowd out
	// tenant lifecycle events on TenantStream.
	OLTEventONTDownSubject       Subject = "olt.event.ont.down"
	OLTEventONTUpSubject         Subject = "olt.event.ont.up"
	OLTEventDyingGaspSubject     Subject = "olt.event.dying_gasp"
	OLTEventONTDeactivatedSubject Subject = "olt.event.ont.deactivated"
	OLTEventLinkDownSubject      Subject = "olt.event.link.down"
	OLTEventLinkUpSubject        Subject = "olt.event.link.up"
	OLTEventLOSSubject           Subject = "olt.event.los"
	OLTEventLOSRecoveredSubject  Subject = "olt.event.los.recovered"
	OLTEventAlarmActiveSubject   Subject = "olt.event.alarm.active"
	OLTEventAlarmClearedSubject  Subject = "olt.event.alarm.cleared"
	OLTEventColdStartSubject     Subject = "olt.event.cold_start"
	OLTEventAuthFailureSubject   Subject = "olt.event.auth_failure"
	OLTEventTrapUnknownSubject   Subject = "olt.event.trap.unknown"

	// OLT Manager service-level events — sync lifecycle, capability
	// detection, ONT reconciliation, alarm reconciliation, health
	// changes. Distinct stream (OLTManagerStream) so high-velocity trap
	// events on OLTEventStream don't crowd out reconciliation results.
	OLTManagerSyncStartedSubject        Subject = "oltmanager.sync.started"
	OLTManagerSyncCompletedSubject      Subject = "oltmanager.sync.completed"
	OLTManagerSyncFailedSubject         Subject = "oltmanager.sync.failed"
	// Per-phase progress events for multi-step initial / full sync.
	// Run-level Started/Completed/Failed bracket the whole job; the Phase
	// subjects narrate progress inside it (Phase A identity, B chassis,
	// C ports/optics, D VLANs, E ONTs, …). The SSE bridge in olt-manager
	// fans these out to the dashboard so operators see U2000-style phase
	// timelines while a fresh OLT is discovered.
	OLTManagerSyncPhaseStartedSubject   Subject = "oltmanager.sync.phase.started"
	OLTManagerSyncPhaseProgressSubject  Subject = "oltmanager.sync.phase.progress"
	OLTManagerSyncPhaseCompletedSubject Subject = "oltmanager.sync.phase.completed"
	OLTManagerSyncPhaseFailedSubject    Subject = "oltmanager.sync.phase.failed"
	OLTManagerCapabilityDetectedSubject Subject = "oltmanager.capability.detected"
	OLTManagerONTDiscoveredSubject      Subject = "oltmanager.ont.discovered"
	OLTManagerONTUpdatedSubject         Subject = "oltmanager.ont.updated"
	OLTManagerONTDeletedSubject         Subject = "oltmanager.ont.deleted"
	// Provisioning lifecycle events emitted by olt-manager when an
	// operator drives the Discover → Profile → Service Port → Execute
	// wizard. Distinct from the *.discovered/updated/deleted sync-side
	// events because they are operator-triggered and carry the full
	// register / service-port draft for downstream audit + dashboard
	// fan-out.
	OLTManagerONTRegisteredSubject         Subject = "oltmanager.ont.registered"
	OLTManagerONTActivatedSubject          Subject = "oltmanager.ont.activated"
	OLTManagerServicePortCreatedSubject    Subject = "oltmanager.serviceport.created"
	OLTManagerServicePortDeletedSubject    Subject = "oltmanager.serviceport.deleted"
	OLTManagerAlarmReconciledSubject    Subject = "oltmanager.alarm.reconciled"
	OLTManagerHealthChangedSubject      Subject = "oltmanager.health.changed"

	// Inventory Service Events
	InventoryDeviceCreatedSubject Subject = "inventory.device.created"
	InventoryDeviceUpdatedSubject Subject = "inventory.device.updated"
	InventoryDeviceDeletedSubject Subject = "inventory.device.deleted"

	// Inventory Stock Events
	InventoryStockTransferredSubject Subject = "inventory.stock.transferred"
	InventoryStockAdjustedSubject    Subject = "inventory.stock.adjusted"
	InventoryStockProvisionedSubject Subject = "inventory.stock.provisioned"
	InventoryStockLowSubject         Subject = "inventory.stock.low"

	// Inventory Asset Lifecycle Events
	InventoryAssetAssignedSubject         Subject = "inventory.asset.assigned"
	InventoryAssetInstalledSubject        Subject = "inventory.asset.installed"
	InventoryAssetReturnedSubject         Subject = "inventory.asset.returned"
	InventoryAssetFaultySubject           Subject = "inventory.asset.faulty"
	InventoryAssetRMASubject              Subject = "inventory.asset.rma"
	InventoryAssetScrappedSubject         Subject = "inventory.asset.scrapped"
	// New asset lifecycle subjects — full ONT/CPE flow: creation,
	// stock-tenant transfer, condition re-grading via QC, and saga
	// compensation when subscriber-side cpe_requested can't be
	// satisfied by inventory-side authority.
	InventoryAssetCreatedSubject          Subject = "inventory.asset.created"
	InventoryAssetTenantAssignedSubject   Subject = "inventory.asset.tenant_assigned"
	InventoryAssetConditionChangedSubject Subject = "inventory.asset.condition_changed"
	InventoryAssetAssignmentFailedSubject Subject = "inventory.asset.assignment_failed"

	// Inventory Inward Events
	InventoryInwardCreatedSubject Subject = "inventory.inward.created"
	InventoryInwardPostedSubject  Subject = "inventory.inward.posted"

	// Tenant Provider Binding Events
	TenantKYCProviderBindingCreatedSubject     Subject = "tenant.kycbinding.created"
	TenantKYCProviderBindingUpdatedSubject     Subject = "tenant.kycbinding.updated"
	TenantKYCProviderBindingDeletedSubject     Subject = "tenant.kycbinding.deleted"
	TenantAppMessagingBindingCreatedSubject    Subject = "tenant.appmessagingbinding.created"
	TenantAppMessagingBindingUpdatedSubject    Subject = "tenant.appmessagingbinding.updated"
	TenantAppMessagingBindingDeletedSubject    Subject = "tenant.appmessagingbinding.deleted"
	TenantCDNProviderBindingCreatedSubject     Subject = "tenant.cdnbinding.created"
	TenantCDNProviderBindingUpdatedSubject     Subject = "tenant.cdnbinding.updated"
	TenantCDNProviderBindingDeletedSubject     Subject = "tenant.cdnbinding.deleted"
	TenantStorageProviderBindingCreatedSubject Subject = "tenant.storagebinding.created"
	TenantStorageProviderBindingUpdatedSubject Subject = "tenant.storagebinding.updated"
	TenantStorageProviderBindingDeletedSubject Subject = "tenant.storagebinding.deleted"
	TenantESignProviderBindingCreatedSubject   Subject = "tenant.esignbinding.created"
	TenantESignProviderBindingUpdatedSubject   Subject = "tenant.esignbinding.updated"
	TenantESignProviderBindingDeletedSubject   Subject = "tenant.esignbinding.deleted"
	SMSProviderCreatedSubject                  Subject = "tenant.smsprovider.created"
	SMSProviderUpdatedSubject                  Subject = "tenant.smsprovider.updated"
	SMSProviderDeletedSubject                  Subject = "tenant.smsprovider.deleted"
	MailProviderCreatedSubject                 Subject = "tenant.mailprovider.created"
	MailProviderUpdatedSubject                 Subject = "tenant.mailprovider.updated"
	MailProviderDeletedSubject                 Subject = "tenant.mailprovider.deleted"

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

	// Department/Team Events (collection-based: full document)
	DeptTeamCreatedSubject Subject = "deptteam.created"
	DeptTeamUpdatedSubject Subject = "deptteam.updated"
	DeptTeamDeletedSubject Subject = "deptteam.deleted"

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

	// Subscriber CPE assignment-saga events. Apply to BOTH broadband and
	// SMB subscriptions (discriminated via SubscriptionType field on the
	// payload). Published by subscriber-service whenever an operator
	// requests / releases a CPE binding; consumed by inventory-service
	// which makes the authoritative atomic transition on the asset.
	SubscriberCpeRequestedSubject Subject = "subscriber.cpe.requested"
	SubscriberCpeReleasedSubject  Subject = "subscriber.cpe.released"

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
	// PlanTenantPricing events removed — feature deprecated.
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
	// Product Template Events (consumed by inventory-service for cache sync)
	ProductCreatedSubject Subject = "product.created"
	ProductUpdatedSubject Subject = "product.updated"
	ProductDeletedSubject Subject = "product.deleted"

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

	// Venue Service Events (3 simple events — remote services read Status field for decisions)
	VenueOrderCreatedSubject Subject = "venue.order.created"
	VenueOrderUpdatedSubject Subject = "venue.order.updated"
	VenueOrderDeletedSubject Subject = "venue.order.deleted"
	VenueMenuCreatedSubject  Subject = "venue.menu.created"
	VenueMenuUpdatedSubject  Subject = "venue.menu.updated"
	VenueMenuDeletedSubject  Subject = "venue.menu.deleted"

	// Billing Service Events
	BillingPaymentCompletedSubject   Subject = "billing.payment.completed"
	BillingPaymentFailedSubject      Subject = "billing.payment.failed"
	BillingInvoiceCreatedSubject     Subject = "billing.invoice.created"
	BillingCreditNoteIssuedSubject   Subject = "billing.creditnote.issued"
	BillingDebitNoteIssuedSubject    Subject = "billing.debitnote.issued"
	BillingDunningReminderSubject    Subject = "billing.dunning.reminder"
	BillingDunningWarningSubject     Subject = "billing.dunning.warning"
	BillingDunningSuspensionSubject  Subject = "billing.dunning.suspension"
	BillingDunningTerminationSubject Subject = "billing.dunning.termination"

	// Billing CDC Events (Change Data Capture from PostgreSQL WAL)
	BillingInvoiceUpdatedSubject      Subject = "billing.invoice.updated"
	BillingPaymentUpdatedSubject      Subject = "billing.payment.updated"
	BillingSubscriptionCreatedSubject Subject = "billing.subscription.created"
	BillingSubscriptionUpdatedSubject Subject = "billing.subscription.updated"
	BillingCommissionCreatedSubject   Subject = "billing.commission.created"
	BillingCommissionUpdatedSubject   Subject = "billing.commission.updated"
	BillingTransactionCreatedSubject  Subject = "billing.ledger.transaction.created"
	BillingTransactionUpdatedSubject  Subject = "billing.ledger.transaction.updated"
	BillingReceiptCreatedSubject      Subject = "billing.receipt.created"
	BillingReceiptUpdatedSubject      Subject = "billing.receipt.updated"
	BillingRefundCreatedSubject       Subject = "billing.refund.created"
	BillingRefundUpdatedSubject       Subject = "billing.refund.updated"

	// Billing Price Book CDC Events (from PostgreSQL WAL CDC)
	BillingPriceBookCreatedSubject Subject = "billing.pricebook.created"
	BillingPriceBookUpdatedSubject Subject = "billing.pricebook.updated"
	BillingPriceBookDeletedSubject Subject = "billing.pricebook.deleted"

	// License lifecycle (issued/transitioned/renewed/etc.)
	LicenseIssuedSubject              Subject = "license.issued"
	LicenseEntitlementChangedSubject  Subject = "license.entitlement.changed"
	LicenseSuspendedSubject           Subject = "license.suspended"
	LicenseReinstatedSubject          Subject = "license.reinstated"
	LicenseRenewedSubject             Subject = "license.renewed"
	LicenseExpiredSubject             Subject = "license.expired"
	LicenseTerminatedSubject          Subject = "license.terminated"

	// Token wallet events (metered SKU consumption)
	LicenseTokenTopupSubject       Subject = "license.token.topup"
	LicenseTokenConsumedSubject    Subject = "license.token.consumed"
	LicenseTokenLowBalanceSubject  Subject = "license.token.lowbalance"
	LicenseTokenExhaustedSubject   Subject = "license.token.exhausted"

	// Installation lifecycle
	LicenseInstallationEnrolledSubject       Subject = "license.installation.enrolled"
	LicenseInstallationOfflineSubject        Subject = "license.installation.offline"
	LicenseInstallationRecoveredSubject      Subject = "license.installation.recovered"
	LicenseInstallationDecommissionedSubject Subject = "license.installation.decommissioned"

	// JWS revocation broadcast (consumed by all in-cluster verifiers)
	LicenseJWSRevokedSubject Subject = "license.jws.revoked"

	// Subscriber-tier (PER_SUBSCRIBER pricing) — count crossed the
	// licensed quantity. billing-service consumes to bill overage.
	LicenseTierExceededSubject Subject = "license.tier.exceeded"
	LicenseTierRecoveredSubject Subject = "license.tier.recovered"
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

	// Audit Events (Global - any service can publish)
	AuditUserActionSubject       Subject = "audit.user.action"
	AuditAuthActionSubject       Subject = "audit.auth.action"
	AuditTenantActionSubject     Subject = "audit.tenant.action"
	AuditSubscriberActionSubject Subject = "audit.subscriber.action"
	AuditPlanActionSubject       Subject = "audit.plan.action"
	AuditInventoryActionSubject  Subject = "audit.inventory.action"
	AuditTicketActionSubject     Subject = "audit.ticket.action"
	AuditSystemActionSubject     Subject = "audit.system.action"
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
			PaymentGatewayCreatedSubject,
			PaymentGatewayUpdateSubject,
			PaymentGatewayDeleteSubject,
			ExternalRadiusCreatedSubject,
			ExternalRadiusUpdateSubject,
			ExternalRadiusDeleteSubject,
			DeviceCreatedSubject,
			DeviceUpdatedSubject,
			DeviceDeletedSubject,
			// Tenant Provider Binding Events
			CDNProviderCreatedSubject,
			CDNProviderUpdateSubject,
			CDNProviderDeleteSubject,
			TenantCDNProviderBindingCreatedSubject,
			TenantCDNProviderBindingUpdatedSubject,
			TenantCDNProviderBindingDeletedSubject,
			TenantStorageProviderBindingCreatedSubject,
			TenantStorageProviderBindingUpdatedSubject,
			TenantStorageProviderBindingDeletedSubject,
			TenantKYCProviderBindingCreatedSubject,
			TenantKYCProviderBindingUpdatedSubject,
			TenantKYCProviderBindingDeletedSubject,
			TenantAppMessagingBindingCreatedSubject,
			TenantAppMessagingBindingUpdatedSubject,
			TenantAppMessagingBindingDeletedSubject,
			TenantESignProviderBindingCreatedSubject,
			TenantESignProviderBindingUpdatedSubject,
			TenantESignProviderBindingDeletedSubject,
			SMSProviderCreatedSubject,
			SMSProviderUpdatedSubject,
			SMSProviderDeletedSubject,
			MailProviderCreatedSubject,
			MailProviderUpdatedSubject,
			MailProviderDeletedSubject,
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
			// Department/Team Events
			DeptTeamCreatedSubject,
			DeptTeamUpdatedSubject,
			DeptTeamDeletedSubject,
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
	InventoryStream: {
		Name:        InventoryStream,
		Description: "Stream for inventory service events",
		Subjects: []Subject{
			// Device events
			InventoryDeviceCreatedSubject,
			InventoryDeviceUpdatedSubject,
			InventoryDeviceDeletedSubject,
			// Stock events
			InventoryStockTransferredSubject,
			InventoryStockAdjustedSubject,
			InventoryStockProvisionedSubject,
			InventoryStockLowSubject,
			// Asset lifecycle events
			InventoryAssetCreatedSubject,
			InventoryAssetTenantAssignedSubject,
			InventoryAssetConditionChangedSubject,
			InventoryAssetAssignedSubject,
			InventoryAssetAssignmentFailedSubject,
			InventoryAssetInstalledSubject,
			InventoryAssetReturnedSubject,
			InventoryAssetFaultySubject,
			InventoryAssetRMASubject,
			InventoryAssetScrappedSubject,
			// Inward events
			InventoryInwardCreatedSubject,
			InventoryInwardPostedSubject,
		},
	},

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
			// CPE assignment-saga events (cross-domain — broadband + SMB)
			SubscriberCpeRequestedSubject,
			SubscriberCpeReleasedSubject,
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
			// Product Template Events (consumed by inventory-service)
			ProductCreatedSubject,
			ProductUpdatedSubject,
			ProductDeletedSubject,
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

	AuditGlobalStream: {
		Name:        AuditGlobalStream,
		Description: "Global stream for audit trail events from all services",
		Subjects: []Subject{
			AuditUserActionSubject,
			AuditAuthActionSubject,
			AuditTenantActionSubject,
			AuditSubscriberActionSubject,
			AuditPlanActionSubject,
			AuditInventoryActionSubject,
			AuditTicketActionSubject,
			AuditSystemActionSubject,
		},
	},

	VenueStream: {
		Name:        VenueStream,
		Description: "Stream for venue service events (orders, menus)",
		Subjects: []Subject{
			VenueOrderCreatedSubject,
			VenueOrderUpdatedSubject,
			VenueOrderDeletedSubject,
			VenueMenuCreatedSubject,
			VenueMenuUpdatedSubject,
			VenueMenuDeletedSubject,
		},
	},

	BillingStream: {
		Name:        BillingStream,
		Description: "Stream for billing service events (payments, invoices, subscriptions, commissions, CDC)",
		Subjects: []Subject{
			BillingPaymentCompletedSubject,
			BillingPaymentFailedSubject,
			BillingInvoiceCreatedSubject,
			BillingCreditNoteIssuedSubject,
			BillingDebitNoteIssuedSubject,
			BillingDunningReminderSubject,
			BillingDunningWarningSubject,
			BillingDunningSuspensionSubject,
			BillingDunningTerminationSubject,
			// CDC subjects
			BillingInvoiceUpdatedSubject,
			BillingPaymentUpdatedSubject,
			BillingSubscriptionCreatedSubject,
			BillingSubscriptionUpdatedSubject,
			BillingCommissionCreatedSubject,
			BillingCommissionUpdatedSubject,
			BillingTransactionCreatedSubject,
			BillingTransactionUpdatedSubject,
			BillingReceiptCreatedSubject,
			BillingReceiptUpdatedSubject,
			BillingRefundCreatedSubject,
			BillingRefundUpdatedSubject,
			// Price Book CDC subjects
			BillingPriceBookCreatedSubject,
			BillingPriceBookUpdatedSubject,
			BillingPriceBookDeletedSubject,
		},
	},

	OLTEventStream: {
		Name:        OLTEventStream,
		Description: "Real-time OLT runtime events from olt-manager (SNMP traps + active polling). Higher velocity than TenantStream; isolated so trap storms don't crowd out lifecycle events.",
		Subjects: []Subject{
			OLTEventONTDownSubject,
			OLTEventONTUpSubject,
			OLTEventDyingGaspSubject,
			OLTEventONTDeactivatedSubject,
			OLTEventLinkDownSubject,
			OLTEventLinkUpSubject,
			OLTEventLOSSubject,
			OLTEventLOSRecoveredSubject,
			OLTEventAlarmActiveSubject,
			OLTEventAlarmClearedSubject,
			OLTEventColdStartSubject,
			OLTEventAuthFailureSubject,
			OLTEventTrapUnknownSubject,
		},
	},

	OLTManagerStream: {
		Name:        OLTManagerStream,
		Description: "Service-level events from olt-manager-service: OLT lifecycle (created/updated/deleted), sync lifecycle, capability detection, ONT reconciliation, alarm reconciliation, health changes. olt-manager owns the source-of-truth for OLT records and is the publisher of olt.* lifecycle events.",
		Subjects: []Subject{
			// OLT lifecycle (olt-manager is the publisher; was on TenantStream pre-cutover)
			OLTCreatedSubject,
			OLTUpdatedSubject,
			OLTDeletedSubject,
			// Service-level events emitted during sync runs / reconciliation
			OLTManagerSyncStartedSubject,
			OLTManagerSyncCompletedSubject,
			OLTManagerSyncFailedSubject,
			OLTManagerSyncPhaseStartedSubject,
			OLTManagerSyncPhaseProgressSubject,
			OLTManagerSyncPhaseCompletedSubject,
			OLTManagerSyncPhaseFailedSubject,
			OLTManagerCapabilityDetectedSubject,
			OLTManagerONTDiscoveredSubject,
			OLTManagerONTUpdatedSubject,
			OLTManagerONTDeletedSubject,
			OLTManagerONTRegisteredSubject,
			OLTManagerONTActivatedSubject,
			OLTManagerServicePortCreatedSubject,
			OLTManagerServicePortDeletedSubject,
			OLTManagerAlarmReconciledSubject,
			OLTManagerHealthChangedSubject,
		},
	},

	LicenseStream: {
		Name:        LicenseStream,
		Description: "Stream for license-service events: licenses, entitlements, wallets, installations, JWS revocation",
		Subjects: []Subject{
			LicenseIssuedSubject,
			LicenseEntitlementChangedSubject,
			LicenseSuspendedSubject,
			LicenseReinstatedSubject,
			LicenseRenewedSubject,
			LicenseExpiredSubject,
			LicenseTerminatedSubject,
			LicenseTokenTopupSubject,
			LicenseTokenConsumedSubject,
			LicenseTokenLowBalanceSubject,
			LicenseTokenExhaustedSubject,
			LicenseInstallationEnrolledSubject,
			LicenseInstallationOfflineSubject,
			LicenseInstallationRecoveredSubject,
			LicenseInstallationDecommissionedSubject,
			LicenseJWSRevokedSubject,
			LicenseTierExceededSubject,
			LicenseTierRecoveredSubject,
		},
	},
}
