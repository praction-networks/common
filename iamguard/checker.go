package iamguard

import (
	"context"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/common/logger"
)

// Config controls one Check() pass.
type Config struct {
	// Service is the canonical service name used to filter seed policies and
	// APISIX entries (e.g. "auth-service", "olt-manager"). Required.
	Service string

	// Router is the chi router whose routes will be walked. Required.
	Router chi.Router

	// SeedPolicies is the slice of policies the service knows about. The
	// guard filters down to Service when comparing. Required.
	SeedPolicies []SeedPolicy

	// APISIXAdminURL is the APISIX admin endpoint (e.g.
	// https://apisix-admin.internal). When set, the guard fetches deployed
	// routes from the live API at boot — the smart path for runtime, since
	// the deployed state IS the source of truth and there is nothing to
	// mount, nothing to bake into the image, no per-environment yaml.
	//
	// Used together with APISIXAdminKey for the X-API-KEY header.
	APISIXAdminURL string

	// APISIXAdminKey is the value of the X-API-KEY header APISIX admin
	// expects. Pulled from a K8s Secret in production; from an env var in
	// dev. Empty string skips the header — only useful when the admin
	// endpoint is firewalled to an internal mesh.
	APISIXAdminKey string

	// APISIXServicePathPrefix narrows the APISIX route set down to the
	// caller's own surface when filtering by label is not possible. e.g.
	// pass "auth/" inside auth-service so the per-service guard does not
	// flag every tenant-service or olt-manager APISIX route as
	// `apisix_without_chi`. Optional — when empty + admin API in use, all
	// returned routes are considered.
	APISIXServicePathPrefix string

	// APISIXRoutesDir is the directory containing the K8S/apisix/routes/*.yaml
	// files. Used by the offline scan (CI lint, dev). Ignored when
	// APISIXAdminURL is set — admin API wins. Optional.
	APISIXRoutesDir string

	// PublicPaths lists chi routes that intentionally bypass Casbin (login,
	// password-reset, internal health). Each entry is "METHOD path" in the
	// canonical RouteKey.String() form (e.g. "POST auth/login"). Default
	// public paths (metrics, healthz) are always included on top of this.
	PublicPaths map[string]struct{}

	// IdentityScopedPaths lists routes that are authenticated but gated by
	// identity (the user's own userAccess slice or own profile), NOT by RBAC.
	// They are excluded from `*_without_policy` drift kinds. Entries follow
	// the same "METHOD path" form as PublicPaths.
	//
	// Example for auth-service:
	//
	//	{
	//	  "GET auth/me":              {},
	//	  "POST auth/tenant/switch":  {},
	//	  "GET auth/users/me/preferences": {},
	//	}
	IdentityScopedPaths map[string]struct{}

	// Metrics is the per-service Prometheus bundle. If nil, NewMetrics(Service)
	// is called on first invocation.
	Metrics *Metrics
}

