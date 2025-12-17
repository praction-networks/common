package planservice

import "time"

// PlanItemTenantPricingEvent represents tenant-specific pricing for a plan item in events
type PlanItemTenantPricingEvent struct {
	ItemID          string           `json:"itemId" bson:"itemId"`                     // Reference to PlanItem.ItemID
	ProductID       string           `json:"productId" bson:"productId"`              // Reference to PlanItem.ProductID
	Unit            Unit             `json:"unit" bson:"unit"`                         // Reference to PlanItem.Unit
	PricingModel    PricingModelType `json:"pricingModel" bson:"pricingModel"`         // "FLAT", "TIERED", "VOLUME", "BUNDLE"
	OverridePrice   *float64          `json:"overridePrice,omitempty" bson:"overridePrice,omitempty"` // Flat price override (for FLAT model)
	PricingTiers    []PricingTierEvent `json:"pricingTiers,omitempty" bson:"pricingTiers,omitempty"`     // Tiered pricing override (for TIERED model)
	VolumeDiscounts []VolumeDiscountEvent `json:"volumeDiscounts,omitempty" bson:"volumeDiscounts,omitempty"` // Volume discounts override (for VOLUME model)
	DiscountPercent *float64          `json:"discountPercent,omitempty" bson:"discountPercent,omitempty"` // Percentage discount from base plan price
	DiscountAmount  *float64          `json:"discountAmount,omitempty" bson:"discountAmount,omitempty"`   // Fixed amount discount from base plan price
	Metadata        map[string]any    `json:"metadata,omitempty" bson:"metadata,omitempty"`               // Additional pricing metadata
}

// PricingTierEvent represents a pricing tier in events
type PricingTierEvent struct {
	MinQty  float64 `json:"minQty" bson:"minQty"`
	MaxQty  *float64 `json:"maxQty,omitempty" bson:"maxQty,omitempty"`
	Price   float64 `json:"price" bson:"price"`
	Unit    Unit    `json:"unit" bson:"unit"`
}

// VolumeDiscountEvent represents a volume discount in events
type VolumeDiscountEvent struct {
	MinQty     float64 `json:"minQty" bson:"minQty"`
	MaxQty     *float64 `json:"maxQty,omitempty" bson:"maxQty,omitempty"`
	DiscountPct float64 `json:"discountPct" bson:"discountPct"`
	DiscountAmt *float64 `json:"discountAmt,omitempty" bson:"discountAmt,omitempty"`
}

// PlanTenantPricingCreatedEvent represents a plan tenant pricing creation event
type PlanTenantPricingCreatedEvent struct {
	ID            string                      `json:"id" bson:"_id"`
	PlanID        string                      `json:"planId" bson:"planId"`                    // Reference to PlanModel.ID
	TenantID      string                      `json:"tenantId" bson:"tenantId"`                // Tenant identifier
	ItemPricing   []PlanItemTenantPricingEvent `json:"itemPricing" bson:"itemPricing"`          // Tenant-specific pricing for each plan item
	EffectiveFrom time.Time                   `json:"effectiveFrom" bson:"effectiveFrom"`      // When this pricing becomes effective
	EffectiveTo   *time.Time                  `json:"effectiveTo,omitempty" bson:"effectiveTo,omitempty"` // When this pricing expires (nil = no expiry)
	Status        Status                      `json:"status" bson:"status"`                    // ACTIVE, INACTIVE, etc.
	CreatedAt     time.Time                   `json:"createdAt" bson:"createdAt"`
	CreatedBy     string                      `json:"createdBy" bson:"createdBy"`
	Version       int                         `json:"version" bson:"version"`
}

// PlanTenantPricingUpdatedEvent represents a plan tenant pricing update event
type PlanTenantPricingUpdatedEvent struct {
	ID            string                      `json:"id" bson:"_id"`
	PlanID        string                      `json:"planId" bson:"planId"`
	TenantID      string                      `json:"tenantId" bson:"tenantId"`
	ItemPricing   []PlanItemTenantPricingEvent `json:"itemPricing,omitempty" bson:"itemPricing,omitempty"`
	EffectiveFrom *time.Time                  `json:"effectiveFrom,omitempty" bson:"effectiveFrom,omitempty"`
	EffectiveTo   *time.Time                  `json:"effectiveTo,omitempty" bson:"effectiveTo,omitempty"`
	Status        Status                      `json:"status,omitempty" bson:"status,omitempty"`
	UpdatedAt     time.Time                   `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy     string                      `json:"updatedBy" bson:"updatedBy"`
	Version       int                         `json:"version" bson:"version"`
}

// PlanTenantPricingDeletedEvent represents a plan tenant pricing deletion event
type PlanTenantPricingDeletedEvent struct {
	ID       string    `json:"id" bson:"_id"`
	PlanID   string    `json:"planId" bson:"planId"`
	TenantID string    `json:"tenantId" bson:"tenantId"`
	Version  int       `json:"version" bson:"version"`
}

