package tenantevent

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TenantInsertEventModel struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Code        string `bson:"code" json:"code"`
	Type        string `bson:"type" json:"type"`
	Fqdn        string `bson:"fqdn,omitempty" json:"fqdn,omitempty"`
	Environment string `bson:"environment,omitempty" json:"environment,omitempty"`
	EntType     string `bson:"entType,omitempty" json:"entType,omitempty"`

	// Hierarchy Management
	ParentTenantID string   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"`
	ChildIDs       []string `bson:"childIds,omitempty" json:"childIds,omitempty"`
	Path           string   `bson:"path" json:"path"`
	PathDepth      int      `bson:"pathDepth" json:"pathDepth"`
	Ancestors      []string `bson:"ancestors,omitempty" json:"ancestors,omitempty"`
	Level          int      `bson:"level" json:"level"`
	IsLeaf         bool     `bson:"isLeaf" json:"isLeaf"`
	ChildrenCount  int      `bson:"childrenCount" json:"childrenCount"`

	// Hierarchy Constraints
	MaxDepth          int      `bson:"maxDepth,omitempty" json:"maxDepth,omitempty"`
	AllowedChildTypes []string `bson:"allowedChildTypes,omitempty" json:"allowedChildTypes,omitempty"`

	DefaultEmail     string          `bson:"defaultEmail" json:"defaultEmail"`
	DefaultPhone     string          `bson:"defaultPhone" json:"defaultPhone"`
	PermanentAddress AddressModel    `bson:"permanentAddress" json:"permanentAddress"`
	CurrentAddress   AddressModel    `bson:"currentAddress" json:"currentAddress"`
	TenantGST        []GSTModel      `bson:"tenantGST,omitempty" json:"tenantGST,omitempty"`
	TenantPAN        PANModel        `bson:"tenantPAN,omitempty" json:"tenantPAN,omitempty"`
	TenantTAN        TANModel        `bson:"tenantTAN,omitempty" json:"tenantTAN,omitempty"`
	TenantCIN        CINModel        `bson:"tenantCIN,omitempty" json:"tenantCIN,omitempty"`
	EnabledFeatures  EnabledFeatures `bson:"enabledFeatures,omitempty" json:"enabledFeatures,omitempty"`
	OLTs             []string        `bson:"olts,omitempty" json:"olts,omitempty"`
	TenantMFAPolicy  TenantMFAPolicy `bson:"tenantMFAPolicy,omitempty" json:"tenantMFAPolicy,omitempty"`

	// Note: Provider configurations (KYC, SMS, Mail, App Messaging) are now managed through
	// separate binding collections and published via separate provider binding events.
	// See: tenant.provider.binding.event.model.go

	IsSystem bool `bson:"isSystem" json:"isSystem"`
	IsActive bool `bson:"isActive" json:"isActive"`
	Version  int  `bson:"version,omitempty" json:"version,omitempty"`
}

type TenantUpdateEventModel struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	Code        string `bson:"code,omitempty" json:"code,omitempty"`
	Type        string `bson:"type,omitempty" json:"type,omitempty"`
	Fqdn        string `bson:"fqdn,omitempty" json:"fqdn,omitempty"`
	Environment string `bson:"environment,omitempty" json:"environment,omitempty"`
	EntType     string `bson:"entType,omitempty" json:"entType,omitempty"`

	// Hierarchy Management
	ParentTenantID string   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"`
	ChildIDs       []string `bson:"childIds,omitempty" json:"childIds,omitempty"`
	Path           string   `bson:"path,omitempty" json:"path,omitempty"`
	PathDepth      int      `bson:"pathDepth,omitempty" json:"pathDepth,omitempty"`
	Ancestors      []string `bson:"ancestors,omitempty" json:"ancestors,omitempty"`
	Level          int      `bson:"level,omitempty" json:"level,omitempty"`
	IsLeaf         bool     `bson:"isLeaf,omitempty" json:"isLeaf,omitempty"`
	ChildrenCount  int      `bson:"childrenCount,omitempty" json:"childrenCount,omitempty"`

	// Hierarchy Constraints
	MaxDepth          int      `bson:"maxDepth,omitempty" json:"maxDepth,omitempty"`
	AllowedChildTypes []string `bson:"allowedChildTypes,omitempty" json:"allowedChildTypes,omitempty"`

	DefaultEmail     string           `bson:"defaultEmail,omitempty" json:"defaultEmail,omitempty"`
	DefaultPhone     string           `bson:"defaultPhone,omitempty" json:"defaultPhone,omitempty"`
	PermanentAddress *AddressModel    `bson:"permanentAddress,omitempty" json:"permanentAddress,omitempty"`
	CurrentAddress   *AddressModel    `bson:"currentAddress,omitempty" json:"currentAddress,omitempty"`
	TenantGST        []GSTModel       `bson:"tenantGST,omitempty" json:"tenantGST,omitempty"`
	TenantPAN        *PANModel        `bson:"tenantPAN,omitempty" json:"tenantPAN,omitempty"`
	TenantTAN        *TANModel        `bson:"tenantTAN,omitempty" json:"tenantTAN,omitempty"`
	TenantCIN        *CINModel        `bson:"tenantCIN,omitempty" json:"tenantCIN,omitempty"`
	EnabledFeatures  *EnabledFeatures `bson:"enabledFeatures,omitempty" json:"enabledFeatures,omitempty"`
	OLTs             []string         `bson:"olts,omitempty" json:"olts,omitempty"`

	// Note: Provider configurations (KYC, SMS, Mail, App Messaging) are now managed through
	// separate binding collections and published via separate provider binding events.
	// See: tenant.provider.binding.event.model.go

	IsSystem *bool `bson:"isSystem,omitempty" json:"isSystem,omitempty"`
	IsActive *bool `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Version  int   `bson:"version,omitempty" json:"version,omitempty"`
}

type GSTModel struct {
	State      string    `bson:"state,omitempty" json:"state,omitempty"`
	GSTIN      string    `bson:"gstin,omitempty" json:"gstin,omitempty"`
	IsVerified bool      `bson:"isVerified,omitempty" json:"isVerified,omitempty"`
	VerifiedAt time.Time `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
}

