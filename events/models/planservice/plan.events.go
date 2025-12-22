package planservice

import (
	"time"
)

type PlanModel struct {
	ID            string       `json:"id" bson:"_id"`
	OwnerTenantID *string      `json:"ownerTenantId,omitempty" bson:"ownerTenantId,omitempty"` // Tenant who created/owns this (has edit rights)
	TenantIDs     []string     `json:"tenantIds,omitempty" bson:"tenantIds,omitempty"`         // Tenants who can access/use this plan (ISP, reseller, distributor, partner)
	Scope         CatalogScope `json:"scope" bson:"scope"`                                     // GLOBAL or TENANT

	Code                string                   `json:"code" bson:"code"`
	Name                string                   `json:"name" bson:"name"`
	Description         string                   `json:"description,omitempty" bson:"description,omitempty"`
	PlanType           PlanType     `json:"planType" bson:"planType"`
	PlanSubType         PlanSubType  `json:"planSubType,omitempty" bson:"planSubType,omitempty"`
	BillingCycle        BillingCycle `json:"billingCycle" bson:"billingCycle"`                                   // Default/primary cycle
	BillingCyclePricing []BillingCyclePricing    `json:"billingCyclePricing,omitempty" bson:"billingCyclePricing,omitempty"` // Multiple cycles with pricing
	ValidityDays        *int                     `json:"validityDays,omitempty" bson:"validityDays,omitempty"`               // Optional override (mainly for ONE_TIME)
	RenewalPolicy       *RenewalPolicy           `json:"renewalPolicy,omitempty" bson:"renewalPolicy,omitempty"`
	ContractTerms       *ContractTerms           `json:"contractTerms,omitempty" bson:"contractTerms,omitempty"` // Contract and commitment terms
	SLAs                []SLA                    `json:"slas,omitempty" bson:"slas,omitempty"`                   // Service Level Agreements
	CreditPolicy        *CreditPolicy            `json:"creditPolicy,omitempty" bson:"creditPolicy,omitempty"`   // Credit limits and financial controls
	Items               []PlanItem               `json:"items" bson:"items"`
	AttachedPromotions  []string                 `json:"attachedPromotions,omitempty" bson:"attachedPromotions,omitempty"`
	Status              Status       `json:"status" bson:"status"`
	NASAttributes      NASAttributes            `json:"nasAttributes,omitempty" bson:"nasAttributes,omitempty"` // RADIUS attributes per NAS vendor type
	CreatedAt           time.Time                `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time                `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy           string                   `json:"updatedBy" bson:"updatedBy"`
	CreatedBy           string                   `json:"createdBy" bson:"createdBy"`
	Version             int                      `json:"version" bson:"version"`
}

// PlanUpdatedEvent represents a plan update event
type PlanUpdatedEvent struct {
	ID                 string         `json:"id" bson:"_id"`
	Code               string         `json:"code,omitempty" bson:"code,omitempty"`
	Name               string         `json:"name,omitempty" bson:"name,omitempty"`
	Description        string         `json:"description,omitempty" bson:"description,omitempty"`
	PlanType           PlanType       `json:"planType,omitempty" bson:"planType,omitempty"`
	BillingCycle       BillingCycle   `json:"billingCycle,omitempty" bson:"billingCycle,omitempty"`
	Status             Status         `json:"status,omitempty" bson:"status,omitempty"`
	Version            int            `json:"version" bson:"version"`
}

// PlanDeletedEvent represents a plan deletion event
type PlanDeletedEvent struct {
	ID       string    `json:"id" bson:"_id"`
	Code     string    `json:"code" bson:"code"`
	Version  int       `json:"version" bson:"version"`
}

type PlanItem struct {
	ItemID              string                       `json:"itemId" bson:"itemId"`
	ProductID           string                       `json:"productId" bson:"productId"`
	Unit                Unit             `json:"unit" bson:"unit"`
	Qty                 int                          `json:"qty" bson:"qty"`
	PriceStrategy       PriceStrategy    `json:"priceStrategy" bson:"priceStrategy"`
	PricingModel        PricingModelType `json:"pricingModel" bson:"pricingModel"`                                   // "FLAT", "TIERED", "VOLUME", "BUNDLE"
	OverridePrice       *float64                     `json:"overridePrice,omitempty" bson:"overridePrice,omitempty"`             // Flat price (for FLAT model)
	// PricingTiers        []PricingTier                `json:"pricingTiers,omitempty" bson:"pricingTiers,omitempty"`               // Tiered pricing (for TIERED model)
	// VolumeDiscounts     []VolumeDiscount             `json:"volumeDiscounts,omitempty" bson:"volumeDiscounts,omitempty"`         // Volume discounts (for VOLUME model)
	// ResourceAllocations []ResourceAllocation         `json:"resourceAllocations,omitempty" bson:"resourceAllocations,omitempty"` // Resource allocation tracking
	// ResourceConstraints []ResourceConstraint         `json:"resourceConstraints,omitempty" bson:"resourceConstraints,omitempty"` // Resource constraints
	Role                PlanItemRole     `json:"role" bson:"role"`                                                   // "BASE", "ADDON", "BUNDLE", "OPTIONAL"
	Metadata            map[string]any               `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
// BillingCyclePricing defines pricing for a specific billing cycle
type BillingCyclePricing struct {
	BillingCycle BillingCycle `json:"billingCycle" bson:"billingCycle"`                   // Quantity + Unit (e.g., {quantity: 1, unit: "MONTH"})
	BasePrice    float64                  `json:"basePrice" bson:"basePrice"`                         // Base price for this cycle
	DiscountPct  *float64                 `json:"discountPct,omitempty" bson:"discountPct,omitempty"` // Optional discount percentage
	IsDefault    bool                     `json:"isDefault" bson:"isDefault"`                         // Mark one as default
	IsActive     bool                     `json:"isActive" bson:"isActive"`                           // Enable/disable this cycle
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
	CreditLimit               *float64                     `json:"creditLimit,omitempty" bson:"creditLimit,omitempty"`                             // Credit limit amount
	CreditLimitType           CreditLimitType  `json:"creditLimitType" bson:"creditLimitType"`                                         // "HARD", "SOFT", "NONE"
	CreditCheckRequired       bool                         `json:"creditCheckRequired" bson:"creditCheckRequired"`                                 // Whether credit check is required
	DepositRequired           *float64                     `json:"depositRequired,omitempty" bson:"depositRequired,omitempty"`                     // Deposit amount required
	DepositRefundable         bool                         `json:"depositRefundable" bson:"depositRefundable"`                                     // Whether deposit is refundable
	PaymentTerms              PaymentTermsType `json:"paymentTerms" bson:"paymentTerms"`                                               // "PREPAID", "POSTPAID", "HYBRID"
	PaymentMethodRestrictions []string                     `json:"paymentMethodRestrictions,omitempty" bson:"paymentMethodRestrictions,omitempty"` // Restricted payment methods
	BillingCycleRestriction   *string                      `json:"billingCycleRestriction,omitempty" bson:"billingCycleRestriction,omitempty"`     // Restricted billing cycles
	AutoSuspendOnLimit        bool                         `json:"autoSuspendOnLimit" bson:"autoSuspendOnLimit"`                                   // Auto-suspend when credit limit reached
	GracePeriodDays           int                          `json:"gracePeriodDays,omitempty" bson:"gracePeriodDays,omitempty"`                     // Grace period after limit exceeded
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