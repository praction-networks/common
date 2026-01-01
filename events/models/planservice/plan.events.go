package planservice

import (
	"time"
)

type PlanCreatedEvent struct {
	ID                   string                `json:"id" bson:"_id"`
	OwnerTenantID        *string               `json:"ownerTenantId,omitempty" bson:"ownerTenantId,omitempty"` // Tenant who created/owns this (has edit rights)
	TenantIDs            []string              `json:"tenantIds,omitempty" bson:"tenantIds,omitempty"`         // Tenants who can access/use this plan (ISP, reseller, distributor, partner)
	Scope                CatalogScope          `json:"scope" bson:"scope"`                                     // GLOBAL or TENANT
	Code                 string                `json:"code" bson:"code"`
	Name                 string                `json:"name" bson:"name"`
	PlanType             PlanType              `json:"planType" bson:"planType"`
	PlanSubType          PlanSubType           `json:"planSubType,omitempty" bson:"planSubType,omitempty"`
	BillingCyclePricing  []BillingCyclePricing `json:"billingCyclePricing" bson:"billingCyclePricing"`                       // Billing cycles with pricing (fixed or percentage-based)
	HotspotBillingConfig *HotspotBillingConfig `json:"hotspotBillingConfig,omitempty" bson:"hotspotBillingConfig,omitempty"` // HOTSPOT-specific: session duration, connection frequency, post-expiration action
	ValidityDays         *int                  `json:"validityDays,omitempty" bson:"validityDays,omitempty"`                 // Optional override (mainly for ONE_TIME)
	RenewalPolicy        *RenewalPolicy        `json:"renewalPolicy,omitempty" bson:"renewalPolicy,omitempty"`
	ContractTerms        *ContractTerms        `json:"contractTerms,omitempty" bson:"contractTerms,omitempty"` // Contract and commitment terms
	SLAs                 []SLA                 `json:"slas,omitempty" bson:"slas,omitempty"`                   // Service Level Agreements
	CreditPolicy         *CreditPolicy         `json:"creditPolicy,omitempty" bson:"creditPolicy,omitempty"`   // Credit limits and financial controls
	Items                []PlanItem            `json:"items" bson:"items"`
	AttachedPromotions   []string              `json:"attachedPromotions,omitempty" bson:"attachedPromotions,omitempty"`
	Status               Status                `json:"status" bson:"status"`
	NASAttributes        NASAttributes         `json:"nasAttributes,omitempty" bson:"nasAttributes,omitempty"`   // RADIUS attributes per NAS vendor type
	PlanDefaultFor       []string              `json:"planDefaultFor,omitempty" bson:"planDefaultFor,omitempty"` // Whether this is the default plan for the tenant
	Version              int                   `json:"version" bson:"version"`
}

type PlanUpdateEvent struct {
	ID                   string                `json:"id" bson:"_id"`
	OwnerTenantID        *string               `json:"ownerTenantId,omitempty" bson:"ownerTenantId,omitempty"` // Tenant who created/owns this (has edit rights)
	TenantIDs            []string              `json:"tenantIds,omitempty" bson:"tenantIds,omitempty"`         // Tenants who can access/use this plan (ISP, reseller, distributor, partner)
	Scope                CatalogScope          `json:"scope" bson:"scope"`                                     // GLOBAL or TENANT
	Code                 string                `json:"code" bson:"code"`
	Name                 string                `json:"name" bson:"name"`
	PlanType             PlanType              `json:"planType" bson:"planType"`
	PlanSubType          PlanSubType           `json:"planSubType,omitempty" bson:"planSubType,omitempty"`
	BillingCyclePricing  []BillingCyclePricing `json:"billingCyclePricing" bson:"billingCyclePricing"`                       // Billing cycles with pricing (fixed or percentage-based)
	HotspotBillingConfig *HotspotBillingConfig `json:"hotspotBillingConfig,omitempty" bson:"hotspotBillingConfig,omitempty"` // HOTSPOT-specific: session duration, connection frequency, post-expiration action
	ValidityDays         *int                  `json:"validityDays,omitempty" bson:"validityDays,omitempty"`                 // Optional override (mainly for ONE_TIME)
	RenewalPolicy        *RenewalPolicy        `json:"renewalPolicy,omitempty" bson:"renewalPolicy,omitempty"`
	ContractTerms        *ContractTerms        `json:"contractTerms,omitempty" bson:"contractTerms,omitempty"` // Contract and commitment terms
	SLAs                 []SLA                 `json:"slas,omitempty" bson:"slas,omitempty"`                   // Service Level Agreements
	CreditPolicy         *CreditPolicy         `json:"creditPolicy,omitempty" bson:"creditPolicy,omitempty"`   // Credit limits and financial controls
	Items                []PlanItem            `json:"items" bson:"items"`
	AttachedPromotions   []string              `json:"attachedPromotions,omitempty" bson:"attachedPromotions,omitempty"`
	Status               Status                `json:"status" bson:"status"`
	NASAttributes        NASAttributes         `json:"nasAttributes,omitempty" bson:"nasAttributes,omitempty"`   // RADIUS attributes per NAS vendor type
	PlanDefaultFor       []string              `json:"planDefaultFor,omitempty" bson:"planDefaultFor,omitempty"` // Whether this is the default plan for the tenant
	Version              int                   `json:"version" bson:"version"`
}

