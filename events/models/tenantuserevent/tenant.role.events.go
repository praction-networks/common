package tenantuserevent

type TenantRoleCreateEvent struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string   `json:"name" bson:"name" validate:"required"`
	PolicyIDs []string `json:"policyIds" bson:"policyIds" validate:"required,dive"` // Associated policy IDs
	IsSystem  bool     `json:"isSystem" bson:"isSystem"`
	Version   int      `json:"version" bson:"version"`
}

type TenantRoleUpdateEvent struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string   `json:"name" bson:"name" validate:"required"`
	PolicyIDs []string `json:"policyIds" bson:"policyIds" validate:"required,dive"` // Associated policy IDs
	IsSystem  bool     `json:"isSystem" bson:"isSystem"`                            // System roles cannot be modified
	IsActive  bool     `json:"isActive" bson:"isActive"`
	IsVisible bool     `json:"isVisible" bson:"isVisible"` // Visible to users for assignment
}

type TenantRoleDeleteEvent struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
}
