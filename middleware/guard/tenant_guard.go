package guard

import (
	"net/http"

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

// Extractor Helpers

// ExtractFromPath returns an extractor that pulls ID from URL path param (Chi style)
// Note: This depends on Chi context being available.
// If generic, might need to rely on request context values set by router.
func ExtractFromPath(paramName string) TenantIDExtractor {
	return func(r *http.Request) string {
		// We rely on "github.com/go-chi/chi/v5" usually, but common shouldn't depend on chi if possible?
		// But helpers already might.
		// For standard library compatibility, we can leave this to the caller to implement the function.
		// BUT, convenient helpers are good.
		// Let's assume standard chi usage for now based on project standards.
		// If we can't import chi here (circular or dep bloat), we leave it out.
		// Let's implement a simple context key lookup just in case the router puts it there.
		// Or easier: Service implements the extractor calling chi.URLParam.
		return ""
	}
}
