package planservice

import "time"

// PromotionConditionEvent represents a promotion condition in events
type PromotionConditionEvent struct {
	Type     string      `json:"type" bson:"type"` // "PLAN_TYPE", "CUSTOMER_SEGMENT", "QUANTITY", "DATE_RANGE", "BUNDLE"
	Field    string      `json:"field,omitempty" bson:"field,omitempty"`
	Operator string      `json:"operator" bson:"operator"` // "EQ", "GT", "LT", "IN", "BETWEEN", "CONTAINS"
	Value    interface{} `json:"value" bson:"value"`
}

// PromotionActionEvent represents a promotion action in events
type PromotionActionEvent struct {
	Type        string   `json:"type" bson:"type"`     // "DISCOUNT_PERCENT", "DISCOUNT_AMOUNT", "FREE_ITEM", "UPGRADE"
	Target      string   `json:"target" bson:"target"` // "PLAN", "ITEM", "TOTAL"
	Value       float64  `json:"value" bson:"value"`
	ProductID   *string  `json:"productId,omitempty" bson:"productId,omitempty"`
	MaxDiscount *float64 `json:"maxDiscount,omitempty" bson:"maxDiscount,omitempty"` // Cap on discount
}

// PromotionCreatedEvent represents a promotion creation event
type PromotionCreatedEvent struct {
	ID          string                    `json:"id" bson:"_id"`
	Code        string                    `json:"code" bson:"code"`
	Name        string                    `json:"name" bson:"name"`
	Type        string                    `json:"type" bson:"type"`
	Scope       string                    `json:"scope" bson:"scope"`
	Conditions  []PromotionConditionEvent `json:"conditions" bson:"conditions"`   // Structured conditions (required)
	Actions     []PromotionActionEvent    `json:"actions" bson:"actions"`         // Structured actions (required)
	Stackable   bool                      `json:"stackable" bson:"stackable"`     // Can stack with other promotions
	Exclusivity string                    `json:"exclusivity" bson:"exclusivity"` // "EXCLUSIVE", "STACKABLE"
	Priority    int                       `json:"priority" bson:"priority"`       // Higher priority = evaluated first
	ValidFrom   time.Time                 `json:"validFrom" bson:"validFrom"`
	ValidTo     *time.Time                `json:"validTo,omitempty" bson:"validTo,omitempty"`
	Status      Status                    `json:"status" bson:"status"`
	Version     int                       `json:"version" bson:"version"`
}

// PromotionUpdatedEvent represents a promotion update event
type PromotionUpdatedEvent struct {
	ID          string                    `json:"id" bson:"_id"`
	Code        string                    `json:"code,omitempty" bson:"code,omitempty"`
	Name        string                    `json:"name,omitempty" bson:"name,omitempty"`
	Type        string                    `json:"type,omitempty" bson:"type,omitempty"`
	Scope       string                    `json:"scope,omitempty" bson:"scope,omitempty"`
	Conditions  []PromotionConditionEvent `json:"conditions,omitempty" bson:"conditions,omitempty"`
	Actions     []PromotionActionEvent    `json:"actions,omitempty" bson:"actions,omitempty"`
	Stackable   *bool                     `json:"stackable,omitempty" bson:"stackable,omitempty"`
	Exclusivity string                    `json:"exclusivity,omitempty" bson:"exclusivity,omitempty"`
	Priority    *int                      `json:"priority,omitempty" bson:"priority,omitempty"`
	ValidFrom   *time.Time                `json:"validFrom,omitempty" bson:"validFrom,omitempty"`
	ValidTo     *time.Time                `json:"validTo,omitempty" bson:"validTo,omitempty"`
	Status      Status                    `json:"status,omitempty" bson:"status,omitempty"`
	Version     int                       `json:"version" bson:"version"`
}

// PromotionDeletedEvent represents a promotion deletion event
type PromotionDeletedEvent struct {
	ID      string `json:"id" bson:"_id"`
	Code    string `json:"code" bson:"code"`
	Version int    `json:"version" bson:"version"`
}
