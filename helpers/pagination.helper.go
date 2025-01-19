package helpers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/praction-networks/common/appError"
)

type PaginatedFeedQuery struct {
	Limit          int                          `json:"limit" validate:"gte=1, lte=5000"`
	Offset         int                          `json:"offset" validate:"gte=0"`
	Sort           map[string]string            `json:"sort"`            // Dynamic sorting: field -> order (asc/desc)
	Filters        map[string]map[string]string `json:"filters"`         // Dynamic search: field -> operator -> value
	IncludeFields  []string                     `json:"include_fields"`  // Fields to include in response
	ExcludeFields  []string                     `json:"exclude_fields"`  // Fields to exclude in response
	PaginationMeta bool                         `json:"pagination_meta"` // Return pagination metadata
	DistinctField  string                       `json:"distinct_field"`  // Field for distinct results
}

// Parse parses query parameters into the PaginatedFeedQuery struct.
func (fq *PaginatedFeedQuery) Parse(r *http.Request) error {
	qs := r.URL.Query()

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
		if len(values) > 0 && key != "limit" && key != "offset" && key != "sort" {
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
