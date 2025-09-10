package helpers

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/praction-networks/common/appError"
	"go.mongodb.org/mongo-driver/bson"
)

type PaginatedFeedQuery struct {
	Limit          int                          `json:"limit" validate:"gte=1, lte=10000"` // Maximum items per page
	Offset         int                          `json:"offset" validate:"gte=0"`           // Offset for pagination
	Sort           map[string]string            `json:"sort"`                              // Dynamic sorting: field -> order (asc/desc)
	Filters        map[string]map[string]string `json:"filters"`                           // Dynamic search: field -> operator -> value
	IncludeFields  []string                     `json:"include_fields"`                    // Fields to include in response
	ExcludeFields  []string                     `json:"exclude_fields"`                    // Fields to exclude in response
	PaginationMeta bool                         `json:"pagination_meta"`                   // Return pagination metadata
	DistinctField  string                       `json:"distinct_field"`                    // Field for distinct results
	Search         string                       `json:"search"`                            // Full-text search query
	DateFrom       *time.Time                   `json:"date_from"`                         // Date range filter - from
	DateTo         *time.Time                   `json:"date_to"`                           // Date range filter - to
}

// PaginationConfig holds configuration for pagination behavior
type PaginationConfig struct {
	MaxLimit            int      `json:"max_limit"`             // Maximum allowed limit (default: 1000)
	DefaultLimit        int      `json:"default_limit"`         // Default limit (default: 10)
	AllowedSortFields   []string `json:"allowed_sort_fields"`   // Whitelist of sortable fields
	AllowedFilterFields []string `json:"allowed_filter_fields"` // Whitelist of filterable fields
	DateField           string   `json:"date_field"`            // Default date field for date range filtering
	SearchFields        []string `json:"search_fields"`         // Fields to search in for full-text search
}

// DefaultPaginationConfig returns a sensible default configuration
func DefaultPaginationConfig() *PaginationConfig {
	return &PaginationConfig{
		MaxLimit:            1000,
		DefaultLimit:        50,
		AllowedSortFields:   []string{"createdAt", "updatedAt", "name", "id"},
		AllowedFilterFields: []string{"isActive", "type", "environment", "status"},
		DateField:           "createdAt",
		SearchFields:        []string{"name", "description", "code"},
	}
}

// Parse parses query parameters into the PaginatedFeedQuery struct.
func (fq *PaginatedFeedQuery) Parse(r *http.Request) error {
	return fq.ParseWithConfig(r, DefaultPaginationConfig())
}