// PlanDeletedEvent represents a plan deletion event
type PlanDeletedEvent struct {
	ID             string   `json:"id" bson:"_id"`
	Code           string   `json:"code" bson:"code"`
	PlanDefaultFor []string `json:"planDefaultFor,omitempty" bson:"planDefaultFor,omitempty"`
	Version        int      `json:"version" bson:"version"`
}

type PlanItem struct {
	ProductID   string       `json:"productId" bson:"productId"`                     // Product ID reference
	ProductCode string       `json:"productCode" bson:"productCode"`                 // Product code for display/reference
	BasePrice   *float64     `json:"basePrice,omitempty" bson:"basePrice,omitempty"` // Base price for this product in the plan
	Role        PlanItemRole `json:"role" bson:"role"`                               // "BASE", "ADDON", "BUNDLE", "OPTIONAL"
}

// BillingCyclePricing defines pricing for a specific billing cycle
type BillingCyclePricing struct {
	// Max Broadband Subscriptions defines the maximum number of broadband subscriptions that can be created for the plan
	MaxSimultaneousSessions   int          `json:"maxSimultaneousSessions" bson:"maxSimultaneousSessions"`
	MaxBroadbandSubscriptions int          `json:"maxBroadbandSubscriptions" bson:"maxBroadbandSubscriptions"`
	BillingCycle              BillingCycle `json:"billingCycle" bson:"billingCycle"`                   // Quantity + Unit (e.g., {quantity: 1, unit: "MONTH"})
	BasePrice                 float64      `json:"basePrice" bson:"basePrice"`                         // Base price for this cycle
	DiscountPct               *float64     `json:"discountPct,omitempty" bson:"discountPct,omitempty"` // Optional discount percentage
	IsDefault                 bool         `json:"isDefault" bson:"isDefault"`                         // Mark one as default
	IsActive                  bool         `json:"isActive" bson:"isActive"`                           // Enable/disable this cycle
}

// HotspotBillingConfig defines connection time management for HOTSPOT plan subtype
// Handles session duration, connection frequency, and post-expiration behavior
// Example: 4-hour sessions, once per day, disconnect after expiry
type HotspotBillingConfig struct {
	// MaxConcurrentSession defines the maximum number of concurrent sessions that can be connected to the plan
	MaxConcurrentSession int `json:"maxConcurrentSession" bson:"maxConcurrentSession"`

	// MaxDevices defines the maximum number of devices that can be connected to the plan
	MaxDevices int `json:"maxDevices" bson:"maxDevices"`


	// SessionDuration defines how long each connection session lasts
	// Example: {quantity: 4, unit: "HOUR"} = 4-hour sessions
	// Example: {quantity: 1, unit: "HOUR"} = 1-hour sessions
	// Example: {quantity: 30, unit: "MINUTE"} = 30-minute sessions
	SessionDuration BillingCycle `json:"sessionDuration" bson:"sessionDuration"`

	// ConnectionQuota defines how many times per period user can connect (dynamic)
	// Example: {quantity: 2, unit: "DAY"} = 2 times per day
	// Example: {quantity: 1, unit: "WEEK"} = 1 time per week
	// Example: {quantity: 3, unit: "MONTH"} = 3 times per month
	// When both quantity and unit are nil, it means unlimited connections within validity period
	ConnectionQuota HotspotConnectionQuota `json:"connectionQuota" bson:"connectionQuota"`

	// PostExpirationAction defines what happens when session expires
	// DISCONNECT: Immediately disconnect user (most common)
	// REDIRECT: Redirect to payment/upgrade page
	// ALLOW_RENEW: Allow user to manually purchase new session
	PostExpirationAction HotspotPostExpirationAction `json:"postExpirationAction" bson:"postExpirationAction"`

	// RedirectURL is used when PostExpirationAction is REDIRECT
	// URL to redirect user for payment/upgrade
	RedirectURL *string `json:"redirectUrl,omitempty" bson:"redirectUrl,omitempty"`
}

type RenewalPolicy struct {
	AutoRenew         bool  `json:"autoRenew" bson:"autoRenew"`                                     // Enable automatic renewal
	GraceDays         *int  `json:"graceDays,omitempty" bson:"graceDays,omitempty"`                 // Grace period in days for payment (only when autoRenew is true)
	SuspendAfterGrace *bool `json:"suspendAfterGrace,omitempty" bson:"suspendAfterGrace,omitempty"` // Suspend after grace period (only when autoRenew is true)
}

