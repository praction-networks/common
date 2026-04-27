package planservice

import "github.com/praction-networks/common/models/money"

// ProductTaxSnapshot carries the resolved tax rate + GST components for a
// taxable product at the moment an event is published. Downstream services
// (billing) read this directly without having to fetch the tax-rate model
// separately. When an admin edits the underlying TaxRateModel, plan-service
// re-publishes every product referencing it so the snapshot stays current
// — the snapshot is a denormalised cache, not a point-in-time pin.
type ProductTaxSnapshot struct {
	TaxRateID string  `json:"taxRateId"`          // TaxRateModel._id in plan-service
	TaxGroup  string  `json:"taxGroup,omitempty"` // e.g. "GST"
	Model     string  `json:"model,omitempty"`    // e.g. "GST"
	CGST      float64 `json:"cgst"`               // central GST %
	SGST      float64 `json:"sgst"`               // state GST %
	IGST      float64 `json:"igst"`               // inter-state GST %
	// TotalRate = CGST+SGST for intra-state, or IGST for inter-state.
	// Cached here so billing never has to recompute.
	TotalRate float64 `json:"totalRate"`
}

// ProductTemplateEvent represents a product template lifecycle event published by plan-service
type ProductTemplateEvent struct {
	ProductID     string   `json:"productId"`
	ProductCode   string   `json:"productCode"`
	ProductName   string   `json:"productName"`
	ProductType   string   `json:"productType"` // SERVICE, GOODS, BUNDLE
	BasePrice     *money.Money `json:"basePrice,omitempty"` // paise
	Status        string   `json:"status"`
	OwnerTenantID *string  `json:"ownerTenantId,omitempty"` // Required: ISP or Enterprise tenant that owns this product
	TenantIDs     []string `json:"tenantIds,omitempty"`     // Owner is auto-included
	IsActive      bool     `json:"isActive"`

	// Tax — product-driven and authoritative. IsTaxable + TaxRateID come
	// straight from ProductModel; Tax is the resolved snapshot (CGST/SGST
	// /IGST %) so billing doesn't have to cache TaxRateModel separately.
	// When IsTaxable=false both TaxRateID and Tax are nil.
	IsTaxable bool                `json:"isTaxable"`
	TaxRateID *string             `json:"taxRateId,omitempty"`
	Tax       *ProductTaxSnapshot `json:"tax,omitempty"`

	// PriceIsInclusive — declares whether Product.Price already contains
	// GST (true) or GST is added on top at billing time (false). Flows
	// through plan → price book → invoice so the interpretation stays
	// consistent catalogue-through-billing.
	PriceIsInclusive bool `json:"priceIsInclusive"`

	Version int `json:"version"`
}
