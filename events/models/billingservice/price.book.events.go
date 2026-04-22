package billingservice

import (
	"time"

	"github.com/praction-networks/common/events/models/planservice"
)

// PriceBookItemCDC represents a price book item in CDC events from billing-service
type PriceBookItemCDC struct {
	ProductID       string                       `json:"productId"`
	Unit            planservice.Unit             `json:"unit"`
	PricingModel    planservice.PricingModelType `json:"pricingModel"`
	Amount          *float64                     `json:"amount,omitempty"`
	DiscountPercent *float64                     `json:"discountPercent,omitempty"`
	PricingTiers    []PricingTierCDC             `json:"pricingTiers,omitempty"`
	VolumeDiscounts []VolumeDiscountCDC          `json:"volumeDiscounts,omitempty"`
	BillingCycle    *planservice.BillingCycle    `json:"billingCycle,omitempty"`
}

// PricingTierCDC represents a pricing tier in CDC events
type PricingTierCDC struct {
	MinQty    float64  `json:"minQty"`
	MaxQty    *float64 `json:"maxQty,omitempty"`
	UnitPrice float64  `json:"unitPrice"`
	FlatFee   float64  `json:"flatFee"`
}

// VolumeDiscountCDC represents a volume discount in CDC events
type VolumeDiscountCDC struct {
	MinQty      float64  `json:"minQty"`
	MaxQty      *float64 `json:"maxQty,omitempty"`
	DiscountPct float64  `json:"discountPct"`
	DiscountAmt float64  `json:"discountAmt"`
}

// PriceBookCDCEvent is the payload published by billing-service CDC for price book changes.
// This is consumed by plan-service to maintain a local cache of price books.
//
// Leaf-tenant pricing fields (BasePrice, CommissionAmount, CommissionType, AgrPolicy)
// are billing-owned — billing is the source of truth for commission/AGR configuration
// and plan-service stores them as a read-only projection for catalog display.
type PriceBookCDCEvent struct {
	ID               string             `json:"id"`
	Scope            string             `json:"scope"`
	PlanID           *string            `json:"planId,omitempty"`
	TenantID         *string            `json:"tenantId,omitempty"`
	TargetTenantType *string            `json:"targetTenantType,omitempty"`
	Type             string             `json:"type"` // WHOLESALE or RETAIL
	Country          string             `json:"country"`
	State            *string            `json:"state,omitempty"`
	Currency         string             `json:"currency"`
	EffectiveFrom    time.Time          `json:"effectiveFrom"`
	EffectiveTo      *time.Time         `json:"effectiveTo,omitempty"`
	Items            []PriceBookItemCDC `json:"items"`
	PlanLevelPrice   *float64           `json:"planLevelPrice,omitempty"`
	TargetTenantID   *string            `json:"targetTenantId,omitempty"`

	// Leaf-tenant pricing (billing-owned, read-only on plan side).
	BasePrice        *int64  `json:"basePrice,omitempty"`        // customer sell price in paise
	CommissionAmount *int64  `json:"commissionAmount,omitempty"` // FIXED paise or PERCENT basis points
	CommissionType   *string `json:"commissionType,omitempty"`   // FIXED | PERCENT
	AgrPolicy        string  `json:"agrPolicy,omitempty"`        // NONE | ONLINE_ONLY | ALL

	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
	Version   int       `json:"version"`
}
