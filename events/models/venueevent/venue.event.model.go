package venueevent

import "time"

// OrderStatus represents the status of an order in venue events
type OrderStatus string

// Order status constants — remote services use these to decide behavior
const (
	OrderStatusPlaced    OrderStatus = "PLACED"
	OrderStatusAccepted  OrderStatus = "ACCEPTED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusReady     OrderStatus = "READY"
	OrderStatusServed    OrderStatus = "SERVED"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderEventModel represents an order event published by venue-service.
// Remote services read the Status field to decide how to handle the event.
type OrderEventModel struct {
	ID             string           `json:"id" bson:"_id"`
	TenantID       string           `json:"tenantId" bson:"tenantId"`
	TableID        string           `json:"tableId" bson:"tableId"`
	SubscriberID   string           `json:"subscriberId,omitempty" bson:"subscriberId,omitempty"`
	OrderNumber    string           `json:"orderNumber" bson:"orderNumber"`
	DisplayNumber  int              `json:"displayNumber" bson:"displayNumber"`
	Status         OrderStatus      `json:"status" bson:"status"`
	Items          []OrderItemEvent `json:"items" bson:"items"`
	Subtotal       int64            `json:"subtotal" bson:"subtotal"`
	DiscountAmount int64            `json:"discountAmount,omitempty" bson:"discountAmount,omitempty"`
	TaxAmount      int64            `json:"taxAmount,omitempty" bson:"taxAmount,omitempty"`
	TotalAmount    int64            `json:"totalAmount" bson:"totalAmount"`
	Currency       string           `json:"currency,omitempty" bson:"currency,omitempty"`
	Version        int              `json:"version" bson:"version"`
	CreatedAt      time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt" bson:"updatedAt"`
}

// OrderItemEvent represents an individual item in an order event
type OrderItemEvent struct {
	ID         string `json:"id" bson:"id"`
	MenuItemID string `json:"menuItemId" bson:"menuItemId"`
	Name       string `json:"name" bson:"name"`
	Quantity   int    `json:"quantity" bson:"quantity"`
	UnitPrice  int64  `json:"unitPrice" bson:"unitPrice"`
	TotalPrice int64  `json:"totalPrice" bson:"totalPrice"`
}

// MenuEventModel represents a menu event published by venue-service
type MenuEventModel struct {
	ID       string `json:"id" bson:"_id"`
	TenantID string `json:"tenantId" bson:"tenantId"`
	Name     string `json:"name" bson:"name"`
	IsActive bool   `json:"isActive" bson:"isActive"`
	Version  int    `json:"version" bson:"version"`
}
