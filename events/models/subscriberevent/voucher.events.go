package subscriberevent

import "time"

// VoucherCreatedEvent represents a voucher instance creation event
// Published when a new voucher is inserted into MongoDB
type VoucherCreatedEvent struct {
	ID               string             `json:"id" bson:"_id"`
	OwnerTenantID    string             `json:"ownerTenantId" bson:"ownerTenantId"`
	Scope            VoucherScope       `json:"scope" bson:"scope"`
	AllowedTenantIDs []string           `json:"allowedTenantIds,omitempty" bson:"allowedTenantIds,omitempty"`
	TemplateID       string             `json:"templateId" bson:"templateId"`
	PlanID           string             `json:"planId" bson:"planId"`
	VoucherCode      string             `json:"voucherCode" bson:"voucherCode"`
	Status           string             `json:"status" bson:"status"`
	BatchID          *string            `json:"batchId,omitempty" bson:"batchId,omitempty"`
	BatchStatus      VoucherBatchStatus `json:"batchStatus" bson:"batchStatus"`
	IssuedAt         *time.Time         `json:"issuedAt,omitempty" bson:"issuedAt,omitempty"`
	ExpiresAt        time.Time          `json:"expiresAt" bson:"expiresAt"`
	Distribution     string             `json:"distributionType" bson:"distributionType"`
	Version          int                `json:"version" bson:"version"`
}

// VoucherUpdatedEvent represents a voucher instance update event
// Published for any voucher update (status change, batch status change, usage, expiry, revocation, extension)
// Contains the full document state from MongoDB so consumers can upsert
type VoucherUpdatedEvent struct {
	ID               string             `json:"id" bson:"_id"`
	OwnerTenantID    string             `json:"ownerTenantId" bson:"ownerTenantId"`
	Scope            VoucherScope       `json:"scope" bson:"scope"`
	AllowedTenantIDs []string           `json:"allowedTenantIds,omitempty" bson:"allowedTenantIds,omitempty"`
	TemplateID       string             `json:"templateId" bson:"templateId"`
	PlanID           string             `json:"planId" bson:"planId"`
	VoucherCode      string             `json:"voucherCode" bson:"voucherCode"`
	Status           string             `json:"status" bson:"status"`
	BatchID          *string            `json:"batchId,omitempty" bson:"batchId,omitempty"`
	BatchStatus      VoucherBatchStatus `json:"batchStatus" bson:"batchStatus"`
	ExpiresAt        time.Time          `json:"expiresAt" bson:"expiresAt"`
	Version          int                `json:"version" bson:"version"`
}

// VoucherDeletedEvent represents a voucher instance deletion event
// Published when a voucher is deleted from MongoDB
type VoucherDeletedEvent struct {
	ID            string `json:"id" bson:"_id"`
	OwnerTenantID string `json:"ownerTenantId" bson:"ownerTenantId"`
	VoucherCode   string `json:"voucherCode" bson:"voucherCode"`
	Version       int    `json:"version" bson:"version"`
}
