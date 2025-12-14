package subscriberevent

type FormType string

const (
	FormTypeSubscriber FormType = "subscriber" // User details form: name, mobile, email, address, KYC (ISP/Reseller/Partner/Distributor)
	FormTypeBroadband  FormType = "broadband"  // System/Account form: PPPoE, MAC, session details, account config (ISP/Reseller/Partner/Distributor)
	FormTypeHotspot    FormType = "hotspot"    // System/Account form: auth settings, device bindings, MAC, session details (Enterprise/Branch)
	FormTypeLogin      FormType = "login"      // User details form: name, mobile, email, address, KYC for login (Enterprise/Branch)
	FormTypeSignup     FormType = "signup"     // User details form: name, mobile, email, address, KYC for registration (Enterprise/Branch)
)

// AuthMethod represents an authentication method for login/signup forms
type AuthMethod string

const (
	AuthMethodOTP      AuthMethod = "OTP"
	AuthMethodPassword AuthMethod = "PASSWORD"
	AuthMethodSocial   AuthMethod = "SOCIAL"
	AuthMethodVoucher  AuthMethod = "VOUCHER"
	AuthMethodSSO      AuthMethod = "SSO"
	AuthMethodQRCode   AuthMethod = "QRCODE"
	AuthMethodOneTap   AuthMethod = "ONETAP"
)

// FormConfigCreatedEvent represents a form configuration creation event
// FormConfig includes fields, layout, sections, settings, and toggles
// This must match the structure in subscriber-service/internal/models/form.config.model.go
// Note: CreatedBy/UpdatedBy are service-level metadata and not included in events
type FormConfigCreatedEvent struct {
	ID       string            `json:"id" bson:"id"`
	TenantID string            `json:"tenantId" bson:"tenantId"`
	FormType FormType          `json:"formType" bson:"formType"` // subscriber, broadband, hotspot, login, signup
	Fields   []FieldConfigItem `json:"fields" bson:"fields"`     // Fields stored directly in FormConfig
	Sections []FormSection     `json:"sections,omitempty" bson:"sections,omitempty"`
	Layout   *FormLayout       `json:"layout,omitempty" bson:"layout,omitempty"`
	Settings *FormSettings     `json:"settings,omitempty" bson:"settings,omitempty"`
	Toggles  []FormToggle      `json:"toggles,omitempty" bson:"toggles,omitempty"`
	Version  int               `json:"version" bson:"version"`
}

// FormConfigUpdatedEvent represents a form configuration update event
// Note: UpdatedBy is service-level metadata and not included in events
type FormConfigUpdatedEvent struct {
	ID       string            `json:"id" bson:"id"`
	TenantID string            `json:"tenantId" bson:"tenantId"`
	FormType string            `json:"formType" bson:"formType"`
	Fields   []FieldConfigItem `json:"fields" bson:"fields"` // Fields stored directly in FormConfig
	Sections []FormSection     `json:"sections,omitempty" bson:"sections,omitempty"`
	Layout   *FormLayout       `json:"layout,omitempty" bson:"layout,omitempty"`
	Settings *FormSettings     `json:"settings,omitempty" bson:"settings,omitempty"`
	Toggles  []FormToggle      `json:"toggles,omitempty" bson:"toggles,omitempty"`
	Version  int               `json:"version" bson:"version"`
}

// FormConfigDeletedEvent represents a form configuration deletion event
type FormConfigDeletedEvent struct {
	ID       string `json:"id" bson:"id"`
	TenantID string `json:"tenantId" bson:"tenantId"`
	FormType string `json:"formType" bson:"formType"`
}

// FormSection represents a form section/group in events
// This must match the structure in subscriber-service/internal/models/form.config.model.go
type FormSection struct {
	ID            string   `json:"id" bson:"id"`                                           // Unique section ID
	Title         string   `json:"title" bson:"title"`                                     // Section title/heading
	Description   string   `json:"description,omitempty" bson:"description,omitempty"`     // Section description
	FieldKeys     []string `json:"fieldKeys" bson:"fieldKeys"`                             // Fields belonging to this section
	DisplayOrder  int      `json:"displayOrder" bson:"displayOrder"`                       // Order of section in form
	IsCollapsible bool     `json:"isCollapsible,omitempty" bson:"isCollapsible,omitempty"` // Can section be collapsed?
	IsCollapsed   bool     `json:"isCollapsed,omitempty" bson:"isCollapsed,omitempty"`     // Is section collapsed by default?
}

