package productservice

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
	StatusDraft    Status = "DRAFT"
	StautsAssigned Status = "ASSIGNED"
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
	UnitYear   Unit = "YEAR"
	UnitMonth  Unit = "MONTH"
	UnitDay    Unit = "DAY"
	UnitHour   Unit = "HOUR"
	UnitMinute Unit = "MINUTE"
	UnitGB     Unit = "GB"
	UnitMB     Unit = "MB"
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