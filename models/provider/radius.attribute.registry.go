package provider

import (
	"fmt"
	"regexp"
	"strings"
)

// RadiusAttributeDataType defines the expected data type of a RADIUS attribute value.
type RadiusAttributeDataType string

const (
	RadiusDataTypeString  RadiusAttributeDataType = "string"
	RadiusDataTypeNumber  RadiusAttributeDataType = "number"
	RadiusDataTypeIP      RadiusAttributeDataType = "ip_address"
	RadiusDataTypeBoolean RadiusAttributeDataType = "boolean"
	RadiusDataTypeEnum    RadiusAttributeDataType = "enum"
)

// RadiusAttributeSchema defines a single RADIUS attribute and its validation rules.
type RadiusAttributeSchema struct {
	Key         string                  `json:"key"`
	Label       string                  `json:"label"`
	Description string                  `json:"description"`
	DataType    RadiusAttributeDataType `json:"dataType"`
	Placeholder string                  `json:"placeholder,omitempty"`
	Required    bool                    `json:"required,omitempty"`
	Options     []FieldOption           `json:"options,omitempty"` // Used if DataType is 'enum'
	Pattern     string                  `json:"pattern,omitempty"` // Regex pattern for validation
	MinValue    *float64                `json:"minValue,omitempty"`
	MaxValue    *float64                `json:"maxValue,omitempty"`
}

// VendorDictionaryInfo holds a collection of RADIUS attributes for a specific vendor dictionary.
type VendorDictionaryInfo struct {
	Value       string                  `json:"value"` // e.g., "MIKROTIK"
	Label       string                  `json:"label"` // e.g., "MikroTik"
	Description string                  `json:"description"`
	Attributes  []RadiusAttributeSchema `json:"attributes"`
}

// RadiusDictionaryRegistry is the single source of truth for all supported RADIUS dictionaries.
var RadiusDictionaryRegistry = map[string]VendorDictionaryInfo{
	"STANDARD": {
		Value:       "STANDARD",
		Label:       "Standard RADIUS",
		Description: "Standard RADIUS parameters defined by RFCs.",
		Attributes: []RadiusAttributeSchema{
			{
				Key:         "Framed-IP-Address",
				Label:       "Framed IP Address",
				Description: "Address to be configured for the user",
				DataType:    RadiusDataTypeIP,
				Required:    false,
			},
			{
				Key:         "Framed-Pool",
				Label:       "Framed Pool",
				Description: "IP pool from which to assign an IP address",
				DataType:    RadiusDataTypeString,
				Required:    false,
			},
		},
	},
	"MIKROTIK": {
		Value:       "MIKROTIK",
		Label:       "MikroTik",
		Description: "MikroTik RouterOS specific RADIUS attributes.",
		Attributes: []RadiusAttributeSchema{
			{
				Key:         "Mikrotik-Rate-Limit",
				Label:       "Rate Limit",
				Description: "Format: rx-rate[/tx-rate] [rx-burst-rate[/tx-burst-rate] [rx-burst-threshold[/tx-burst-threshold] [rx-burst-time[/tx-burst-time] [priority] [rx-rate-min[/tx-rate-min]]]]",
				DataType:    RadiusDataTypeString,
				Required:    false,
				Pattern:     `^\d+[kM](/\d+[kM])?.*`, // basic validation for speed formats like 10M/10M
			},
			{
				Key:         "Mikrotik-Address-List",
				Label:       "Address List",
				Description: "Firewall address list to which user is added",
				DataType:    RadiusDataTypeString,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Group",
				Label:       "User Group",
				Description: "MikroTik User group name",
				DataType:    RadiusDataTypeString,
				Required:    false,
			},
		},
	},
	"CISCO": {
		Value:       "CISCO",
		Label:       "Cisco",
		Description: "Cisco specific RADIUS attributes.",
		Attributes: []RadiusAttributeSchema{
			{
				Key:         "Cisco-AVPair",
				Label:       "Cisco AVPair",
				Description: "Arbitrary AV-Pairs used for advanced Cisco configs",
				DataType:    RadiusDataTypeString,
				Required:    false,
			},
		},
	},
}

