// common/helpers/querybuilder.go
package helpers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/metrics"
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
	// Optional: Query caching configuration
	CacheConfig *CacheConfig
	// Optional: Rate limiting configuration
	RateLimitConfig *RateLimitConfig
	// Optional: Full-text search configuration
	FullTextSearch *FullTextSearchConfig
	// Optional: Geospatial search configuration
	GeoSpatial *GeoSpatialConfig
}

// CacheConfig configures query result caching
type CacheConfig struct {
	Enabled   bool          // Enable caching
	TTL       time.Duration // Cache TTL
	MaxSize   int           // Maximum cache entries
	CacheImpl QueryCache    // Cache implementation (optional, uses in-memory by default)
}

// RateLimitConfig configures rate limiting
type RateLimitConfig struct {
	Enabled     bool          // Enable rate limiting
	MaxRequests int           // Maximum requests per window
	Window      time.Duration // Time window
	LimiterImpl RateLimiter   // Rate limiter implementation (optional, uses in-memory by default)
}

// FullTextSearchConfig configures MongoDB full-text search
type FullTextSearchConfig struct {
	Enabled         bool     // Enable full-text search
	DefaultLanguage string   // Default language for text search (e.g., "english")
	TextIndexFields []string // Fields that have text indexes
}

// GeoSpatialConfig configures geospatial queries
type GeoSpatialConfig struct {
	Enabled       bool    // Enable geospatial queries
	LocationField string  // Field name containing location (GeoJSON Point)
	MaxDistance   float64 // Maximum distance in meters for $geoNear
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
	TextSearch     *TextSearchQuery
	GeoQuery       *GeoQuery
}

// Output: ready-to-use Mongo query
type MongoQuery struct {
	Filter        bson.M
	FindOptions   *options.FindOptions
	DistinctField string
	WithMeta      bool
	Limit         int
	Offset        int
	// Aggregation pipeline (optional, for complex queries)
	Pipeline []bson.M
	// Full-text search query
	TextSearch *TextSearchQuery
	// Geospatial query
	GeoQuery *GeoQuery
}

// TextSearchQuery represents a MongoDB full-text search query
type TextSearchQuery struct {
	Search        string // Search term
	Language      string // Language for stemming
	CaseSensitive bool   // Case sensitive search
}

// GeoQuery represents a geospatial query
type GeoQuery struct {
	Type        string    // "near", "within", "intersects"
	Coordinates []float64 // [longitude, latitude] for Point
	MaxDistance float64   // Maximum distance in meters (for $geoNear)
	MinDistance float64   // Minimum distance in meters (for $geoNear)
	Geometry    bson.M    // GeoJSON geometry (for $geoWithin, $geoIntersects)
}

/* ==============================
   PUBLIC ENTRYPOINT
   ============================== */

// QueryCache interface for caching query results
type QueryCache interface {
	Get(ctx context.Context, key string) (MongoQuery, bool)
	Set(ctx context.Context, key string, query MongoQuery, ttl time.Duration) error
	SetAsync(ctx context.Context, key string, query MongoQuery, ttl time.Duration) // Non-blocking async set
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// RateLimiter interface for rate limiting queries
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	Reset(ctx context.Context, key string) error
}

// In-memory query cache implementation
type inMemoryQueryCache struct {
	cache   map[string]cacheEntry
	mu      sync.RWMutex
	maxSize int
}

type cacheEntry struct {
	query     MongoQuery
	expiresAt time.Time
}

func NewInMemoryQueryCache(maxSize int) QueryCache {
	return &inMemoryQueryCache{
		cache:   make(map[string]cacheEntry),
		maxSize: maxSize,
	}
}

func (c *inMemoryQueryCache) Get(ctx context.Context, key string) (MongoQuery, bool) {
	// Fast path: check context cancellation first (non-blocking check)
	select {
	case <-ctx.Done():
		return MongoQuery{}, false
	default:
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return MongoQuery{}, false
	}

	if time.Now().After(entry.expiresAt) {
		// Expired, delete it (async cleanup to avoid blocking)
		go func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			// Double-check after acquiring write lock
			if entry, stillExists := c.cache[key]; stillExists && time.Now().After(entry.expiresAt) {
				delete(c.cache, key)
			}
		}()
		return MongoQuery{}, false
	}

	return entry.query, true
}

