package tenantevent

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TenantInsertEventModel struct {
	ID               string                         `json:"id"`
	Name             string                         `json:"name"`
	Code             string                         `json:"code"`
	Type             string                         `json:"type"`
	Fqdn             string                         `json:"fqdn,omitempty"`
	Environment      string                         `json:"environment,omitempty"`
	EntType          string                         `json:"entType,omitempty"`
	ParentTenantID   string                         `json:"parentTenantID,omitempty"`
	DefaultEmail     string                         `json:"defaultEmail"`
	DefaultPhone     string                         `json:"defaultPhone"`
	PermanentAddress AddressModel                   `json:"permanentAddress"`
	CurrentAddress   AddressModel                   `json:"currentAddress"`
	TenantGST        []GSTModel                     `json:"tenantGST"`
	TenantPAN        PANModel                       `json:"tenantPAN"`
	TenantTAN        TANModel                       `json:"tenantTAN"`
	TenantCIN        CINModel                       `json:"tenantCIN"`
	EnabledFeatures  EnabledFeatures                `json:"enabledFeatures,omitempty"`
	OLTs             []string                       `json:"olts,omitempty"`
	KYCProvider      ProvidersModel                 `json:"kycProvider,omitempty"`
	PaymentGateway   ProvidersModel                 `json:"paymentGateway,omitempty"`
	SMSProvider      ProvidersModel                 `json:"smsProvider,omitempty"`
	MailProvider     ProvidersModel                 `json:"mailProvider,omitempty"`
	AppsMessanger    []AppMessagingProvidersModel   `json:"appsMessanger,omitempty"`
	ExternalRadius   []ExternalRadiusProvidersModel `json:"externalRadius,omitempty"`
	IsActive         bool                           `json:"isActive"`
	Version          int                            `json:"version,omitempty"`
}

type TenantUpdateEventModel struct {
	ID               string                         `json:"id"`
	Name             string                         `json:"name"`
	Code             string                         `json:"code"`
	Type             string                         `json:"type"`
	Fqdn             string                         `json:"fqdn,omitempty"`
	Environment      string                         `json:"environment,omitempty"`
	EntType          string                         `json:"entType,omitempty"`
	ParentTenantID   string                         `json:"parentTenantID,omitempty"`
	ChildIDs         []string                       `json:"childIDs,omitempty"`
	DefaultEmail     string                         `json:"defaultEmail"`
	DefaultPhone     string                         `json:"defaultPhone"`
	PermanentAddress AddressModel                   `json:"permanentAddress"`
	CurrentAddress   AddressModel                   `json:"currentAddress"`
	TenantGST        []GSTModel                     `json:"tenantGST"`
	TenantPAN        PANModel                       `json:"tenantPAN"`
	TenantTAN        TANModel                       `json:"tenantTAN"`
	TenantCIN        CINModel                       `json:"tenantCIN"`
	EnabledFeatures  EnabledFeatures                `json:"enabledFeatures,omitempty"`
	OLTs             []string                       `json:"olts,omitempty"`
	KYCProvider      ProvidersModel                 `json:"kycProvider,omitempty"`
	PaymentGateway   ProvidersModel                 `json:"paymentGateway,omitempty"`
	SMSProvider      ProvidersModel                 `json:"smsProvider,omitempty"`
	MailProvider     ProvidersModel                 `json:"mailProvider,omitempty"`
	AppsMessanger    []AppMessagingProvidersModel   `json:"appsMessanger,omitempty"`
	ExternalRadius   []ExternalRadiusProvidersModel `json:"externalRadius,omitempty"`
	IsActive         bool                           `json:"isActive"`
	Version          int                            `json:"version,omitempty"`
}