// Check runs one boot-time guard pass and returns the report.
//
// Failure modes are intentionally soft: the report is always returned, even
// when some sources cannot be loaded — the missing source simply contributes
// no drift entries for that side of the comparison. The caller decides
// whether soft failure (log + metric) or hard failure (panic, refuse to
// serve) is appropriate for its environment.
func Check(cfg Config) (*Report, error) {
	if cfg.Service == "" {
		return nil, errMissingService
	}
	if cfg.Router == nil {
		return nil, errMissingRouter
	}
	if cfg.Metrics == nil {
		cfg.Metrics = NewMetrics(cfg.Service)
	}

	report := &Report{
		Service: cfg.Service,
		Counts:  map[DriftKind]int{},
	}

	chiRoutes, err := WalkChiRoutes(cfg.Router, SkipFunc(cfg.PublicPaths))
	if err != nil {
		return nil, err
	}

	seedRoutes := WalkSeedPolicies(cfg.SeedPolicies, cfg.Service)

	// APISIX side: prefer the admin API (runtime source of truth, env-
	// agnostic, no file mount required) and fall back to the local yaml
	// directory only when the admin URL is unset. The yaml path is mainly
	// for offline scans (CI lint, dev). Both paths feed the same comparator
	// downstream.
	var apisixRoutes []APISIXRoute
	switch {
	case cfg.APISIXAdminURL != "":
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		apisixRoutes, err = FetchAPISIXRoutesFromAdmin(ctx, cfg.APISIXAdminURL, cfg.APISIXAdminKey)
		cancel()
		if err != nil {
			logger.Warn("iamguard: APISIX admin fetch failed, skipping APISIX side of check", err, "admin_url", cfg.APISIXAdminURL)
			apisixRoutes = nil
		}
	case cfg.APISIXRoutesDir != "":
		apisixRoutes, err = WalkAPISIXRoutes(cfg.APISIXRoutesDir)
		if err != nil {
			logger.Warn("iamguard: APISIX walk failed, skipping APISIX side of check", err, "routes_dir", cfg.APISIXRoutesDir)
			apisixRoutes = nil
		}
	default:
		logger.Warn("iamguard: APISIX source not configured, skipping APISIX side of check")
	}

	if len(apisixRoutes) > 0 {
		// Filter by labels.service first; routes deployed before the label
		// rollout fall through to the path-prefix filter so older clusters
		// still scope correctly.
		filtered := FilterAPISIXByService(apisixRoutes, cfg.Service)
		if len(filtered) == 0 && cfg.APISIXServicePathPrefix != "" {
			filtered = FilterAPISIXByPathPrefix(apisixRoutes, cfg.APISIXServicePathPrefix)
		}
		apisixRoutes = filtered
	}

	chiSet := make(map[RouteKey]ChiRoute, len(chiRoutes))
	for _, r := range chiRoutes {
		chiSet[r.Key] = r
	}
	seedSet := make(map[RouteKey]seedRoute, len(seedRoutes))
	for _, p := range seedRoutes {
		seedSet[p.Key] = p
	}
	apisixSet := make(map[RouteKey]APISIXRoute, len(apisixRoutes))
	for _, a := range apisixRoutes {
		apisixSet[a.Key] = a
	}

	isIdentity := func(k RouteKey) bool {
		if cfg.IdentityScopedPaths == nil {
			return false
		}
		_, ok := cfg.IdentityScopedPaths[k.String()]
		return ok
	}

	// chi ↔ seed (skip identity-scoped: handler gates on userAccess, not RBAC).
	for k, r := range chiSet {
		if isIdentity(k) {
			continue
		}
		if _, ok := seedSet[k]; !ok {
			report.add(DriftEntry{Kind: ChiWithoutPolicy, Route: k, Hint: r.HandlerName})
		}
	}
	for k, p := range seedSet {
		if isIdentity(k) {
			continue
		}
		if _, ok := chiSet[k]; !ok {
			report.add(DriftEntry{Kind: PolicyWithoutChi, Route: k, Hint: p.PermissionKey})
		}
	}

	// chi ↔ APISIX (does NOT skip identity-scoped — both routing layers must
	// serve every endpoint regardless of authorization model).
	if len(apisixRoutes) > 0 {
		for k, r := range chiSet {
			if _, ok := apisixSet[k]; !ok {
				report.add(DriftEntry{Kind: ChiWithoutAPISIX, Route: k, Hint: r.HandlerName})
			}
		}
		for k, a := range apisixSet {
			if _, ok := chiSet[k]; !ok {
				report.add(DriftEntry{Kind: APISIXWithoutChi, Route: k, Hint: a.Name})
			}
		}

		for k, a := range apisixSet {
			if a.Public || isIdentity(k) {
				continue
			}
			if _, ok := seedSet[k]; !ok {
				report.add(DriftEntry{Kind: APISIXWithoutPolicy, Route: k, Hint: a.Name})
			}
		}
		for k, p := range seedSet {
			if isIdentity(k) {
				continue
			}
			if _, ok := apisixSet[k]; !ok {
				report.add(DriftEntry{Kind: PolicyWithoutAPISIX, Route: k, Hint: p.PermissionKey})
			}
		}
	}

	sortReport(report)
	emitMetrics(report, cfg.Metrics)
	return report, nil
}

// LogReport emits a structured per-kind summary plus one WARN line per drift
// entry. Output is line-oriented so it is greppable from container logs.
//
// The common logger interprets a leading `nil` as a misaligned key when no
// real error is present, so this routine passes kv pairs directly (no nil
// placeholder) when there is nothing to surface as an error.
func LogReport(rep *Report) {
	if rep == nil {
		return
	}
	if !rep.HasDrift() {
		logger.Info("iamguard: no route drift detected", "service", rep.Service)
		return
	}

	logger.Warn("iamguard: route drift detected",
		"service", rep.Service,
		"total", len(rep.Drift),
		"chi_without_policy", rep.Counts[ChiWithoutPolicy],
		"policy_without_chi", rep.Counts[PolicyWithoutChi],
		"chi_without_apisix", rep.Counts[ChiWithoutAPISIX],
		"apisix_without_chi", rep.Counts[APISIXWithoutChi],
		"apisix_without_policy", rep.Counts[APISIXWithoutPolicy],
		"policy_without_apisix", rep.Counts[PolicyWithoutAPISIX],
	)

	for _, d := range rep.Drift {
		logger.Warn("iamguard: drift entry",
			"kind", string(d.Kind),
			"method", d.Route.Method,
			"path", d.Route.Path,
			"hint", d.Hint,
		)
	}
}

func (r *Report) add(e DriftEntry) {
	r.Drift = append(r.Drift, e)
	r.Counts[e.Kind]++
}

func sortReport(r *Report) {
	sort.Slice(r.Drift, func(i, j int) bool {
		a, b := r.Drift[i], r.Drift[j]
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.Route.Path != b.Route.Path {
			return a.Route.Path < b.Route.Path
		}
		return a.Route.Method < b.Route.Method
	})
}

func emitMetrics(r *Report, m *Metrics) {
	if m == nil {
		return
	}
	kinds := []DriftKind{
		ChiWithoutPolicy, PolicyWithoutChi,
		ChiWithoutAPISIX, APISIXWithoutChi,
		APISIXWithoutPolicy, PolicyWithoutAPISIX,
	}
	for _, k := range kinds {
		m.RouteDriftLastPass.WithLabelValues(string(k)).Set(float64(r.Counts[k]))
	}
	for _, d := range r.Drift {
		m.RouteDriftTotal.WithLabelValues(string(d.Kind)).Inc()
	}
}
