package guard

import (
	"net/http"
	"strings"

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
				logger.Warn("Access denied: Tenant not found in hierarchy cache during system check", "tenant_id", contextTenantID)
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
// Note: For system users WITHOUT a tenant context, this middleware is SKIPPED (no list injected).
// For system users WITH a tenant context, accessible tenants are computed for that context.
// Handlers should check if GetAccessibleTenants returns nil (global access case) and handle accordingly.
func AccessibleTenantsMiddleware(cache hierarchy.TenantHierarchyCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get context tenant ID first
			contextTenantID := helpers.GetTenantID(ctx)

			// System users (SuperAdmin / scope=system in JWT): ALWAYS grant global
			// access regardless of X-Tenant-ID value. Zero-Trust IAM design:
			// system scope bypasses tenant scoping everywhere. X-Tenant-ID is
			// informational only for system users — the dashboard sends it to
			// satisfy forward-auth requirements, but the guard ignores it.
			if helpers.IsSystemUser(ctx) {
				logger.Debug("System user - granting global access (X-Tenant-ID ignored for scoping)",
					"contextTenantID", contextTenantID)
				next.ServeHTTP(w, r)
				return
			}

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
				"totalAccessible", len(accessibleTenants),
				"isSystemUser", helpers.IsSystemUser(ctx))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TenantSetupGuardMiddleware creates a middleware that blocks API calls from tenants
// that haven't completed their post-login setup wizard.
//
// Rules:
//  1. System Users → ALLOW (SuperAdmin always has access)
//  2. System Tenant → ALLOW (system tenant is always fully set up)
//  3. Tenant with setupComplete=true → ALLOW
//  4. Specific allowed paths (GET /tenants/*, PATCH /tenants/*) → ALLOW (wizard needs these)
//  5. Auth endpoints → ALLOW (token refresh, logout, etc.)
//  6. Otherwise → DENY with HTTP 403 + message to complete setup
func TenantSetupGuardMiddleware(cache hierarchy.TenantHierarchyCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. System Users: Always allow
			if helpers.IsSystemUser(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Get context tenant
			contextTenantID := helpers.GetTenantID(r.Context())
			if contextTenantID == "" {
				// No tenant context — skip (auth endpoints, health checks, etc.)
				next.ServeHTTP(w, r)
				return
			}

			// 3. Look up tenant in cache. GetFresh consults the optional
			//    TenantHierarchyProvider when the cached entry is stale
			//    (SetupComplete=false). This breaks the dashboard ↔ /setup
			//    loop that fires when a tenant flips SetupComplete=true but
			//    the in-memory cache misses the NATS tenant.updated event
			//    (consumer offline, replica behind, etc.).
			tenantData, exists := cache.GetFresh(r.Context(), contextTenantID)
			if !exists {
				// Not in cache — allow through (cache might not be populated yet)
				logger.Debug("Tenant not found in hierarchy cache, skipping setup guard", "tenant_id", contextTenantID)
				next.ServeHTTP(w, r)
				return
			}

			// 4. System tenant or setup already complete — allow
			if tenantData.IsSystem || tenantData.SetupComplete {
				next.ServeHTTP(w, r)
				return
			}

			// 5. Allow specific endpoints needed by the setup wizard:
			//    - GET/PATCH on tenant endpoints (wizard reads and updates tenant)
			//    - Auth endpoints (token refresh, policies, etc.)
			path := r.URL.Path
			if isSetupAllowedEndpoint(r.Method, path) {
				next.ServeHTTP(w, r)
				return
			}

			// 6. Setup not complete — reject
			logger.Warn("API call rejected: tenant setup not complete",
				"tenant_id", contextTenantID,
				"method", r.Method,
				"path", path)

			helpers.HandleAppError(w, appError.New(
				appError.InvalidOperation,
				"Tenant setup is incomplete. Please complete the setup wizard before accessing other features.",
				http.StatusForbidden,
				nil,
			))
		})
	}
}

// isSetupAllowedEndpoint returns true for endpoints that the setup wizard needs
func isSetupAllowedEndpoint(method, path string) bool {
	// Auth / access control endpoints (always allowed)
	if strings.Contains(path, "/auth/") ||
		strings.Contains(path, "/access/") {
		return true
	}

	// Setup wizard only needs to read/update its OWN tenant by ID
	// Pattern: /api/v1/tenant/{tenantId} or /api/v1/tenant/{tenantId}/theme
	// Do NOT allow listing all tenants (GET /api/v1/tenant) — that should be blocked until setup is done
	if strings.Contains(path, "/tenant/") &&
		(method == http.MethodGet || method == http.MethodPatch || method == http.MethodPut) {
		// Verify this is a specific tenant path (has an ID segment after /tenant/)
		// NOT /tenant/hierarchy, /tenant/children etc. that are bulk operations
		parts := strings.Split(path, "/tenant/")
		if len(parts) > 1 {
			remaining := parts[len(parts)-1]
			// Allow: {tenantId}, {tenantId}/theme, {tenantId}/domain/verify
			// Block: hierarchy, children, descendants, etc.
			if remaining != "" && !strings.Contains(remaining, "hierarchy") &&
				!strings.Contains(remaining, "children") &&
				!strings.Contains(remaining, "descendants") &&
				!strings.Contains(remaining, "ancestors") &&
				!strings.Contains(remaining, "siblings") {
				return true
			}
		}
	}

	// Domain verification (setup wizard triggers POST to re-verify domain/SSL)
	if strings.Contains(path, "/domain/verify") && method == http.MethodPost {
		return true
	}

	// Tenant-user reads (dashboard fetches current user profile during setup)
	if strings.Contains(path, "/tenant-users") &&
		(method == http.MethodGet || method == http.MethodPatch) {
		return true
	}

	// Health / readiness probes
	if strings.HasSuffix(path, "/healthz") ||
		strings.HasSuffix(path, "/readyz") ||
		strings.HasSuffix(path, "/health") {
		return true
	}

	return false
}
