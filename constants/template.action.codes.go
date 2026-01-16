package constants

import "strings"

// TemplateActionCode represents a template action code following the <DOMAIN>.<ACTION>[.<VARIANT>] format
type TemplateActionCode string

// All supported template action codes
const (
	// -----------------------------------------------------
	// 1. Authentication & Security Actions
	// -----------------------------------------------------
	TemplateActionCodeOTPLogin            TemplateActionCode = "AUTH.OTP_LOGIN"
	TemplateActionCodePasswordRecovery    TemplateActionCode = "AUTH.PASSWORD_RECOVERY"
	TemplateActionCodePasswordChanged     TemplateActionCode = "AUTH.PASSWORD_CHANGED"
	TemplateActionCodeNewDeviceLogin      TemplateActionCode = "AUTH.NEW_DEVICE_LOGIN"
	TemplateActionCodeAccountVerification TemplateActionCode = "AUTH.ACCOUNT_VERIFICATION"
	TemplateActionCodeMFAEnabled          TemplateActionCode = "AUTH.MFA_ENABLED"
	TemplateActionCodeMFADisabled         TemplateActionCode = "AUTH.MFA_DISABLED"
	TemplateActionCodeSuspiciousLogin     TemplateActionCode = "AUTH.SUSPICIOUS_LOGIN"
	TemplateActionCodeLockedOut           TemplateActionCode = "AUTH.LOCKED_OUT"
	TemplateActionCodeUnlocked            TemplateActionCode = "AUTH.UNLOCKED"

	// -----------------------------------------------------
	// 2. User Onboarding & Profile Actions
	// -----------------------------------------------------
	TemplateActionCodeWelcome              TemplateActionCode = "USER.WELCOME"
	TemplateActionCodeEmailVerification    TemplateActionCode = "USER.EMAIL_VERIFICATION"
	TemplateActionCodeMobileVerification   TemplateActionCode = "USER.MOBILE_VERIFICATION"
	TemplateActionCodeWhatsAppVerification TemplateActionCode = "USER.WHATSAPP_VERIFICATION"
	TemplateActionCodeProfileUpdated       TemplateActionCode = "USER.PROFILE_UPDATED"
	TemplateActionCodeKYCRequired          TemplateActionCode = "USER.KYC_REQUIRED"
	TemplateActionCodeKYCApproved          TemplateActionCode = "USER.KYC_APPROVED"
	TemplateActionCodeKYCRejected          TemplateActionCode = "USER.KYC_REJECTED"
	TemplateActionCodeDeactivated          TemplateActionCode = "USER.DEACTIVATED"
	TemplateActionCodeReinstated           TemplateActionCode = "USER.REINSTATED"

	// -----------------------------------------------------
	// 3. Tenant / Organization Management Actions
	// -----------------------------------------------------
	TemplateActionCodeTenantInviteAdmin  TemplateActionCode = "TENANT.INVITE_ADMIN"
	TemplateActionCodeTenantInviteUser   TemplateActionCode = "TENANT.INVITE_USER"
	TemplateActionCodeTenantCreated      TemplateActionCode = "TENANT.CREATED"
	TemplateActionCodeTenantSuspended    TemplateActionCode = "TENANT.SUSPENDED"
	TemplateActionCodeTenantReinstated   TemplateActionCode = "TENANT.REINSTATED"
	TemplateActionCodeTenantDomainMapped TemplateActionCode = "TENANT.DOMAIN_MAPPED"
	TemplateActionCodeTenantPlanChanged  TemplateActionCode = "TENANT.PLAN_CHANGED"

	// -----------------------------------------------------
	// 4. Billing, Payments & Finance Actions
	// -----------------------------------------------------
	TemplateActionCodeBillingInvoiceGenerated           TemplateActionCode = "BILLING.INVOICE_GENERATED"
	TemplateActionCodeBillingPaymentReceipt             TemplateActionCode = "BILLING.PAYMENT_RECEIPT"
	TemplateActionCodeBillingPaymentFailed              TemplateActionCode = "BILLING.PAYMENT_FAILED"
	TemplateActionCodeBillingSubscriptionExpiring       TemplateActionCode = "BILLING.SUBSCRIPTION_EXPIRING"
	TemplateActionCodeBillingSubscriptionExpired        TemplateActionCode = "BILLING.SUBSCRIPTION_EXPIRED"
	TemplateActionCodeBillingOrderConfirmation          TemplateActionCode = "BILLING.ORDER_CONFIRMATION"
	TemplateActionCodeBillingRefundIssued               TemplateActionCode = "BILLING.REFUND_ISSUED"
	TemplateActionCodeBillingUpcomingRenewal            TemplateActionCode = "BILLING.UPCOMING_RENEWAL"
	TemplateActionCodeBillingAutoDebitFailed            TemplateActionCode = "BILLING.AUTO_DEBIT_FAILED"
	TemplateActionCodeBillingLowBalance                 TemplateActionCode = "BILLING.LOW_BALANCE"
	TemplateActionCodeBillingSubscriptionRenewal        TemplateActionCode = "BILLING.SUBSCRIPTION_RENEWAL"
	TemplateActionCodeBillingSubscriptionRenewalFailed  TemplateActionCode = "BILLING.SUBSCRIPTION_RENEWAL_FAILED"
	TemplateActionCodeBillingSubscriptionRenewalSuccess TemplateActionCode = "BILLING.SUBSCRIPTION_RENEWAL_SUCCESS"
	TemplateActionCodeBillingInvoiceOverdue             TemplateActionCode = "BILLING.INVOICE_OVERDUE"
	TemplateActionCodeBillingInvoicePaid                TemplateActionCode = "BILLING.INVOICE_PAID"
	TemplateActionCodeBillingInvoiceGeneratedOverdue    TemplateActionCode = "BILLING.INVOICE_GENERATED_OVERDUE"
	TemplateActionCodeBillingPaymentInitiated           TemplateActionCode = "BILLING.PAYMENT_INITIATED"
	TemplateActionCodeBillingPaymentProcessing          TemplateActionCode = "BILLING.PAYMENT_PROCESSING"
	TemplateActionCodeBillingPaymentReversed            TemplateActionCode = "BILLING.PAYMENT_REVERSED"
	TemplateActionCodeBillingPaymentRefundRequested     TemplateActionCode = "BILLING.REFUND_REQUESTED"
	TemplateActionCodeBillingLateFeeApplied             TemplateActionCode = "BILLING.LATE_FEE_APPLIED"
	TemplateActionCodeBillingCreditApplied              TemplateActionCode = "BILLING.CREDIT_APPLIED"
	TemplateActionCodeBillingAutoDebitSuccess           TemplateActionCode = "BILLING.AUTO_DEBIT_SUCCESS"



	// -----------------------------------------------------
	// 5. Network Management System (NMS) — SNMP / Poller / Alarms Actions
	// -----------------------------------------------------
	TemplateActionCodeNMSDeviceDown        TemplateActionCode = "NMS.DEVICE_DOWN"
	TemplateActionCodeNMSDeviceUp          TemplateActionCode = "NMS.DEVICE_UP"
	TemplateActionCodeNMSMaintenanceAlert  TemplateActionCode = "NMS.MAINTENANCE_ALERT"
	TemplateActionCodeNMSServiceUpdate     TemplateActionCode = "NMS.SERVICE_UPDATE"
	TemplateActionCodeNMSDeviceRebooted    TemplateActionCode = "NMS.DEVICE_REBOOTED"
	TemplateActionCodeNMSConfigPushSuccess TemplateActionCode = "NMS.CONFIG_PUSH_SUCCESS"
	TemplateActionCodeNMSConfigPushFailed  TemplateActionCode = "NMS.CONFIG_PUSH_FAILED"
	TemplateActionCodeNMSHighCPU           TemplateActionCode = "NMS.HIGH_CPU"
	TemplateActionCodeNMSHighMemory        TemplateActionCode = "NMS.HIGH_MEMORY"
	TemplateActionCodeNMSPortDown          TemplateActionCode = "NMS.PORT_DOWN"
	TemplateActionCodeNMSPortUp            TemplateActionCode = "NMS.PORT_UP"
	TemplateActionCodeNMSLinkFlap          TemplateActionCode = "NMS.LINK_FLAP"
	TemplateActionCodeNMSOpticalPowerLow   TemplateActionCode = "NMS.OPTICAL_POWER_LOW"
	TemplateActionCodeNMSOpticalPowerHigh  TemplateActionCode = "NMS.OPTICAL_POWER_HIGH"
	TemplateActionCodeNMSDeviceUnreachable TemplateActionCode = "NMS.DEVICE_UNREACHABLE"
	TemplateActionCodeNMSDeviceAuthFailed  TemplateActionCode = "NMS.DEVICE_AUTH_FAILED"

	// -----------------------------------------------------
	// 6. ACS / TR-069 / GenieACS Actions
	// -----------------------------------------------------
	TemplateActionCodeACSDeviceInform           TemplateActionCode = "ACS.DEVICE_INFORM"
	TemplateActionCodeACSFirmwareUpgradeSuccess TemplateActionCode = "ACS.FIRMWARE_UPGRADE_SUCCESS"
	TemplateActionCodeACSFirmwareUpgradeFailed  TemplateActionCode = "ACS.FIRMWARE_UPGRADE_FAILED"
	TemplateActionCodeACSParameterChanged       TemplateActionCode = "ACS.PARAMETER_CHANGED"
	TemplateActionCodeACSWANDown                TemplateActionCode = "ACS.WAN_DOWN"
	TemplateActionCodeACSWANUp                  TemplateActionCode = "ACS.WAN_UP"
	TemplateActionCodeACSWIFIDisabled           TemplateActionCode = "ACS.WIFI_DISABLED"
	TemplateActionCodeACSWIFIEnabled            TemplateActionCode = "ACS.WIFI_ENABLED"

	// -----------------------------------------------------
	// 7. Inventory & Asset Management Actions
	// -----------------------------------------------------
	TemplateActionCodeInventoryLowStock       TemplateActionCode = "INVENTORY.LOW_STOCK"
	TemplateActionCodeInventoryOutOfStock     TemplateActionCode = "INVENTORY.OUT_OF_STOCK"
	TemplateActionCodeInventorySerialAdded    TemplateActionCode = "INVENTORY.SERIAL_ADDED"
	TemplateActionCodeInventorySerialAssigned TemplateActionCode = "INVENTORY.SERIAL_ASSIGNED"
	TemplateActionCodeInventorySerialReturned TemplateActionCode = "INVENTORY.SERIAL_RETURNED"
	TemplateActionCodeInventoryPOCreated      TemplateActionCode = "INVENTORY.PO_CREATED"
	TemplateActionCodeInventoryPOReceived     TemplateActionCode = "INVENTORY.PO_RECEIVED"
	TemplateActionCodeInventoryExpiryAlert    TemplateActionCode = "INVENTORY.EXPIRY_ALERT"

	// -----------------------------------------------------
	// 8. Support / CRM Actions
	// -----------------------------------------------------
	TemplateActionCodeSupportTicketCreated   TemplateActionCode = "SUPPORT.TICKET_CREATED"
	TemplateActionCodeSupportTicketUpdated   TemplateActionCode = "SUPPORT.TICKET_UPDATED"
	TemplateActionCodeSupportTicketAssigned  TemplateActionCode = "SUPPORT.TICKET_ASSIGNED"
	TemplateActionCodeSupportTicketClosed    TemplateActionCode = "SUPPORT.TICKET_CLOSED"
	TemplateActionCodeSupportFeedbackRequest TemplateActionCode = "SUPPORT.FEEDBACK_REQUEST"

	// -----------------------------------------------------
	// 9. Marketing / Communication Actions
	// -----------------------------------------------------
	TemplateActionCodeMarketingPromotion    TemplateActionCode = "MARKETING.PROMOTION"
	TemplateActionCodeMarketingNewsletter   TemplateActionCode = "MARKETING.NEWSLETTER"
	TemplateActionCodeMarketingAnnouncement TemplateActionCode = "MARKETING.ANNOUNCEMENT"
	TemplateActionCodeMarketingEventInvite  TemplateActionCode = "MARKETING.EVENT_INVITE"
	TemplateActionCodeMarketingSurvey       TemplateActionCode = "MARKETING.SURVEY"

	// -----------------------------------------------------
	// 10. DevOps / Infrastructure Alerts Actions
	// -----------------------------------------------------
	TemplateActionCodeInfraDeploymentFailed  TemplateActionCode = "INFRA.DEPLOYMENT_FAILED"
	TemplateActionCodeInfraDeploymentSuccess TemplateActionCode = "INFRA.DEPLOYMENT_SUCCESS"
	TemplateActionCodeInfraPodRestarting     TemplateActionCode = "INFRA.POD_RESTARTING"
	TemplateActionCodeInfraCPUSpike          TemplateActionCode = "INFRA.CPU_SPIKE"
	TemplateActionCodeInfraDiskFull          TemplateActionCode = "INFRA.DISK_FULL"
	TemplateActionCodeInfraSSLExpiring       TemplateActionCode = "INFRA.SSL_EXPIRING"

	// -----------------------------------------------------
	// 11. Leads & Pre-Sales Actions
	// -----------------------------------------------------
	TemplateActionCodeLeadCreated       TemplateActionCode = "LEAD.CREATED"
	TemplateActionCodeLeadAssigned      TemplateActionCode = "LEAD.ASSIGNED"
	TemplateActionCodeLeadFollowupDue   TemplateActionCode = "LEAD.FOLLOWUP_DUE"
	TemplateActionCodeLeadFollowupDone  TemplateActionCode = "LEAD.FOLLOWUP_DONE"
	TemplateActionCodeLeadStatusUpdated TemplateActionCode = "LEAD.STATUS_UPDATED"
	TemplateActionCodeLeadConverted     TemplateActionCode = "LEAD.CONVERTED"
	TemplateActionCodeLeadLost          TemplateActionCode = "LEAD.LOST"
	TemplateActionCodeLeadReopened      TemplateActionCode = "LEAD.REOPENED"

	// -----------------------------------------------------
	// 11. Hotspots Logic Actions
	// -----------------------------------------------------
	TemplateActionCodeHotSpotLoginOTP TemplateActionCode = "HOTSPOT.SUBSCRIBER_LOGIN_OTP"

	// -----------------------------------------------------
	// 12. Feasibility Check & Site Survey Actions
	// -----------------------------------------------------
	TemplateActionCodeFeasibilityRequested              TemplateActionCode = "FEASIBILITY.REQUESTED"
	TemplateActionCodeFeasibilityAssigned               TemplateActionCode = "FEASIBILITY.ASSIGNED"
	TemplateActionCodeFeasibilitySurveyScheduled        TemplateActionCode = "FEASIBILITY.SURVEY_SCHEDULED"
	TemplateActionCodeFeasibilitySurveyRescheduled      TemplateActionCode = "FEASIBILITY.SURVEY_RESCHEDULED"
	TemplateActionCodeFeasibilitySurveyCompleted        TemplateActionCode = "FEASIBILITY.SURVEY_COMPLETED"
	TemplateActionCodeFeasibilityPassed                 TemplateActionCode = "FEASIBILITY.PASSED"
	TemplateActionCodeFeasibilityFailed                 TemplateActionCode = "FEASIBILITY.FAILED"
	TemplateActionCodeFeasibilityAdditionalInfoRequired TemplateActionCode = "FEASIBILITY.ADDITIONAL_INFO_REQUIRED"
	TemplateActionCodeFeasibilityCancelled              TemplateActionCode = "FEASIBILITY.CANCELLED"

	// -----------------------------------------------------
	// 13. Installation, Provisioning & Activation Actions
	// -----------------------------------------------------
	TemplateActionCodeInstallationRequested    TemplateActionCode = "INSTALLATION.REQUESTED"
	TemplateActionCodeInstallationScheduled    TemplateActionCode = "INSTALLATION.SCHEDULED"
	TemplateActionCodeInstallationRescheduled  TemplateActionCode = "INSTALLATION.RESCHEDULED"
	TemplateActionCodeInstallationAssigned     TemplateActionCode = "INSTALLATION.ASSIGNED"
	TemplateActionCodeInstallationInProgress   TemplateActionCode = "INSTALLATION.IN_PROGRESS"
	TemplateActionCodeInstallationCompleted    TemplateActionCode = "INSTALLATION.COMPLETED"
	TemplateActionCodeInstallationFailed       TemplateActionCode = "INSTALLATION.FAILED"
	TemplateActionCodeInstallationCancelled    TemplateActionCode = "INSTALLATION.CANCELLED"
	TemplateActionCodeProvisioningRequested    TemplateActionCode = "PROVISIONING.REQUESTED"
	TemplateActionCodeProvisioningCompleted    TemplateActionCode = "PROVISIONING.COMPLETED"
	TemplateActionCodeProvisioningFailed       TemplateActionCode = "PROVISIONING.FAILED"
	TemplateActionCodeActivationPendingPayment TemplateActionCode = "ACTIVATION.PENDING_PAYMENT"
	TemplateActionCodeActivationCompleted      TemplateActionCode = "ACTIVATION.COMPLETED"
	TemplateActionCodeActivationFailed         TemplateActionCode = "ACTIVATION.FAILED"

	// -----------------------------------------------------
	// 14. Customer Lifecycle & Churn / Retention Actions
	// -----------------------------------------------------
	TemplateActionCodeCustomerWelcomeCallScheduled TemplateActionCode = "CUSTOMER.WELCOME_CALL_SCHEDULED"
	TemplateActionCodeCustomerWelcomeCallCompleted TemplateActionCode = "CUSTOMER.WELCOME_CALL_COMPLETED"
	TemplateActionCodeCustomerUpgradeOffer         TemplateActionCode = "CUSTOMER.UPGRADE_OFFER"
	TemplateActionCodeCustomerDowngradeRequested   TemplateActionCode = "CUSTOMER.DOWNGRADE_REQUESTED"
	TemplateActionCodeCustomerRelocationRequested  TemplateActionCode = "CUSTOMER.RELOCATION_REQUESTED"
	TemplateActionCodeCustomerRelocationCompleted  TemplateActionCode = "CUSTOMER.RELOCATION_COMPLETED"
	TemplateActionCodeCustomerRenewalReminder      TemplateActionCode = "CUSTOMER.RENEWAL_REMINDER"
	TemplateActionCodeCustomerChurnRisk            TemplateActionCode = "CUSTOMER.CHURN_RISK"
	TemplateActionCodeCustomerDisconnected         TemplateActionCode = "CUSTOMER.DISCONNECTED"

	// -----------------------------------------------------
	// 15. System Templates (System-level notifications)
	// These templates are owned by System tenants only
	// Format: System.<DOMAIN>.<ACTION>
	// -----------------------------------------------------
	SystemTemplateActionCodeOTPLogin            TemplateActionCode = "System.AUTH.OTP_LOGIN"
	SystemTemplateActionCodePasswordRecovery    TemplateActionCode = "System.AUTH.PASSWORD_RECOVERY"
	SystemTemplateActionCodePasswordChanged     TemplateActionCode = "System.AUTH.PASSWORD_CHANGED"
	SystemTemplateActionCodeNewDeviceLogin      TemplateActionCode = "System.AUTH.NEW_DEVICE_LOGIN"
	SystemTemplateActionCodeAccountVerification TemplateActionCode = "System.AUTH.ACCOUNT_VERIFICATION"
	SystemTemplateActionCodeMFAEnabled          TemplateActionCode = "System.AUTH.MFA_ENABLED"
	SystemTemplateActionCodeMFADisabled         TemplateActionCode = "System.AUTH.MFA_DISABLED"
	SystemTemplateActionCodeSuspiciousLogin     TemplateActionCode = "System.AUTH.SUSPICIOUS_LOGIN"
	SystemTemplateActionCodeLockedOut           TemplateActionCode = "System.AUTH.LOCKED_OUT"
	SystemTemplateActionCodeUnlocked            TemplateActionCode = "System.AUTH.UNLOCKED"
)

