package logengineevent

import "time"

// TenantLogConfigEvent represents tenant-level log engine configuration
type TenantLogConfigEvent struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`

	// Log collection settings
	LogCollectionEnabled bool `json:"logCollectionEnabled"`

	// Retention settings
	RetentionDays    int    `json:"retentionDays"`    // How long to keep logs
	RetentionPolicy  string `json:"retentionPolicy"`  // HOT, WARM, COLD tiering
	ArchiveAfterDays int    `json:"archiveAfterDays"` // Move to cold storage after N days

	// Compaction settings
	CompactionEnabled    bool `json:"compactionEnabled"`
	CompactionIntervalHr int  `json:"compactionIntervalHr"`

	// Query settings
	QueryEnabled         bool `json:"queryEnabled"`
	MaxQueryRangeDays    int  `json:"maxQueryRangeDays"`    // Max time range for queries
	MaxResultsPerQuery   int  `json:"maxResultsPerQuery"`

	// Alert/notification settings
	AlertsEnabled bool   `json:"alertsEnabled"`
	AlertWebhook  string `json:"alertWebhook,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RetentionPolicy constants
const (
	RetentionPolicyStandard  = "STANDARD"   // Keep in hot storage
	RetentionPolicyTiered    = "TIERED"     // Move to warm/cold over time
	RetentionPolicyArchive   = "ARCHIVE"    // Move to cold storage quickly
	RetentionPolicyCompliance = "COMPLIANCE" // Extended retention for compliance
)

