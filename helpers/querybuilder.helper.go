// common/helpers/querybuilder.go
package helpers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/praction-networks/common/appError"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* ==============================
   PUBLIC CONFIG & MODELS
   ============================== */

// Pagination-only
type PaginationConfig struct {
	MaxLimit     int
	DefaultLimit int
}

// Sorting
type SortingConfig struct {
	AllowedFields []string
}

// Search (regex-OR across fields)
type SearchConfig struct {
	Fields []string
}

// Field typing
type Kind int

const (
	KindString Kind = iota
	KindInt
	KindFloat
	KindBool
	KindTime
	KindObjectID
)

type Decoder func(raw string) (any, error)

type FieldSpec struct {
	Kind            Kind
	Layouts         []string // time parsing layouts (fallback order)
	Decoder         Decoder  // optional custom coercion (overrides Kind)
	Alias           string   // API field -> DB field
	CaseInsensitive bool     // default regex $options:"i"
	AllowedOps      []string // restrict operators per field; empty => allow defaults
	AllowRawRegex   bool     // allow `regex` op (vs only contains|startsWith|endsWith)
}

// Filtering config
type FilteringConfig struct {
	AllowedFields []string // whitelist of filterable fields (API names)
	DateField     string   // API field used for date range (date_from/date_to)
	FieldSpecs    map[string]FieldSpec
}

// Projection (include/exclude/distinct whitelists)
type ProjectionConfig struct {
	AllowInclude  bool
	AllowExclude  bool
	AllowedFields []string // whitelist for include/exclude/distinct; empty => no extra check
}

// Security/caps
type Limits struct {
	MaxFilters       int // e.g., 25
	MaxSorts         int // e.g., 5
	MaxIncludeFields int // e.g., 50
	MaxExcludeFields int // e.g., 50
	MaxValueLength   int // e.g., 512
	MaxSearchLength  int // e.g., 256
	MaxRegexLength   int // e.g., 256
	FindMaxTimeMS    int // e.g., 3000
	MaxQueryDepth    int // e.g., 3 - Prevent deep nesting
	MaxOrConditions  int // e.g., 10 - Limit $or complexity
	MaxQuerySize     int // e.g., 8192 - Total query string size limit
}

// Master policy passed from feature (e.g., Tenants)
type QueryPolicy struct {
	Pagination  PaginationConfig
	Sorting     SortingConfig
	Filtering   FilteringConfig
	Search      SearchConfig
	Projection  ProjectionConfig
	TimeLayouts []string // global fallback for time parsing
	Limits      Limits
}

// Sensible defaults (optional helper)
func DefaultPolicy() QueryPolicy {
	return QueryPolicy{
		Pagination:  PaginationConfig{MaxLimit: 1000, DefaultLimit: 50},
		Sorting:     SortingConfig{AllowedFields: []string{"createdAt", "updatedAt"}},
		Filtering:   FilteringConfig{DateField: "createdAt", FieldSpecs: map[string]FieldSpec{}},
		Search:      SearchConfig{Fields: []string{"name", "description"}},
		Projection:  ProjectionConfig{AllowInclude: true, AllowExclude: true},
		TimeLayouts: []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"},
		Limits: Limits{
			MaxFilters:       25,
			MaxSorts:         5,
			MaxIncludeFields: 50,
			MaxExcludeFields: 50,
			MaxValueLength:   512,
			MaxSearchLength:  256,
			MaxRegexLength:   256,
			FindMaxTimeMS:    3000,
			MaxQueryDepth:    3,
			MaxOrConditions:  10,
			MaxQuerySize:     8192,
		},
	}
}

// Input model parsed from query params
type SortField struct {
	Field string
	Order string // asc|desc
}

type PaginatedFeedQuery struct {
	Limit          int
	Offset         int
	Sort           []SortField
	Filters        map[string]map[string]string
	IncludeFields  []string
	ExcludeFields  []string
	PaginationMeta bool
	DistinctField  string
	Search         string
	DateFrom       *time.Time
	DateTo         *time.Time
}

// Output: ready-to-use Mongo query
type MongoQuery struct {
	Filter        bson.M
	FindOptions   *options.FindOptions
	DistinctField string
	WithMeta      bool
	Limit         int
	Offset        int
}

