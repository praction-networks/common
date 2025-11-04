package tenantevent

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TenantInsertEventModel struct {
	ID               string                         `bson:"_id" json:"id"`
	Name             string                         `bson:"name" json:"name"`
	Code             string                         `bson:"code" json:"code"`
	Type             string                         `bson:"type" json:"type"`
	Fqdn             string                         `bson:"fqdn,omitempty" json:"fqdn,omitempty"`
	Environment      string                         `bson:"environment,omitempty" json:"environment,omitempty"`
	EntType          string                         `bson:"entType,omitempty" json:"entType,omitempty"`
	ParentTenantID   string                         `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"`
	DefaultEmail     string                         `bson:"defaultEmail" json:"defaultEmail"`
	DefaultPhone     string                         `bson:"defaultPhone" json:"defaultPhone"`
	PermanentAddress AddressModel                   `bson:"permanentAddress" json:"permanentAddress"`
	CurrentAddress   AddressModel                   `bson:"currentAddress" json:"currentAddress"`
	TenantGST        []GSTModel                     `bson:"tenantGST" json:"tenantGST"`
	TenantPAN        PANModel                       `bson:"tenantPAN" json:"tenantPAN"`
	TenantTAN        TANModel                       `bson:"tenantTAN" json:"tenantTAN"`
	TenantCIN        CINModel                       `bson:"tenantCIN" json:"tenantCIN"`
	EnabledFeatures  EnabledFeatures                `bson:"enabledFeatures,omitempty" json:"enabledFeatures,omitempty"`
	OLTs             []string                       `bson:"olts,omitempty" json:"olts,omitempty"`
	KYCProvider      ProvidersModel                 `bson:"kycProvider,omitempty" json:"kycProvider,omitempty"`
	PaymentGateway   ProvidersModel                 `bson:"paymentGateway,omitempty" json:"paymentGateway,omitempty"`
	SMSProvider      ProvidersModel                 `bson:"smsProvider,omitempty" json:"smsProvider,omitempty"`
	MailProvider     ProvidersModel                 `bson:"mailProvider,omitempty" json:"mailProvider,omitempty"`
	AppsMessanger    []AppMessagingProvidersModel   `bson:"appsMessanger,omitempty" json:"appsMessanger,omitempty"`
	ExternalRadius   []ExternalRadiusProvidersModel `bson:"externalRadius,omitempty" json:"externalRadius,omitempty"`
	IsActive         bool                           `bson:"isActive" json:"isActive"`
	Version          int                            `bson:"version,omitempty" json:"version,omitempty"`
}

type TenantUpdateEventModel struct {
	ID               string                         `bson:"_id" json:"id"`
	Name             string                         `bson:"name" json:"name"`
	Code             string                         `bson:"code" json:"code"`
	Type             string                         `bson:"type" json:"type"`
	Fqdn             string                         `bson:"fqdn,omitempty" json:"fqdn,omitempty"`
	Environment      string                         `bson:"environment,omitempty" json:"environment,omitempty"`
	EntType          string                         `bson:"entType,omitempty" json:"entType,omitempty"`
	ParentTenantID   string                         `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"`
	ChildIDs         []string                       `bson:"childIDs,omitempty" json:"childIDs,omitempty"`
	DefaultEmail     string                         `bson:"defaultEmail" json:"defaultEmail"`
	DefaultPhone     string                         `bson:"defaultPhone" json:"defaultPhone"`
	PermanentAddress AddressModel                   `bson:"permanentAddress" json:"permanentAddress"`
	CurrentAddress   AddressModel                   `bson:"currentAddress" json:"currentAddress"`
	TenantGST        []GSTModel                     `bson:"tenantGST" json:"tenantGST"`
	TenantPAN        PANModel                       `bson:"tenantPAN" json:"tenantPAN"`
	TenantTAN        TANModel                       `bson:"tenantTAN" json:"tenantTAN"`
	TenantCIN        CINModel                       `bson:"tenantCIN" json:"tenantCIN"`
	EnabledFeatures  EnabledFeatures                `bson:"enabledFeatures,omitempty" json:"enabledFeatures,omitempty"`
	OLTs             []string                       `bson:"olts,omitempty" json:"olts,omitempty"`
	KYCProvider      ProvidersModel                 `bson:"kycProvider,omitempty" json:"kycProvider,omitempty"`
	PaymentGateway   ProvidersModel                 `bson:"paymentGateway,omitempty" json:"paymentGateway,omitempty"`
	SMSProvider      ProvidersModel                 `bson:"smsProvider,omitempty" json:"smsProvider,omitempty"`
	MailProvider     ProvidersModel                 `bson:"mailProvider,omitempty" json:"mailProvider,omitempty"`
	AppsMessanger    []AppMessagingProvidersModel   `bson:"appsMessanger,omitempty" json:"appsMessanger,omitempty"`
	ExternalRadius   []ExternalRadiusProvidersModel `bson:"externalRadius,omitempty" json:"externalRadius,omitempty"`
	IsActive         bool                           `bson:"isActive" json:"isActive"`
	Version          int                            `bson:"version,omitempty" json:"version,omitempty"`
}

