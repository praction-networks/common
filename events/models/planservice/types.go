package planservice

type Status string

const (
	StatusActive         Status = "ACTIVE"
	StatusInactive       Status = "INACTIVE"
	StatusDraft          Status = "DRAFT"
	StautsAssigned       Status = "ASSIGNED"
	StatusPendingApproval Status = "PENDING_APPROVAL"
	StatusArchived       Status = "ARCHIVED"
	StatusDeprecated     Status = "DEPRECATED" // Old plan, no new subscriptions
	StatusSunset         Status = "SUNSET"     // Being phased out
)


type ProductType string

const (
	ProductTypeService ProductType = "SERVICE"
	ProductTypeGoods   ProductType = "GOODS"
	ProductTypeInstlalation ProductType = "INSTLALATION"
	ProductTypeOther   ProductType = "OTHER"
)

type Unit string

const (
	// Time-based units
	UnitYear   Unit = "YEAR"
	UnitMonth  Unit = "MONTH"
	UnitDay    Unit = "DAY"
	UnitHour   Unit = "HOUR"
	UnitMinute Unit = "MINUTE"
	
	// Data units
	UnitGB     Unit = "GB"
	UnitMB     Unit = "MB"
	UnitTB     Unit = "TB"     // Terabytes
	UnitKB     Unit = "KB"     // Kilobytes
	
	// Bandwidth units
	UnitMbps   Unit = "MBPS"   // Megabits per second
	UnitKbps   Unit = "KBPS"   // Kilobits per second
	UnitGbps   Unit = "GBPS"   // Gigabits per second
	
	// Voice/Communication units
	UnitMinutes Unit = "MINUTES" // Voice minutes (separate from time)
	UnitSMS     Unit = "SMS"     // SMS count
	UnitMMS     Unit = "MMS"     // MMS count
	UnitCall    Unit = "CALL"    // Call count
	
	// Hardware/Resource units
	UnitSim        Unit = "SIM"        // SIM card count
	UnitDevice     Unit = "DEVICE"      // Device count
	UnitLocation   Unit = "LOCATION"   // Site/location count
	UnitConcurrent Unit = "CONCURRENT"  // Concurrent connections
	
	// Generic
	UnitCount  Unit = "COUNT"
)

type BillingModel string

const (
	BillingModelRecurring BillingModel = "RECURRING"
	BillingModelOneTime   BillingModel = "ONE_TIME"
	BillingModelUsage     BillingModel = "USAGE"
)

type TaxMode string

const (
	TaxModeAuto  TaxMode = "AUTO"
	TaxModeFixed TaxMode = "FIXED"
)

type PriceBookScope string

const (
	PriceBookScopeGlobal PriceBookScope = "GLOBAL"
	PriceBookScopeTenant PriceBookScope = "TENANT"
	PriceBookScopeZone PriceBookScope = "ZONE"
)

type PlanType string

const (
	PlanTypePrepaid  PlanType = "PREPAID"
	PlanTypePostpaid PlanType = "POSTPAID"
	PlanTypeHybrid   PlanType = "HYBRID"
)

type BillingCycle string

const (
	BillingCycleHourly    BillingCycle = "HOURLY"
	BillingCycleDaily     BillingCycle = "DAILY"
	BillingCycleWeekly    BillingCycle = "WEEKLY"
	BillingCycleMonthly   BillingCycle = "MONTHLY"
	BillingCycleQuarterly BillingCycle = "QUARTERLY"
	BillingCycleSemiAnnually BillingCycle = "SEMI_ANNUALLY"
	BillingCycleYearly    BillingCycle = "YEARLY"
	BillingCycleOneTime   BillingCycle = "ONE_TIME"
)

type PriceStrategy string

const (
	PriceStrategyDefault  PriceStrategy = "DEFAULT"
	PriceStrategyOverride PriceStrategy = "OVERRIDE"
	PriceStrategyDerived  PriceStrategy = "DERIVED"
)

type OverrideTargetType string

const (
	OverrideTargetProduct OverrideTargetType = "PRODUCT"
	OverrideTargetPlan    OverrideTargetType = "PLAN"
	OverrideTargetPriceBook OverrideTargetType = "PRICE_BOOK"
)

// PricingModelType represents the type of pricing model
type PricingModelType string

const (
	PricingModelTypeFlat   PricingModelType = "FLAT"   // Single flat price
	PricingModelTypeTiered PricingModelType = "TIERED" // Tiered pricing (quantity-based)
	PricingModelTypeVolume PricingModelType = "VOLUME"  // Volume discount pricing
	PricingModelTypeBundle PricingModelType = "BUNDLE" // Bundle pricing
)