func (c *inMemoryQueryCache) Set(ctx context.Context, key string, query MongoQuery, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict oldest entries if cache is full
	if len(c.cache) >= c.maxSize {
		c.evictOldest()
	}

	c.cache[key] = cacheEntry{
		query:     query,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

// SetAsync performs a non-blocking async cache set operation
func (c *inMemoryQueryCache) SetAsync(ctx context.Context, key string, query MongoQuery, ttl time.Duration) {
	go func() {
		// Use a separate context with timeout to prevent goroutine leaks
		setCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Make a copy of the query to avoid race conditions
		queryCopy := query

		// Attempt to set in cache (non-blocking for caller)
		select {
		case <-setCtx.Done():
			// Timeout - cache set failed, but that's OK (non-critical)
			return
		default:
			_ = c.Set(setCtx, key, queryCopy, ttl)
		}
	}()
}

func (c *inMemoryQueryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
	return nil
}

func (c *inMemoryQueryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]cacheEntry)
	return nil
}

func (c *inMemoryQueryCache) evictOldest() {
	oldestKey := ""
	oldestTime := time.Now()

	for key, entry := range c.cache {
		if entry.expiresAt.Before(oldestTime) {
			oldestTime = entry.expiresAt
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}

// In-memory rate limiter implementation (token bucket)
type inMemoryRateLimiter struct {
	buckets     map[string]*rateBucket
	mu          sync.RWMutex
	maxRequests int
	window      time.Duration
}

type rateBucket struct {
	tokens     int
	lastRefill time.Time
}

func NewInMemoryRateLimiter(maxRequests int, window time.Duration) RateLimiter {
	return &inMemoryRateLimiter{
		buckets:     make(map[string]*rateBucket),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *inMemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	// Fast path: check context cancellation first (non-blocking check)
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	now := time.Now()

	if !exists {
		rl.buckets[key] = &rateBucket{
			tokens:     rl.maxRequests - 1,
			lastRefill: now,
		}
		return true, nil
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.lastRefill)
	if elapsed >= rl.window {
		bucket.tokens = rl.maxRequests
		bucket.lastRefill = now
	} else {
		// Refill proportionally
		refillAmount := int(float64(rl.maxRequests) * elapsed.Seconds() / rl.window.Seconds())
		if refillAmount > 0 {
			bucket.tokens = min(rl.maxRequests, bucket.tokens+refillAmount)
			bucket.lastRefill = now
		}
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true, nil
	}

	return false, nil
}

func (rl *inMemoryRateLimiter) Reset(ctx context.Context, key string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.buckets, key)
	return nil
}

// BuildFromRequest builds a MongoDB query from HTTP request with caching and rate limiting
func BuildFromRequest(r *http.Request, policy QueryPolicy) (MongoQuery, error) {
	ctx := r.Context()

	// Rate limiting check
	if policy.RateLimitConfig != nil && policy.RateLimitConfig.Enabled {
		limiter := policy.RateLimitConfig.LimiterImpl
		if limiter == nil {
			limiter = NewInMemoryRateLimiter(policy.RateLimitConfig.MaxRequests, policy.RateLimitConfig.Window)
		}

		// Use client IP or user ID as rate limit key
		rateLimitKey := getRateLimitKey(r)
		allowed, err := limiter.Allow(ctx, rateLimitKey)
		if err != nil {
			return MongoQuery{}, appError.New(appError.InternalServerError, "Rate limit check failed", http.StatusInternalServerError, err)
		}
		if !allowed {
			metrics.QueryRateLimitHits.WithLabelValues(rateLimitKey).Inc()
			return MongoQuery{}, appError.New(appError.RateLimitExceeded, "Rate limit exceeded", http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded for key: %s", rateLimitKey))
		}
	}

	// Check cache if enabled
	if policy.CacheConfig != nil && policy.CacheConfig.Enabled {
		cache := policy.CacheConfig.CacheImpl
		if cache == nil {
			cache = NewInMemoryQueryCache(policy.CacheConfig.MaxSize)
		}

		cacheKey := generateCacheKey(r, policy)
		if cached, found := cache.Get(ctx, cacheKey); found {
			metrics.QueryCacheHits.Inc()
			return cached, nil
		}
		metrics.QueryCacheMisses.Inc()
	}

	// Parse and build query
	startTime := time.Now()
	var fq PaginatedFeedQuery
	if err := parseInto(&fq, r, policy); err != nil {
		metrics.QueryBuildErrors.WithLabelValues("parse").Inc()
		return MongoQuery{}, err
	}

	query, err := buildMongoQuery(fq, policy)
	if err != nil {
		metrics.QueryBuildErrors.WithLabelValues("build").Inc()
		return MongoQuery{}, err
	}

	// Record query build duration (non-blocking - Prometheus metrics are atomic)
	duration := time.Since(startTime)
	metrics.QueryBuildDuration.Observe(duration.Seconds())
	metrics.QueryBuildsTotal.Inc()

	// Cache the query if enabled (async/non-blocking)
	if policy.CacheConfig != nil && policy.CacheConfig.Enabled {
		cache := policy.CacheConfig.CacheImpl
		if cache == nil {
			cache = NewInMemoryQueryCache(policy.CacheConfig.MaxSize)
		}
		cacheKey := generateCacheKey(r, policy)
		ttl := policy.CacheConfig.TTL
		if ttl == 0 {
			ttl = 5 * time.Minute // Default TTL
		}

		// Use async set to avoid blocking the request
		// Cache failures are non-critical - query was built successfully
		if asyncCache, ok := cache.(interface {
			SetAsync(context.Context, string, MongoQuery, time.Duration)
		}); ok {
			asyncCache.SetAsync(ctx, cacheKey, query, ttl)
		} else {
			// Fallback to sync set if async not supported
			go func() {
				_ = cache.Set(context.Background(), cacheKey, query, ttl)
			}()
		}
	}

	return query, nil
}

// generateCacheKey creates a cache key from request and policy
func generateCacheKey(r *http.Request, _ QueryPolicy) string {
	key := r.Method + ":" + r.URL.Path + "?" + r.URL.RawQuery
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// getRateLimitKey extracts a key for rate limiting (IP or user ID)
func getRateLimitKey(r *http.Request) string {
	// Try to get user ID from header first
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return "user:" + userID
	}
	// Fall back to IP address
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
		}
	}
	return "ip:" + ip
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
		"text_search": true, "text_language": true, "text_case_sensitive": true,
		"lat": true, "lon": true, "geo_type": true, "max_distance": true, "min_distance": true,
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
		// Note: No sanitization needed here - MongoDB driver handles injection prevention
		// when using bson.M with proper type coercion. Field names are validated separately.
		// Values are properly escaped in buildFilter() using regexp.QuoteMeta for regex ops.
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
		// Note: No sanitization needed - search terms are properly escaped using regexp.QuoteMeta
		// in buildFilter() when building the $regex query. This preserves legitimate special characters.
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

	// --- Full-text search ---
	if v := qs.Get("text_search"); v != "" && p.FullTextSearch != nil && p.FullTextSearch.Enabled {
		fq.TextSearch = &TextSearchQuery{
			Search:        strings.TrimSpace(v),
			Language:      qs.Get("text_language"),
			CaseSensitive: qs.Get("text_case_sensitive") == "true",
		}
		if fq.TextSearch.Language == "" {
			fq.TextSearch.Language = p.FullTextSearch.DefaultLanguage
		}
		if p.Limits.MaxSearchLength > 0 && len(fq.TextSearch.Search) > p.Limits.MaxSearchLength {
			fq.TextSearch.Search = enforceLenCap(fq.TextSearch.Search, p.Limits.MaxSearchLength)
		}
	}

	// --- Geospatial queries ---
	if p.GeoSpatial != nil && p.GeoSpatial.Enabled {
		if latStr := qs.Get("lat"); latStr != "" {
			if lonStr := qs.Get("lon"); lonStr != "" {
				lat, err1 := strconv.ParseFloat(latStr, 64)
				lon, err2 := strconv.ParseFloat(lonStr, 64)
				if err1 == nil && err2 == nil {
					fq.GeoQuery = &GeoQuery{
						Type:        qs.Get("geo_type"),  // "near", "within", "intersects"
						Coordinates: []float64{lon, lat}, // GeoJSON format: [longitude, latitude]
						MaxDistance: p.GeoSpatial.MaxDistance,
					}
					if maxDistStr := qs.Get("max_distance"); maxDistStr != "" {
						if maxDist, err := strconv.ParseFloat(maxDistStr, 64); err == nil {
							fq.GeoQuery.MaxDistance = maxDist
						}
					}
					if minDistStr := qs.Get("min_distance"); minDistStr != "" {
						if minDist, err := strconv.ParseFloat(minDistStr, 64); err == nil {
							fq.GeoQuery.MinDistance = minDist
						}
					}
					if fq.GeoQuery.Type == "" {
						fq.GeoQuery.Type = "near" // Default to $geoNear
					}
				}
			}
		}
	}

	// Validate query complexity
	if err := validateQueryComplexity(*fq, p); err != nil {
		return err
	}

	// Record query complexity metric
	complexity := calculateQueryComplexity(*fq)
	metrics.QueryComplexity.Observe(float64(complexity))

	return nil
}

// calculateQueryComplexity calculates a complexity score for the query
func calculateQueryComplexity(fq PaginatedFeedQuery) int {
	complexity := 0
	complexity += len(fq.Filters) * 2 // Each filter adds complexity
	complexity += len(fq.Sort)        // Each sort adds complexity
	if fq.Search != "" {
		complexity += len(strings.Fields(fq.Search)) * 2 // Search terms add complexity
	}
	if fq.TextSearch != nil {
		complexity += 5 // Full-text search is more complex
	}
	if fq.GeoQuery != nil {
		complexity += 5 // Geospatial queries are more complex
	}
	if fq.DateFrom != nil || fq.DateTo != nil {
		complexity += 2 // Date ranges add complexity
	}
	return complexity
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

	query := MongoQuery{
		Filter:        filter,
		FindOptions:   opts,
		DistinctField: fq.DistinctField,
		WithMeta:      fq.PaginationMeta,
		Limit:         fq.Limit,
		Offset:        fq.Offset,
		TextSearch:    fq.TextSearch,
		GeoQuery:      fq.GeoQuery,
	}

	// Build aggregation pipeline for complex queries (full-text search, geospatial)
	if fq.TextSearch != nil || fq.GeoQuery != nil {
		pipeline := buildAggregationPipeline(fq, p, filter, opts)
		query.Pipeline = pipeline
	}

	return query, nil
}

// buildAggregationPipeline builds MongoDB aggregation pipeline for complex queries
func buildAggregationPipeline(fq PaginatedFeedQuery, p QueryPolicy, filter bson.M, _ *options.FindOptions) []bson.M {
	pipeline := []bson.M{}

	// $geoNear must be first stage if present
	if fq.GeoQuery != nil && p.GeoSpatial != nil {
		geoStage := bson.M{
			"$geoNear": bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": fq.GeoQuery.Coordinates,
				},
				"distanceField": "distance",
				"spherical":     true,
			},
		}
		if fq.GeoQuery.MaxDistance > 0 {
			geoStage["$geoNear"].(bson.M)["maxDistance"] = fq.GeoQuery.MaxDistance
		}
		if fq.GeoQuery.MinDistance > 0 {
			geoStage["$geoNear"].(bson.M)["minDistance"] = fq.GeoQuery.MinDistance
		}
		if p.GeoSpatial.LocationField != "" {
			geoStage["$geoNear"].(bson.M)["key"] = p.GeoSpatial.LocationField
		}
		pipeline = append(pipeline, geoStage)
	}

	// $match stage for filters
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	// $text search stage (must come after $geoNear if present, but before $match)
	if fq.TextSearch != nil && p.FullTextSearch != nil {
		textStage := bson.M{
			"$match": bson.M{
				"$text": bson.M{
					"$search": fq.TextSearch.Search,
				},
			},
		}
		if fq.TextSearch.Language != "" {
			textStage["$match"].(bson.M)["$text"].(bson.M)["$language"] = fq.TextSearch.Language
		}
		if fq.TextSearch.CaseSensitive {
			textStage["$match"].(bson.M)["$text"].(bson.M)["$caseSensitive"] = true
		}
		// Insert text search before $match if no $geoNear, otherwise after
		if fq.GeoQuery == nil {
			pipeline = append([]bson.M{textStage}, pipeline...)
		} else {
			// Insert after $geoNear
			pipeline = append(pipeline[:1], append([]bson.M{textStage}, pipeline[1:]...)...)
		}
	}

	// $sort stage
	if len(fq.Sort) > 0 {
		sortDoc := bson.M{}
		for _, s := range fq.Sort {
			dir := 1
			if s.Order == "desc" {
				dir = -1
			}
			sortDoc[dbField(s.Field, p)] = dir
		}
		pipeline = append(pipeline, bson.M{"$sort": sortDoc})
	}

	// $skip and $limit stages
	if fq.Offset > 0 {
		pipeline = append(pipeline, bson.M{"$skip": fq.Offset})
	}
	if fq.Limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": min(fq.Limit, p.Pagination.MaxLimit)})
	}

	// $project stage for field projection
	if len(fq.IncludeFields) > 0 || len(fq.ExcludeFields) > 0 {
		proj := bson.M{}
		if len(fq.IncludeFields) > 0 {
			for _, f := range fq.IncludeFields {
				proj[dbField(f, p)] = 1
			}
		} else {
			for _, f := range fq.ExcludeFields {
				proj[dbField(f, p)] = 0
			}
		}
		pipeline = append(pipeline, bson.M{"$project": proj})
	}

	return pipeline
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