// GetRadiusDictionaryConfig returns the full list of supported vendors and their attribute schemas for frontend rendering.
func GetRadiusDictionaryConfig() []VendorDictionaryInfo {
	dictionaries := make([]VendorDictionaryInfo, 0, len(RadiusDictionaryRegistry))
	for _, info := range RadiusDictionaryRegistry {
		dictionaries = append(dictionaries, info)
	}
	return dictionaries
}

// ValidateRadiusAttributes validates a map of key-value attributes against the known registry schemas.
func ValidateRadiusAttributes(attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}

	// Build a flat lookup map for quick validation of any supported key
	schemaMap := make(map[string]RadiusAttributeSchema)
	for _, vendor := range RadiusDictionaryRegistry {
		for _, attr := range vendor.Attributes {
			schemaMap[attr.Key] = attr
		}
	}

	for key, rawValue := range attributes {
		schema, exists := schemaMap[key]
		if !exists {
			// If not strictly dropping unknown keys, we might want to just skip
			// but for strong validation as a single source of truth, we should fail or just warning.
			// The user said: "our registory will be more advanced like it will key their expected values and validation rule if required". 
			// If it's explicitly given but not unknown, what happens? Let's ensure strict validation.
			return fmt.Errorf("RADIUS attribute '%s' is not defined in any supported vendor dictionary", key)
		}

		if err := validateAttributeValue(key, rawValue, schema); err != nil {
			return err
		}
	}

	return nil
}

func validateAttributeValue(key string, rawValue interface{}, schema RadiusAttributeSchema) error {
	if rawValue == nil {
		if schema.Required {
			return fmt.Errorf("RADIUS attribute '%s' is required but missing", key)
		}
		return nil
	}

	switch schema.DataType {
	case RadiusDataTypeString:
		val, ok := rawValue.(string)
		if !ok {
			return fmt.Errorf("RADIUS attribute '%s' must be a string", key)
		}
		if schema.Pattern != "" {
			regex, err := regexp.Compile(schema.Pattern)
			if err == nil && !regex.MatchString(val) {
				return fmt.Errorf("RADIUS attribute '%s' value '%s' does not match required pattern", key, val)
			}
		}

	case RadiusDataTypeNumber:
		// Accept float64 or int
		var val float64
		switch v := rawValue.(type) {
		case float64:
			val = v
		case int:
			val = float64(v)
		case int32:
			val = float64(v)
		case int64:
			val = float64(v)
		default:
			return fmt.Errorf("RADIUS attribute '%s' must be a number", key)
		}

		if schema.MinValue != nil && val < *schema.MinValue {
			return fmt.Errorf("RADIUS attribute '%s' must be at least %v", key, *schema.MinValue)
		}
		if schema.MaxValue != nil && val > *schema.MaxValue {
			return fmt.Errorf("RADIUS attribute '%s' must be at most %v", key, *schema.MaxValue)
		}

	case RadiusDataTypeIP:
		val, ok := rawValue.(string)
		if !ok {
			return fmt.Errorf("RADIUS attribute '%s' must be an IP address string", key)
		}
		// simple regex for basic IPv4 validation
		ipPattern := `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`
		matched, _ := regexp.MatchString(ipPattern, val)
		if !matched {
			return fmt.Errorf("RADIUS attribute '%s' must be a valid IPv4 address", key)
		}

	case RadiusDataTypeBoolean:
		_, ok := rawValue.(bool)
		if !ok {
			// Some clients might send string "true" / "false"
			valStr, okStr := rawValue.(string)
			if okStr && (strings.ToLower(valStr) == "true" || strings.ToLower(valStr) == "false") {
				return nil
			}
			return fmt.Errorf("RADIUS attribute '%s' must be a boolean", key)
		}

	case RadiusDataTypeEnum:
		valStr, ok := rawValue.(string)
		if !ok {
			return fmt.Errorf("RADIUS attribute '%s' must be a string (enum option)", key)
		}
		validOption := false
		for _, opt := range schema.Options {
			if opt.Value == valStr {
				validOption = true
				break
			}
		}
		if !validOption {
			return fmt.Errorf("RADIUS attribute '%s' value '%s' is not a valid option", key, valStr)
		}
	}

	return nil
}
