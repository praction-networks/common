package tenantevent

import "fmt"

// PolicyAccount is the per-tenant account-management policy bucket.
//
// Source: backend-contract §11.5, Sweep 4 Q2 (emailChangeAllowed) plus the
// 7 deferred-but-listed §11.5 keys land here so admin UI can configure them
// before each consumer service wires them up.
//
// Bool fields use *bool because false ≠ "not set" matters for sparse PATCH
// semantics — a nil pointer means "client didn't send the field"; a pointer
// to false means "client explicitly set it to false". MergeTenantPolicy
// preserves this distinction.
type PolicyAccount struct {
	EmailChangeAllowed   *bool    `json:"emailChangeAllowed,omitempty"   bson:"emailChangeAllowed,omitempty"`
	AllowMobileChange    *bool    `json:"allowMobileChange,omitempty"    bson:"allowMobileChange,omitempty"`
	AllowWhatsappChange  *bool    `json:"allowWhatsappChange,omitempty"  bson:"allowWhatsappChange,omitempty"`
	FuelMode             string   `json:"fuelMode,omitempty"             bson:"fuelMode,omitempty"` // "" | "FLAT" | "PER_KM"
	FuelRatePerKm        *float64 `json:"fuelRatePerKm,omitempty"        bson:"fuelRatePerKm,omitempty"`
	LeaveCategories      []string `json:"leaveCategories,omitempty"      bson:"leaveCategories,omitempty"`
	VehicleRequired      *bool    `json:"vehicleRequired,omitempty"      bson:"vehicleRequired,omitempty"`
	DocumentVaultEnabled *bool    `json:"documentVaultEnabled,omitempty" bson:"documentVaultEnabled,omitempty"`
}

func (a PolicyAccount) Validate() error {
	switch a.FuelMode {
	case "", "FLAT", "PER_KM":
		return nil
	default:
		return fmt.Errorf("FuelMode must be FLAT or PER_KM (got %q)", a.FuelMode)
	}
}
