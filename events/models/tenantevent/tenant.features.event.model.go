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
	EnabledFeaturesFieldCoreIsAnalyticsEnabled    = "Core.IsAnalyticsEnabled"
	EnabledFeaturesFieldCoreIsI9ShieldEnabled     = "Core.IsI9ShieldEnabled"
	EnabledFeaturesFieldCoreIsAIEnabled           = "Core.IsAIEnabled"

	// Notification features
	EnabledFeaturesFieldNotificationsIsUserMailNotificationEnabled     = "Notifications.IsUserMailNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsUserSMSNotificationEnabled      = "Notifications.IsUserSMSNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsUserWhatsappNotificationEnabled = "Notifications.IsUserWhatsappNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsUserTelegramNotificationEnabled = "Notifications.IsUserTelegramNotificationEnabled"
	EnabledFeaturesFieldNotificationsIsPushNotificationEnabled         = "Notifications.IsPushNotificationEnabled"

	// Hotspot auth features (management toggles removed — handled at infra level)
	EnabledFeaturesFieldHotspotIsOTPAuthEnabled       = "Hotspot.IsOTPAuthEnabled"
	EnabledFeaturesFieldHotspotIsSocialLoginEnabled   = "Hotspot.IsSocialLoginEnabled"
	EnabledFeaturesFieldHotspotIsVoucherLoginEnabled  = "Hotspot.IsVoucherLoginEnabled"
	EnabledFeaturesFieldHotspotIsPasswordLoginEnabled = "Hotspot.IsPasswordLoginEnabled"

	// VAS features
	EnabledFeaturesFieldVASIsIPTVEnabled            = "VAS.IsIPTVEnabled"
	EnabledFeaturesFieldVASIsOTTEnabled             = "VAS.IsOTTEnabled"
	EnabledFeaturesFieldVASIsVoiceServiceEnabled    = "VAS.IsVoiceServiceEnabled"
	EnabledFeaturesFieldVASIsVPNEnabled             = "VAS.IsVPNEnabled"
	EnabledFeaturesFieldVASIsCloudBackupEnabled     = "VAS.IsCloudBackupEnabled"
	EnabledFeaturesFieldVASIsFirewallServiceEnabled = "VAS.IsFirewallServiceEnabled"
	EnabledFeaturesFieldVASIsDNSSecurityEnabled     = "VAS.IsDNSSecurityEnabled"

	// Marketing features (simplified — campaigns/automation removed)
	EnabledFeaturesFieldMarketingIsPromoCampaignEnabled   = "Marketing.IsPromoCampaignEnabled"
	EnabledFeaturesFieldMarketingIsCouponsEnabled         = "Marketing.IsCouponsEnabled"
	EnabledFeaturesFieldMarketingIsLoyaltyProgramEnabled  = "Marketing.IsLoyaltyProgramEnabled"
	EnabledFeaturesFieldMarketingIsReferralProgramEnabled = "Marketing.IsReferralProgramEnabled"

	// Billing features
	EnabledFeaturesFieldBillingIsOnlinePaymentEnabled       = "Billing.IsOnlinePaymentEnabled"
	EnabledFeaturesFieldBillingIsMultiSplitPaymentEnabled   = "Billing.IsMultiSplitPaymentEnabled"
	EnabledFeaturesFieldBillingIsInvoiceAutoReminderEnabled = "Billing.IsInvoiceAutoReminderEnabled"
	EnabledFeaturesFieldBillingIsPaymentSplitEnabled        = "Billing.IsPaymentSplitEnabled"

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
	Hotspot       HotspotFeatures      `json:"hotspot" bson:"hotspot"`             // Hotspot & captive portal (auth methods only)
	VAS           VASFeatures          `json:"vas" bson:"vas"`                     // Value-added services
	Marketing     MarketingFeatures    `json:"marketing" bson:"marketing"`         // Marketing & promotions
	Billing       BillingFeatures      `json:"billing" bson:"billing"`             // Billing & collections
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

// HotspotFeatures - Authentication methods only (management toggles removed)
type HotspotFeatures struct {
	IsOTPAuthEnabled       bool `json:"isOtpAuthEnabled" bson:"isOtpAuthEnabled"`
	IsSocialLoginEnabled   bool `json:"isSocialLoginEnabled" bson:"isSocialLoginEnabled"`
	IsVoucherLoginEnabled  bool `json:"isVoucherLoginEnabled" bson:"isVoucherLoginEnabled"`
	IsPasswordLoginEnabled bool `json:"isPasswordLoginEnabled" bson:"isPasswordLoginEnabled"`
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

// MarketingFeatures - Promotions & retention (campaigns/automation removed)
type MarketingFeatures struct {
	IsPromoCampaignEnabled   bool `json:"isPromoCampaignEnabled" bson:"isPromoCampaignEnabled"`
	IsCouponsEnabled         bool `json:"isCouponsEnabled" bson:"isCouponsEnabled"`
	IsLoyaltyProgramEnabled  bool `json:"isLoyaltyProgramEnabled" bson:"isLoyaltyProgramEnabled"`
	IsReferralProgramEnabled bool `json:"isReferralProgramEnabled" bson:"isReferralProgramEnabled"`
}

// BillingFeatures - Billing & collections
type BillingFeatures struct {
	IsOnlinePaymentEnabled       bool `json:"isOnlinePaymentEnabled" bson:"isOnlinePaymentEnabled"`
	IsMultiSplitPaymentEnabled   bool `json:"isMultiSplitPaymentEnabled" bson:"isMultiSplitPaymentEnabled"`
	IsInvoiceAutoReminderEnabled bool `json:"isInvoiceAutoReminderEnabled" bson:"isInvoiceAutoReminderEnabled"`
	IsPaymentSplitEnabled        bool `json:"isPaymentSplitEnabled" bson:"isPaymentSplitEnabled"`
}

// AnalyticsFeatures - Analytics & feedback
type AnalyticsFeatures struct {
	IsUsageAnalyticsEnabled   bool `json:"isUsageAnalyticsEnabled" bson:"isUsageAnalyticsEnabled"`
	IsRevenueAnalyticsEnabled bool `json:"isRevenueAnalyticsEnabled" bson:"isRevenueAnalyticsEnabled"`
	IsNetworkAnalyticsEnabled bool `json:"isNetworkAnalyticsEnabled" bson:"isNetworkAnalyticsEnabled"`
	IsCustomerFeedbackEnabled bool `json:"isCustomerFeedbackEnabled" bson:"isCustomerFeedbackEnabled"`
	IsChurnPredictionEnabled  bool `json:"isChurnPredictionEnabled" bson:"isChurnPredictionEnabled"`
}

// SecurityFeatures - Security & compliance (I9 Shield)
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
