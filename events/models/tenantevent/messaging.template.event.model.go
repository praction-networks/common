package tenantevent

import "time"

// TenantTemplateConfig represents tenant-specific configuration for a template
// This allows enabling/disabling a template for a specific tenant
// Note: Sender ID and DLT Template ID are SMS-specific and stored in MessagingSMSTemplate, not tenant-specific
type TenantTemplateConfig struct {
	TenantID string `bson:"tenantID" json:"tenantID"`
	IsActive bool   `bson:"isActive" json:"isActive"` // Can disable template for specific tenant
}

// MessagingTemplateInsertEventModel represents a messaging template creation event
type MessagingTemplateInsertEventModel struct {
	ID               string                 `bson:"_id" json:"id"`
	Name             string                 `bson:"name" json:"name"`
	Code             string                 `bson:"code" json:"code"`
	Channel          string                 `bson:"channel" json:"channel"`
	Language         string                 `bson:"language" json:"language"`
	Description      string                 `bson:"description,omitempty" json:"description,omitempty"`
	IsActive         bool                   `bson:"isActive" json:"isActive"`
	IsSystemTemplate bool                   `bson:"isSystemTemplate" json:"isSystemTemplate"` // If true, this template is for system notifications only. Can only be assigned by system users.
	Tags             []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	AssignedTo       []string               `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"`
	TenantConfigs    []TenantTemplateConfig `bson:"tenantConfigs,omitempty" json:"tenantConfigs,omitempty"` // Tenant-specific template configurations (array) - only IsActive flag
	AssignableBy     []string               `bson:"assignableBy,omitempty" json:"assignableBy,omitempty"`   // Which tenant types can assign this template. Ignored if IsSystemTemplate is true.

	// Channel-specific template data (contains content and variables)
	SMS      *MessagingSMSTemplate      `bson:"sms,omitempty" json:"sms,omitempty"`
	Email    *MessagingEmailTemplate    `bson:"email,omitempty" json:"email,omitempty"`
	Telegram *MessagingTelegramTemplate `bson:"telegram,omitempty" json:"telegram,omitempty"`
	WhatsApp *MessagingWhatsAppTemplate `bson:"whatsapp,omitempty" json:"whatsapp,omitempty"`

	Version int `bson:"version" json:"version"`
}

// MessagingTemplateUpdateEventModel represents a messaging template update event
type MessagingTemplateUpdateEventModel struct {
	ID               string                 `bson:"_id" json:"id"`
	Name             string                 `bson:"name" json:"name"`
	Code             string                 `bson:"code" json:"code"`
	Channel          string                 `bson:"channel" json:"channel"`
	Language         string                 `bson:"language" json:"language"`
	Description      string                 `bson:"description,omitempty" json:"description,omitempty"`
	IsActive         bool                   `bson:"isActive" json:"isActive"`
	IsSystemTemplate bool                   `bson:"isSystemTemplate" json:"isSystemTemplate"` // If true, this template is for system notifications only. Can only be assigned by system users.
	Tags             []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	AssignedTo       []string               `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"`
	TenantConfigs    []TenantTemplateConfig `bson:"tenantConfigs,omitempty" json:"tenantConfigs,omitempty"` // Tenant-specific template configurations (array) - only IsActive flag
	AssignableBy     []string               `bson:"assignableBy,omitempty" json:"assignableBy,omitempty"`   // Which tenant types can assign this template. Ignored if IsSystemTemplate is true.

	// Channel-specific template data (contains content and variables)
	SMS      *MessagingSMSTemplate      `bson:"sms,omitempty" json:"sms,omitempty"`
	Email    *MessagingEmailTemplate    `bson:"email,omitempty" json:"email,omitempty"`
	Telegram *MessagingTelegramTemplate `bson:"telegram,omitempty" json:"telegram,omitempty"`
	WhatsApp *MessagingWhatsAppTemplate `bson:"whatsapp,omitempty" json:"whatsapp,omitempty"`

	Version   int        `bson:"version" json:"version"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"` // Soft delete timestamp
}

// MessagingTemplateDeleteEventModel represents a messaging template deletion event
type MessagingTemplateDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}

