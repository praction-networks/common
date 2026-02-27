package tenantevent

// Field constants for EnabledFeatures - used for tracking changed fields
const (
	// Core features
	EnabledFeaturesFieldCoreIsHotspotEnabled      = "Core.IsHotspotEnabled"
	EnabledFeaturesFieldCoreIsUserKYCEnabled      = "Core.IsUserKYCEnabled"
	EnabledFeaturesFieldCoreIsNotificationEnabled = "Core.IsNotificationEnabled"
	EnabledFeaturesFieldCoreIsUserPortalEnabled   = "Core.IsUserPortalEnabled"
	EnabledFeaturesFieldCoreIsVasEnabled          = "Core.IsVasEnabled"
	EnabledFeaturesFieldCoreIsMarketingEnabled    = "Core.IsMarketingEnabled"
	EnabledFeaturesFieldCoreIsBillingEnabled      = "Core.IsBillingEnabled"
	EnabledFeaturesFieldCoreIsOperationsEnabled   = "Core.IsOperationsEnabled"
	EnabledFeaturesFieldCoreIsAnalyticsEnabled    = "Core.IsAnalyticsEnabled"
	EnabledFeaturesFieldCoreIsI9ShieldEnabled     = "Core.IsI9ShieldEnabled"
	EnabledFeaturesFieldCoreIsAIEnabled           = "Core.IsAIEnabled"

	// Notification features
	EnabledFeaturesFieldNotificationsIsUserMailNotificationEnabled     = "Notifications.IsUserMailNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsUserSMSNotificationEnabled      = "Notifications.IsUserSMSNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsUserWhatsappNotificationEnabled = "Notifications.IsUserWhatsappNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsUserTelegramNotificationEnabled = "Notifications.IsUserTelegramNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsPushNotificationEnabled         = "Notifications.IsPushNotificationEnabled"

	// Hotspot features
	EnabledFeaturesFieldHotspotIsOTPAuthEnabled           = "Hotspot.IsOTPAuthEnabled"
	EnabledFeaturesFieldHotspotIsSocialLoginEnabled       = "Hotspot.IsSocialLoginEnabled"
	EnabledFeaturesFieldHotspotIsVoucherLoginEnabled      = "Hotspot.IsVoucherLoginEnabled"
	EnabledFeaturesFieldHotspotIsPasswordLoginEnabled     = "Hotspot.IsPasswordLoginEnabled"
	EnabledFeaturesFieldHotspotIsGuestWiFiEnabled         = "Hotspot.IsGuestWiFiEnabled"
	EnabledFeaturesFieldHotspotIsSessionManagementEnabled = "Hotspot.IsSessionManagementEnabled"
	EnabledFeaturesFieldHotspotIsCustomBrandingEnabled    = "Hotspot.IsCustomBrandingEnabled"
	EnabledFeaturesFieldHotspotIsAdvertisementEnabled     = "Hotspot.IsAdvertisementEnabled"
	EnabledFeaturesFieldHotspotIsUsageLimitEnabled        = "Hotspot.IsUsageLimitEnabled"
	EnabledFeaturesFieldHotspotIsMultipleSSIDEnabled      = "Hotspot.IsMultipleSSIDEnabled"

	// VAS features
	EnabledFeaturesFieldVASIsIPTVEnabled            = "VAS.IsIPTVEnabled"
	EnabledFeaturesFieldVASIsOTTEnabled             = "VAS.IsOTTEnabled"
	EnabledFeaturesFieldVASIsVoiceServiceEnabled    = "VAS.IsVoiceServiceEnabled"
	EnabledFeaturesFieldVASIsVPNEnabled             = "VAS.IsVPNEnabled"
	EnabledFeaturesFieldVASIsCloudBackupEnabled     = "VAS.IsCloudBackupEnabled"
	EnabledFeaturesFieldVASIsFirewallServiceEnabled = "VAS.IsFirewallServiceEnabled"
	EnabledFeaturesFieldVASIsDNSSecurityEnabled     = "VAS.IsDNSSecurityEnabled"

	// Marketing features
	EnabledFeaturesFieldMarketingIsCampaginServiceEnabled     = "Marketing.IsCampaginServiceEnabled"
	EnabledFeaturesFieldMarketingIsPromoCampaignEnabled       = "Marketing.IsPromoCampaignEnabled"
	EnabledFeaturesFieldMarketingIsEmailCampaignEnabled       = "Marketing.IsEmailCampaignEnabled"
	EnabledFeaturesFieldMarketingIsSMSCampaignEnabled         = "Marketing.IsSMSCampaignEnabled"
	EnabledFeaturesFieldMarketingIsMarketingAutomationEnabled = "Marketing.IsMarketingAutomationEnabled"
	EnabledFeaturesFieldMarketingIsCouponsEnabled             = "Marketing.IsCouponsEnabled"
	EnabledFeaturesFieldMarketingIsLoyaltyProgramEnabled      = "Marketing.IsLoyaltyProgramEnabled"
	EnabledFeaturesFieldMarketingIsReferralProgramEnabled     = "Marketing.IsReferralProgramEnabled"

	// Billing features
	EnabledFeaturesFieldBillingIsOnlinePaymentEnabled       = "Billing.IsOnlinePaymentEnabled"
	EnabledFeaturesFieldBillingIsMultiSplitPaymentEnabled   = "Billing.IsMultiSplitPaymentEnabled"
	EnabledFeaturesFieldBillingIsInvoiceAutoReminderEnabled = "Billing.IsInvoiceAutoReminderEnabled"
	EnabledFeaturesFieldBillingIsPaymentSplitEnabled        = "Billing.IsPaymentSplitEnabled"

	// Operations features
	EnabledFeaturesFieldOperationsIsTicketingEnabled      = "Operations.IsTicketingEnabled"
	EnabledFeaturesFieldOperationsIsCRMIntegrationEnabled = "Operations.IsCRMIntegrationEnabled"
	EnabledFeaturesFieldOperationsIsERPIntegrationEnabled = "Operations.IsERPIntegrationEnabled"
	EnabledFeaturesFieldOperationsIsInventorySyncEnabled  = "Operations.IsInventorySyncEnabled"
	EnabledFeaturesFieldOperationsIsPartnerPortalEnabled  = "Operations.IsPartnerPortalEnabled"

	// Analytics features
	EnabledFeaturesFieldAnalyticsIsUsageAnalyticsEnabled   = "Analytics.IsUsageAnalyticsEnabled"
	EnabledFeaturesFieldAnalyticsIsRevenueAnalyticsEnabled = "Analytics.IsRevenueAnalyticsEnabled"
	EnabledFeaturesFieldAnalyticsIsNetworkAnalyticsEnabled = "Analytics.IsNetworkAnalyticsEnabled"
	EnabledFeaturesFieldAnalyticsIsCustomerFeedbackEnabled = "Analytics.IsCustomerFeedbackEnabled"
	EnabledFeaturesFieldAnalyticsIsChurnPredictionEnabled  = "Analytics.IsChurnPredictionEnabled"

	// Inventory features
	EnabledFeaturesFieldInventoryIsInventoryEnabled = "Inventory.IsInventoryEnabled"

	// I9Shield features (Security)
	EnabledFeaturesFieldI9ShieldIs2FAEnabled                 = "I9Shield.Is2FAEnabled"
	EnabledFeaturesFieldI9ShieldIsAuditLoggingEnabled        = "I9Shield.IsAuditLoggingEnabled"
	EnabledFeaturesFieldI9ShieldIsDataRetentionPolicyEnabled = "I9Shield.IsDataRetentionPolicyEnabled"
	EnabledFeaturesFieldI9ShieldIsGDPRComplianceEnabled      = "I9Shield.IsGDPRComplianceEnabled"

	// AI features
	EnabledFeaturesFieldAIIsAIAssistantEnabled           = "AI.IsAIAssistantEnabled"
	EnabledFeaturesFieldAIIsAIBasedSupportEnabled        = "AI.IsAIBasedSupportEnabled"
	EnabledFeaturesFieldAIIsNetworkOptimizationAIEnabled = "AI.IsNetworkOptimizationAIEnabled"
	EnabledFeaturesFieldAIIsPredictiveMaintenanceEnabled = "AI.IsPredictiveMaintenanceEnabled"
)

