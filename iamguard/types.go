package iamguard

import "strings"

// RouteKey is the canonical identifier for a route across chi / APISIX / seed.
// Path is the normalised template (see NormalisePath) so that:
//
//	chi          : /api/v1/auth/roles/{id}
//	APISIX uri   : /api/v1/auth/roles/*
//	seed Resource: auth/roles/{id}
//
// all collapse to the same RouteKey.
type RouteKey struct {
	Path   string // canonical path, no leading slash, no /api/v1 prefix, params as "{*}"
	Method string // upper-case verb, e.g. "GET"
}

// String returns "METHOD path" — handy for log lines and map keys.
func (k RouteKey) String() string {
	return k.Method + " " + k.Path
}

// SeedPolicy is the minimal projection of a service's Casbin seed policy the
// guard needs. Callers convert their internal models.Policy slice into this
// shape via ConvertSeedPolicies (see seed_walk.go) so this package stays
// decoupled from every service's models package.
type SeedPolicy struct {
	// Service is the canonical service name (e.g. "auth-service",
	// "olt-manager"). The guard filters by this when comparing.
	Service string

	// Resource is the template path used to compute the RouteKey (e.g.
	// "auth/users/{id}"). It will be normalised via NormalisePath.
	Resource string

	// Action is the HTTP verb (GET / POST / PUT / PATCH / DELETE).
	// Wildcard "*" entries are filtered out — they do not map to a single
	// HTTP method and would generate false drift.
	Action string

	// PermissionKey is the dotted UI key (e.g. "auth.users.view"). Used as
	// the Hint in drift entries.
	PermissionKey string

	// Name is the human-readable policy name. Used as a fallback Hint.
	Name string
}

// DriftKind categorises a single drift entry.
type DriftKind string

const (
	// ChiWithoutPolicy: chi handler exists but no Casbin seed policy covers
	// it. Symptom: every non-system role hits 403 at forward-auth time.
	ChiWithoutPolicy DriftKind = "chi_without_policy"

	// PolicyWithoutChi: a Casbin seed policy points at a path no chi handler
	// serves. Symptom: dead permission key in the UI; orphan policy.
	PolicyWithoutChi DriftKind = "policy_without_chi"

	// ChiWithoutAPISIX: chi handler exists but no APISIX route forwards
	// traffic to it. Symptom: 404 at the gateway.
	ChiWithoutAPISIX DriftKind = "chi_without_apisix"

	// APISIXWithoutChi: APISIX route forwards to a path no chi handler
	// serves. Symptom: 404 from upstream after forward-auth passes.
	APISIXWithoutChi DriftKind = "apisix_without_chi"

	// APISIXWithoutPolicy: APISIX route is protected (public=false) but no
	// Casbin seed policy covers it. Same 403 class as ChiWithoutPolicy.
	APISIXWithoutPolicy DriftKind = "apisix_without_policy"

	// PolicyWithoutAPISIX: seed policy exists but no APISIX route exposes it.
	// Symptom: permission key visible in UI but the route is unreachable.
	PolicyWithoutAPISIX DriftKind = "policy_without_apisix"
)

// DriftEntry is one row in the guard report.
type DriftEntry struct {
	Kind  DriftKind
	Route RouteKey
	// Hint is a short human-readable explanation. For ChiWithoutPolicy this
	// is the chi handler name; for PolicyWithoutChi it is the Casbin
	// PermissionKey or policy Name.
	Hint string
}

// Report is the result of one guard pass.
type Report struct {
	Service string
	Drift   []DriftEntry
	Counts  map[DriftKind]int
}

// HasDrift returns true when at least one drift entry was found.
func (r *Report) HasDrift() bool {
	return r != nil && len(r.Drift) > 0
}

// NormalisePath converts any route template (chi, APISIX uri, seed Resource)
// into a comparable canonical form:
//
//   - strip "/api/v1/" / "/api/v2/" / "/api/" prefix
//   - trim leading + trailing "/"
//   - lower-case
//   - replace every "{...}" path-parameter with "{*}" so chi's "{id}",
//     and the seed's "{factorId}" collide on the same canonical slot
//   - replace bare "*" path segments with "{*}" so APISIX's wildcard syntax
//     (`/api/v1/auth/roles/*`) lines up with chi/seed's bracketed params
//     (`auth/roles/{id}`).
func NormalisePath(p string) string {
	p = strings.TrimPrefix(p, "/api/v1/")
	p = strings.TrimPrefix(p, "/api/v2/")
	p = strings.TrimPrefix(p, "/api/")
	p = strings.Trim(p, "/")
	p = strings.ToLower(p)

	// First pass: collapse "{...}" into a single placeholder.
	var b strings.Builder
	b.Grow(len(p))
	inParam := false
	for _, r := range p {
		switch {
		case r == '{':
			inParam = true
			b.WriteString("{*}")
		case r == '}':
			inParam = false
		case inParam:
			// skip the parameter name itself
		default:
			b.WriteRune(r)
		}
	}

	// Second pass: rewrite bare "*" path segments into "{*}" so APISIX
	// wildcard syntax aligns with chi/seed bracketed params.
	parts := strings.Split(b.String(), "/")
	for i, part := range parts {
		if part == "*" {
			parts[i] = "{*}"
		}
	}
	return strings.Join(parts, "/")
}

// NormaliseMethod upper-cases an HTTP verb for comparison.
func NormaliseMethod(m string) string {
	return strings.ToUpper(strings.TrimSpace(m))
}