// String returns the string representation of the action code
func (c TemplateActionCode) String() string {
	return string(c)
}

// IsValid checks if the action code is valid
func (c TemplateActionCode) IsValid() bool {
	_, exists := GetAllActionCodes()[c]
	return exists
}

// GetDescription returns the human-readable description for an action code
func (c TemplateActionCode) GetDescription() string {
	descriptions := GetActionCodeDescriptions()
	if desc, exists := descriptions[c]; exists {
		return desc
	}
	return ""
}

// GetAllActionCodes returns a map of all valid action codes
func GetAllActionCodes() map[TemplateActionCode]bool {
	return map[TemplateActionCode]bool{
		TemplateActionCodeOTPLogin:                          true,
		TemplateActionCodePasswordRecovery:                  true,
		TemplateActionCodePasswordChanged:                   true,
		TemplateActionCodeNewDeviceLogin:                    true,
		TemplateActionCodeAccountVerification:               true,
		TemplateActionCodeMFAEnabled:                        true,
		TemplateActionCodeMFADisabled:                       true,
		TemplateActionCodeSuspiciousLogin:                   true,
		TemplateActionCodeLockedOut:                         true,
		TemplateActionCodeUnlocked:                          true,
		TemplateActionCodeWelcome:                           true,
		TemplateActionCodeEmailVerification:                 true,
		TemplateActionCodeMobileVerification:                true,
		TemplateActionCodeWhatsAppVerification:              true,
		TemplateActionCodeProfileUpdated:                    true,
		TemplateActionCodeKYCRequired:                       true,
		TemplateActionCodeKYCApproved:                       true,
		TemplateActionCodeKYCRejected:                       true,
		TemplateActionCodeDeactivated:                       true,
		TemplateActionCodeReinstated:                        true,
		TemplateActionCodeTenantInviteAdmin:                 true,
		TemplateActionCodeTenantInviteUser:                  true,
		TemplateActionCodeTenantCreated:                     true,
		TemplateActionCodeTenantSuspended:                   true,
		TemplateActionCodeTenantReinstated:                  true,
		TemplateActionCodeTenantDomainMapped:                true,
		TemplateActionCodeTenantPlanChanged:                 true,
		TemplateActionCodeBillingInvoiceGenerated:           true,
		TemplateActionCodeBillingPaymentReceipt:             true,
		TemplateActionCodeBillingPaymentFailed:              true,
		TemplateActionCodeBillingSubscriptionExpiring:       true,
		TemplateActionCodeBillingSubscriptionExpired:        true,
		TemplateActionCodeBillingOrderConfirmation:          true,
		TemplateActionCodeBillingRefundIssued:               true,
		TemplateActionCodeBillingUpcomingRenewal:            true,
		TemplateActionCodeBillingAutoDebitFailed:            true,
		TemplateActionCodeBillingLowBalance:                 true,
		TemplateActionCodeBillingSubscriptionRenewal:        true,
		TemplateActionCodeBillingSubscriptionRenewalFailed:  true,
		TemplateActionCodeBillingSubscriptionRenewalSuccess: true,
		TemplateActionCodeBillingInvoiceOverdue:             true,
		TemplateActionCodeBillingInvoicePaid:                true,
		TemplateActionCodeBillingInvoiceGeneratedOverdue:    true,
		TemplateActionCodeBillingPaymentInitiated:           true,
		TemplateActionCodeBillingPaymentProcessing:          true,
		TemplateActionCodeBillingPaymentReversed:            true,
		TemplateActionCodeBillingPaymentRefundRequested:     true,
		TemplateActionCodeBillingLateFeeApplied:             true,
		TemplateActionCodeBillingCreditApplied:              true,
		TemplateActionCodeBillingAutoDebitSuccess:           true,
		TemplateActionCodeNMSDeviceDown:                     true,
		TemplateActionCodeNMSDeviceUp:                       true,
		TemplateActionCodeNMSMaintenanceAlert:               true,
		TemplateActionCodeNMSServiceUpdate:                  true,
		TemplateActionCodeNMSDeviceRebooted:                 true,
		TemplateActionCodeNMSConfigPushSuccess:              true,
		TemplateActionCodeNMSConfigPushFailed:               true,
		TemplateActionCodeNMSHighCPU:                        true,
		TemplateActionCodeNMSHighMemory:                     true,
		TemplateActionCodeNMSPortDown:                       true,
		TemplateActionCodeNMSPortUp:                         true,
		TemplateActionCodeNMSLinkFlap:                       true,
		TemplateActionCodeNMSOpticalPowerLow:                true,
		TemplateActionCodeNMSOpticalPowerHigh:               true,
		TemplateActionCodeNMSDeviceUnreachable:              true,
		TemplateActionCodeNMSDeviceAuthFailed:               true,
		TemplateActionCodeACSDeviceInform:                   true,
		TemplateActionCodeACSFirmwareUpgradeSuccess:         true,
		TemplateActionCodeACSFirmwareUpgradeFailed:          true,
		TemplateActionCodeACSParameterChanged:               true,
		TemplateActionCodeACSWANDown:                        true,
		TemplateActionCodeACSWANUp:                          true,
		TemplateActionCodeACSWIFIDisabled:                   true,
		TemplateActionCodeACSWIFIEnabled:                    true,
		TemplateActionCodeInventoryLowStock:                 true,
		TemplateActionCodeInventoryOutOfStock:               true,
		TemplateActionCodeInventorySerialAdded:              true,
		TemplateActionCodeInventorySerialAssigned:           true,
		TemplateActionCodeInventorySerialReturned:           true,
		TemplateActionCodeInventoryPOCreated:                true,
		TemplateActionCodeInventoryPOReceived:               true,
		TemplateActionCodeInventoryExpiryAlert:              true,
		TemplateActionCodeSupportTicketCreated:              true,
		TemplateActionCodeSupportTicketUpdated:              true,
		TemplateActionCodeSupportTicketAssigned:             true,
		TemplateActionCodeSupportTicketClosed:               true,
		TemplateActionCodeSupportFeedbackRequest:            true,
		TemplateActionCodeMarketingPromotion:                true,
		TemplateActionCodeMarketingNewsletter:               true,
		TemplateActionCodeMarketingAnnouncement:             true,
		TemplateActionCodeMarketingEventInvite:              true,
		TemplateActionCodeMarketingSurvey:                   true,
		TemplateActionCodeInfraDeploymentFailed:             true,
		TemplateActionCodeInfraDeploymentSuccess:            true,
		TemplateActionCodeInfraPodRestarting:                true,
		TemplateActionCodeInfraCPUSpike:                     true,
		TemplateActionCodeInfraDiskFull:                     true,
		TemplateActionCodeInfraSSLExpiring:                  true,
		TemplateActionCodeLeadCreated:                       true,
		TemplateActionCodeLeadAssigned:                      true,
		TemplateActionCodeLeadFollowupDue:                   true,
		TemplateActionCodeLeadFollowupDone:                  true,
		TemplateActionCodeLeadStatusUpdated:                 true,
		TemplateActionCodeLeadConverted:                     true,
		TemplateActionCodeLeadLost:                          true,
		TemplateActionCodeLeadReopened:                      true,
		TemplateActionCodeFeasibilityRequested:              true,
		TemplateActionCodeFeasibilityAssigned:               true,
		TemplateActionCodeFeasibilitySurveyScheduled:        true,
		TemplateActionCodeFeasibilitySurveyRescheduled:      true,
		TemplateActionCodeFeasibilitySurveyCompleted:        true,
		TemplateActionCodeFeasibilityPassed:                 true,
		TemplateActionCodeFeasibilityFailed:                 true,
		TemplateActionCodeFeasibilityAdditionalInfoRequired: true,
		TemplateActionCodeFeasibilityCancelled:              true,
		TemplateActionCodeInstallationRequested:             true,
		TemplateActionCodeInstallationScheduled:             true,
		TemplateActionCodeInstallationRescheduled:           true,
		TemplateActionCodeInstallationAssigned:              true,
		TemplateActionCodeInstallationInProgress:            true,
		TemplateActionCodeInstallationCompleted:             true,
		TemplateActionCodeInstallationFailed:                true,
		TemplateActionCodeInstallationCancelled:             true,
		TemplateActionCodeProvisioningRequested:             true,
		TemplateActionCodeProvisioningCompleted:             true,
		TemplateActionCodeProvisioningFailed:                true,
		TemplateActionCodeActivationPendingPayment:          true,
		TemplateActionCodeActivationCompleted:               true,
		TemplateActionCodeActivationFailed:                  true,
		TemplateActionCodeCustomerWelcomeCallScheduled:      true,
		TemplateActionCodeCustomerWelcomeCallCompleted:      true,
		TemplateActionCodeCustomerUpgradeOffer:              true,
		TemplateActionCodeCustomerDowngradeRequested:        true,
		TemplateActionCodeCustomerRelocationRequested:       true,
		TemplateActionCodeCustomerRelocationCompleted:       true,
		TemplateActionCodeCustomerRenewalReminder:           true,
		TemplateActionCodeCustomerChurnRisk:                 true,
		TemplateActionCodeCustomerDisconnected:              true,
		TemplateActionCodeHotSpotLoginOTP:                   true,
		// System templates (System.* prefix)
		SystemTemplateActionCodeOTPLogin:                     true,
		SystemTemplateActionCodePasswordRecovery:             true,
		SystemTemplateActionCodePasswordChanged:              true,
		SystemTemplateActionCodeNewDeviceLogin:               true,
		SystemTemplateActionCodeAccountVerification:          true,
		SystemTemplateActionCodeMFAEnabled:                   true,
		SystemTemplateActionCodeMFADisabled:                  true,
		SystemTemplateActionCodeSuspiciousLogin:              true,
		SystemTemplateActionCodeLockedOut:                    true,
		SystemTemplateActionCodeUnlocked:                     true,
	}
}