type GSTModel struct {
	State      string    `bson:"state,omitempty" json:"state,omitempty"`
	GSTIN      string    `bson:"gstin,omitempty" json:"gstin,omitempty"`
	IsVerified bool      `bson:"isVerified,omitempty" json:"isVerified,omitempty"`
	VerifiedAt time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
}

type TenantDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}

type ProvidersModel struct {
	DefaultProviderID string   `bson:"defaultProviderID,omitempty" json:"defaultProviderID,omitempty"`
	Providers         []string `bson:"providers,omitempty" json:"providers,omitempty"`
}

type AppMessagingProvidersModel struct {
	MessageProviderID string `bson:"messageProviderID,omitempty" json:"messageProviderID,omitempty"`
	MessageProvider   string `bson:"messageProvider,omitempty" json:"messageProvider,omitempty"`
}

type ExternalRadiusProvidersModel struct {
	ExternalRadiusProviderID string   `bson:"externalRadiusProviderID,omitempty" json:"externalRadiusProviderID,omitempty"`
	ExternalRadiusProviders  []string `bson:"externalRadiusProviders,omitempty" json:"externalRadiusProviders,omitempty"`
}

type PANModel struct {
	PAN        string    `bson:"pan,omitempty" json:"pan,omitempty"`
	IsVerified bool      `bson:"isVerified,omitempty" json:"isVerified,omitempty"`
	VerifiedAt time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
}

type TANModel struct {
	TAN        string    `bson:"tan,omitempty" json:"tan,omitempty"`
	IsVerified bool      `bson:"isVerified,omitempty" json:"isVerified,omitempty"`
	VerifiedAt time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
}

type CINModel struct {
	CIN        string    `bson:"cin,omitempty" json:"cin,omitempty"`
	IsVerified bool      `bson:"isVerified,omitempty" json:"isVerified,omitempty"`
	VerifiedAt time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
}

type NotificationGateways struct {
	SMS      []string `bson:"sms,omitempty" json:"sms,omitempty"`
	Mail     []string `bson:"mail,omitempty" json:"mail,omitempty"`
	WhatsApp []string `bson:"whatsapp,omitempty" json:"whatsapp,omitempty"`
	Telegram []string `bson:"telegram,omitempty" json:"telegram,omitempty"`
}

type ExternalRadiusSettings struct {
	JazeID  string `bson:"jazeEnabled,omitempty" json:"jazeEnabled,omitempty"`
	IPACTID string `bson:"ipactEnabled,omitempty" json:"ipactEnabled,omitempty"`
}

// EnabledFeatures defines all feature toggles for a Domain/Tenant.
// NOTE: Flat structure maintained for backward-compatibility with existing JSON/BSON.
type EnabledFeatures struct {
	// ===== Core User Experience =====
	IsUserPortalEnabled    bool `bson:"isUserPortalEnabled" json:"isUserPortalEnabled"`       // Enable user portal
	IsUserKYCEnabled       bool `bson:"isUserKYCEnabled" json:"isUserKYCEnabled"`             // Enable KYC for users
	IsOnlinePaymentEnabled bool `bson:"isOnlinePaymentEnabled" json:"isOnlinePaymentEnabled"` // Enable online payment

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