/* ==============================
   PUBLIC ENTRYPOINT
   ============================== */

func BuildFromRequest(r *http.Request, policy QueryPolicy) (MongoQuery, error) {
	var fq PaginatedFeedQuery
	if err := parseInto(&fq, r, policy); err != nil {
		return MongoQuery{}, err
	}
	return buildMongoQuery(fq, policy)
}

/* ==============================
   PARSING
   ============================== */

func parseInto(fq *PaginatedFeedQuery, r *http.Request, p QueryPolicy) error {
	qs := r.URL.Query()

	// Validate query string size
	if p.Limits.MaxQuerySize > 0 && len(r.URL.RawQuery) > p.Limits.MaxQuerySize {
		return badReqWithField("query", fmt.Sprintf("query string too large, maximum allowed: %d bytes", p.Limits.MaxQuerySize))
	}

	// --- Pagination (PaginationConfig only) ---
	if v := qs.Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil || l < 1 {
			return badReqWithField("limit", "invalid limit value, must be a positive integer")
		}
		if p.Pagination.MaxLimit > 0 && l > p.Pagination.MaxLimit {
			return badReqWithField("limit", fmt.Sprintf("limit %d exceeds maximum allowed %d", l, p.Pagination.MaxLimit))
		}
		fq.Limit = l
	} else {
		if p.Pagination.DefaultLimit > 0 {
			fq.Limit = p.Pagination.DefaultLimit
		} else {
			fq.Limit = 50
		}
	}
	if v := qs.Get("offset"); v != "" {
		o, err := strconv.Atoi(v)
		if err != nil || o < 0 {
			return badReqWithField("offset", "invalid offset value, must be a non-negative integer")
		}
		fq.Offset = o
	}

	// --- Sorting ---
	if v := qs.Get("sort"); v != "" {
		for _, item := range splitCSV(v) {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) != 2 {
				return badReqWithField("sort", "invalid sort format, use field:asc|desc")
			}
			field := strings.TrimSpace(parts[0])
			order := strings.TrimSpace(parts[1])

			if !validAPIFieldName(field) {
				return badReqWithField("sort", "invalid sort field name: "+field)
			}
			if !inSlice(field, p.Sorting.AllowedFields) {
				return badReqWithField("sort", "sort field not allowed: "+field)
			}
			if order != "asc" && order != "desc" {
				return badReqWithField("sort", "invalid sort order: "+order+", must be 'asc' or 'desc'")
			}
			fq.Sort = append(fq.Sort, SortField{Field: field, Order: order})
		}
	}
	if p.Limits.MaxSorts > 0 && len(fq.Sort) > p.Limits.MaxSorts {
		return badReqWithField("sort", fmt.Sprintf("too many sort fields, maximum allowed: %d", p.Limits.MaxSorts))
	}

	// --- Filters (any unknown key becomes a filter) ---
	fq.Filters = map[string]map[string]string{}
	known := map[string]bool{
		"limit": true, "offset": true, "sort": true, "include": true, "exclude": true,
		"pagination_meta": true, "distinct": true, "search": true, "date_from": true, "date_to": true,
	}
	for key, vals := range qs {
		if known[key] || len(vals) == 0 {
			continue
		}
		field, op := parseFilterKey(key)
		if !validAPIFieldName(field) {
			return badReqWithField("filter", "invalid filter field name: "+field)
		}
		if !inSlice(field, p.Filtering.AllowedFields) {
			return badReqWithField("filter", "filter field not allowed: "+field)
		}
		val := vals[0]
		if p.Limits.MaxValueLength > 0 && len(val) > p.Limits.MaxValueLength {
			val = enforceLenCap(val, p.Limits.MaxValueLength)
		}
		// Sanitize input to prevent NoSQL injection
		val = sanitizeInput(val)
		if fq.Filters[field] == nil {
			fq.Filters[field] = map[string]string{}
		}
		fq.Filters[field][op] = val
	}
	if p.Limits.MaxFilters > 0 && len(fq.Filters) > p.Limits.MaxFilters {
		return badReqWithField("filter", fmt.Sprintf("too many filters, maximum allowed: %d", p.Limits.MaxFilters))
	}

	// --- Projection ---
	if v := qs.Get("include"); v != "" {
		if !p.Projection.AllowInclude {
			return badReqWithField("include", "include projection not allowed")
		}
		fq.IncludeFields = splitCSV(v)
	}
	if v := qs.Get("exclude"); v != "" {
		if !p.Projection.AllowExclude {
			return badReqWithField("exclude", "exclude projection not allowed")
		}
		fq.ExcludeFields = splitCSV(v)
	}
	if len(fq.IncludeFields) > 0 && len(fq.ExcludeFields) > 0 {
		return badReqWithField("projection", "cannot combine 'include' and 'exclude' projections")
	}
	if p.Limits.MaxIncludeFields > 0 && len(fq.IncludeFields) > p.Limits.MaxIncludeFields {
		return badReqWithField("include", fmt.Sprintf("too many include fields, maximum allowed: %d", p.Limits.MaxIncludeFields))
	}
	if p.Limits.MaxExcludeFields > 0 && len(fq.ExcludeFields) > p.Limits.MaxExcludeFields {
		return badReqWithField("exclude", fmt.Sprintf("too many exclude fields, maximum allowed: %d", p.Limits.MaxExcludeFields))
	}
	// Validate projection field names & whitelists
	for _, f := range fq.IncludeFields {
		if !validAPIFieldName(f) {
			return badReqWithField("include", "invalid include field name: "+f)
		}
		if !mustBeAllowed(f, p.Projection.AllowedFields) {
			return badReqWithField("include", "include field not allowed: "+f)
		}
	}
	for _, f := range fq.ExcludeFields {
		if !validAPIFieldName(f) {
			return badReqWithField("exclude", "invalid exclude field name: "+f)
		}
		if !mustBeAllowed(f, p.Projection.AllowedFields) {
			return badReqWithField("exclude", "exclude field not allowed: "+f)
		}
	}

	// --- Misc flags ---
	fq.PaginationMeta = qs.Get("pagination_meta") == "true"
	if v := qs.Get("distinct"); v != "" {
		if !validAPIFieldName(v) {
			return badReqWithField("distinct", "invalid distinct field name: "+v)
		}
		if !mustBeAllowed(v, p.Projection.AllowedFields) {
			return badReqWithField("distinct", "distinct field not allowed: "+v)
		}
		fq.DistinctField = strings.TrimSpace(v)
	}

	// --- Search ---
	if v := qs.Get("search"); v != "" {
		s := strings.TrimSpace(v)
		if p.Limits.MaxSearchLength > 0 && len(s) > p.Limits.MaxSearchLength {
			s = enforceLenCap(s, p.Limits.MaxSearchLength)
		}
		// Sanitize search input to prevent NoSQL injection
		s = sanitizeInput(s)
		fq.Search = s
	}

	// --- Date range ---
	if v := qs.Get("date_from"); v != "" {
		t, err := parseTimeMulti(v, firstNonEmpty(p.TimeLayouts))
		if err != nil {
			return badReqWithField("date_from", "invalid date_from format, expected: "+strings.Join(firstNonEmpty(p.TimeLayouts), " or "))
		}
		fq.DateFrom = &t
	}
	if v := qs.Get("date_to"); v != "" {
		t, err := parseTimeMulti(v, firstNonEmpty(p.TimeLayouts))
		if err != nil {
			return badReqWithField("date_to", "invalid date_to format, expected: "+strings.Join(firstNonEmpty(p.TimeLayouts), " or "))
		}
		fq.DateTo = &t
	}

	// Validate query complexity
	if err := validateQueryComplexity(*fq, p); err != nil {
		return err
	}

	return nil
}