// ActivationMethodType represents activation method
type ActivationMethodType string

const (
	ActivationMethodTypeImmediate ActivationMethodType = "IMMEDIATE" // Activate immediately
	ActivationMethodTypeScheduled ActivationMethodType = "SCHEDULED" // Schedule activation
	ActivationMethodTypeManual    ActivationMethodType = "MANUAL"   // Manual activation required
)

// ProrationMethodType represents proration calculation method
type ProrationMethodType string

const (
	ProrationMethodTypeDaily       ProrationMethodType = "DAILY"       // Daily proration
	ProrationMethodTypeHourly      ProrationMethodType = "HOURLY"      // Hourly proration
	ProrationMethodTypeProportional ProrationMethodType = "PROPORTIONAL" // Proportional proration
)

// ProrationRoundingType represents proration rounding method
type ProrationRoundingType string

const (
	ProrationRoundingTypeUp      ProrationRoundingType = "UP"      // Round up
	ProrationRoundingTypeDown    ProrationRoundingType = "DOWN"   // Round down
	ProrationRoundingTypeNearest ProrationRoundingType = "NEAREST" // Round to nearest
)

// TrialType represents trial period type
type TrialType string

const (
	TrialTypeFree       TrialType = "FREE"       // Free trial
	TrialTypeDiscounted TrialType = "DISCOUNTED" // Discounted trial
)

// RenewalMethodType represents renewal method
type RenewalMethodType string

const (
	RenewalMethodTypeAuto    RenewalMethodType = "AUTO"    // Automatic renewal
	RenewalMethodTypeManual  RenewalMethodType = "MANUAL"  // Manual renewal required
	RenewalMethodTypeOptional RenewalMethodType = "OPTIONAL" // Optional renewal
)

// UpgradeType represents upgrade option type
type UpgradeType string

const (
	UpgradeTypeAuto       UpgradeType = "AUTO"       // Automatic upgrade
	UpgradeTypeOptional   UpgradeType = "OPTIONAL"   // Optional upgrade
	UpgradeTypeRecommended UpgradeType = "RECOMMENDED" // Recommended upgrade
)

// CreditLimitType represents credit limit type
type CreditLimitType string

const (
	CreditLimitTypeHard CreditLimitType = "HARD" // Hard limit - service stops
	CreditLimitTypeSoft CreditLimitType = "SOFT" // Soft limit - warnings only
	CreditLimitTypeNone CreditLimitType = "NONE" // No credit limit
)

// PaymentTermsType represents payment terms
type PaymentTermsType string

const (
	PaymentTermsTypePrepaid  PaymentTermsType = "PREPAID"  // Prepaid payment
	PaymentTermsTypePostpaid PaymentTermsType = "POSTPAID" // Postpaid payment
	PaymentTermsTypeHybrid   PaymentTermsType = "HYBRID"   // Hybrid payment
)

// ResourceType represents resource allocation type
type ResourceType string

const (
	ResourceTypeIPAddress ResourceType = "IP_ADDRESS"
	ResourceTypeVLAN      ResourceType = "VLAN"
	ResourceTypePort      ResourceType = "PORT"
	ResourceTypeEquipment ResourceType = "EQUIPMENT"
	ResourceTypeBandwidth ResourceType = "BANDWIDTH"
	ResourceTypeCircuit   ResourceType = "CIRCUIT"
)

// AllocationType represents resource allocation method
type AllocationType string

const (
	AllocationTypeDedicated AllocationType = "DEDICATED" // Dedicated resource
	AllocationTypeShared    AllocationType = "SHARED"    // Shared resource
	AllocationTypePool      AllocationType = "POOL"      // Pool-based allocation
)

// MigrationType represents plan migration type
type MigrationType string

const (
	MigrationTypeUpgrade   MigrationType = "UPGRADE"   // Upgrade migration
	MigrationTypeDowngrade MigrationType = "DOWNGRADE" // Downgrade migration
	MigrationTypeLateral   MigrationType = "LATERAL"   // Lateral migration
	MigrationTypeCustom    MigrationType = "CUSTOM"    // Custom migration
)

// PlanItemRole represents the role of an item in a plan
type PlanItemRole string

const (
	PlanItemRoleBase     PlanItemRole = "BASE"     // Base plan item
	PlanItemRoleAddon    PlanItemRole = "ADDON"    // Add-on item
	PlanItemRoleBundle   PlanItemRole = "BUNDLE"   // Bundle item
	PlanItemRoleOptional PlanItemRole = "OPTIONAL" // Optional item
)

