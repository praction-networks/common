package provider

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"strconv"
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
	RadiusDataTypeInteger RadiusAttributeDataType = "integer"
	RadiusDataTypeIPv6    RadiusAttributeDataType = "ipv6_address"
	RadiusDataTypeOctet   RadiusAttributeDataType = "octet_string"
)

// RadiusAttributeCategory groups attributes by functional area.
type RadiusAttributeCategory string

const (
	RadiusCategoryAuthentication RadiusAttributeCategory = "authentication"
	RadiusCategoryAuthorization  RadiusAttributeCategory = "authorization"
	RadiusCategoryAccounting     RadiusAttributeCategory = "accounting"
	RadiusCategoryQoS            RadiusAttributeCategory = "qos"
	RadiusCategoryFiltering      RadiusAttributeCategory = "filtering"
	RadiusCategoryRouting        RadiusAttributeCategory = "routing"
	RadiusCategoryTunneling      RadiusAttributeCategory = "tunneling"
	RadiusCategoryVendorSpecific RadiusAttributeCategory = "vendor_specific"
	RadiusCategoryGeneral        RadiusAttributeCategory = "general"
	RadiusCategorySession        RadiusAttributeCategory = "session"
	RadiusCategoryMobility       RadiusAttributeCategory = "mobility"
)



// RadiusAttributeSchema defines a single RADIUS attribute and its validation rules.
type RadiusAttributeSchema struct {
	Key            string                  `json:"key"`
	AttributeID    int                     `json:"attributeId,omitempty"` // Numeric RADIUS attribute ID
	Label          string                  `json:"label"`
	Description    string                  `json:"description"`
	DataType       RadiusAttributeDataType `json:"dataType"`
	Category       RadiusAttributeCategory `json:"category,omitempty"`
	Placeholder    string                  `json:"placeholder,omitempty"`
	Required       bool                    `json:"required,omitempty"`
	Options        []FieldOption           `json:"options,omitempty"` // Used if DataType is 'enum'
	Pattern        string                  `json:"pattern,omitempty"` // Regex pattern for validation
	MinValue       *float64                `json:"minValue,omitempty"`
	MaxValue       *float64                `json:"maxValue,omitempty"`
	MinLength      *int                    `json:"minLength,omitempty"`
	MaxLength      *int                    `json:"maxLength,omitempty"`
	Examples       []string                `json:"examples,omitempty"` // Example values
	Deprecated     bool                    `json:"deprecated,omitempty"`
	MultiValue     bool                    `json:"multiValue,omitempty"`     // Attribute may appear multiple times in one packet
	CoAApplicable  bool                    `json:"coaApplicable,omitempty"`  // Usable in Change-of-Authorization packets
	DMApplicable   bool                    `json:"dmApplicable,omitempty"`   // Usable in Disconnect-Message packets
}

// radiusKeyIndex maps every attribute Key → schema across all vendors (built by init).
// radiusVendorKeyIndex maps vendor → (attribute Key → schema) (built by init).
// radiusCategoryIndex maps category → []schema (built by init).
// radiusKeyVendors maps attribute Key → []vendor keys that define it (built by init).
var (
	radiusKeyIndex       map[string]RadiusAttributeSchema
	radiusVendorKeyIndex map[string]map[string]RadiusAttributeSchema
	radiusCategoryIndex  map[RadiusAttributeCategory][]RadiusAttributeSchema
	radiusKeyVendors     map[string][]string
)

// VendorDictionaryInfo holds a collection of RADIUS attributes for a specific vendor dictionary.
type VendorDictionaryInfo struct {
	Value       string                  `json:"value"` // e.g., "MIKROTIK"
	Label       string                  `json:"label"` // e.g., "MikroTik"
	Description string                  `json:"description"`
	VendorID    int                     `json:"vendorId,omitempty"` // Vendor-Specific attribute ID
	Attributes  []RadiusAttributeSchema `json:"attributes"`
}

