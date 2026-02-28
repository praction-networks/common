package tenantuserevent

import "time"

// DeptTeamCreateEvent is published when a department/team is created.
// Consuming services pick the fields they need from the full document.
type DeptTeamCreateEvent struct {
	ID          string   `json:"id" bson:"_id"`
	TenantID    string   `json:"tenantId" bson:"tenantId"`
	Name        string   `json:"name" bson:"name"`
	Code        string   `json:"code" bson:"code"`
	Type        string   `json:"type" bson:"type"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	ParentID    string   `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Path        string   `json:"path,omitempty" bson:"path,omitempty"`
	Ancestors   []string `json:"ancestors,omitempty" bson:"ancestors,omitempty"`
	Level       int      `json:"level" bson:"level"`
	HeadUserID  string   `json:"headUserId,omitempty" bson:"headUserId,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Color       string   `json:"color,omitempty" bson:"color,omitempty"`
	MemberCount int      `json:"memberCount" bson:"memberCount"`
	IsActive    bool     `json:"isActive" bson:"isActive"`
	Version     int      `json:"version" bson:"version"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// DeptTeamUpdateEvent is published when a department/team is updated.
// Contains the full document after update — consuming services decide what to use.
type DeptTeamUpdateEvent struct {
	ID          string   `json:"id" bson:"_id"`
	TenantID    string   `json:"tenantId" bson:"tenantId"`
	Name        string   `json:"name" bson:"name"`
	Code        string   `json:"code" bson:"code"`
	Type        string   `json:"type" bson:"type"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	ParentID    string   `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Path        string   `json:"path,omitempty" bson:"path,omitempty"`
	Ancestors   []string `json:"ancestors,omitempty" bson:"ancestors,omitempty"`
	Level       int      `json:"level" bson:"level"`
	HeadUserID  string   `json:"headUserId,omitempty" bson:"headUserId,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Color       string   `json:"color,omitempty" bson:"color,omitempty"`
	MemberCount int      `json:"memberCount" bson:"memberCount"`
	IsActive    bool     `json:"isActive" bson:"isActive"`
	Version     int      `json:"version" bson:"version"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// DeptTeamDeleteEvent is published when a department/team is deleted.
type DeptTeamDeleteEvent struct {
	ID       string `json:"id" bson:"_id"`
	TenantID string `json:"tenantId,omitempty" bson:"tenantId,omitempty"`
	Version  int    `json:"version" bson:"version"`
}
