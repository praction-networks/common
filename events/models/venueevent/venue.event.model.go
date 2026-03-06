package venueevent

import "time"

// OrderEventModel represents an order event published by venue-service
type OrderEventModel struct {
	ID             string           `json:"id"`
	TenantID       string           `json:"tenantId"`
	TableID        string           `json:"tableId"`
	SubscriberID   string           `json:"subscriberId,omitempty"`
	OrderNumber    string           `json:"orderNumber"`
	DisplayNumber  int              `json:"displayNumber"`
	Status         string           `json:"status"`
	Items          []OrderItemEvent `json:"items"`
	Subtotal       int64            `json:"subtotal"`
	DiscountAmount int64            `json:"discountAmount,omitempty"`
	TaxAmount      int64            `json:"taxAmount,omitempty"`
	TotalAmount    int64            `json:"totalAmount"`
	Currency       string           `json:"currency,omitempty"`
	Version        int              `json:"version"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

// OrderItemEvent represents an individual item in an order event
type OrderItemEvent struct {
	ID         string `json:"id"`
	MenuItemID string `json:"menuItemId"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int64  `json:"unitPrice"`
	TotalPrice int64  `json:"totalPrice"`
}

// MenuEventModel represents a menu event published by venue-service
type MenuEventModel struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
	Version  int    `json:"version"`
}
