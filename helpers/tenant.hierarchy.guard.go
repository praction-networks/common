package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
)

// ==================== CONSTANTS ====================

const (
	// TenantHierarchyCachePrefix is the Redis cache key prefix for tenant hierarchy data
	TenantHierarchyCachePrefix = "tenant:hierarchy:"

	// TenantHierarchyCacheTTL is the cache TTL for tenant hierarchy data (5 minutes)
	// Tenant hierarchy doesn't change frequently, so 5 minutes is a good balance
	TenantHierarchyCacheTTL = 5 * time.Minute

	// SystemTenantCacheTTL is the cache TTL for system tenant flag (10 minutes)
	// System tenant flag rarely changes, so longer TTL is acceptable
	SystemTenantCacheTTL = 10 * time.Minute
)

// ==================== INTERFACES ====================

// TenantHierarchyProvider interface allows each service to provide tenant data
// Services can implement this using their own tenant repository/cache/event service
type TenantHierarchyProvider interface {
	// GetTenantAncestors returns the list of ancestor tenant IDs for a given tenant
	GetTenantAncestors(ctx context.Context, tenantID string) ([]string, error)

	// GetTenantByID returns tenant hierarchy data (for caching and system flag check)
	GetTenantByID(ctx context.Context, tenantID string) (*TenantHierarchyData, error)
}

// TenantHierarchyData contains minimal hierarchy info needed for validation
type TenantHierarchyData struct {
	ID        string   `json:"id"`
	Ancestors []string `json:"ancestors"` // Array of ancestor tenant IDs
	Level     int      `json:"level"`     // Optional: hierarchy level
	IsSystem  bool     `json:"isSystem"`  // System tenant flag
}

// RedisClientInterface defines the minimal Redis interface needed for caching
// This allows services to use any Redis client implementation
type RedisClientInterface interface {
	Get(ctx context.Context, key string) RedisStringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) RedisStatusCmd
	Del(ctx context.Context, keys ...string) RedisIntCmd
}

// RedisStringCmd represents a Redis string command result
type RedisStringCmd interface {
	Result() (string, error)
}

// RedisStatusCmd represents a Redis status command result
type RedisStatusCmd interface {
	Err() error
}

// RedisIntCmd represents a Redis integer command result
type RedisIntCmd interface {
	Err() error
}

// ==================== REDIS CACHE HELPERS ====================

// getTenantHierarchyFromCache retrieves tenant hierarchy data from Redis cache
// Returns nil if cache miss, error, or Redis not available
func getTenantHierarchyFromCache(ctx context.Context, redisClient RedisClientInterface, tenantID string) *TenantHierarchyData {
	if redisClient == nil {
		return nil // Redis not available
	}

	cacheKey := TenantHierarchyCachePrefix + tenantID
	cmd := redisClient.Get(ctx, cacheKey)
	if cmd == nil {
		return nil
	}
	cachedData, err := cmd.Result()
	if err != nil {
		// Cache miss or error - return nil (will fallback to provider)
		return nil
	}

	var tenantData TenantHierarchyData
	if err := json.Unmarshal([]byte(cachedData), &tenantData); err != nil {
		logger.Warn("Failed to unmarshal cached tenant hierarchy data", err, "tenantID", tenantID)
		return nil
	}

	logger.Debug("Tenant hierarchy cache hit", "tenantID", tenantID)
	return &tenantData
}

// setTenantHierarchyCache stores tenant hierarchy data in Redis cache
// Fire-and-forget: doesn't block on cache failure
func setTenantHierarchyCache(ctx context.Context, redisClient RedisClientInterface, tenantID string, data *TenantHierarchyData) {
	if redisClient == nil {
		return // Redis not available
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cacheKey := TenantHierarchyCachePrefix + tenantID
		cacheData, err := json.Marshal(data)
		if err != nil {
			logger.Warn("Failed to marshal tenant hierarchy for cache", err, "tenantID", tenantID)
			return
		}

		cmd := redisClient.Set(cacheCtx, cacheKey, cacheData, TenantHierarchyCacheTTL)
		if cmd == nil {
			return
		}
		if err := cmd.Err(); err != nil {
			logger.Warn("Failed to cache tenant hierarchy data", err, "tenantID", tenantID)
		} else {
			logger.Debug("Cached tenant hierarchy data", "tenantID", tenantID, "ttl", TenantHierarchyCacheTTL)
		}
	}()
}

