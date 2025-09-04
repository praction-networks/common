package helpers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/praction-networks/common/appError"
	"go.mongodb.org/mongo-driver/bson"
)

type PaginatedFeedQuery struct {
	Limit          int                          `json:"limit" validate:"gte=1, lte=5000"` // Maximum items per page
	Offset         int                          `json:"offset" validate:"gte=0"`          // Offset for pagination
	Sort           map[string]string            `json:"sort"`                             // Dynamic sorting: field -> order (asc/desc)
	Filters        map[string]map[string]string `json:"filters"`                          // Dynamic search: field -> operator -> value
	IncludeFields  []string                     `json:"include_fields"`                   // Fields to include in response
	ExcludeFields  []string                     `json:"exclude_fields"`                   // Fields to exclude in response
	PaginationMeta bool                         `json:"pagination_meta"`                  // Return pagination metadata
	DistinctField  string                       `json:"distinct_field"`                   // Field for distinct results
}

// Parse parses query parameters into the PaginatedFeedQuery struct.
func (fq *PaginatedFeedQuery) Parse(r *http.Request) error {
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
	}

	// Parse limit with default value
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return appError.New(appError.InvalidInputError, "invalid 'limit' parameter, must be an integer", http.StatusBadRequest, err)
		}
		fq.Limit = l
	} else {
		fq.Limit = 10 // Default limit
	}

	// Parse offset with default value
	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return appError.New(appError.InvalidInputError, "invalid 'offset' parameter, must be an integer", http.StatusBadRequest, err)
		}
		fq.Offset = o
	} else {
		fq.Offset = 0 // Default offset
	}

	// Parse dynamic sorting
	fq.Sort = make(map[string]string)
	if sortFields := qs.Get("sort"); sortFields != "" {
		for _, field := range strings.Split(sortFields, ",") {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) == 2 {
				fieldName := parts[0]
				order := parts[1]
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

	// Parse dynamic filters
	fq.Filters = make(map[string]map[string]string)
	for key, values := range qs {
		// Skip explicitly whitelisted non-filter parameters
		if validParams[key] {
			continue
		}

		// Process valid filters
		if len(values) > 0 {
			if strings.Contains(key, "[") && strings.Contains(key, "]") {
				field := key[:strings.Index(key, "[")]
				operator := key[strings.Index(key, "[")+1 : strings.Index(key, "]")]
				if fq.Filters[field] == nil {
					fq.Filters[field] = make(map[string]string)
				}
				fq.Filters[field][operator] = values[0]
			} else {
				if fq.Filters[key] == nil {
					fq.Filters[key] = make(map[string]string)
				}
				fq.Filters[key]["eq"] = values[0]
			}
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

// BuildMongoFilter builds MongoDB filter from dynamic filters
// This is a generic utility function that converts PaginatedFeedQuery filters
// into MongoDB query filters for use across all repositories
func BuildMongoFilter(filters map[string]map[string]string) bson.M {
	query := bson.M{}
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
			}
		}
		if len(sub) > 0 {
			query[field] = sub
		}
	}
	return query
}