// EnabledFeatures defines all feature toggles for a Domain/Tenant, organized by category
type EnabledFeatures struct {
	Core          CoreFeatures         `json:"core" bson:"core"`                   // Core user experience features
	Notifications NotificationFeatures `json:"notifications" bson:"notifications"` // Notification channels
	Hotspot       HotspotFeatures      `json:"hotspot" bson:"hotspot"`             // Hotspot & captive portal
	VAS           VASFeatures          `json:"vas" bson:"vas"`                     // Value-added services
	Marketing     MarketingFeatures    `json:"marketing" bson:"marketing"`         // Marketing & campaigns
	Billing       BillingFeatures      `json:"billing" bson:"billing"`             // Billing & collections
	Operations    OperationsFeatures   `json:"operations" bson:"operations"`       // Operations & integrations
	Analytics     AnalyticsFeatures    `json:"analytics" bson:"analytics"`         // Analytics & feedback
	I9Shield      SecurityFeatures     `json:"i9Shield" bson:"i9Shield"`           // I9 Shield features
	Inventory     InventoryFeatures    `json:"inventory" bson:"inventory"`         // Inventory features
	AI            AIFeatures           `json:"ai" bson:"ai"`                       // AI & automation
}

// CoreFeatures - Core user experience
type CoreFeatures struct {
	IsHotspotEnabled      bool `json:"isHotspotEnabled" bson:"isHotspotEnabled"`
	IsUserKYCEnabled      bool `json:"isUserKYCEnabled" bson:"isUserKYCEnabled"`
	IsNotificationEnabled bool `json:"isNotificationEnabled" bson:"isNotificationEnabled"`
	IsUserPortalEnabled   bool `json:"isUserPortalEnabled" bson:"isUserPortalEnabled"`
	IsVasEnabled          bool `json:"isVasEnabled" bson:"isVasEnabled"`
	IsMarketingEnabled    bool `json:"isMarketingEnabled" bson:"isMarketingEnabled"`
	IsBillingEnabled      bool `json:"isBillingEnabled" bson:"isBillingEnabled"`
	IsOperationsEnabled   bool `json:"isOperationsEnabled" bson:"isOperationsEnabled"`
	IsAnalyticsEnabled    bool `json:"isAnalyticsEnabled" bson:"isAnalyticsEnabled"`
	IsI9ShieldEnabled     bool `json:"isI9ShieldEnabled" bson:"isI9ShieldEnabled"`
	IsAIEnabled           bool `json:"isAIEnabled" bson:"isAIEnabled"`
}