// ParseWithConfig parses query parameters with custom configuration
func (fq *PaginatedFeedQuery) ParseWithConfig(r *http.Request, config *PaginationConfig) error {
	qs := r.URL.Query()

	// Valid parameter keys for non-filter fields
	validParams := map[string]bool{
		"limit":           true,
		"offset":          true,
		"sort":            true,
		"pagination_meta": true,
		"include":         true,
		"exclude":         true,
		"distinct":        true,
		"search":          true,
		"date_from":       true,
		"date_to":         true,
	}

	// Parse limit with validation
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return appError.New(appError.InvalidInputError, "invalid 'limit' parameter, must be an integer", http.StatusBadRequest, err)
		}
		if l > config.MaxLimit {
			return appError.New(appError.InvalidInputError,
				"limit exceeds maximum allowed value", http.StatusBadRequest, nil)
		}
		if l < 1 {
			return appError.New(appError.InvalidInputError,
				"limit must be greater than 0", http.StatusBadRequest, nil)
		}
		fq.Limit = l
	} else {
		fq.Limit = config.DefaultLimit
	}

	// Parse offset with default value
	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return appError.New(appError.InvalidInputError, "invalid 'offset' parameter, must be an integer", http.StatusBadRequest, err)
		}
		if o < 0 {
			return appError.New(appError.InvalidInputError,
				"offset must be greater than or equal to 0", http.StatusBadRequest, nil)
		}
		fq.Offset = o
	} else {
		fq.Offset = 0 // Default offset
	}

	// Parse dynamic sorting with field validation
	fq.Sort = make(map[string]string)
	if sortFields := qs.Get("sort"); sortFields != "" {
		for _, field := range strings.Split(sortFields, ",") {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) == 2 {
				fieldName := parts[0]
				order := parts[1]

				// Validate sort field against whitelist
				if !isFieldAllowed(fieldName, config.AllowedSortFields) {
					return appError.New(appError.InvalidInputError,
						"sort field not allowed: "+fieldName, http.StatusBadRequest, nil)
				}

				if order == "asc" || order == "desc" {
					fq.Sort[fieldName] = order
				} else {
					return appError.New(appError.InvalidInputError, "invalid sort order, must be 'asc' or 'desc'", http.StatusBadRequest, nil)
				}
			} else {
				return appError.New(appError.InvalidInputError, "invalid sort format, must be 'field:asc' or 'field:desc'", http.StatusBadRequest, nil)
			}
		}
	}

	// Parse dynamic filters with field validation
	fq.Filters = make(map[string]map[string]string)
	for key, values := range qs {
		// Skip explicitly whitelisted non-filter parameters
		if validParams[key] {
			continue
		}

		// Process valid filters
		if len(values) > 0 {
			var fieldName string
			var operator string

			if strings.Contains(key, "[") && strings.Contains(key, "]") {
				fieldName = key[:strings.Index(key, "[")]
				operator = key[strings.Index(key, "[")+1 : strings.Index(key, "]")]
			} else {
				fieldName = key
				operator = "eq"
			}

			// Validate filter field against whitelist
			if !isFieldAllowed(fieldName, config.AllowedFilterFields) {
				return appError.New(appError.InvalidInputError,
					"filter field not allowed: "+fieldName, http.StatusBadRequest, nil)
			}

			// Validate operator
			if !isValidOperator(operator) {
				return appError.New(appError.InvalidInputError,
					"invalid filter operator: "+operator, http.StatusBadRequest, nil)
			}

			if fq.Filters[fieldName] == nil {
				fq.Filters[fieldName] = make(map[string]string)
			}
			fq.Filters[fieldName][operator] = values[0]
		}
	}

	// Parse included fields
	if include := qs.Get("include"); include != "" {
		fq.IncludeFields = strings.Split(include, ",")
	}

	// Parse excluded fields
	if exclude := qs.Get("exclude"); exclude != "" {
		fq.ExcludeFields = strings.Split(exclude, ",")
	}

	// Parse pagination metadata flag
	if paginationMeta := qs.Get("pagination_meta"); paginationMeta == "true" {
		fq.PaginationMeta = true
	}

	// Parse distinct field
	if distinct := qs.Get("distinct"); distinct != "" {
		fq.DistinctField = distinct
	}

	// Parse search query
	if search := qs.Get("search"); search != "" {
		fq.Search = strings.TrimSpace(search)
	}

	// Parse date range filters
	if dateFrom := qs.Get("date_from"); dateFrom != "" {
		if parsed, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			fq.DateFrom = &parsed
		} else {
			return appError.New(appError.InvalidInputError,
				"invalid date_from format, use RFC3339 (e.g., 2006-01-02T15:04:05Z)", http.StatusBadRequest, err)
		}
	}

	if dateTo := qs.Get("date_to"); dateTo != "" {
		if parsed, err := time.Parse(time.RFC3339, dateTo); err == nil {
			fq.DateTo = &parsed
		} else {
			return appError.New(appError.InvalidInputError,
				"invalid date_to format, use RFC3339 (e.g., 2006-01-02T15:04:05Z)", http.StatusBadRequest, err)
		}
	}

	return nil
}

// Offset returns the parsed offset if set, otherwise 0.
func (fq PaginatedFeedQuery) OffsetValue() int {
	if fq.Offset < 0 {
		return 0
	}
	return fq.Offset
}

func (fq PaginatedFeedQuery) Page() int {
	if fq.Limit == 0 {
		return 1
	}
	return (fq.Offset / fq.Limit) + 1
}

func (fq PaginatedFeedQuery) PageSize() int {
	return fq.Limit
}

// isFieldAllowed checks if a field is in the allowed list
func isFieldAllowed(field string, allowedFields []string) bool {
	for _, allowed := range allowedFields {
		if field == allowed {
			return true
		}
	}
	return false
}

