package authevent

type RoleCreateEventModel struct {
	ID                       string   `bson:"_id" json:"id"`
	Name                     string   `bson:"name" json:"name"`
	DisplayName              string   `bson:"displayName" json:"displayName"`
	Description              string   `bson:"description" json:"description"`
	SuperAdminAssignableOnly bool     `json:"superAdminAssignableOnly,omitempty" bson:"superAdminAssignableOnly,omitempty"`
	AssignedTo               []string `json:"assignedTo" bson:"assignedTo"`
	IsSystem                 bool     `bson:"isSystem" json:"isSystem"`
	IsActive                 bool     `bson:"isActive" json:"isActive"`
	IsVisible                bool     `bson:"isVisible" json:"isVisible"`
	Version                  int      `bson:"version" json:"version"`
}

type RoleUpdateEventModel struct {
	ID                       string   `bson:"_id" json:"id"`
	Name                     string   `bson:"name" json:"name"`
	DisplayName              string   `bson:"displayName" json:"displayName"`
	Description              string   `bson:"description" json:"description"`
	SuperAdminAssignableOnly bool     `json:"superAdminAssignableOnly,omitempty" bson:"superAdminAssignableOnly,omitempty"`
	AssignedTo               []string `json:"assignedTo" bson:"assignedTo"`
	IsSystem                 bool     `bson:"isSystem" json:"isSystem"`
	IsActive                 bool     `bson:"isActive" json:"isActive"`
	IsVisible                bool     `bson:"isVisible" json:"isVisible"`
	Version                  int      `bson:"version" json:"version"`
}

// RoleDeleteEventModel carries the role ID plus the snapshot of users the role
// was assigned to at delete time. Consumers (e.g. auth-service /auth/events SSE)
// match `assignedTo` against the authenticated user ID to push a
// `permissions.invalidated` event so each affected session reloads /auth/me.
type RoleDeleteEventModel struct {
	ID         string   `bson:"_id" json:"id"`
	AssignedTo []string `json:"assignedTo" bson:"assignedTo"`
}