type TenantDeleteEventModel struct {
	ID             string   `bson:"_id" json:"id"`
	ParentTenantID string   `bson:"parentTenantID,omitempty" json:"parentTenantID,omitempty"` // Parent needs to update childrenCount
	Path           string   `bson:"path,omitempty" json:"path,omitempty"`                     // For cascade delete checks
	Ancestors      []string `bson:"ancestors,omitempty" json:"ancestors,omitempty"`           // Notify all ancestors
}

// ==================== NEW HIERARCHY-SPECIFIC EVENTS ====================

// TenantParentChangedEventModel - Published when tenant's parent relationship changes
// Other services use this to update hierarchy-dependent data (permissions, quotas, aggregations)
type TenantParentChangedEventModel struct {
	TenantID     string    `bson:"tenantID" json:"tenantID"`                             // Tenant whose parent changed
	OldParentID  string    `bson:"oldParentID,omitempty" json:"oldParentID,omitempty"`   // Previous parent (empty if was root)
	NewParentID  string    `bson:"newParentID,omitempty" json:"newParentID,omitempty"`   // New parent (empty if becoming root)
	OldPath      string    `bson:"oldPath" json:"oldPath"`                               // Previous path
	NewPath      string    `bson:"newPath" json:"newPath"`                               // New path
	OldAncestors []string  `bson:"oldAncestors,omitempty" json:"oldAncestors,omitempty"` // Previous ancestors
	NewAncestors []string  `bson:"newAncestors,omitempty" json:"newAncestors,omitempty"` // New ancestors
	OldLevel     int       `bson:"oldLevel" json:"oldLevel"`                             // Previous level
	NewLevel     int       `bson:"newLevel" json:"newLevel"`                             // New level
	ChangedAt    time.Time `bson:"changedAt" json:"changedAt"`                           // When change occurred
	ChangedBy    string    `bson:"changedBy,omitempty" json:"changedBy,omitempty"`       // User who made change
}

// TenantChildAddedEventModel - Published when a child is added to a tenant
// Used to update aggregate counts, permissions inheritance, etc.
type TenantChildAddedEventModel struct {
	ParentID      string    `bson:"parentID" json:"parentID"`                   // Parent tenant ID
	ChildID       string    `bson:"childID" json:"childID"`                     // Newly added child ID
	ChildType     string    `bson:"childType" json:"childType"`                 // Type of child (Reseller, Distributor, etc.)
	ChildPath     string    `bson:"childPath" json:"childPath"`                 // Child's path
	ParentPath    string    `bson:"parentPath" json:"parentPath"`               // Parent's path
	NewChildCount int       `bson:"newChildCount" json:"newChildCount"`         // Updated children count
	AddedAt       time.Time `bson:"addedAt" json:"addedAt"`                     // When child was added
	AddedBy       string    `bson:"addedBy,omitempty" json:"addedBy,omitempty"` // User who added child
}