// ContractTerms defines contract and commitment terms for plans
type ContractTerms struct {
	CommitmentPeriod    int        `json:"commitmentPeriod" bson:"commitmentPeriod"` // Months of commitment (0 = no commitment)
	CommitmentStartDate *time.Time `json:"commitmentStartDate,omitempty" bson:"commitmentStartDate,omitempty"`
	CommitmentEndDate   *time.Time `json:"commitmentEndDate,omitempty" bson:"commitmentEndDate,omitempty"`
	EarlyTerminationFee *float64   `json:"earlyTerminationFee,omitempty" bson:"earlyTerminationFee,omitempty"`
	AutoRenewal         bool       `json:"autoRenewal" bson:"autoRenewal"`                                   // Auto-renew at end of commitment
	RenewalNoticeDays   int        `json:"renewalNoticeDays,omitempty" bson:"renewalNoticeDays,omitempty"`   // Notice period before renewal
	MinimumPeriod       int        `json:"minimumPeriod,omitempty" bson:"minimumPeriod,omitempty"`           // Minimum service period (days)
	LockInPeriod        int        `json:"lockInPeriod,omitempty" bson:"lockInPeriod,omitempty"`             // Lock-in period (days)
	CommitmentDiscount  *float64   `json:"commitmentDiscount,omitempty" bson:"commitmentDiscount,omitempty"` // Discount for commitment
}

// SLA defines Service Level Agreements for plans
type SLA struct {
	Metric            string  `json:"metric" bson:"metric"`                   // "UPTIME", "LATENCY", "SPEED", "RESPONSE_TIME"
	TargetValue       float64 `json:"targetValue" bson:"targetValue"`         // e.g., 99.9 for uptime, 50 for latency (ms)
	MeasurementUnit   string  `json:"measurementUnit" bson:"measurementUnit"` // "PERCENT", "MILLISECONDS", "MBPS"
	PenaltyType       string  `json:"penaltyType" bson:"penaltyType"`         // "CREDIT", "REFUND", "SERVICE_CREDIT"
	PenaltyAmount     float64 `json:"penaltyAmount,omitempty" bson:"penaltyAmount,omitempty"`
	MeasurementWindow string  `json:"measurementWindow" bson:"measurementWindow"` // "MONTHLY", "DAILY", "WEEKLY"
	Description       string  `json:"description,omitempty" bson:"description,omitempty"`
}

// CreditPolicy defines credit limits and financial controls for plans
type CreditPolicy struct {
	CreditLimit               *float64         `json:"creditLimit,omitempty" bson:"creditLimit,omitempty"`                             // Credit limit amount
	CreditLimitType           CreditLimitType  `json:"creditLimitType" bson:"creditLimitType"`                                         // "HARD", "SOFT", "NONE"
	CreditCheckRequired       bool             `json:"creditCheckRequired" bson:"creditCheckRequired"`                                 // Whether credit check is required
	DepositRequired           *float64         `json:"depositRequired,omitempty" bson:"depositRequired,omitempty"`                     // Deposit amount required
	DepositRefundable         bool             `json:"depositRefundable" bson:"depositRefundable"`                                     // Whether deposit is refundable
	PaymentTerms              PaymentTermsType `json:"paymentTerms" bson:"paymentTerms"`                                               // "PREPAID", "POSTPAID", "HYBRID"
	PaymentMethodRestrictions []string         `json:"paymentMethodRestrictions,omitempty" bson:"paymentMethodRestrictions,omitempty"` // Restricted payment methods
	BillingCycleRestriction   *string          `json:"billingCycleRestriction,omitempty" bson:"billingCycleRestriction,omitempty"`     // Restricted billing cycles
	AutoSuspendOnLimit        bool             `json:"autoSuspendOnLimit" bson:"autoSuspendOnLimit"`                                   // Auto-suspend when credit limit reached
	GracePeriodDays           int              `json:"gracePeriodDays,omitempty" bson:"gracePeriodDays,omitempty"`                     // Grace period after limit exceeded
}

// NASVendorAttributeValues represents RADIUS attributes for a specific NAS vendor type
// Keys are attribute names (e.g., "sessionTimeout", "bandwidthMaxDown", "rateLimit")
// Values can be string, integer, or other types depending on the attribute
type NASVendorAttributeValues map[string]interface{}

// NASAttributes maps NAS vendor/device types to their RADIUS attribute configurations
// This allows plans to specify different RADIUS attributes for different NAS types
// Structure:
//
//	{
//	  "mikrotik-750": {
//	    "sessionTimeout": 86400,
//	    "idleTimeout": 1800,
//	    "acctInterimInterval": 300,
//	    "bandwidthMaxDown": 20000000,
//	    "bandwidthMaxUp": 20000000,
//	    "rateLimit": "20M/20M",
//	    "replyMessage": "Welcome to Silver Plan"
//	  },
//	  "cisco-ASR": {
//	    "sessionTimeout": 86400,
//	    "idleTimeout": 1800,
//	    ...
//	  },
//	  "juniper-MX": {
//	    ...
//	  },
//	  "other": {
//	    ... (default attributes for unspecified NAS types)
//	  }
//	}
type NASAttributes map[string]NASVendorAttributeValues