/* ==============================
   BUILD MONGO QUERY
   ============================== */

func buildMongoQuery(fq PaginatedFeedQuery, p QueryPolicy) (MongoQuery, error) {
	filter := buildFilter(fq, p)

	opts := options.Find().
		SetLimit(int64(max(1, min(fq.Limit, p.Pagination.MaxLimit)))).
		SetSkip(int64(max(0, fq.Offset)))

	// Sort (already validated)
	if len(fq.Sort) > 0 {
		var d bson.D
		for _, s := range fq.Sort {
			dir := 1
			if s.Order == "desc" {
				dir = -1
			}
			d = append(d, bson.E{Key: dbField(s.Field, p), Value: dir})
		}
		opts.SetSort(d)
	}

	// Projection
	if len(fq.IncludeFields) > 0 {
		proj := bson.M{}
		for _, f := range fq.IncludeFields {
			proj[dbField(f, p)] = 1
		}
		opts.SetProjection(proj)
	} else if len(fq.ExcludeFields) > 0 {
		proj := bson.M{}
		for _, f := range fq.ExcludeFields {
			proj[dbField(f, p)] = 0
		}
		opts.SetProjection(proj)
	}

	// MaxTime to bound query
	if p.Limits.FindMaxTimeMS > 0 {
		opts.SetMaxTime(time.Duration(p.Limits.FindMaxTimeMS) * time.Millisecond)
	}

	return MongoQuery{
		Filter:        filter,
		FindOptions:   opts,
		DistinctField: fq.DistinctField,
		WithMeta:      fq.PaginationMeta,
		Limit:         fq.Limit,
		Offset:        fq.Offset,
	}, nil
}

