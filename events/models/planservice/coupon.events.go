package planservice

import "time"

// CouponCreatedEvent represents a coupon creation event
type CouponCreatedEvent struct {
	ID           string    `json:"id" bson:"_id"`
	Code         string    `json:"code" bson:"code"`
	PromotionID  string    `json:"promotionId" bson:"promotionId"`
	GlobalLimit  *int      `json:"globalLimit,omitempty" bson:"globalLimit,omitempty"`
	PerUserLimit *int      `json:"perUserLimit,omitempty" bson:"perUserLimit,omitempty"`
	ValidFrom    time.Time `json:"validFrom" bson:"validFrom"`
	ValidTo      *time.Time `json:"validTo,omitempty" bson:"validTo,omitempty"`
	Status       Status    `json:"status" bson:"status"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	CreatedBy    string    `json:"createdBy" bson:"createdBy"`
	Version      int       `json:"version" bson:"version"`
}

// CouponUpdatedEvent represents a coupon update event
type CouponUpdatedEvent struct {
	ID           string     `json:"id" bson:"_id"`
	Code         string     `json:"code,omitempty" bson:"code,omitempty"`
	PromotionID  string     `json:"promotionId,omitempty" bson:"promotionId,omitempty"`
	GlobalLimit  *int       `json:"globalLimit,omitempty" bson:"globalLimit,omitempty"`
	PerUserLimit *int       `json:"perUserLimit,omitempty" bson:"perUserLimit,omitempty"`
	ValidFrom    *time.Time `json:"validFrom,omitempty" bson:"validFrom,omitempty"`
	ValidTo      *time.Time `json:"validTo,omitempty" bson:"validTo,omitempty"`
	Status       Status     `json:"status,omitempty" bson:"status,omitempty"`
	UpdatedAt    time.Time  `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy    string     `json:"updatedBy" bson:"updatedBy"`
	Version      int        `json:"version" bson:"version"`
}

// CouponDeletedEvent represents a coupon deletion event
type CouponDeletedEvent struct {
	ID       string    `json:"id" bson:"_id"`
	Code     string    `json:"code" bson:"code"`
	Version  int       `json:"version" bson:"version"`
}