// isValidOperator checks if the operator is valid
func isValidOperator(operator string) bool {
	validOps := []string{"eq", "ne", "gt", "gte", "lt", "lte", "in", "nin", "regex", "exists", "contains", "startsWith", "endsWith"}
	for _, op := range validOps {
		if operator == op {
			return true
		}
	}
	return false
}

// BuildMongoFilter builds MongoDB filter from dynamic filters
// This is a generic utility function that converts PaginatedFeedQuery filters
// into MongoDB query filters for use across all repositories
func BuildMongoFilter(filters map[string]map[string]string) bson.M {
	return BuildMongoFilterWithConfig(filters, "", nil, nil, nil, "")
}

// BuildMongoFilterWithConfig builds MongoDB filter with additional features
func BuildMongoFilterWithConfig(filters map[string]map[string]string, search string, searchFields []string, dateFrom, dateTo *time.Time, dateField string) bson.M {
	query := bson.M{}

	// Process regular filters
	for field, ops := range filters {
		if len(ops) == 1 && ops["eq"] != "" {
			query[field] = ops["eq"]
			continue
		}
		sub := bson.M{}
		for op, val := range ops {
			switch op {
			case "eq":
				sub["$eq"] = val
			case "ne":
				sub["$ne"] = val
			case "gt":
				sub["$gt"] = val
			case "gte":
				sub["$gte"] = val
			case "lt":
				sub["$lt"] = val
			case "lte":
				sub["$lte"] = val
			case "in":
				sub["$in"] = strings.Split(val, ",")
			case "nin":
				sub["$nin"] = strings.Split(val, ",")
			case "regex":
				sub["$regex"] = val
			case "exists":
				sub["$exists"] = val == "true"
			case "contains":
				sub["$regex"] = regexp.QuoteMeta(val)
			case "startsWith":
				sub["$regex"] = "^" + regexp.QuoteMeta(val)
			case "endsWith":
				sub["$regex"] = regexp.QuoteMeta(val) + "$"
			}
		}
		if len(sub) > 0 {
			query[field] = sub
		}
	}

	// Add full-text search
	if search != "" && len(searchFields) > 0 {
		searchConditions := make([]bson.M, 0, len(searchFields))
		for _, field := range searchFields {
			searchConditions = append(searchConditions, bson.M{
				field: bson.M{"$regex": regexp.QuoteMeta(search), "$options": "i"},
			})
		}
		if len(searchConditions) > 0 {
			query["$or"] = searchConditions
		}
	}

	// Add date range filter
	if dateFrom != nil || dateTo != nil {
		dateFilter := bson.M{}
		if dateFrom != nil {
			dateFilter["$gte"] = *dateFrom
		}
		if dateTo != nil {
			dateFilter["$lte"] = *dateTo
		}
		if len(dateFilter) > 0 {
			query[dateField] = dateFilter
		}
	}

	return query
}

// PaginationMetadata holds pagination information for responses
type PaginationMetadata struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
}

// BuildPaginationMetadata creates pagination metadata from query and total count
func BuildPaginationMetadata(query PaginatedFeedQuery, total int64) PaginationMetadata {
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationMetadata{
		Page:       query.Page(),
		PageSize:   query.PageSize(),
		Total:      total,
		TotalPages: totalPages,
		HasNext:    query.Page() < totalPages,
		HasPrev:    query.Page() > 1,
		Offset:     query.OffsetValue(),
		Limit:      query.Limit,
	}
}

// ValidatePaginationConfig validates the pagination configuration
func ValidatePaginationConfig(config *PaginationConfig) error {
	if config.MaxLimit <= 0 {
		return appError.New(appError.InvalidInputError, "MaxLimit must be greater than 0", http.StatusInternalServerError, nil)
	}
	if config.DefaultLimit <= 0 {
		return appError.New(appError.InvalidInputError, "DefaultLimit must be greater than 0", http.StatusInternalServerError, nil)
	}
	if config.DefaultLimit > config.MaxLimit {
		return appError.New(appError.InvalidInputError, "DefaultLimit cannot be greater than MaxLimit", http.StatusInternalServerError, nil)
	}
	return nil
}