type GSTModel struct {
	State      string    `json:"state,omitempty"`
	GSTIN      string    `json:"gstin,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type TenantDeleteEventModel struct {
	ID string `json:"id"`
}

type ProvidersModel struct {
	DefaultProviderID string   `json:"defaultProviderID,omitempty"`
	Providers         []string `json:"providers,omitempty"`
}

type AppMessagingProvidersModel struct {
	MessageProviderID string `json:"messageProviderID,omitempty"`
	MessageProvider   string `json:"messageProvider,omitempty"`
}

type ExternalRadiusProvidersModel struct {
	ExternalRadiusProviderID string   `json:"externalRadiusProviderID,omitempty"`
	ExternalRadiusProviders  []string `json:"externalRadiusProviders,omitempty"`
}

type PANModel struct {
	PAN        string    `json:"pan,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type TANModel struct {
	TAN        string    `json:"tan,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type CINModel struct {
	CIN        string    `json:"cin,omitempty"`
	IsVerified bool      `json:"isVerified,omitempty"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type NotificationGateways struct {
	SMS      []string `json:"sms,omitempty"`
	Mail     []string `json:"mail,omitempty"`
	WhatsApp []string `json:"whatsapp,omitempty"`
	Telegram []string `json:"telegram,omitempty"`
}

type ExternalRadiusSettings struct {
	JazeID  string `json:"jazeEnabled,omitempty"`
	IPACTID string `json:"ipactEnabled,omitempty"`
}

// EnabledFeatures defines all feature toggles for a Domain/Tenant.
// NOTE: Flat structure maintained for backward-compatibility with existing JSON/BSON.
type EnabledFeatures struct {
	// ===== Core User Experience =====
	IsUserPortalEnabled    bool `json:"isUserPortalEnabled" bson:"isUserPortalEnabled"`       // Enable user portal
	IsUserKYCEnabled       bool `json:"isUserKYCEnabled" bson:"isUserKYCEnabled"`             // Enable KYC for users
	IsOnlinePaymentEnabled bool `json:"isOnlinePaymentEnabled" bson:"isOnlinePaymentEnabled"` // Enable online payment

	// ===== Notifications & Messaging =====
	IsUserMailNotificationEnabled     bool `json:"isUserMailNotificationEnabled" bson:"isUserMailNotificationEnabled"`         // Email notifications
	IsUserSMSNotificationEnabled      bool `json:"isUserSMSNotificationEnabled" bson:"isUserSMSNotificationEnabled"`           // SMS notifications
	IsUserWhatsappNotificationEnabled bool `json:"isUserWhatsappNotificationEnabled" bson:"isUserWhatsappNotificationEnabled"` // WhatsApp notifications
	IsUserTelegramNotificationEnabled bool `json:"isUserTelegramNotificationEnabled" bson:"isUserTelegramNotificationEnabled"` // Telegram notifications
	IsPushNotificationEnabled         bool `json:"isPushNotificationEnabled" bson:"isPushNotificationEnabled"`                 // Mobile push notifications
	IsInAppBannerEnabled              bool `json:"isInAppBannerEnabled" bson:"isInAppBannerEnabled"`                           // In-app banner/announcement bar

	// ===== Access / RADIUS Providers =====
	IsJazeeraRadiusProviderEnabled bool `json:"isJazeeraRadiusProviderEnabled" bson:"isJazeeraRadiusProviderEnabled"` // Jazeera RADIUS
	IsIPACTRadiusProviderEnabled   bool `json:"isIPACTRadiusProviderEnabled" bson:"isIPACTRadiusProviderEnabled"`     // IPACT RADIUS
	IsFreeRadiusProviderEnabled    bool `json:"isFreeRadiusProviderEnabled" bson:"isFreeRadiusProviderEnabled"`       // FreeRADIUS

	// ===== Value-Added Services (VAS) =====
	IsIPTVEnabled            bool `json:"isIPTVEnabled" bson:"isIPTVEnabled"`                       // IPTV
	IsOTTEnabled             bool `json:"isOTTEnabled" bson:"isOTTEnabled"`                         // OTT bundles
	IsVoiceServiceEnabled    bool `json:"isVoiceServiceEnabled" bson:"isVoiceServiceEnabled"`       // VoIP/Voice
	IsVPNEnabled             bool `json:"isVPNEnabled" bson:"isVPNEnabled"`                         // Managed VPN
	IsCloudBackupEnabled     bool `json:"isCloudBackupEnabled" bson:"isCloudBackupEnabled"`         // Cloud backup for users
	IsFirewallServiceEnabled bool `json:"isFirewallServiceEnabled" bson:"isFirewallServiceEnabled"` // Managed firewall
	IsDNSFilteringEnabled    bool `json:"isDNSFilteringEnabled" bson:"isDNSFilteringEnabled"`       // DNS filtering / parental control
	IsHotspotLoginEnabled    bool `json:"isHotspotLoginEnabled" bson:"isHotspotLoginEnabled"`       // Captive portal / hotspot
	IsRoamingEnabled         bool `json:"isRoamingEnabled" bson:"isRoamingEnabled"`                 // Wi-Fi roaming across zones

	// ===== Marketing, Campaigns & Retention =====
	IsCampaginServiceEnabled     bool `json:"isCampaginServiceEnabled" bson:"isCampaginServiceEnabled"`         // (legacy) Campaign service
	IsPromoCampaignEnabled       bool `json:"isPromoCampaignEnabled" bson:"isPromoCampaignEnabled"`             // Promotional campaigns (generic)
	IsEmailCampaignEnabled       bool `json:"isEmailCampaignEnabled" bson:"isEmailCampaignEnabled"`             // Email campaigns
	IsSMSCampaignEnabled         bool `json:"isSMSCampaignEnabled" bson:"isSMSCampaignEnabled"`                 // SMS campaigns
	IsMarketingAutomationEnabled bool `json:"isMarketingAutomationEnabled" bson:"isMarketingAutomationEnabled"` // Marketing automation workflows
	IsCouponsEnabled             bool `json:"isCouponsEnabled" bson:"isCouponsEnabled"`                         // Discount coupons
	IsLoyaltyProgramEnabled      bool `json:"isLoyaltyProgramEnabled" bson:"isLoyaltyProgramEnabled"`           // Loyalty/reward points
	IsReferralProgramEnabled     bool `json:"isReferralProgramEnabled" bson:"isReferralProgramEnabled"`         // Referral program

	// ===== Billing & Collections =====
	IsInvoiceAutoReminderEnabled bool `json:"isInvoiceAutoReminderEnabled" bson:"isInvoiceAutoReminderEnabled"` // Auto invoice/payment reminders
	IsPaymentSplitEnabled        bool `json:"isPaymentSplitEnabled" bson:"isPaymentSplitEnabled"`               // Split payments (partners/vendors)

	// ===== Operations & Integrations =====
	IsTicketingEnabled      bool `json:"isTicketingEnabled" bson:"isTicketingEnabled"`           // Helpdesk/ticketing
	IsCRMIntegrationEnabled bool `json:"isCRMIntegrationEnabled" bson:"isCRMIntegrationEnabled"` // CRM integration (Zoho/HubSpot/etc.)
	IsERPIntegrationEnabled bool `json:"isERPIntegrationEnabled" bson:"isERPIntegrationEnabled"` // ERP integration
	IsInventorySyncEnabled  bool `json:"isInventorySyncEnabled" bson:"isInventorySyncEnabled"`   // Inventory sync with ERP
	IsPartnerPortalEnabled  bool `json:"isPartnerPortalEnabled" bson:"isPartnerPortalEnabled"`   // Partner/reseller portal

	// ===== Analytics & Feedback =====
	IsUsageAnalyticsEnabled   bool `json:"isUsageAnalyticsEnabled" bson:"isUsageAnalyticsEnabled"`     // User behavior/usage analytics
	IsRevenueAnalyticsEnabled bool `json:"isRevenueAnalyticsEnabled" bson:"isRevenueAnalyticsEnabled"` // Revenue analytics
	IsNetworkAnalyticsEnabled bool `json:"isNetworkAnalyticsEnabled" bson:"isNetworkAnalyticsEnabled"` // Network performance analytics
	IsCustomerFeedbackEnabled bool `json:"isCustomerFeedbackEnabled" bson:"isCustomerFeedbackEnabled"` // Surveys/feedback
	IsChurnPredictionEnabled  bool `json:"isChurnPredictionEnabled" bson:"isChurnPredictionEnabled"`   // AI churn prediction

	// ===== Security & Compliance =====
	Is2FAEnabled                 bool `json:"is2FAEnabled" bson:"is2FAEnabled"`                                 // Two-factor auth (2FA)
	IsAuditLoggingEnabled        bool `json:"isAuditLoggingEnabled" bson:"isAuditLoggingEnabled"`               // Audit logs
	IsDataRetentionPolicyEnabled bool `json:"isDataRetentionPolicyEnabled" bson:"isDataRetentionPolicyEnabled"` // Data retention enforcement
	IsGDPRComplianceEnabled      bool `json:"isGDPRComplianceEnabled" bson:"isGDPRComplianceEnabled"`           // GDPR/DPDP controls (toggle gates UI/flows)

	// ===== AI & Automation =====
	IsAIAssistantEnabled           bool `json:"isAIAssistantEnabled" bson:"isAIAssistantEnabled"`                     // AI assistant in portal/app
	IsAIBasedSupportEnabled        bool `json:"isAIBasedSupportEnabled" bson:"isAIBasedSupportEnabled"`               // AI chat/agent for support
	IsNetworkOptimizationAIEnabled bool `json:"isNetworkOptimizationAIEnabled" bson:"isNetworkOptimizationAIEnabled"` // AI-driven network tuning
	IsPredictiveMaintenanceEnabled bool `json:"isPredictiveMaintenanceEnabled" bson:"isPredictiveMaintenanceEnabled"` // Predictive maintenance alerts
}

type ISPSettings struct {
	Plans            []string `json:"plans,omitempty" bson:"plans,omitempty"`
	MaxBandwidthMbps int      `json:"maxBandwidthMbps,omitempty" bson:"maxBandwidthMbps,omitempty"`
	IPPoolCIDR       string   `json:"ipPoolCIDR,omitempty" bson:"ipPoolCIDR,omitempty"`
	CoverageArea     bson.M   `json:"coverageArea,omitempty" bson:"coverageArea,omitempty"`
	Latitude         string   `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude        string   `json:"longitude,omitempty" bson:"longitude,omitempty"`
	BillingCycle     string   `json:"billingCycle,omitempty" bson:"billingCycle,omitempty"`
	AutoRenewal      bool     `json:"autoRenewal,omitempty" bson:"autoRenewal,omitempty"`
	SupportContact   string   `json:"supportContact,omitempty" bson:"supportContact,omitempty"`
	AssignedRM       string   `json:"assignedRM,omitempty" bson:"assignedRM,omitempty"`
	DeviceIDs        []string `json:"deviceIDs,omitempty" bson:"deviceIDs,omitempty"`
}

type TenantUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type TenantDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}