// ==================== GUARD FUNCTIONS ====================

// ValidateTenantHierarchyAccess ensures that child tenants cannot access parent tenant resources
// Returns error if access should be denied, nil if allowed
//
// Access Rules:
//   - ✅ System users: Allow all
//   - ✅ Same tenant: Allow
//   - ✅ Parent accessing child: Allow (resource tenant is descendant)
//   - ❌ Child accessing parent: DENY (resource tenant is ancestor)
//   - ❌ Sibling tenants: DENY (no relationship)
//
// Parameters:
//   - ctx: Request context (must contain tenant ID via GetTenantID)
//   - resourceTenantID: The tenant ID that owns the resource being accessed
//   - provider: Implementation of TenantHierarchyProvider to fetch tenant data
//   - redisClient: Optional Redis client for caching (can be nil)
func ValidateTenantHierarchyAccess(
	ctx context.Context,
	resourceTenantID string,
	provider TenantHierarchyProvider,
	redisClient RedisClientInterface,
) error {
	// 1. System users: Allow all
	if IsSystemUser(ctx) {
		return nil
	}

	// 2. Get context tenant ID
	contextTenantID := GetTenantID(ctx)
	if contextTenantID == "" {
		// No tenant context - deny (should not happen in normal flow)
		logger.Warn("No tenant context in hierarchy validation", "resourceTenantID", resourceTenantID)
		return appError.New(
			appError.UnauthorizedAccess,
			"Tenant context is required",
			http.StatusForbidden,
			nil,
		)
	}

	// 3. Same tenant: Allow
	if contextTenantID == resourceTenantID {
		return nil
	}

	// 4. Get context tenant hierarchy data (with Redis cache)
	var contextTenant *TenantHierarchyData
	var err error

	// Try cache first
	if redisClient != nil {
		contextTenant = getTenantHierarchyFromCache(ctx, redisClient, contextTenantID)
	}

	// Cache miss - fetch from provider
	if contextTenant == nil {
		contextTenant, err = provider.GetTenantByID(ctx, contextTenantID)
		if err != nil {
			logger.Error("Failed to fetch context tenant for hierarchy validation", err,
				"contextTenantID", contextTenantID)
			return appError.New(
				appError.DBFetchError,
				"Failed to validate tenant access",
				http.StatusInternalServerError,
				err,
			)
		}

		// Cache for future use
		if redisClient != nil {
			setTenantHierarchyCache(ctx, redisClient, contextTenantID, contextTenant)
		}
	}

	// 5. Get resource tenant hierarchy data (with Redis cache)
	var resourceTenant *TenantHierarchyData

	// Try cache first
	if redisClient != nil {
		resourceTenant = getTenantHierarchyFromCache(ctx, redisClient, resourceTenantID)
	}

	// Cache miss - fetch from provider
	if resourceTenant == nil {
		resourceTenant, err = provider.GetTenantByID(ctx, resourceTenantID)
		if err != nil {
			logger.Error("Failed to fetch resource tenant for hierarchy validation", err,
				"resourceTenantID", resourceTenantID)
			return appError.New(
				appError.EntityNotFound,
				"Resource tenant not found",
				http.StatusNotFound,
				err,
			)
		}

		// Cache for future use
		if redisClient != nil {
			setTenantHierarchyCache(ctx, redisClient, resourceTenantID, resourceTenant)
		}
	}

	// 6. Check if resource tenant is ancestor of context tenant → DENY
	for _, ancestorID := range contextTenant.Ancestors {
		if ancestorID == resourceTenantID {
			logger.Warn("Access denied: child tenant cannot access parent tenant resource",
				"contextTenantID", contextTenantID,
				"resourceTenantID", resourceTenantID)
			return appError.New(
				appError.UnauthorizedAccess,
				"Access denied: child tenant cannot access parent tenant resources",
				http.StatusForbidden,
				nil,
			)
		}
	}

	// 7. Check if resource tenant is descendant of context tenant → ALLOW
	for _, ancestorID := range resourceTenant.Ancestors {
		if ancestorID == contextTenantID {
			logger.Debug("Access allowed: parent tenant accessing child tenant resource",
				"contextTenantID", contextTenantID,
				"resourceTenantID", resourceTenantID)
			return nil
		}
	}

	// 8. No relationship (sibling tenants) → DENY
	logger.Warn("Access denied: no hierarchical relationship between tenants",
		"contextTenantID", contextTenantID,
		"resourceTenantID", resourceTenantID)
	return appError.New(
		appError.UnauthorizedAccess,
		"Access denied: no hierarchical relationship between tenants",
		http.StatusForbidden,
		nil,
	)
}