// GetActionCodeDescriptions returns a map of action codes to their descriptions
func GetActionCodeDescriptions() map[TemplateActionCode]string {
	return map[TemplateActionCode]string{
		TemplateActionCodeOTPLogin:                          "OTP Login - One-Time Password for login",
		TemplateActionCodePasswordRecovery:                  "Password Recovery - Password recovery token",
		TemplateActionCodePasswordChanged:                   "Password Changed - Password change confirmation",
		TemplateActionCodeNewDeviceLogin:                    "New Device Login - Login alert for new device",
		TemplateActionCodeAccountVerification:               "Account Verification - Verification link",
		TemplateActionCodeMFAEnabled:                        "MFA Enabled - Multi-factor authentication activated",
		TemplateActionCodeMFADisabled:                       "MFA Disabled - Multi-factor authentication disabled",
		TemplateActionCodeSuspiciousLogin:                   "Suspicious Login - Unusual login detected",
		TemplateActionCodeLockedOut:                         "Account Locked - Too many failed login attempts",
		TemplateActionCodeUnlocked:                          "Account Unlocked - Unlock confirmation",
		TemplateActionCodeWelcome:                           "Welcome message for new users",
		TemplateActionCodeEmailVerification:                 "Email verification",
		TemplateActionCodeMobileVerification:                "Mobile verification",
		TemplateActionCodeWhatsAppVerification:              "WhatsApp verification",
		TemplateActionCodeProfileUpdated:                    "Profile updated confirmation",
		TemplateActionCodeKYCRequired:                       "KYC Required - KYC submission request",
		TemplateActionCodeKYCApproved:                       "KYC Approved - Verification successful",
		TemplateActionCodeKYCRejected:                       "KYC Rejected - Verification failed",
		TemplateActionCodeDeactivated:                       "User Deactivated",
		TemplateActionCodeReinstated:                        "User Reinstated",
		TemplateActionCodeTenantInviteAdmin:                 "Admin Invitation",
		TemplateActionCodeTenantInviteUser:                  "User Invitation",
		TemplateActionCodeTenantCreated:                     "Tenant Created - New tenant onboarded",
		TemplateActionCodeTenantSuspended:                   "Tenant Suspended - Service disabled",
		TemplateActionCodeTenantReinstated:                  "Tenant Reinstated - Service restored",
		TemplateActionCodeTenantDomainMapped:                "Domain Mapping Updated",
		TemplateActionCodeTenantPlanChanged:                 "Subscription Plan Changed",
		TemplateActionCodeBillingInvoiceGenerated:           "Invoice generated",
		TemplateActionCodeBillingPaymentReceipt:             "Payment receipt confirmation",
		TemplateActionCodeBillingPaymentFailed:              "Payment failed",
		TemplateActionCodeBillingSubscriptionExpiring:       "Subscription expiring",
		TemplateActionCodeBillingSubscriptionExpired:        "Subscription expired",
		TemplateActionCodeBillingOrderConfirmation:          "Order confirmation",
		TemplateActionCodeBillingRefundIssued:               "Refund issued",
		TemplateActionCodeBillingUpcomingRenewal:            "Upcoming subscription renewal",
		TemplateActionCodeBillingAutoDebitFailed:            "Auto debit failed",
		TemplateActionCodeBillingLowBalance:                 "Low balance alert",
		TemplateActionCodeBillingSubscriptionRenewal:        "Subscription renewal",
		TemplateActionCodeBillingSubscriptionRenewalFailed:  "Subscription renewal failed",
		TemplateActionCodeBillingSubscriptionRenewalSuccess: "Subscription renewal successful",
		TemplateActionCodeBillingInvoiceOverdue:             "Invoice overdue",
		TemplateActionCodeBillingInvoicePaid:                "Invoice paid",
		TemplateActionCodeBillingInvoiceGeneratedOverdue:    "Invoice generated overdue",
		TemplateActionCodeBillingPaymentInitiated:           "Payment initiated",
		TemplateActionCodeBillingPaymentProcessing:          "Payment processing",
		TemplateActionCodeBillingPaymentReversed:            "Payment reversed",
		TemplateActionCodeBillingPaymentRefundRequested:     "Payment refund requested",
		TemplateActionCodeBillingLateFeeApplied:             "Late fee applied",
		TemplateActionCodeBillingCreditApplied:              "Credit applied",
		TemplateActionCodeBillingAutoDebitSuccess:           "Auto debit successful",	
		TemplateActionCodeNMSDeviceDown:                     "Device offline alert",
		TemplateActionCodeNMSDeviceUp:                       "Device online alert",
		TemplateActionCodeNMSMaintenanceAlert:               "Maintenance scheduled alert",
		TemplateActionCodeNMSServiceUpdate:                  "Network service update",
		TemplateActionCodeNMSDeviceRebooted:                 "Device reboot detected",
		TemplateActionCodeNMSConfigPushSuccess:              "Configuration push successful",
		TemplateActionCodeNMSConfigPushFailed:               "Configuration push failed",
		TemplateActionCodeNMSHighCPU:                        "High CPU usage alert",
		TemplateActionCodeNMSHighMemory:                     "High memory usage alert",
		TemplateActionCodeNMSPortDown:                       "Port down alert",
		TemplateActionCodeNMSPortUp:                         "Port up alert",
		TemplateActionCodeNMSLinkFlap:                       "Link flapping detected",
		TemplateActionCodeNMSOpticalPowerLow:                "Low optical power detected",
		TemplateActionCodeNMSOpticalPowerHigh:               "High optical power detected",
		TemplateActionCodeNMSDeviceUnreachable:              "Device unreachable",
		TemplateActionCodeNMSDeviceAuthFailed:               "Device SNMP authentication failure",
		TemplateActionCodeACSDeviceInform:                   "ACS Inform received",
		TemplateActionCodeACSFirmwareUpgradeSuccess:         "Firmware upgrade successful",
		TemplateActionCodeACSFirmwareUpgradeFailed:          "Firmware upgrade failed",
		TemplateActionCodeACSParameterChanged:               "Device parameter changed",
		TemplateActionCodeACSWANDown:                        "ONT WAN down alert",
		TemplateActionCodeACSWANUp:                          "ONT WAN up alert",
		TemplateActionCodeACSWIFIDisabled:                   "WiFi disabled event",
		TemplateActionCodeACSWIFIEnabled:                    "WiFi enabled event",
		TemplateActionCodeInventoryLowStock:                 "Low stock alert",
		TemplateActionCodeInventoryOutOfStock:               "Out of stock alert",
		TemplateActionCodeInventorySerialAdded:              "Serial number added to inventory",
		TemplateActionCodeInventorySerialAssigned:           "Serial assigned to customer/device",
		TemplateActionCodeInventorySerialReturned:           "Serial returned to stock",
		TemplateActionCodeInventoryPOCreated:                "Purchase order created",
		TemplateActionCodeInventoryPOReceived:               "Purchase order items received",
		TemplateActionCodeInventoryExpiryAlert:              "Stock expiry alert",
		TemplateActionCodeSupportTicketCreated:              "Support ticket created",
		TemplateActionCodeSupportTicketUpdated:              "Support ticket updated",
		TemplateActionCodeSupportTicketAssigned:             "Ticket assigned to agent",
		TemplateActionCodeSupportTicketClosed:               "Support ticket closed",
		TemplateActionCodeSupportFeedbackRequest:            "Feedback request after ticket resolution",
		TemplateActionCodeMarketingPromotion:                "Promotional message",
		TemplateActionCodeMarketingNewsletter:               "Newsletter",
		TemplateActionCodeMarketingAnnouncement:             "General announcement",
		TemplateActionCodeMarketingEventInvite:              "Event invitation",
		TemplateActionCodeMarketingSurvey:                   "Survey request",
		TemplateActionCodeInfraDeploymentFailed:             "Deployment failed alert",
		TemplateActionCodeInfraDeploymentSuccess:            "Deployment successful",
		TemplateActionCodeInfraPodRestarting:                "Pod restarting / crash loop detected",
		TemplateActionCodeInfraCPUSpike:                     "High CPU usage on node",
		TemplateActionCodeInfraDiskFull:                     "Disk usage critically high",
		TemplateActionCodeInfraSSLExpiring:                  "SSL certificate expiring soon",
		TemplateActionCodeLeadCreated:                       "Lead Created - New lead captured",
		TemplateActionCodeLeadAssigned:                      "Lead Assigned - Assigned to sales owner",
		TemplateActionCodeLeadFollowupDue:                   "Lead Follow-up Due - Reminder to contact lead",
		TemplateActionCodeLeadFollowupDone:                  "Lead Follow-up Completed",
		TemplateActionCodeLeadStatusUpdated:                 "Lead Status Updated - Stage changed",
		TemplateActionCodeLeadConverted:                     "Lead Converted - Converted to customer/opportunity",
		TemplateActionCodeLeadLost:                          "Lead Lost - Not converted",
		TemplateActionCodeLeadReopened:                      "Lead Reopened - Lead reactivated",
		TemplateActionCodeFeasibilityRequested:              "Feasibility Requested - Check feasibility for address",
		TemplateActionCodeFeasibilityAssigned:               "Feasibility Assigned - Task assigned to engineer/team",
		TemplateActionCodeFeasibilitySurveyScheduled:        "Survey Scheduled - Site survey appointment fixed",
		TemplateActionCodeFeasibilitySurveyRescheduled:      "Survey Rescheduled - Appointment changed",
		TemplateActionCodeFeasibilitySurveyCompleted:        "Survey Completed - Field visit done",
		TemplateActionCodeFeasibilityPassed:                 "Feasibility Passed - Service possible at location",
		TemplateActionCodeFeasibilityFailed:                 "Feasibility Failed - Service not possible at location",
		TemplateActionCodeFeasibilityAdditionalInfoRequired: "Feasibility Additional Info Required - Need more details",
		TemplateActionCodeFeasibilityCancelled:              "Feasibility Cancelled",
		TemplateActionCodeInstallationRequested:             "Installation Requested - New connection requested",
		TemplateActionCodeInstallationScheduled:             "Installation Scheduled - Appointment booked",
		TemplateActionCodeInstallationRescheduled:           "Installation Rescheduled",
		TemplateActionCodeInstallationAssigned:              "Installation Assigned - Task assigned to field engineer",
		TemplateActionCodeInstallationInProgress:            "Installation In Progress - Engineer on site / work started",
		TemplateActionCodeInstallationCompleted:             "Installation Completed - Physical work done",
		TemplateActionCodeInstallationFailed:                "Installation Failed - Could not complete installation",
		TemplateActionCodeInstallationCancelled:             "Installation Cancelled by customer or system",
		TemplateActionCodeProvisioningRequested:             "Provisioning Requested - Service/device provisioning",
		TemplateActionCodeProvisioningCompleted:             "Provisioning Completed - Service activated in network",
		TemplateActionCodeProvisioningFailed:                "Provisioning Failed - Error during activation",
		TemplateActionCodeActivationPendingPayment:          "Activation Pending - Awaiting payment/verification",
		TemplateActionCodeActivationCompleted:               "Activation Completed - Service live for customer",
		TemplateActionCodeActivationFailed:                  "Activation Failed - Could not activate service",
		TemplateActionCodeCustomerWelcomeCallScheduled:      "Welcome Call Scheduled - Post-install touchpoint",
		TemplateActionCodeCustomerWelcomeCallCompleted:      "Welcome Call Completed",
		TemplateActionCodeCustomerUpgradeOffer:              "Upgrade Offer - Higher plan / speed proposal",
		TemplateActionCodeCustomerDowngradeRequested:        "Downgrade Requested",
		TemplateActionCodeCustomerRelocationRequested:       "Relocation Requested - Shift connection to new address",
		TemplateActionCodeCustomerRelocationCompleted:       "Relocation Completed",
		TemplateActionCodeCustomerRenewalReminder:           "Renewal Reminder - Plan renewal due",
		TemplateActionCodeCustomerChurnRisk:                 "Churn Risk Alert - At-risk customer detected",
		TemplateActionCodeCustomerDisconnected:              "Customer Disconnected - Permanent disconnection",
		TemplateActionCodeHotSpotLoginOTP:                   "Hotspot Login OTP - One-Time Password for login",
		// System template descriptions
		SystemTemplateActionCodeOTPLogin:                    "System OTP Login - One-Time Password for login (System template)",
		SystemTemplateActionCodePasswordRecovery:           "System Password Recovery - Password recovery token (System template)",
		SystemTemplateActionCodePasswordChanged:            "System Password Changed - Password change confirmation (System template)",
		SystemTemplateActionCodeNewDeviceLogin:             "System New Device Login - Login alert for new device (System template)",
		SystemTemplateActionCodeAccountVerification:        "System Account Verification - Verification link (System template)",
		SystemTemplateActionCodeMFAEnabled:                 "System MFA Enabled - Multi-factor authentication activated (System template)",
		SystemTemplateActionCodeMFADisabled:                "System MFA Disabled - Multi-factor authentication disabled (System template)",
		SystemTemplateActionCodeSuspiciousLogin:           "System Suspicious Login - Unusual login detected (System template)",
		SystemTemplateActionCodeLockedOut:                   "System Account Locked - Too many failed login attempts (System template)",
		SystemTemplateActionCodeUnlocked:                   "System Account Unlocked - Unlock confirmation (System template)",
	}
}

// GetAllActionCodeStrings returns a slice of all action code strings
func GetAllActionCodeStrings() []string {
	codes := GetAllActionCodes()
	result := make([]string, 0, len(codes))
	for code := range codes {
		result = append(result, code.String())
	}
	return result
}

// IsSystemTemplateCode checks if a template action code is a system template
// System templates have the "System." prefix (e.g., "System.AUTH.OTP_LOGIN")
// Tenant templates do not have the "System." prefix (e.g., "AUTH.OTP_LOGIN", "BILLING.INVOICE_GENERATED")
func IsSystemTemplateCode(code string) bool {
	return strings.HasPrefix(code, "System.")
}
