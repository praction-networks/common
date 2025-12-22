package planservice

import "time"

// PriceBookItemEvent represents a price book item in events
type PriceBookItemEvent struct {
	ProductID       string                `json:"productId" bson:"productId"`
	Unit            Unit                  `json:"unit" bson:"unit"`
	PricingModel    PricingModelType      `json:"pricingModel" bson:"pricingModel"`                           // "FLAT", "TIERED", "VOLUME", "BUNDLE"
	Amount          *float64              `json:"amount,omitempty" bson:"amount,omitempty"`                   // Flat price (for FLAT model)
	PricingTiers    []PricingTierEvent    `json:"pricingTiers,omitempty" bson:"pricingTiers,omitempty"`       // Tiered pricing (for TIERED model)
	VolumeDiscounts []VolumeDiscountEvent `json:"volumeDiscounts,omitempty" bson:"volumeDiscounts,omitempty"` // Volume discounts (for VOLUME model)
}

// PriceBookCreatedEvent represents a price book creation event
type PriceBookCreatedEvent struct {
	ID            string               `json:"id" bson:"_id"`
	Scope         PriceBookScope       `json:"scope" bson:"scope"`
	TenantID      *string              `json:"tenantId,omitempty" bson:"tenantId,omitempty"`
	Country       string               `json:"country" bson:"country"`
	State         *string              `json:"state,omitempty" bson:"state,omitempty"`
	Currency      string               `json:"currency" bson:"currency"`
	EffectiveFrom time.Time            `json:"effectiveFrom" bson:"effectiveFrom"`
	EffectiveTo   *time.Time           `json:"effectiveTo,omitempty" bson:"effectiveTo,omitempty"`
	Items         []PriceBookItemEvent `json:"items" bson:"items"`
	Status        Status               `json:"status" bson:"status"`
	CreatedAt     time.Time            `json:"createdAt" bson:"createdAt"`
	CreatedBy     string               `json:"createdBy" bson:"createdBy"`
	Version       int                  `json:"version" bson:"version"`
}

// PriceBookUpdatedEvent represents a price book update event
type PriceBookUpdatedEvent struct {
	ID            string               `json:"id" bson:"_id"`
	Scope         PriceBookScope       `json:"scope,omitempty" bson:"scope,omitempty"`
	TenantID      *string              `json:"tenantId,omitempty" bson:"tenantId,omitempty"`
	Country       string               `json:"country,omitempty" bson:"country,omitempty"`
	State         *string              `json:"state,omitempty" bson:"state,omitempty"`
	Currency      string               `json:"currency,omitempty" bson:"currency,omitempty"`
	EffectiveFrom *time.Time           `json:"effectiveFrom,omitempty" bson:"effectiveFrom,omitempty"`
	EffectiveTo   *time.Time           `json:"effectiveTo,omitempty" bson:"effectiveTo,omitempty"`
	Items         []PriceBookItemEvent `json:"items,omitempty" bson:"items,omitempty"`
	Status        Status               `json:"status,omitempty" bson:"status,omitempty"`
	UpdatedAt     time.Time            `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy     string               `json:"updatedBy" bson:"updatedBy"`
	Version       int                  `json:"version" bson:"version"`
}

// PriceBookDeletedEvent represents a price book deletion event
type PriceBookDeletedEvent struct {
	ID      string `json:"id" bson:"_id"`
	Version int    `json:"version" bson:"version"`
}
