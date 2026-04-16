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

	EnabledFeaturesFieldNotificationsIsPushNotificationEnabled         = "Notifications.IsPushNotificationEnabled"

	// Hotspot auth features (management toggles removed — handled at infra level)
	EnabledFeaturesFieldHotspotIsOTPAuthEnabled       = "Hotspot.IsOTPAuthEnabled"
	EnabledFeaturesFieldHotspotIsSocialLoginEnabled   = "Hotspot.IsSocialLoginEnabled"
	EnabledFeaturesFieldHotspotIsVoucherLoginEnabled  = "Hotspot.IsVoucherLoginEnabled"
	EnabledFeaturesFieldHotspotIsPasswordLoginEnabled = "Hotspot.IsPasswordLoginEnabled"
	EnabledFeaturesFieldHotspotIsInAppBannerEnabled   = "Hotspot.IsInAppBannerEnabled"

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
	EnabledFeaturesFieldI9ShieldIsI9ShieldEnabled            = "I9Shield.IsI9ShieldEnabled"

	// AI features
	EnabledFeaturesFieldAIIsAIAssistantEnabled           = "AI.IsAIAssistantEnabled"
	EnabledFeaturesFieldAIIsAIBasedSupportEnabled        = "AI.IsAIBasedSupportEnabled"
	EnabledFeaturesFieldAIIsNetworkOptimizationAIEnabled = "AI.IsNetworkOptimizationAIEnabled"
	EnabledFeaturesFieldAIIsPredictiveMaintenanceEnabled = "AI.IsPredictiveMaintenanceEnabled"

	// Venue features
	EnabledFeaturesFieldVenueIsVenueEnabled           = "Venue.IsVenueEnabled"
	EnabledFeaturesFieldVenueIsMenuManagementEnabled  = "Venue.IsMenuManagementEnabled"
	EnabledFeaturesFieldVenueIsTableManagementEnabled = "Venue.IsTableManagementEnabled"
	EnabledFeaturesFieldVenueIsZoneManagementEnabled  = "Venue.IsZoneManagementEnabled"
	EnabledFeaturesFieldVenueIsDiscountsEnabled       = "Venue.IsDiscountsEnabled"

	// ACS features
	EnabledFeaturesFieldACSIsACSEnabled             = "ACS.IsACSEnabled"
	EnabledFeaturesFieldACSIsFleetManagementEnabled = "ACS.IsFleetManagementEnabled"
	EnabledFeaturesFieldACSIsRemoteConfigEnabled    = "ACS.IsRemoteConfigEnabled"

	// Network features
	EnabledFeaturesFieldNetworkIsMonitoringEnabled    = "Network.IsMonitoringEnabled"
	EnabledFeaturesFieldNetworkIsTopologyEnabled      = "Network.IsTopologyEnabled"
	EnabledFeaturesFieldNetworkIsOLTManagementEnabled = "Network.IsOLTManagementEnabled"

	// Operations features
	EnabledFeaturesFieldOperationsIsLogEngineEnabled   = "Operations.IsLogEngineEnabled" // Deprecated: moved to Network
	EnabledFeaturesFieldNetworkIsLogEngineEnabled      = "Network.IsLogEngineEnabled"
	EnabledFeaturesFieldOperationsIsAlertsEnabled      = "Operations.IsAlertsEnabled"
	EnabledFeaturesFieldOperationsIsTicketingEnabled   = "Operations.IsTicketingEnabled"
	EnabledFeaturesFieldOperationsIsSubscribersEnabled = "Operations.IsSubscribersEnabled"

	// Catalog features
	EnabledFeaturesFieldCatalogIsProductsEnabled = "Catalog.IsProductsEnabled"
	EnabledFeaturesFieldCatalogIsPlansEnabled    = "Catalog.IsPlansEnabled"

	// Integrations features
	EnabledFeaturesFieldIntegrationsIsIntegrationsEnabled          = "Integrations.IsIntegrationsEnabled"
	EnabledFeaturesFieldIntegrationsIsNotificationTemplatesEnabled = "Integrations.IsNotificationTemplatesEnabled"
	EnabledFeaturesFieldIntegrationsIsSMSResultsEnabled            = "Integrations.IsSMSResultsEnabled"
	EnabledFeaturesFieldIntegrationsIsEmailResultsEnabled          = "Integrations.IsEmailResultsEnabled"

	// System Management features
	EnabledFeaturesFieldSystemIsKubernetesEnabled    = "System.IsKubernetesEnabled"
	EnabledFeaturesFieldSystemIsRolesEnabled         = "System.IsRolesEnabled"
	EnabledFeaturesFieldSystemIsPoliciesEnabled      = "System.IsPoliciesEnabled"
	EnabledFeaturesFieldSystemIsPincodeUploadEnabled = "System.IsPincodeUploadEnabled"
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
	Venue         VenueFeatures        `json:"venue" bson:"venue"`                 // Venue management (menus, tables, zones)
	ACS           ACSFeatures          `json:"acs" bson:"acs"`                     // ACS / CWMP / CPE device management
	Network       NetworkFeatures      `json:"network" bson:"network"`             // Network monitoring & topology
	Operations    OperationsFeatures   `json:"operations" bson:"operations"`       // Operational modules (log engine, alerts, ticketing)
	Catalog       CatalogFeatures      `json:"catalog" bson:"catalog"`             // Product & plan catalog
	Integrations  IntegrationsFeatures `json:"integrations" bson:"integrations"`   // Third-party integrations
	System        SystemFeatures       `json:"system" bson:"system"`               // System management panels
	// Note: KYC policies (KYCFeatures) live on TenantModel.KycPolicy, NOT here.
	// EnabledFeatures only holds boolean on/off toggles.
	// The master toggle Core.IsUserKYCEnabled controls whether KYC is active.
	// KYC gateway config lives in TenantKYCProviderBinding (separate collection).
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

	IsPushNotificationEnabled         bool `json:"isPushNotificationEnabled" bson:"isPushNotificationEnabled"`
}

