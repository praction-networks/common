package planservice

// ProductTemplateEvent represents a product template lifecycle event published by plan-service
type ProductTemplateEvent struct {
	ProductID     string   `json:"productId"`
	ProductCode   string   `json:"productCode"`
	ProductName   string   `json:"productName"`
	ProductType   string   `json:"productType"` // SERVICE, GOODS, BUNDLE
	BasePrice     *float64 `json:"basePrice,omitempty"`
	Status        string   `json:"status"`
	OwnerTenantID *string  `json:"ownerTenantId,omitempty"` // Required: ISP or Enterprise tenant that owns this product
	TenantIDs     []string `json:"tenantIds,omitempty"`     // Owner is auto-included
	IsActive      bool     `json:"isActive"`
	Version       int      `json:"version"`
}