func buildFilter(fq PaginatedFeedQuery, p QueryPolicy) bson.M {
	q := bson.M{}

	// 1) Regular filters with type coercion & operator restrictions
	for apiField, ops := range fq.Filters {
		spec := p.Filtering.FieldSpecs[apiField]
		dbf := dbField(apiField, p)

		if len(ops) == 1 && ops["eq"] != "" {
			q[dbf] = coerceValue(ops["eq"], spec, p)
			continue
		}
		sub := bson.M{}
		for op, raw := range ops {
			if !isOpAllowed(op, spec) {
				continue
			}
			switch op {
			case "eq":
				sub["$eq"] = coerceValue(raw, spec, p)
			case "ne":
				sub["$ne"] = coerceValue(raw, spec, p)
			case "gt":
				sub["$gt"] = coerceValue(raw, spec, p)
			case "gte":
				sub["$gte"] = coerceValue(raw, spec, p)
			case "lt":
				sub["$lt"] = coerceValue(raw, spec, p)
			case "lte":
				sub["$lte"] = coerceValue(raw, spec, p)
			case "in":
				sub["$in"] = coerceSlice(splitCSV(raw), spec, p)
			case "nin":
				sub["$nin"] = coerceSlice(splitCSV(raw), spec, p)
			case "exists":
				sub["$exists"] = parseBool(raw)
			case "regex":
				if !spec.AllowRawRegex {
					continue // raw regex disabled unless explicitly enabled
				}
				pat := raw
				if p.Limits.MaxRegexLength > 0 && len(pat) > p.Limits.MaxRegexLength {
					pat = enforceLenCap(pat, p.Limits.MaxRegexLength)
				}
				// Validate regex pattern for ReDoS protection
				if err := validateRegexPattern(pat); err != nil {
					continue // Skip invalid regex patterns
				}
				sub["$regex"] = pat
				if spec.CaseInsensitive {
					sub["$options"] = "i"
				} else {
					sub["$options"] = "i" // default safe
				}
			case "contains":
				pat := regexp.QuoteMeta(raw)
				if p.Limits.MaxRegexLength > 0 && len(pat) > p.Limits.MaxRegexLength {
					pat = enforceLenCap(pat, p.Limits.MaxRegexLength)
				}
				sub["$regex"] = pat
				sub["$options"] = "i"
			case "startsWith":
				pat := "^" + regexp.QuoteMeta(raw)
				if p.Limits.MaxRegexLength > 0 && len(pat) > p.Limits.MaxRegexLength {
					pat = enforceLenCap(pat, p.Limits.MaxRegexLength)
				}
				sub["$regex"] = pat
				sub["$options"] = "i"
			case "endsWith":
				pat := regexp.QuoteMeta(raw) + "$"
				if p.Limits.MaxRegexLength > 0 && len(pat) > p.Limits.MaxRegexLength {
					pat = enforceLenCap(pat, p.Limits.MaxRegexLength)
				}
				sub["$regex"] = pat
				sub["$options"] = "i"
			}
		}
		if len(sub) > 0 {
			q[dbf] = sub
		}
	}

	// 2) Search OR across fields (safe-quoted)
	if fq.Search != "" && len(p.Search.Fields) > 0 {
		terms := strings.TrimSpace(fq.Search)
		var or []bson.M
		for _, f := range p.Search.Fields {
			or = append(or, bson.M{
				dbField(f, p): bson.M{"$regex": regexp.QuoteMeta(terms), "$options": "i"},
			})
		}
		// Limit $or conditions for performance
		if p.Limits.MaxOrConditions > 0 && len(or) > p.Limits.MaxOrConditions {
			or = or[:p.Limits.MaxOrConditions]
		}
		if len(or) > 0 {
			q["$or"] = or
		}
	}

	// 3) Date range on Filtering.DateField (if present)
	if p.Filtering.DateField != "" && (fq.DateFrom != nil || fq.DateTo != nil) {
		df := bson.M{}
		if fq.DateFrom != nil {
			df["$gte"] = *fq.DateFrom
		}
		if fq.DateTo != nil {
			df["$lte"] = *fq.DateTo
		}
		if len(df) > 0 {
			q[dbField(p.Filtering.DateField, p)] = df
		}
	}

	return q
}

