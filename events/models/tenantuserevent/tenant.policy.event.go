package tenantuserevent

type TenantPolicyCreateEvent struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name" validate:"required"`
	Action   string `json:"action" bson:"action" validate:"required"`
	Service  string `json:"service" bson:"service" validate:"required"`
	Resource string `json:"resource" bson:"resource" validate:"required"`
	IsActive bool   `json:"isActive" bson:"isActive"`
	IsSystem bool   `json:"isSystem" bson:"isSystem"`
	Version  int    `json:"version" bson:"version"`
}

type TenantPolicyUpdateEvent struct {
	ID               string `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string `json:"name" bson:"name" validate:"required"`
	Action           string `json:"action" bson:"action" validate:"required"`
	Service          string `json:"service" bson:"service" validate:"required"`
	Resource         string `json:"resource" bson:"resource" validate:"required"`
	IsAssignedToRole bool   `json:"isAssignedToRole" bson:"isAssignedToRole"`
	IsActive         bool   `json:"isActive" bson:"isActive"`
	IsSystem         bool   `json:"isSystem" bson:"isSystem"`
	Version          int    `json:"version" bson:"version"`
}

type TenantPolicyDeleteEvent struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
}
