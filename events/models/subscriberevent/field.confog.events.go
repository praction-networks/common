package subscriberevent

// FieldConfigItem represents a single field configuration item in events
type FieldConfigItem struct {
	FieldKey          string                `json:"fieldKey" bson:"fieldKey"`
	IsVisible         bool                  `json:"isVisible" bson:"isVisible"`
	IsMandatory       bool                  `json:"isMandatory" bson:"isMandatory"`
	IsReadOnly        bool                  `json:"isReadOnly" bson:"isReadOnly"`
	IsHidden          bool                  `json:"isHidden" bson:"isHidden"`
	DisplayOrder      int                   `json:"displayOrder" bson:"displayOrder"`
	CustomLabel       *string               `json:"customLabel,omitempty" bson:"customLabel,omitempty"`
	CustomPlaceholder *string               `json:"customPlaceholder,omitempty" bson:"customPlaceholder,omitempty"`
	ValidationRules   *FieldValidationRules `json:"validationRules,omitempty" bson:"validationRules,omitempty"`
	ConditionalRules  *ConditionalRules     `json:"conditionalRules,omitempty" bson:"conditionalRules,omitempty"`
}

// FieldValidationRules represents validation rules for a field in events
type FieldValidationRules struct {
	MinLength   *int     `json:"minLength,omitempty" bson:"minLength,omitempty"`
	MaxLength   *int     `json:"maxLength,omitempty" bson:"maxLength,omitempty"`
	Pattern     *string  `json:"pattern,omitempty" bson:"pattern,omitempty"`
	MinValue    *float64 `json:"minValue,omitempty" bson:"minValue,omitempty"`
	MaxValue    *float64 `json:"maxValue,omitempty" bson:"maxValue,omitempty"`
	CustomError *string  `json:"customError,omitempty" bson:"customError,omitempty"`
}

// ConditionalRules represents conditional display/requirement logic in events
type ConditionalRules struct {
	ShowIf     []Condition `json:"showIf,omitempty" bson:"showIf,omitempty"`
	RequiredIf []Condition `json:"requiredIf,omitempty" bson:"requiredIf,omitempty"`
	DisabledIf []Condition `json:"disabledIf,omitempty" bson:"disabledIf,omitempty"`
}

// Condition represents a single condition for conditional logic in events
type Condition struct {
	FieldKey string      `json:"fieldKey" bson:"fieldKey"`
	Operator string      `json:"operator" bson:"operator"`
	Value    interface{} `json:"value" bson:"value"`
}

// FieldConfigCreatedEvent represents a field configuration creation event
type FieldConfigCreatedEvent struct {
	ID       string            `json:"id" bson:"id"`
	TenantID string            `json:"tenantId" bson:"tenantId"`
	FormType string            `json:"formType" bson:"formType"` // subscriber, broadband, hotspot
	Fields   []FieldConfigItem `json:"fields" bson:"fields"`
}

// FieldConfigUpdatedEvent represents a field configuration update event
type FieldConfigUpdatedEvent struct {
	ID       string            `json:"id" bson:"id"`
	TenantID string            `json:"tenantId" bson:"tenantId"`
	FormType string            `json:"formType" bson:"formType"`
	Fields   []FieldConfigItem `json:"fields" bson:"fields"`
	Version  int               `json:"version" bson:"version"`
}

// FieldConfigDeletedEvent represents a field configuration deletion event
type FieldConfigDeletedEvent struct {
	ID       string `json:"id" bson:"id"`
	TenantID string `json:"tenantId" bson:"tenantId"`
	FormType string `json:"formType" bson:"formType"`
}
