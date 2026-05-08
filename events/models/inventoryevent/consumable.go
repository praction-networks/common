package inventoryevent

// Consumable is a catalogue entry. SKU may be tenant-local or global per
// inventory-service config. Cables track in metres, connectors in counts —
// see assets-prd.md §7 for unit semantics.
type Consumable struct {
	ID       string `json:"id"            bson:"_id"`
	TenantID string `json:"tenantId"      bson:"tenantId"`
	Name     string `json:"name"          bson:"name"`
	Unit     string `json:"unit"          bson:"unit"` // "metre" | "count"
	SKU      string `json:"sku,omitempty" bson:"sku,omitempty"`
}