/* ==============================
   HELPERS
   ============================== */

func badReqWithField(field, message string) error {
	return appError.New(appError.InvalidInputError, field+": "+message, http.StatusBadRequest, nil)
}

// validateQueryComplexity performs comprehensive query validation
func validateQueryComplexity(fq PaginatedFeedQuery, p QueryPolicy) error {
	// Validate filter count
	if p.Limits.MaxFilters > 0 && len(fq.Filters) > p.Limits.MaxFilters {
		return badReqWithField("filters", fmt.Sprintf("too many filters, maximum allowed: %d", p.Limits.MaxFilters))
	}

	// Validate sort count
	if p.Limits.MaxSorts > 0 && len(fq.Sort) > p.Limits.MaxSorts {
		return badReqWithField("sort", fmt.Sprintf("too many sort fields, maximum allowed: %d", p.Limits.MaxSorts))
	}

	// Validate include fields count
	if p.Limits.MaxIncludeFields > 0 && len(fq.IncludeFields) > p.Limits.MaxIncludeFields {
		return badReqWithField("include", fmt.Sprintf("too many include fields, maximum allowed: %d", p.Limits.MaxIncludeFields))
	}

	// Validate exclude fields count
	if p.Limits.MaxExcludeFields > 0 && len(fq.ExcludeFields) > p.Limits.MaxExcludeFields {
		return badReqWithField("exclude", fmt.Sprintf("too many exclude fields, maximum allowed: %d", p.Limits.MaxExcludeFields))
	}

	// Validate search length
	if p.Limits.MaxSearchLength > 0 && len(fq.Search) > p.Limits.MaxSearchLength {
		return badReqWithField("search", fmt.Sprintf("search term too long, maximum allowed: %d characters", p.Limits.MaxSearchLength))
	}

	// Check for injection attempts in all inputs
	if err := detectInjectionAttempts(fq); err != nil {
		return err
	}

	return nil
}

// detectInjectionAttempts scans all query inputs for potential injection attacks
func detectInjectionAttempts(fq PaginatedFeedQuery) error {
	// Check search term for injection patterns
	if fq.Search != "" && containsInjectionPattern(fq.Search) {
		return badReqWithField("search", "potentially malicious input detected")
	}

	// Check filter values for injection patterns
	for field, ops := range fq.Filters {
		for _, value := range ops {
			if containsInjectionPattern(value) {
				return badReqWithField("filter", fmt.Sprintf("potentially malicious input detected in field '%s'", field))
			}
		}
	}

	// Check include/exclude fields for injection patterns
	for _, field := range fq.IncludeFields {
		if containsInjectionPattern(field) {
			return badReqWithField("include", "potentially malicious field name detected")
		}
	}

	for _, field := range fq.ExcludeFields {
		if containsInjectionPattern(field) {
			return badReqWithField("exclude", "potentially malicious field name detected")
		}
	}

	// Check distinct field
	if fq.DistinctField != "" && containsInjectionPattern(fq.DistinctField) {
		return badReqWithField("distinct", "potentially malicious field name detected")
	}

	return nil
}