// HotspotFeatures - Authentication methods only (management toggles removed)
type HotspotFeatures struct {
	IsOTPAuthEnabled       bool `json:"isOtpAuthEnabled" bson:"isOtpAuthEnabled"`
	IsSocialLoginEnabled   bool `json:"isSocialLoginEnabled" bson:"isSocialLoginEnabled"`
	IsVoucherLoginEnabled  bool `json:"isVoucherLoginEnabled" bson:"isVoucherLoginEnabled"`
	IsPasswordLoginEnabled bool `json:"isPasswordLoginEnabled" bson:"isPasswordLoginEnabled"`
	IsInAppBannerEnabled   bool `json:"isInAppBannerEnabled" bson:"isInAppBannerEnabled"`
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
	IsI9ShieldEnabled            bool `json:"isI9ShieldEnabled" bson:"isI9ShieldEnabled"`
}

// AIFeatures - AI & automation
type AIFeatures struct {
	IsAIAssistantEnabled           bool `json:"isAiAssistantEnabled" bson:"isAiAssistantEnabled"`
	IsAIBasedSupportEnabled        bool `json:"isAiBasedSupportEnabled" bson:"isAiBasedSupportEnabled"`
	IsNetworkOptimizationAIEnabled bool `json:"isNetworkOptimizationAIEnabled" bson:"isNetworkOptimizationAIEnabled"`
	IsPredictiveMaintenanceEnabled bool `json:"isPredictiveMaintenanceEnabled" bson:"isPredictiveMaintenanceEnabled"`
}

// VenueFeatures - Venue management (menus, tables, zones, discounts)
type VenueFeatures struct {
	IsVenueEnabled           bool `json:"isVenueEnabled" bson:"isVenueEnabled"`
	IsMenuManagementEnabled  bool `json:"isMenuManagementEnabled" bson:"isMenuManagementEnabled"`
	IsTableManagementEnabled bool `json:"isTableManagementEnabled" bson:"isTableManagementEnabled"`
	IsZoneManagementEnabled  bool `json:"isZoneManagementEnabled" bson:"isZoneManagementEnabled"`
	IsDiscountsEnabled       bool `json:"isDiscountsEnabled" bson:"isDiscountsEnabled"`
}

// ACSFeatures - ACS / CWMP / CPE device management
type ACSFeatures struct {
	IsACSEnabled             bool `json:"isAcsEnabled" bson:"isAcsEnabled"`
	IsFleetManagementEnabled bool `json:"isFleetManagementEnabled" bson:"isFleetManagementEnabled"`
	IsRemoteConfigEnabled    bool `json:"isRemoteConfigEnabled" bson:"isRemoteConfigEnabled"`
}

// NetworkFeatures - Network monitoring & topology
type NetworkFeatures struct {
	IsMonitoringEnabled    bool `json:"isMonitoringEnabled" bson:"isMonitoringEnabled"`
	IsTopologyEnabled      bool `json:"isTopologyEnabled" bson:"isTopologyEnabled"`
	IsOLTManagementEnabled bool `json:"isOltManagementEnabled" bson:"isOltManagementEnabled"`
	IsLogEngineEnabled     bool `json:"isLogEngineEnabled" bson:"isLogEngineEnabled"`
}

// OperationsFeatures - Operational modules
type OperationsFeatures struct {
	IsAlertsEnabled      bool `json:"isAlertsEnabled" bson:"isAlertsEnabled"`
	IsTicketingEnabled   bool `json:"isTicketingEnabled" bson:"isTicketingEnabled"`
	IsSubscribersEnabled bool `json:"isSubscribersEnabled" bson:"isSubscribersEnabled"`
}

// CatalogFeatures - Product & plan catalog
type CatalogFeatures struct {
	IsProductsEnabled bool `json:"isProductsEnabled" bson:"isProductsEnabled"`
	IsPlansEnabled    bool `json:"isPlansEnabled" bson:"isPlansEnabled"`
}

// IntegrationsFeatures - Third-party integration modules
type IntegrationsFeatures struct {
	IsIntegrationsEnabled          bool `json:"isIntegrationsEnabled" bson:"isIntegrationsEnabled"`
	IsNotificationTemplatesEnabled bool `json:"isNotificationTemplatesEnabled" bson:"isNotificationTemplatesEnabled"`
	IsSMSResultsEnabled            bool `json:"isSmsResultsEnabled" bson:"isSmsResultsEnabled"`
	IsEmailResultsEnabled          bool `json:"isEmailResultsEnabled" bson:"isEmailResultsEnabled"`
}

// SystemFeatures - System management panel visibility controls
type SystemFeatures struct {
	IsKubernetesEnabled    bool `json:"isKubernetesEnabled" bson:"isKubernetesEnabled"`
	IsRolesEnabled         bool `json:"isRolesEnabled" bson:"isRolesEnabled"`
	IsPoliciesEnabled      bool `json:"isPoliciesEnabled" bson:"isPoliciesEnabled"`
	IsPincodeUploadEnabled bool `json:"isPincodeUploadEnabled" bson:"isPincodeUploadEnabled"`
}