// InventoryFeatures - Inventory features
type InventoryFeatures struct {
	IsInventoryEnabled bool `json:"isInventoryEnabled" bson:"isInventoryEnabled"`
}

// NotificationFeatures - Notification channels
type NotificationFeatures struct {
	IsUserMailNotificationEnabled     bool `json:"isUserMailNotificationEnabled" bson:"isUserMailNotificationEnabled"`
	IsUserSMSNotificationEnabled      bool `json:"isUserSMSNotificationEnabled" bson:"isUserSMSNotificationEnabled"`
	IsUserWhatsappNotificationEnabled bool `json:"isUserWhatsappNotificationEnabled" bson:"isUserWhatsappNotificationEnabled"`
	IsUserTelegramNotificationEnabled bool `json:"isUserTelegramNotificationEnabled" bson:"isUserTelegramNotificationEnabled"`
	IsPushNotificationEnabled         bool `json:"isPushNotificationEnabled" bson:"isPushNotificationEnabled"`
}

// HotspotFeatures - Hotspot & captive portal
type HotspotFeatures struct { // Master hotspot toggle
	IsOTPAuthEnabled           bool `json:"isOtpAuthEnabled" bson:"isOtpAuthEnabled"`                     // OTP-based authentication
	IsSocialLoginEnabled       bool `json:"isSocialLoginEnabled" bson:"isSocialLoginEnabled"`             // Social media login (Google, Facebook)
	IsVoucherLoginEnabled      bool `json:"isVoucherLoginEnabled" bson:"isVoucherLoginEnabled"`           // Voucher/scratch card access
	IsPasswordLoginEnabled     bool `json:"isPasswordLoginEnabled" bson:"isPasswordLoginEnabled"`         // Password-based authentication
	IsGuestWiFiEnabled         bool `json:"isGuestWiFiEnabled" bson:"isGuestWiFiEnabled"`                 // Guest WiFi portal
	IsSessionManagementEnabled bool `json:"isSessionManagementEnabled" bson:"isSessionManagementEnabled"` // Session tracking & management
	IsCustomBrandingEnabled    bool `json:"isCustomBrandingEnabled" bson:"isCustomBrandingEnabled"`       // Custom portal branding
	IsAdvertisementEnabled     bool `json:"isAdvertisementEnabled" bson:"isAdvertisementEnabled"`         // Advertisement injection
	IsUsageLimitEnabled        bool `json:"isUsageLimitEnabled" bson:"isUsageLimitEnabled"`               // Data/time usage limits
	IsMultipleSSIDEnabled      bool `json:"isMultipleSSIDEnabled" bson:"isMultipleSSIDEnabled"`           // Multiple SSID support
}

// VASFeatures - Value-added services
type VASFeatures struct {
	IsIPTVEnabled            bool `json:"isIptvEnabled" bson:"isIptvEnabled"`
	IsOTTEnabled             bool `json:"isOttEnabled" bson:"isOttEnabled"`
	IsVoiceServiceEnabled    bool `json:"isVoiceServiceEnabled" bson:"isVoiceServiceEnabled"`
	IsVPNEnabled             bool `json:"isVpnEnabled" bson:"isVpnEnabled"`
	IsCloudBackupEnabled     bool `json:"isCloudBackupEnabled" bson:"isCloudBackupEnabled"`
	IsFirewallServiceEnabled bool `json:"isFirewallServiceEnabled" bson:"isFirewallServiceEnabled"`
	IsDNSSecurityEnabled     bool `json:"isDnsSecurityEnabled" bson:"isDnsSecurityEnabled"`
}

