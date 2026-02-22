package planservice

// ProductTemplateEvent represents a product template lifecycle event published by plan-service
type ProductTemplateEvent struct {
	ProductID   string   `json:"productId"`
	ProductCode string   `json:"productCode"`
	ProductName string   `json:"productName"`
	ProductType string   `json:"productType"` // SERVICE, GOODS, BUNDLE
	BasePrice   *float64 `json:"basePrice,omitempty"`
	Status      string   `json:"status"`
	TenantIDs   []string `json:"tenantIds,omitempty"`
	Scope       string   `json:"scope,omitempty"` // GLOBAL or TENANT
	IsActive    bool     `json:"isActive"`
	Version     int      `json:"version"`
}