// FormLayout represents form layout settings in events
// This must match the structure in subscriber-service/internal/models/form.config.model.go
type FormLayout struct {
	Columns               int            `json:"columns,omitempty" bson:"columns,omitempty"`                             // Number of columns (1, 2, 3)
	ColumnGap             int            `json:"columnGap,omitempty" bson:"columnGap,omitempty"`                         // Gap between columns (px)
	RowGap                int            `json:"rowGap,omitempty" bson:"rowGap,omitempty"`                               // Gap between rows (px)
	LabelPosition         string         `json:"labelPosition,omitempty" bson:"labelPosition,omitempty"`                 // "top", "left", "right" (default: "top")
	LabelWidth            int            `json:"labelWidth,omitempty" bson:"labelWidth,omitempty"`                       // Label width for left/right positioning (px)
	ResponsiveBreakpoints map[string]int `json:"responsiveBreakpoints,omitempty" bson:"responsiveBreakpoints,omitempty"` // Breakpoints for responsive layout
}

// FormSettings represents form-level settings in events
// This must match the structure in subscriber-service/internal/models/form.config.model.go
type FormSettings struct {
	SubmitButtonText string `json:"submitButtonText,omitempty" bson:"submitButtonText,omitempty"` // Custom submit button text
	CancelButtonText string `json:"cancelButtonText,omitempty" bson:"cancelButtonText,omitempty"` // Custom cancel button text
	ShowCancelButton bool   `json:"showCancelButton,omitempty" bson:"showCancelButton,omitempty"` // Show cancel button?
	RedirectURL      string `json:"redirectURL,omitempty" bson:"redirectURL,omitempty"`           // Redirect after successful submission
	SuccessMessage   string `json:"successMessage,omitempty" bson:"successMessage,omitempty"`     // Custom success message
	ErrorMessage     string `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`         // Custom error message
	ShowProgressBar  bool   `json:"showProgressBar,omitempty" bson:"showProgressBar,omitempty"`   // Show progress indicator?
	AutoSave         bool   `json:"autoSave,omitempty" bson:"autoSave,omitempty"`                 // Auto-save form data?
	AutoSaveInterval int    `json:"autoSaveInterval,omitempty" bson:"autoSaveInterval,omitempty"` // Auto-save interval (seconds)
}

// FormToggle represents a form-level toggle that controls conditional field visibility in events
// This must match the structure in subscriber-service/internal/models/form.config.model.go
// Example: Toggle between "mobile" and "email" - when toggle is ON, show mobile field; when OFF, show email field
type FormToggle struct {
	ID           string         `json:"id" bson:"id"`                                         // Unique toggle ID (e.g., "contact-method-toggle")
	Label        string         `json:"label" bson:"label"`                                   // Toggle label (e.g., "Contact Method")
	Description  string         `json:"description,omitempty" bson:"description,omitempty"`   // Toggle description
	Type         string         `json:"type" bson:"type"`                                     // "radio", "toggle", "select" (default: "toggle")
	Options      []ToggleOption `json:"options" bson:"options"`                               // Options for radio/select, or ON/OFF states for toggle
	DefaultValue string         `json:"defaultValue,omitempty" bson:"defaultValue,omitempty"` // Default selected value
	DisplayOrder int            `json:"displayOrder" bson:"displayOrder"`                     // Order in form
	FieldGroups  []FieldGroup   `json:"fieldGroups" bson:"fieldGroups"`                       // Field groups controlled by this toggle
}

// ToggleOption represents an option for a toggle (for radio/select types) in events
type ToggleOption struct {
	Value       string `json:"value" bson:"value"`                                 // Option value (e.g., "mobile", "email")
	Label       string `json:"label" bson:"label"`                                 // Option label (e.g., "Mobile Number", "Email Address")
	Description string `json:"description,omitempty" bson:"description,omitempty"` // Option description
}

// FieldGroup represents a group of fields that are shown/hidden together based on toggle state in events
type FieldGroup struct {
	ID          string   `json:"id" bson:"id"`                   // Unique group ID
	ToggleValue string   `json:"toggleValue" bson:"toggleValue"` // Toggle value that activates this group (e.g., "mobile", "email", "on", "off")
	FieldKeys   []string `json:"fieldKeys" bson:"fieldKeys"`     // Fields in this group (e.g., ["primaryMobile", "alternateMobile"])
	IsVisible   bool     `json:"isVisible" bson:"isVisible"`     // Should this group be visible when toggle value matches?
}
