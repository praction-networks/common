package tenantuserevent

import "time"

// OrgUnitAddress mirrors the address fields carried in OrgUnit events.
type OrgUnitAddress struct {
	Line1      string `json:"line1,omitempty"      bson:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"      bson:"line2,omitempty"`
	City       string `json:"city,omitempty"       bson:"city,omitempty"`
	State      string `json:"state,omitempty"      bson:"state,omitempty"`
	Country    string `json:"country,omitempty"    bson:"country,omitempty"`
	PostalCode string `json:"postalCode,omitempty" bson:"postalCode,omitempty"`
}

// OrgUnitContact mirrors the contact fields carried in OrgUnit events.
type OrgUnitContact struct {
	Email          string `json:"email,omitempty"          bson:"email,omitempty"`
	Phone          string `json:"phone,omitempty"          bson:"phone,omitempty"`
	AlternatePhone string `json:"alternatePhone,omitempty" bson:"alternatePhone,omitempty"`
}

// OrgUnitTax mirrors the tax fields carried in OrgUnit events.
type OrgUnitTax struct {
	GSTIN           string `json:"gstin,omitempty"           bson:"gstin,omitempty"`
	PAN             string `json:"pan,omitempty"             bson:"pan,omitempty"`
	TaxRegistration string `json:"taxRegistration,omitempty" bson:"taxRegistration,omitempty"`
}

// OrgUnitBusinessHours mirrors the business-hours fields carried in OrgUnit events.
type OrgUnitBusinessHours struct {
	Timezone  string `json:"timezone,omitempty"  bson:"timezone,omitempty"`
	OpenTime  string `json:"openTime,omitempty"  bson:"openTime,omitempty"`
	CloseTime string `json:"closeTime,omitempty" bson:"closeTime,omitempty"`
	WorkDays  []int  `json:"workDays,omitempty"  bson:"workDays,omitempty"`
}

// OrgUnitCreateEvent is published when an organisational unit is created.
// Consuming services pick the fields they need from the full document.
type OrgUnitCreateEvent struct {
	ID          string `json:"id"                    bson:"_id"`
	TenantID    string `json:"tenantId"              bson:"tenantId"`
	Type        string `json:"type"                  bson:"type"`
	Name        string `json:"name"                  bson:"name"`
	Code        string `json:"code"                  bson:"code"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`

	// Hierarchy (materialized path)
	ParentID  string   `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Path      string   `json:"path"               bson:"path"`
	Ancestors []string `json:"ancestors"          bson:"ancestors"`
	Level     int      `json:"level"              bson:"level"`

	// Leadership
	HeadUserID  string `json:"headUserId,omitempty" bson:"headUserId,omitempty"`
	MemberCount int64  `json:"memberCount"          bson:"memberCount"`

	// Operational
	IsHQ          bool                  `json:"isHQ"                    bson:"isHQ"`
	Currency      string                `json:"currency,omitempty"      bson:"currency,omitempty"`
	BusinessHours *OrgUnitBusinessHours `json:"businessHours,omitempty" bson:"businessHours,omitempty"`

	Address *OrgUnitAddress `json:"address,omitempty" bson:"address,omitempty"`
	Contact *OrgUnitContact `json:"contact,omitempty" bson:"contact,omitempty"`
	Tax     *OrgUnitTax     `json:"tax,omitempty"     bson:"tax,omitempty"`

	Tags  []string `json:"tags,omitempty"  bson:"tags,omitempty"`
	Color string   `json:"color,omitempty" bson:"color,omitempty"`

	IsActive  bool      `json:"isActive"              bson:"isActive"`
	Version   int64     `json:"version"               bson:"version"`
	CreatedAt time.Time `json:"createdAt"             bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"             bson:"updatedAt"`
	CreatedBy string    `json:"createdBy,omitempty"   bson:"createdBy,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"   bson:"updatedBy,omitempty"`
}

// OrgUnitUpdateEvent is published when an organisational unit is updated.
// Contains the full document after update — consuming services decide what to use.
type OrgUnitUpdateEvent struct {
	ID          string `json:"id"                    bson:"_id"`
	TenantID    string `json:"tenantId"              bson:"tenantId"`
	Type        string `json:"type"                  bson:"type"`
	Name        string `json:"name"                  bson:"name"`
	Code        string `json:"code"                  bson:"code"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`

	// Hierarchy (materialized path)
	ParentID  string   `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Path      string   `json:"path"               bson:"path"`
	Ancestors []string `json:"ancestors"          bson:"ancestors"`
	Level     int      `json:"level"              bson:"level"`

	// Leadership
	HeadUserID  string `json:"headUserId,omitempty" bson:"headUserId,omitempty"`
	MemberCount int64  `json:"memberCount"          bson:"memberCount"`

	// Operational
	IsHQ          bool                  `json:"isHQ"                    bson:"isHQ"`
	Currency      string                `json:"currency,omitempty"      bson:"currency,omitempty"`
	BusinessHours *OrgUnitBusinessHours `json:"businessHours,omitempty" bson:"businessHours,omitempty"`

	Address *OrgUnitAddress `json:"address,omitempty" bson:"address,omitempty"`
	Contact *OrgUnitContact `json:"contact,omitempty" bson:"contact,omitempty"`
	Tax     *OrgUnitTax     `json:"tax,omitempty"     bson:"tax,omitempty"`

	Tags  []string `json:"tags,omitempty"  bson:"tags,omitempty"`
	Color string   `json:"color,omitempty" bson:"color,omitempty"`

	IsActive  bool      `json:"isActive"              bson:"isActive"`
	Version   int64     `json:"version"               bson:"version"`
	CreatedAt time.Time `json:"createdAt"             bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"             bson:"updatedAt"`
	CreatedBy string    `json:"createdBy,omitempty"   bson:"createdBy,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"   bson:"updatedBy,omitempty"`
}

// OrgUnitDeleteEvent is published when an organisational unit is deleted.
type OrgUnitDeleteEvent struct {
	ID       string `json:"id"                bson:"_id"`
	TenantID string `json:"tenantId,omitempty" bson:"tenantId,omitempty"`
	Version  int64  `json:"version"           bson:"version"`
}
