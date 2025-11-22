package notification

type TemplateCode struct {
	ID          string   `json:"id" bson:"_id"`
	Code        string   `json:"code" bson:"code" validate:"required,template_code_format"`
	Description string   `json:"description" bson:"description" validate:"required,min=3,max=200"`
	IsActive    bool     `json:"isActive" bson:"isActive"`
	UsedBy      []string `json:"usedBy,omitempty" bson:"usedBy,omitempty"`
	Version     int      `json:"version" bson:"version"`
}
