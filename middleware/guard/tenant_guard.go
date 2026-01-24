package guard

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/caching/hierarchy"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
)

// TenantIDExtractor is a function that extracts the relevant tenant ID from the request
type TenantIDExtractor func(r *http.Request) string

// TenantHierarchyGuardMiddleware creates a middleware that enforces hierarchy access rules
// Rules:
// 1. System Users (SuperAdmin) -> ALLOW
// 2. Resource Tenant == Context Tenant -> ALLOW (Self access)
// 3. Resource Tenant is a descendant of Context Tenant -> ALLOW (Parent accessing Child)
// 4. Otherwise -> DENY
func TenantHierarchyGuardMiddleware(
	cache hierarchy.TenantHierarchyCache,
	extractTenantID TenantIDExtractor,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. System Users: Allow implicit access
			// We check context user roles usually set by auth middleware
			if helpers.IsSystemUser(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Extract resource tenant ID
			resourceTenantID := extractTenantID(r)
			if resourceTenantID == "" {
				// If no tenant ID found, we might skip or deny.
				// For safety in a guard, if we can't key the resource, we should probably proceed
				// only if it's not a tenant-scoped resource, but this middleware IS for tenant-scoped.
				// However, maybe the extractor logic has a fallback?
				// Similar to the one in tenant-service, we'll log and skip if empty (assume route mismatch or system route)
				// But strictly speaking, if applied, it should match.
				logger.Debug("Skipping hierarchy guard: no tenant ID extracted")
				next.ServeHTTP(w, r)
				return
			}

			// 3. Get Context Tenant (who is calling?)
			contextTenantID := helpers.GetTenantID(r.Context())
			if contextTenantID == "" {
				helpers.HandleAppError(w, appError.New(
					appError.UnauthorizedAccess,
					"Tenant context is required",
					http.StatusForbidden,
					nil,
				))
				return
			}

			// 4. Same Tenant Check
			if contextTenantID == resourceTenantID {
				next.ServeHTTP(w, r)
				return
			}

			// 5. Hierarchy Check (Is Resource a Child of Context?)
			// We check if resourceTenant is a child of contextTenantID
			// (i.e. Context Tenant is an Ancestor of Resource Tenant)
			if cache.IsChild(contextTenantID, resourceTenantID) {
				// Allowed: Parent accessing child resource
				next.ServeHTTP(w, r)
				return
			}

			// 6. Access Denied
			// If we are here: Not System, Not Self, Not Parent -> Deny
			logger.Warn("Access denied: Tenant Hierarchy Violation",
				"contextTenant", contextTenantID,
				"resourceTenant", resourceTenantID)

			helpers.HandleAppError(w, appError.New(
				appError.UnauthorizedAccess,
				"Access denied: You do not have permission to access resources of this tenant",
				http.StatusForbidden,
				nil,
			))
		})
	}
}

// SystemLevelGuardMiddleware creates a middleware that protects system-level routes
// Rules:
// 1. System Users (SuperAdmin) -> ALLOW
// 2. System Tenant -> ALLOW
// 3. Otherwise -> DENY
func SystemLevelGuardMiddleware(cache hierarchy.TenantHierarchyCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. System Users: Allow
			if helpers.IsSystemUser(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Get Context Tenant
			contextTenantID := helpers.GetTenantID(r.Context())
			if contextTenantID == "" {
				helpers.HandleAppError(w, appError.New(
					appError.UnauthorizedAccess,
					"Tenant context is required",
					http.StatusForbidden,
					nil,
				))
				return
			}

			// 3. Check if Context Tenant is System Tenant
			// We use the cache to look up the tenant details
			tenantData, exists := cache.Get(contextTenantID)
			if !exists {
				// If not in cache, we cannot verify if it's system.
				// For security, we must deny or fallback.
				// Assuming cache is source of truth for "System" flag.
				// However, if cache is partial, this might be issue.
				// But we load FULL cache at startup.
				logger.Warn("Access denied: Tenant not found in hierarchy cache during system check", "tenantID", contextTenantID)
				helpers.HandleAppError(w, appError.New(
					appError.UnauthorizedAccess,
					"Access denied: Tenant verification failed",
					http.StatusForbidden,
					nil,
				))
				return
			}

			if tenantData.IsSystem {
				// Allowed: System Tenant
				next.ServeHTTP(w, r)
				return
			}

			// 4. Access Denied
			logger.Warn("Access denied: System-level resource requires system tenant or system user",
				"contextTenant", contextTenantID)

			helpers.HandleAppError(w, appError.New(
				appError.UnauthorizedAccess,
				"Access denied: system-level resources are only accessible by system tenant or system users",
				http.StatusForbidden,
				nil,
			))
		})
	}
}

// Extractor Helpers

// ExtractFromPath returns an extractor that pulls ID from URL path param (Chi style)
func ExtractFromPath(paramName string) TenantIDExtractor {
	return func(r *http.Request) string {
		return chi.URLParam(r, paramName)
	}
}

// AccessibleTenantsMiddleware injects the list of accessible tenant IDs into the request context
// This list includes the context tenant + all its descendants (children, grandchildren, etc.)
// Handlers can use helpers.GetAccessibleTenants(ctx) to retrieve this list for query filtering
//
// Usage:
//
//	r.Use(guard.AccessibleTenantsMiddleware(cache))
//
// Then in handler/repository:
//
//	accessibleTenants := helpers.GetAccessibleTenants(ctx)
//	filter := bson.M{"tenantId": bson.M{"$in": accessibleTenants}}
//
// Note: For system users, this middleware is SKIPPED (no list injected).
// Handlers should check if GetAccessibleTenants returns nil (system user case) and handle accordingly.
func AccessibleTenantsMiddleware(cache hierarchy.TenantHierarchyCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// System users: Skip injection (they have access to all tenants)
			if helpers.IsSystemUser(ctx) {
				next.ServeHTTP(w, r)
				return
			}

			// Get context tenant ID
			contextTenantID := helpers.GetTenantID(ctx)
			if contextTenantID == "" {
				// No tenant context - skip injection
				next.ServeHTTP(w, r)
				return
			}

			// Compute accessible tenants: self + all descendants
			descendants := cache.GetDescendants(contextTenantID)
			accessibleTenants := append([]string{contextTenantID}, descendants...)

			// Inject into context
			ctx = helpers.SetAccessibleTenants(ctx, accessibleTenants)

			logger.Debug("Injected accessible tenants into context",
				"contextTenant", contextTenantID,
				"totalAccessible", len(accessibleTenants))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