// ValidateSystemLevelAccess ensures that system-level resources are only accessible by system tenant/system user
// Returns error if access should be denied, nil if allowed
//
// Access Rules:
//   - ✅ System users: Allow (helpers.IsSystemUser(ctx) == true)
//   - ✅ System tenant: Allow (tenant.IsSystem == true)
//   - ❌ All others: DENY
//
// Parameters:
//   - ctx: Request context (must contain tenant ID via GetTenantID)
//   - provider: Implementation of TenantHierarchyProvider to fetch tenant data
//   - redisClient: Optional Redis client for caching (can be nil)
func ValidateSystemLevelAccess(
	ctx context.Context,
	provider TenantHierarchyProvider,
	redisClient RedisClientInterface,
) error {
	// 1. System users: Allow
	if IsSystemUser(ctx) {
		return nil
	}

	// 2. Get context tenant ID
	contextTenantID := GetTenantID(ctx)
	if contextTenantID == "" {
		logger.Warn("No tenant context in system-level validation")
		return appError.New(
			appError.UnauthorizedAccess,
			"Tenant context is required",
			http.StatusForbidden,
			nil,
		)
	}

	// 3. Get context tenant data (with Redis cache)
	var contextTenant *TenantHierarchyData
	var err error

	// Try cache first
	if redisClient != nil {
		contextTenant = getTenantHierarchyFromCache(ctx, redisClient, contextTenantID)
	}

	// Cache miss - fetch from provider
	if contextTenant == nil {
		contextTenant, err = provider.GetTenantByID(ctx, contextTenantID)
		if err != nil {
			logger.Error("Failed to fetch context tenant for system-level validation", err,
				"contextTenantID", contextTenantID)
			return appError.New(
				appError.DBFetchError,
				"Failed to validate tenant access",
				http.StatusInternalServerError,
				err,
			)
		}

		// Cache for future use
		if redisClient != nil {
			setTenantHierarchyCache(ctx, redisClient, contextTenantID, contextTenant)
		}
	}

	// 4. Check if context tenant is system tenant
	if contextTenant.IsSystem {
		logger.Debug("Access allowed: system tenant accessing system-level resource",
			"contextTenantID", contextTenantID)
		return nil
	}

	// 5. DENY: Not system user and not system tenant
	logger.Warn("Access denied: system-level resource requires system tenant or system user",
		"contextTenantID", contextTenantID)
	return appError.New(
		appError.UnauthorizedAccess,
		"Access denied: system-level resources are only accessible by system tenant or system users",
		http.StatusForbidden,
		nil,
	)
}

// InvalidateTenantHierarchyCache invalidates Redis cache for a tenant
// Call this when tenant hierarchy changes (parent changed, deleted, etc.)
//
// Parameters:
//   - ctx: Request context
//   - redisClient: Redis client (can be nil if Redis not available)
//   - tenantID: Tenant ID to invalidate
func InvalidateTenantHierarchyCache(ctx context.Context, redisClient RedisClientInterface, tenantID string) error {
	if redisClient == nil {
		return nil // Redis not available - no-op
	}

	cacheKey := TenantHierarchyCachePrefix + tenantID
	if err := redisClient.Del(ctx, cacheKey).Err(); err != nil {
		logger.Warn("Failed to invalidate tenant hierarchy cache", err, "tenantID", tenantID)
		return err
	}

	logger.Debug("Invalidated tenant hierarchy cache", "tenantID", tenantID)
	return nil
}