// containsInjectionPattern checks if input contains common injection patterns
func containsInjectionPattern(input string) bool {
	// Common NoSQL injection patterns
	injectionPatterns := []string{
		"$where", "$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin", "$exists", "$regex",
		"$or", "$and", "$not", "$nor", "$all", "$elemMatch", "$size", "$type",
		"javascript:", "this.", "function", "eval", "exec", "script",
		"sleep(", "waitfor", "delay", "benchmark", "load_file", "into outfile",
		"union", "select", "insert", "update", "delete", "drop", "create", "alter",
		"<script", "</script", "onload=", "onerror=", "onclick=", "onmouseover=",
		"document.cookie", "window.location", "alert(", "confirm(", "prompt(",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range injectionPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

// sanitizeInput removes potentially dangerous characters and patterns to prevent NoSQL injection
func sanitizeInput(input string) string {
	// Remove or escape dangerous characters
	dangerousChars := map[string]string{
		"$":  "", // Remove MongoDB operators
		".":  "", // Remove dot notation
		"(":  "", // Remove parentheses
		")":  "",
		"[":  "", // Remove brackets
		"]":  "",
		"{":  "", // Remove braces
		"}":  "",
		"\\": "", // Remove backslashes
		"/":  "", // Remove forward slashes
		":":  "", // Remove colons
		";":  "", // Remove semicolons
		"=":  "", // Remove equals
		"<":  "", // Remove comparison operators
		">":  "",
		"!":  "", // Remove exclamation
		"@":  "", // Remove at symbol
		"#":  "", // Remove hash
		"%":  "", // Remove percent
		"^":  "", // Remove caret
		"&":  "", // Remove ampersand
		"*":  "", // Remove asterisk
		"|":  "", // Remove pipe
		"~":  "", // Remove tilde
		"`":  "", // Remove backtick
		"'":  "", // Remove single quotes
		"\"": "", // Remove double quotes
		"?":  "", // Remove question mark
		"+":  "", // Remove plus
		"-":  "", // Remove minus
	}

	result := input
	for char, replacement := range dangerousChars {
		result = strings.ReplaceAll(result, char, replacement)
	}

	// Remove common NoSQL injection patterns
	injectionPatterns := []string{
		"$where", "$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin", "$exists", "$regex",
		"$or", "$and", "$not", "$nor", "$all", "$elemMatch", "$size", "$type",
		"javascript:", "this.", "function", "eval", "exec", "script",
		"sleep(", "waitfor", "delay", "benchmark", "load_file", "into outfile",
		"union", "select", "insert", "update", "delete", "drop", "create", "alter",
	}

	lowerResult := strings.ToLower(result)
	for _, pattern := range injectionPatterns {
		if strings.Contains(lowerResult, pattern) {
			result = strings.ReplaceAll(result, pattern, "")
		}
	}

	return strings.TrimSpace(result)
}

// validateRegexPattern checks for potentially dangerous regex patterns that could cause ReDoS
func validateRegexPattern(pattern string) error {
	// Check for common ReDoS patterns
	dangerousPatterns := []string{
		"(a+)+", "(a*)*", "(a|a)+", "(a|a)*",
		"(a+)*", "(a*)+", "(.+)+", "(.*)*",
		"(.+)*", "(.*)+",
	}

	for _, dangerous := range dangerousPatterns {
		if strings.Contains(pattern, dangerous) {
			return fmt.Errorf("potentially dangerous regex pattern detected")
		}
	}

	// Test if the regex compiles without issues
	_, err := regexp.Compile(pattern)
	return err
}

func parseFilterKey(key string) (field, operator string) {
	operator = "eq"
	if i := strings.Index(key, "["); i != -1 && strings.HasSuffix(key, "]") {
		field = key[:i]
		operator = key[i+1 : len(key)-1]
	} else {
		field = key
	}
	return
}

func validAPIFieldName(s string) bool {
	if s == "" {
		return false
	}
	// NoSQL injection prevention - block dangerous characters
	dangerousChars := []string{"$", ".", "(", ")", "[", "]", "{", "}", "\\", "/", ":", ";", "=", "<", ">", "!", "@", "#", "%", "^", "&", "*", "|", "~", "`", "'", "\"", "?", "+", "-"}
	for _, char := range dangerousChars {
		if strings.Contains(s, char) {
			return false
		}
	}
	// Block common NoSQL injection patterns
	dangerousPatterns := []string{
		"$where", "$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin", "$exists", "$regex",
		"$or", "$and", "$not", "$nor", "$all", "$elemMatch", "$size", "$type",
		"javascript:", "this.", "function", "eval", "exec", "script",
	}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(s), pattern) {
			return false
		}
	}
	return true
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func inSlice(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func mustBeAllowed(field string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, a := range allowed {
		if a == field {
			return true
		}
	}
	return false
}

func enforceLenCap(s string, cap int) string {
	if cap <= 0 {
		return s
	}
	if len(s) > cap {
		return s[:cap]
	}
	return s
}

func parseTimeMulti(v string, layouts []string) (time.Time, error) {
	for _, l := range layouts {
		if t, err := time.Parse(l, v); err == nil {
			return t, nil
		}
	}
	return time.Parse(time.RFC3339, v)
}

func firstNonEmpty(v []string) []string {
	if len(v) > 0 {
		return v
	}
	return []string{time.RFC3339}
}

func parseBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "t", "true", "yes", "y":
		return true
	default:
		return false
	}
}