// TenantChildRemovedEventModel - Published when a child is removed from a tenant
// Used to update aggregate counts, cleanup permissions, etc.
type TenantChildRemovedEventModel struct {
	ParentID      string    `bson:"parentID" json:"parentID"`                       // Parent tenant ID
	ChildID       string    `bson:"childID" json:"childID"`                         // Removed child ID
	ChildType     string    `bson:"childType" json:"childType"`                     // Type of removed child
	NewChildCount int       `bson:"newChildCount" json:"newChildCount"`             // Updated children count
	IsNowLeaf     bool      `bson:"isNowLeaf" json:"isNowLeaf"`                     // True if parent has no more children
	RemovedAt     time.Time `bson:"removedAt" json:"removedAt"`                     // When child was removed
	RemovedBy     string    `bson:"removedBy,omitempty" json:"removedBy,omitempty"` // User who removed child
}

// TenantHierarchyRecomputedEventModel - Published when hierarchy fields are recalculated
// Used to sync hierarchy data across services
type TenantHierarchyRecomputedEventModel struct {
	TenantID     string    `bson:"tenantID" json:"tenantID"`                 // Tenant whose hierarchy was recomputed
	Path         string    `bson:"path" json:"path"`                         // New path
	PathDepth    int       `bson:"pathDepth" json:"pathDepth"`               // New depth
	Ancestors    []string  `bson:"ancestors" json:"ancestors"`               // New ancestors
	Level        int       `bson:"level" json:"level"`                       // New level
	RecomputedAt time.Time `bson:"recomputedAt" json:"recomputedAt"`         // When recomputation occurred
	Reason       string    `bson:"reason,omitempty" json:"reason,omitempty"` // Why recomputation happened (parent_changed, migration, etc.)
}

type ProvidersModel struct {
	DefaultProviderID string                      `bson:"defaultProviderID,omitempty" json:"defaultProviderID,omitempty"`
	UseTemplate       bool                        `bson:"useTemplate" json:"useTemplate"`
	UseParent         bool                        `bson:"useParent,omitempty" json:"useParent,omitempty"`
	TenantConfig      *TenantPaymentGatewayConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Providers         []string                    `bson:"providers,omitempty" json:"providers,omitempty"`
}

type TenantPaymentGatewayConfig struct {
	DefaultProviderID string         `bson:"defaultProviderID" json:"defaultProviderID"`
	Providers         []string       `bson:"providers" json:"providers"`
	Gateway           PaymentGateway `bson:"gateway" json:"gateway"`
	Metadata          map[string]any `bson:"metadata" json:"metadata"`
}

type SMSProviderTenantConfig struct {
	ProviderID   string                   `bson:"providerID,omitempty" json:"providerID,omitempty"`
	Channel      string                   `bson:"channel,omitempty" json:"channel,omitempty"` // SMS channel (e.g., "OTP", "Promotional", "Transactional")
	UseTemplate  bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent    bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	TenantConfig *TenantSMSProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority     int                      `bson:"priority" json:"priority"`
	IsActive     bool                     `bson:"isActive" json:"isActive"`
	FailoverOn   bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries   int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight       int                      `bson:"weight,omitempty" json:"weight,omitempty"`
}

type TenantSMSProviderConfig struct {
	Provider SMSType        `bson:"provider" json:"provider"`
	Metadata map[string]any `bson:"metadata" json:"metadata"`
}

type MailServerProviderTenantConfig struct {
	ProviderID   string                  `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate  bool                    `bson:"useTemplate" json:"useTemplate"`
	UseParent    bool                    `bson:"useParent,omitempty" json:"useParent,omitempty"`
	TenantConfig *TenantMailServerConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority     int                     `bson:"priority" json:"priority"`
	IsActive     bool                    `bson:"isActive" json:"isActive"`
	FailoverOn   bool                    `bson:"failoverOn" json:"failoverOn"`
	MaxRetries   int                     `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight       int                     `bson:"weight,omitempty" json:"weight,omitempty"`
}

type TenantMailServerConfig struct {
	SortCode        MailProviderType `bson:"sort_code" json:"sort_code"`
	SMTPConfig      *SMTPConfig      `bson:"smtp,omitempty" json:"smtp,omitempty"`
	SendGridConfig  *SendGridConfig  `bson:"sendgrid,omitempty" json:"sendgrid,omitempty"`
	MailgunConfig   *MailgunConfig   `bson:"mailgun,omitempty" json:"mailgun,omitempty"`
	PostalConfig    *PostalConfig    `bson:"postal,omitempty" json:"postal,omitempty"`
	MailchimpConfig *MailchimpConfig `bson:"mailchimp,omitempty" json:"mailchimp,omitempty"`
}