// RadiusDictionaryRegistry is the single source of truth for all supported RADIUS dictionaries.
// Each vendor dictionary is defined in its own radius.vendor.<name>.go file.
var RadiusDictionaryRegistry = map[string]VendorDictionaryInfo{
	"STANDARD": radiusVendorStandard,
	"MIKROTIK": radiusVendorMikrotik,
	"CISCO": radiusVendorCisco,
	"JUNIPER": radiusVendorJuniper,
	"ARUBA": radiusVendorAruba,
	"HUAWEI": radiusVendorHuawei,
	"FORTINET": radiusVendorFortinet,
	"RUCKUS": radiusVendorRuckus,
	"UBIQUITI": radiusVendorUbiquiti,
	"CAMBIUM": radiusVendorCambium,
	"AEROHIVE": radiusVendorAerohive,
	"NOKIA_SR_OS": radiusVendorNokiaSrOs,
	"ERICSSON_SSR": radiusVendorEricssonSsr,
	"A10_NETWORKS": radiusVendorA10Networks,
	"NETELASTIC": radiusVendorNetelastic,
	"HP_PROCURVE": radiusVendorHpProcurve,
	"EXTREME_XOS": radiusVendorExtremeXos,
	"BROCADE_ICX": radiusVendorBrocadeIcx,
	"PALO_ALTO": radiusVendorPaloAlto,
	"ZTE": radiusVendorZte,
	"CALIX": radiusVendorCalix,
	"ADTRAN": radiusVendorAdtran,
	"CASA_SYSTEMS": radiusVendorCasaSystems,
	"STARENT": radiusVendorStarent,
}

// init builds the package-level lookup indexes from RadiusDictionaryRegistry.
// This runs once at startup so every subsequent lookup is O(1) map access.
// It also detects duplicate keys across vendors and panics early to prevent
// silent data-loss from map overwrites.
func init() {
	radiusKeyIndex = make(map[string]RadiusAttributeSchema, 512)
	radiusVendorKeyIndex = make(map[string]map[string]RadiusAttributeSchema, len(RadiusDictionaryRegistry))
	radiusCategoryIndex = make(map[RadiusAttributeCategory][]RadiusAttributeSchema, 16)
	radiusKeyVendors = make(map[string][]string, 512)

	// First pass: build per-vendor maps and collect key→[]vendor mapping.
	for vendor, info := range RadiusDictionaryRegistry {
		vendorMap := make(map[string]RadiusAttributeSchema, len(info.Attributes))
		for _, attr := range info.Attributes {
			vendorMap[attr.Key] = attr
			radiusKeyVendors[attr.Key] = append(radiusKeyVendors[attr.Key], vendor)
		}
		radiusVendorKeyIndex[vendor] = vendorMap
	}

	// Second pass: detect cross-vendor duplicates (same key, different vendor, both non-standard).
	// Standard attributes are intentionally shared, so skip STANDARD vs STANDARD conflicts.
	for key, vendors := range radiusKeyVendors {
		if len(vendors) > 1 {
			// Allow the same key to appear in STANDARD + one vendor (intentional re-use like NAS-IP-Address).
			// Panic only if two non-STANDARD vendors both define the same key, which is always a data error.
			nonStd := 0
			for _, v := range vendors {
				if v != "STANDARD" {
					nonStd++
				}
			}
			if nonStd > 1 {
				panic(fmt.Sprintf("radius registry: duplicate key %q defined in multiple non-standard vendors: %v", key, vendors))
			}
		}
	}

	// Third pass: build the flat key index and category index from vendor maps.
	for _, vendorMap := range radiusVendorKeyIndex {
		for _, attr := range vendorMap {
			// STANDARD wins over vendor when a key appears in both (vendor re-exports standard attrs).
			if _, exists := radiusKeyIndex[attr.Key]; !exists {
				radiusKeyIndex[attr.Key] = attr
			}
			if attr.Category != "" {
				radiusCategoryIndex[attr.Category] = append(radiusCategoryIndex[attr.Category], attr)
			}
		}
	}

	// Sort each category slice by Key for deterministic output.
	for cat := range radiusCategoryIndex {
		slice := radiusCategoryIndex[cat]
		sort.Slice(slice, func(i, j int) bool { return slice[i].Key < slice[j].Key })
		radiusCategoryIndex[cat] = slice
	}
}

// GetRadiusDictionaryConfig returns the full list of supported vendors and their attribute schemas,
// sorted by vendor label, for frontend rendering.
func GetRadiusDictionaryConfig() []VendorDictionaryInfo {
	dictionaries := make([]VendorDictionaryInfo, 0, len(RadiusDictionaryRegistry))
	for _, info := range RadiusDictionaryRegistry {
		dictionaries = append(dictionaries, info)
	}
	sort.Slice(dictionaries, func(i, j int) bool {
		return dictionaries[i].Label < dictionaries[j].Label
	})
	return dictionaries
}