func coerceValue(raw string, spec FieldSpec, p QueryPolicy) any {
	if spec.Decoder != nil {
		if v, err := spec.Decoder(raw); err == nil {
			return v
		}
	}
	switch spec.Kind {
	case KindInt:
		if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return i
		}
	case KindFloat:
		if f, err := strconv.ParseFloat(raw, 64); err == nil {
			return f
		}
	case KindBool:
		return parseBool(raw)
	case KindTime:
		layouts := spec.Layouts
		if len(layouts) == 0 {
			layouts = p.TimeLayouts
		}
		for _, l := range layouts {
			if t, err := time.Parse(l, raw); err == nil {
				return t
			}
		}
	case KindObjectID:
		if oid, err := primitive.ObjectIDFromHex(raw); err == nil {
			return oid
		}
	}
	return raw
}

func coerceSlice(raws []string, spec FieldSpec, p QueryPolicy) []any {
	out := make([]any, 0, len(raws))
	for _, s := range raws {
		out = append(out, coerceValue(s, spec, p))
	}
	return out
}

func isOpAllowed(op string, spec FieldSpec) bool {
	if len(spec.AllowedOps) == 0 {
		// Default allowed operators
		switch op {
		case "eq", "ne", "gt", "gte", "lt", "lte", "in", "nin", "exists", "contains", "startsWith", "endsWith":
			return true
		case "regex":
			// only allow if explicitly enabled via AllowRawRegex
			return spec.AllowRawRegex
		default:
			return false
		}
	}
	for _, a := range spec.AllowedOps {
		if a == op {
			return true
		}
	}
	return false
}

func dbField(apiField string, p QueryPolicy) string {
	if spec, ok := p.Filtering.FieldSpecs[apiField]; ok && spec.Alias != "" {
		// Sanitize the alias to prevent injection
		return sanitizeFieldName(spec.Alias)
	}
	// Sanitize the field name to prevent injection
	return sanitizeFieldName(apiField)
}

// sanitizeFieldName ensures field names are safe for MongoDB queries
func sanitizeFieldName(fieldName string) string {
	// Remove any remaining dangerous characters
	dangerousChars := []string{"$", ".", "(", ")", "[", "]", "{", "}", "\\", "/", ":", ";", "=", "<", ">", "!", "@", "#", "%", "^", "&", "*", "|", "~", "`", "'", "\"", "?", "+", "-"}
	result := fieldName
	for _, char := range dangerousChars {
		result = strings.ReplaceAll(result, char, "")
	}
	return strings.TrimSpace(result)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