type KYCProviderTenantConfig struct {
	ProviderID   string                   `bson:"providerID,omitempty" json:"providerID,omitempty"`
	UseTemplate  bool                     `bson:"useTemplate" json:"useTemplate"`
	UseParent    bool                     `bson:"useParent,omitempty" json:"useParent,omitempty"`
	TenantConfig *TenantKYCProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	Priority     int                      `bson:"priority" json:"priority"`
	IsActive     bool                     `bson:"isActive" json:"isActive"`
	FailoverOn   bool                     `bson:"failoverOn" json:"failoverOn"`
	MaxRetries   int                      `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight       int                      `bson:"weight,omitempty" json:"weight,omitempty"`
}

type TenantKYCProviderConfig struct {
	Provider KYCType        `bson:"provider" json:"provider"`
	Metadata map[string]any `bson:"metadata" json:"metadata"`
}

type AppMessagingProvidersModel struct {
	MessageProviderID string                    `bson:"messageProviderID,omitempty" json:"messageProviderID,omitempty"`
	Channel           string                    `bson:"channel,omitempty" json:"channel,omitempty"` // Messaging channel (WhatsApp, Telegram, etc.)
	UseTemplate       bool                      `bson:"useTemplate" json:"useTemplate"`
	UseParent         bool                      `bson:"useParent,omitempty" json:"useParent,omitempty"`
	TenantConfig      *TenantAppMessagingConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	MessageProvider   MessagingProvider         `bson:"messageProvider,omitempty" json:"messageProvider,omitempty"`
	Priority          int                       `bson:"priority,omitempty" json:"priority,omitempty"`
	IsActive          bool                      `bson:"isActive,omitempty" json:"isActive,omitempty"`
	FailoverOn        bool                      `bson:"failoverOn,omitempty" json:"failoverOn,omitempty"`
	MaxRetries        int                       `bson:"maxRetries,omitempty" json:"maxRetries,omitempty"`
	Weight            int                       `bson:"weight,omitempty" json:"weight,omitempty"`
}

type TenantAppMessagingConfig struct {
	MessageProvider      MessagingProvider `bson:"messageProvider" json:"messageProvider"`
	Metadata             map[string]any    `bson:"metadata" json:"metadata"`
	PhoneNumberID        []string          `bson:"phoneNumberID,omitempty" json:"phoneNumberID,omitempty"` // Business Phone Number ID(s) - array of phone numbers for sending messages
	AccessToken          string            `bson:"accessToken,omitempty" json:"accessToken,omitempty"`     // Meta Graph API Access Token - for sending messages and template submissions
	AllowedMessagesTypes []string          `bson:"allowedMessagesTypes,omitempty" json:"allowedMessagesTypes,omitempty"`
}

type ExternalRadiusProvidersModel struct {
	ExternalRadiusProviderID string                      `bson:"externalRadiusProviderID,omitempty" json:"externalRadiusProviderID,omitempty"`
	UseTemplate              bool                        `bson:"useTemplate" json:"useTemplate"`
	UseParent                bool                        `bson:"useParent,omitempty" json:"useParent,omitempty"`
	TenantConfig             *TenantRadiusProviderConfig `bson:"tenantConfig,omitempty" json:"tenantConfig,omitempty"`
	ExternalRadiusProviders  RadiusProvider              `bson:"externalRadiusProviders,omitempty" json:"externalRadiusProviders,omitempty"`
}

type TenantRadiusProviderConfig struct {
	Provider RadiusProvider `bson:"provider" json:"provider"`
	Metadata map[string]any `bson:"metadata" json:"metadata"`
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

// Lightweight enum aliases to mirror tenant-service models without a direct dependency
type MailProviderType string
type MessagingProvider string
type RadiusProvider string
type KYCType string
type PaymentGateway string

// ProviderTenantConfig represents a common configuration for providers assigned to tenants
// Used for SMS, Mail, and other providers that support priority-based routing with failover
type ProviderTenantConfig struct {
	ProviderID string `json:"providerID" bson:"providerID" validate:"required,isCuid2"`
	Priority   int    `json:"priority" bson:"priority" validate:"required,min=1"`                                 // 1 = Primary, 2 = Secondary, 3 = Tertiary, etc.
	IsActive   bool   `json:"isActive" bson:"isActive"`                                                           // Can disable provider for specific tenant
	FailoverOn bool   `json:"failoverOn" bson:"failoverOn"`                                                       // Enable automatic failover to next provider
	MaxRetries int    `json:"maxRetries,omitempty" bson:"maxRetries,omitempty" validate:"omitempty,min=1,max=10"` // Max retries before failover (default: 3)
	Weight     int    `json:"weight,omitempty" bson:"weight,omitempty" validate:"omitempty,min=1,max=100"`        // For load balancing (1-100)
}
