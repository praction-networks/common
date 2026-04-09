package constants

// TenantType represents the type of a tenant in the system
type TenantType string

const (
	TenantTypeSystem      TenantType = "System"      // Root super tenant - manages all other tenants
	TenantTypeDistributor TenantType = "Distributor" // Distributor distributes services to Partner
	TenantTypeReseller    TenantType = "Reseller"    // Reseller resells services to Distributor
	TenantTypePartner     TenantType = "Partner"     // Partner partners with ISP to sell services
	TenantTypeISP         TenantType = "ISP"         // ISP provides services to end users
	TenantTypeEnterprise  TenantType = "Enterprise"  // Enterprise provides services to own end customers
	TenantTypeBranch      TenantType = "Branch"      // Branch of an Enterprise
)

// IsRootTenantType returns true if the tenant type is a root-level type (ISP or Enterprise)
func IsRootTenantType(tenantType string) bool {
	return TenantType(tenantType) == TenantTypeISP || TenantType(tenantType) == TenantTypeEnterprise
}

// Role name constants
const (
	RoleTenantAdmin = "tenant_admin"
	RoleSuperAdmin  = "super_admin"
)