// MarketingFeatures - Marketing, campaigns & retention
type MarketingFeatures struct {
	IsCampaginServiceEnabled     bool `json:"isCampaginServiceEnabled" bson:"isCampaginServiceEnabled"`
	IsPromoCampaignEnabled       bool `json:"isPromoCampaignEnabled" bson:"isPromoCampaignEnabled"`
	IsEmailCampaignEnabled       bool `json:"isEmailCampaignEnabled" bson:"isEmailCampaignEnabled"`
	IsSMSCampaignEnabled         bool `json:"isSmsCampaignEnabled" bson:"isSmsCampaignEnabled"`
	IsMarketingAutomationEnabled bool `json:"isMarketingAutomationEnabled" bson:"isMarketingAutomationEnabled"`
	IsCouponsEnabled             bool `json:"isCouponsEnabled" bson:"isCouponsEnabled"`
	IsLoyaltyProgramEnabled      bool `json:"isLoyaltyProgramEnabled" bson:"isLoyaltyProgramEnabled"`
	IsReferralProgramEnabled     bool `json:"isReferralProgramEnabled" bson:"isReferralProgramEnabled"`
}

// BillingFeatures - Billing & collections
type BillingFeatures struct {
	IsOnlinePaymentEnabled       bool `json:"isOnlinePaymentEnabled" bson:"isOnlinePaymentEnabled"`
	IsMultiSplitPaymentEnabled   bool `json:"isMultiSplitPaymentEnabled" bson:"isMultiSplitPaymentEnabled"`
	IsInvoiceAutoReminderEnabled bool `json:"isInvoiceAutoReminderEnabled" bson:"isInvoiceAutoReminderEnabled"`
	IsPaymentSplitEnabled        bool `json:"isPaymentSplitEnabled" bson:"isPaymentSplitEnabled"`
}

// OperationsFeatures - Operations & integrations
type OperationsFeatures struct {
	IsTicketingEnabled      bool `json:"isTicketingEnabled" bson:"isTicketingEnabled"`
	IsCRMIntegrationEnabled bool `json:"isCrmIntegrationEnabled" bson:"isCrmIntegrationEnabled"`
	IsERPIntegrationEnabled bool `json:"isErpIntegrationEnabled" bson:"isErpIntegrationEnabled"`
	IsInventorySyncEnabled  bool `json:"isInventorySyncEnabled" bson:"isInventorySyncEnabled"`
	IsPartnerPortalEnabled  bool `json:"isPartnerPortalEnabled" bson:"isPartnerPortalEnabled"`
}

// AnalyticsFeatures - Analytics & feedback
type AnalyticsFeatures struct {
	IsUsageAnalyticsEnabled   bool `json:"isUsageAnalyticsEnabled" bson:"isUsageAnalyticsEnabled"`
	IsRevenueAnalyticsEnabled bool `json:"isRevenueAnalyticsEnabled" bson:"isRevenueAnalyticsEnabled"`
	IsNetworkAnalyticsEnabled bool `json:"isNetworkAnalyticsEnabled" bson:"isNetworkAnalyticsEnabled"`
	IsCustomerFeedbackEnabled bool `json:"isCustomerFeedbackEnabled" bson:"isCustomerFeedbackEnabled"`
	IsChurnPredictionEnabled  bool `json:"isChurnPredictionEnabled" bson:"isChurnPredictionEnabled"`
}

// SecurityFeatures - Security & compliance
type SecurityFeatures struct {
	Is2FAEnabled                 bool `json:"is2FAEnabled" bson:"is2FAEnabled"`
	IsAuditLoggingEnabled        bool `json:"isAuditLoggingEnabled" bson:"isAuditLoggingEnabled"`
	IsDataRetentionPolicyEnabled bool `json:"isDataRetentionPolicyEnabled" bson:"isDataRetentionPolicyEnabled"`
	IsGDPRComplianceEnabled      bool `json:"isGdprComplianceEnabled" bson:"isGdprComplianceEnabled"`
}

// AIFeatures - AI & automation
type AIFeatures struct {
	IsAIAssistantEnabled           bool `json:"isAiAssistantEnabled" bson:"isAiAssistantEnabled"`
	IsAIBasedSupportEnabled        bool `json:"isAiBasedSupportEnabled" bson:"isAiBasedSupportEnabled"`
	IsNetworkOptimizationAIEnabled bool `json:"isNetworkOptimizationAIEnabled" bson:"isNetworkOptimizationAIEnabled"`
	IsPredictiveMaintenanceEnabled bool `json:"isPredictiveMaintenanceEnabled" bson:"isPredictiveMaintenanceEnabled"`
}