// ValidateRadiusAttributes validates a map of key-value attributes against all known vendor schemas.
// The lookup uses the package-level index built by init(), so it is O(1) per key.
func ValidateRadiusAttributes(attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}
	for key, rawValue := range attributes {
		schema, exists := radiusKeyIndex[key]
		if !exists {
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
	case RadiusDataTypeString, RadiusDataTypeIP, RadiusDataTypeIPv6:
		val, ok := rawValue.(string)
		if !ok {
			return fmt.Errorf("RADIUS attribute '%s' must be a string", key)
		}

		if schema.MinLength != nil && len(val) < *schema.MinLength {
			return fmt.Errorf("RADIUS attribute '%s' must be at least %d characters", key, *schema.MinLength)
		}
		if schema.MaxLength != nil && len(val) > *schema.MaxLength {
			return fmt.Errorf("RADIUS attribute '%s' must not exceed %d characters", key, *schema.MaxLength)
		}

		if schema.Pattern != "" {
			regex, err := regexp.Compile(schema.Pattern)
			if err == nil && !regex.MatchString(val) {
				return fmt.Errorf("RADIUS attribute '%s' value '%s' does not match required pattern", key, val)
			}
		}

		switch schema.DataType {
		case RadiusDataTypeIP:
			if net.ParseIP(val) == nil || strings.Contains(val, ":") {
				return fmt.Errorf("RADIUS attribute '%s' must be a valid IPv4 address", key)
			}
		case RadiusDataTypeIPv6:
			if net.ParseIP(val) == nil || !strings.Contains(val, ":") {
				return fmt.Errorf("RADIUS attribute '%s' must be a valid IPv6 address", key)
			}
		}

	case RadiusDataTypeInteger, RadiusDataTypeNumber:
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
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("RADIUS attribute '%s' must be a number", key)
			}
			val = parsed
		default:
			return fmt.Errorf("RADIUS attribute '%s' must be a number", key)
		}

		if schema.DataType == RadiusDataTypeInteger && val != float64(int64(val)) {
			return fmt.Errorf("RADIUS attribute '%s' must be an integer", key)
		}

		if schema.MinValue != nil && val < *schema.MinValue {
			return fmt.Errorf("RADIUS attribute '%s' must be at least %v", key, *schema.MinValue)
		}
		if schema.MaxValue != nil && val > *schema.MaxValue {
			return fmt.Errorf("RADIUS attribute '%s' must be at most %v", key, *schema.MaxValue)
		}

	case RadiusDataTypeBoolean:
		_, ok := rawValue.(bool)
		if !ok {
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

// ---------------------------------------------------------------------------
// Key validation helpers
// ---------------------------------------------------------------------------

// IsValidKey reports whether key is defined in any vendor's dictionary.
func IsValidKey(key string) bool {
	_, ok := radiusKeyIndex[key]
	return ok
}

// IsValidKeyForVendor reports whether key is defined in the given vendor's dictionary.
// vendor is compared case-insensitively (e.g. "mikrotik" == "MIKROTIK").
func IsValidKeyForVendor(vendor, key string) bool {
	vendorMap, ok := radiusVendorKeyIndex[strings.ToUpper(vendor)]
	if !ok {
		return false
	}
	_, ok = vendorMap[key]
	return ok
}

// GetAttributeSchema returns the schema for key from any vendor's dictionary.
func GetAttributeSchema(key string) (RadiusAttributeSchema, bool) {
	schema, ok := radiusKeyIndex[key]
	return schema, ok
}

// GetAttributeSchemaForVendor returns the schema for key restricted to a specific vendor.
func GetAttributeSchemaForVendor(vendor, key string) (RadiusAttributeSchema, bool) {
	vendorMap, ok := radiusVendorKeyIndex[strings.ToUpper(vendor)]
	if !ok {
		return RadiusAttributeSchema{}, false
	}
	schema, ok := vendorMap[key]
	return schema, ok
}

// ---------------------------------------------------------------------------
// Vendor-scoped validation
// ---------------------------------------------------------------------------

// ValidateRadiusAttributesForVendor validates attributes restricted to a specific vendor's dictionary.
// Use this when the caller already knows which vendor dictionary should be authoritative.
func ValidateRadiusAttributesForVendor(vendor string, attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}
	vendorUpper := strings.ToUpper(vendor)
	vendorMap, ok := radiusVendorKeyIndex[vendorUpper]
	if !ok {
		return fmt.Errorf("vendor '%s' is not in the RADIUS dictionary registry; valid vendors: %s",
			vendor, strings.Join(GetValidVendors(), ", "))
	}
	for key, rawValue := range attributes {
		schema, exists := vendorMap[key]
		if !exists {
			return fmt.Errorf("RADIUS attribute '%s' is not defined for vendor '%s'", key, vendor)
		}
		if err := validateAttributeValue(key, rawValue, schema); err != nil {
			return err
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Registry discovery helpers
// ---------------------------------------------------------------------------

// GetValidVendors returns a sorted list of all vendor keys in the registry.
func GetValidVendors() []string {
	vendors := make([]string, 0, len(RadiusDictionaryRegistry))
	for k := range RadiusDictionaryRegistry {
		vendors = append(vendors, k)
	}
	sort.Strings(vendors)
	return vendors
}

// GetValidKeysForVendor returns a sorted list of all attribute keys for the given vendor.
// Returns an error if the vendor is unknown.
func GetValidKeysForVendor(vendor string) ([]string, error) {
	vendorMap, ok := radiusVendorKeyIndex[strings.ToUpper(vendor)]
	if !ok {
		return nil, fmt.Errorf("vendor '%s' is not in the RADIUS dictionary registry", vendor)
	}
	keys := make([]string, 0, len(vendorMap))
	for k := range vendorMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

// GetAttributesByCategory returns all attribute schemas across all vendors that belong to
// category, sorted by Key. Uses the pre-built category index — O(1) lookup.
func GetAttributesByCategory(category RadiusAttributeCategory) []RadiusAttributeSchema {
	return radiusCategoryIndex[category] // already sorted by init()
}

// GetVendorsForKey returns all vendor keys that define the given attribute key, sorted.
// Most attributes return a single vendor; standard attrs shared by a vendor return two.
func GetVendorsForKey(attrKey string) []string {
	vendors := radiusKeyVendors[attrKey]
	if len(vendors) == 0 {
		return nil
	}
	out := make([]string, len(vendors))
	copy(out, vendors)
	sort.Strings(out)
	return out
}

// GetCoAAttributes returns all attributes marked as applicable in Change-of-Authorization packets,
// sorted by Key.
func GetCoAAttributes() []RadiusAttributeSchema {
	var result []RadiusAttributeSchema
	for _, attr := range radiusKeyIndex {
		if attr.CoAApplicable {
			result = append(result, attr)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result
}

// GetDMAttributes returns all attributes marked as applicable in Disconnect-Message packets,
// sorted by Key.
func GetDMAttributes() []RadiusAttributeSchema {
	var result []RadiusAttributeSchema
	for _, attr := range radiusKeyIndex {
		if attr.DMApplicable {
			result = append(result, attr)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result
}

// ---------------------------------------------------------------------------
// Multi-error validation
// ---------------------------------------------------------------------------

// ValidateRadiusAttributesAll validates every attribute and collects all errors instead of
// stopping at the first one. Returns nil if all attributes are valid.
func ValidateRadiusAttributesAll(attributes map[string]interface{}) []error {
	if attributes == nil {
		return nil
	}
	var errs []error
	for key, rawValue := range attributes {
		schema, exists := radiusKeyIndex[key]
		if !exists {
			errs = append(errs, fmt.Errorf("RADIUS attribute '%s' is not defined in any supported vendor dictionary", key))
			continue
		}
		if err := validateAttributeValue(key, rawValue, schema); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ValidateRadiusAttributesForVendorAll is the multi-error variant of ValidateRadiusAttributesForVendor.
func ValidateRadiusAttributesForVendorAll(vendor string, attributes map[string]interface{}) []error {
	if attributes == nil {
		return nil
	}
	vendorUpper := strings.ToUpper(vendor)
	vendorMap, ok := radiusVendorKeyIndex[vendorUpper]
	if !ok {
		return []error{fmt.Errorf("vendor '%s' is not in the RADIUS dictionary registry; valid vendors: %s",
			vendor, strings.Join(GetValidVendors(), ", "))}
	}
	var errs []error
	for key, rawValue := range attributes {
		schema, exists := vendorMap[key]
		if !exists {
			errs = append(errs, fmt.Errorf("RADIUS attribute '%s' is not defined for vendor '%s'", key, vendor))
			continue
		}
		if err := validateAttributeValue(key, rawValue, schema); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

// SearchAttributes performs a case-insensitive substring search across Key, Label, and
// Description of all attributes in all vendors. Results are sorted by Key.
// Pass vendor="" to search all vendors, or a specific vendor key to restrict the search.
func SearchAttributes(query, vendor string) []RadiusAttributeSchema {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}

	var source map[string]RadiusAttributeSchema
	if vendor == "" {
		source = radiusKeyIndex
	} else {
		vm, ok := radiusVendorKeyIndex[strings.ToUpper(vendor)]
		if !ok {
			return nil
		}
		source = vm
	}

	var result []RadiusAttributeSchema
	for _, attr := range source {
		if strings.Contains(strings.ToLower(attr.Key), q) ||
			strings.Contains(strings.ToLower(attr.Label), q) ||
			strings.Contains(strings.ToLower(attr.Description), q) {
			result = append(result, attr)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result
}

// ---------------------------------------------------------------------------
// Internal helpers for schema definitions
// ---------------------------------------------------------------------------

func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int             { return &v }