// MessagingTemplateVariable represents a template variable definition
type MessagingTemplateVariable struct {
	Key          string      `bson:"key" json:"key"`
	Placeholder  *int        `bson:"placeholder,omitempty" json:"placeholder,omitempty"`
	Type         string      `bson:"type" json:"type"`
	Required     bool        `bson:"required" json:"required"`
	DefaultValue interface{} `bson:"default_value,omitempty" json:"default_value,omitempty"`
	Description  string      `bson:"description,omitempty" json:"description,omitempty"`
}

// MessagingSMSTemplate represents SMS-specific template data
// Sender ID and DLT Template ID are SMS-specific (not tenant-specific)
// These values are used when sending SMS for any tenant using this template
type MessagingSMSTemplate struct {
	Content       string                      `bson:"content" json:"content"`
	Variables     []MessagingTemplateVariable `bson:"variables" json:"variables"`
	DLTTemplateID string                      `bson:"dlt_template_id,omitempty" json:"dlt_template_id,omitempty"` // DLT template ID for SMS (SMS-specific, not tenant-specific)
	SenderID      string                      `bson:"sender_id,omitempty" json:"sender_id,omitempty"`             // Sender ID for SMS (SMS-specific, not tenant-specific)
	MaxLength     int                         `bson:"max_length,omitempty" json:"max_length,omitempty"`           // SMS character limit (default: 160)
}

// MessagingEmailTemplate represents Email-specific template data
type MessagingEmailTemplate struct {
	Subject     string                      `bson:"subject" json:"subject"`
	TextBody    string                      `bson:"text_body,omitempty" json:"text_body,omitempty"`
	HTMLBody    string                      `bson:"html_body" json:"html_body"`
	PreviewText string                      `bson:"preview_text,omitempty" json:"preview_text,omitempty"` // Email preview text
	Variables   []MessagingTemplateVariable `bson:"variables" json:"variables"`
}

// MessagingTelegramTemplate represents Telegram-specific template data
type MessagingTelegramTemplate struct {
	Body      string                      `bson:"body" json:"body"`
	ParseMode string                      `bson:"parse_mode,omitempty" json:"parse_mode,omitempty"`
	Variables []MessagingTemplateVariable `bson:"variables" json:"variables"`
}

// MessagingWhatsAppTemplate represents Meta-compatible WhatsApp template structure
type MessagingWhatsAppTemplate struct {
	Namespace    string                       `bson:"namespace,omitempty" json:"namespace,omitempty"`
	Name         string                       `bson:"name" json:"name"`
	Category     string                       `bson:"category" json:"category"`
	LanguageCode string                       `bson:"language_code" json:"language_code"`
	Status       string                       `bson:"status" json:"status"`
	Components   []MessagingWhatsAppComponent `bson:"components" json:"components"`
	ProviderMeta map[string]any               `bson:"provider_meta,omitempty" json:"provider_meta,omitempty"`
	Variables    []MessagingTemplateVariable  `bson:"variables" json:"variables"`
}

// MessagingWhatsAppComponent represents a component in WhatsApp template
type MessagingWhatsAppComponent struct {
	Type    string                    `bson:"type" json:"type"`
	SubType string                    `bson:"sub_type,omitempty" json:"sub_type,omitempty"` // quick_reply, url
	Text    string                    `bson:"text,omitempty" json:"text,omitempty"`
	Format  string                    `bson:"format,omitempty" json:"format,omitempty"` // IMAGE, VIDEO, DOCUMENT
	Index   *int                      `bson:"index,omitempty" json:"index,omitempty"`
	Buttons []MessagingWhatsAppButton `bson:"buttons,omitempty" json:"buttons,omitempty"`
	Example *MessagingWhatsAppExample `bson:"example,omitempty" json:"example,omitempty"`
}

// MessagingWhatsAppButton represents a button in WhatsApp template
type MessagingWhatsAppButton struct {
	Type  string `bson:"type" json:"type"`
	Text  string `bson:"text" json:"text"`
	URL   string `bson:"url,omitempty" json:"url,omitempty"`
	Phone string `bson:"phone,omitempty" json:"phone,omitempty"`
	Index int    `bson:"index" json:"index"`
}

// MessagingWhatsAppExample represents example values for WhatsApp template
type MessagingWhatsAppExample struct {
	HeaderTextExamples [][]string `bson:"header_text_examples,omitempty" json:"header_text_examples,omitempty"`
	BodyTextExamples   [][]string `bson:"body_text_examples,omitempty" json:"body_text_examples,omitempty"`
	URLButtonExamples  [][]string `bson:"url_button_examples,omitempty" json:"url_button_examples,omitempty"`
}
