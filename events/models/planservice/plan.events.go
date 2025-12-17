package planservice

import "time"

// PlanCreatedEvent represents a plan creation event
type PlanCreatedEvent struct {
	ID                 string    `json:"id" bson:"_id"`
	Code               string    `json:"code" bson:"code"`
	Name               string    `json:"name" bson:"name"`
	Description        string    `json:"description,omitempty" bson:"description,omitempty"`
	PlanType           PlanType  `json:"planType" bson:"planType"`
	BillingCycle       BillingCycle `json:"billingCycle" bson:"billingCycle"`
	Status             Status    `json:"status" bson:"status"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
	CreatedBy          string    `json:"createdBy" bson:"createdBy"`
	Version            int       `json:"version" bson:"version"`
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
	UpdatedAt          time.Time      `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy          string         `json:"updatedBy" bson:"updatedBy"`
	Version            int            `json:"version" bson:"version"`
}

// PlanDeletedEvent represents a plan deletion event
type PlanDeletedEvent struct {
	ID       string    `json:"id" bson:"_id"`
	Code     string    `json:"code" bson:"code"`
	Version  int       `json:"version" bson:"version"`
}

