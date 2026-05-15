package iamguard

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// apisixRouteFile mirrors the structure of every yaml under K8S/apisix/routes/.
// Only the fields the guard needs are captured so minor schema additions
// (rate_profile, plugins, priority, etc.) do not break parsing.
type apisixRouteFile struct {
	Service string        `yaml:"service"`
	Routes  []apisixEntry `yaml:"routes"`
}

type apisixEntry struct {
	Name    string   `yaml:"name"`
	URI     string   `yaml:"uri"`
	URIs    []string `yaml:"uris"`
	Methods []string `yaml:"methods"`
	Public  bool     `yaml:"public"`
}

// APISIXRoute is the in-memory record kept for each parsed APISIX entry.
type APISIXRoute struct {
	Key     RouteKey
	Service string
	Name    string
	Public  bool
}

// WalkAPISIXRoutes parses every *.yaml under routesDir and returns one
// APISIXRoute per (uri, method) pair. routesDir is typically the path to
// K8S/apisix/routes, supplied at runtime via the APISIX_ROUTES_DIR env var.
//
// If routesDir is empty or unreadable a descriptive error is returned so the
// caller can decide whether to fail or skip — the boot-time integration in
// each service logs a warning and runs the chi↔seed comparison without
// APISIX coverage.
func WalkAPISIXRoutes(routesDir string) ([]APISIXRoute, error) {
	if routesDir == "" {
		return nil, fmt.Errorf("apisix routes dir not configured (set APISIX_ROUTES_DIR)")
	}

	info, err := os.Stat(routesDir)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", routesDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", routesDir)
	}

	out := make([]APISIXRoute, 0, 256)

	err = filepath.WalkDir(routesDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read %s: %w", path, readErr)
		}
		var parsed apisixRouteFile
		if yamlErr := yaml.Unmarshal(raw, &parsed); yamlErr != nil {
			return fmt.Errorf("parse %s: %w", path, yamlErr)
		}
		for _, entry := range parsed.Routes {
			uris := entry.URIs
			if entry.URI != "" {
				uris = append(uris, entry.URI)
			}
			for _, uri := range uris {
				methods := entry.Methods
				if len(methods) == 0 {
					methods = []string{"GET"}
				}
				for _, m := range methods {
					out = append(out, APISIXRoute{
						Key: RouteKey{
							Path:   NormalisePath(uri),
							Method: NormaliseMethod(m),
						},
						Service: parsed.Service,
						Name:    entry.Name,
						Public:  entry.Public,
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterAPISIXByService returns only the entries whose service matches the
// given prefix. Pass "" to keep all entries.
//
// "Service" is populated by the yaml walker (top-level `service:` key) and
// by the admin walker when the APISIX route carries `labels.service`. When
// neither source supplies a value (older routes deployed before the label
// was introduced), use FilterAPISIXByPathPrefix instead so the per-service
// guard can still scope its check.
func FilterAPISIXByService(routes []APISIXRoute, service string) []APISIXRoute {
	if service == "" {
		return routes
	}
	out := make([]APISIXRoute, 0, len(routes))
	for _, r := range routes {
		if r.Service == service {
			out = append(out, r)
		}
	}
	return out
}

// FilterAPISIXByPathPrefix narrows the APISIX route list to entries whose
// normalised path starts with the given prefix (e.g. "auth/" for
// auth-service, "olt/" for olt-manager). Use this when running the guard
// against the APISIX admin API at runtime — APISIX returns routes for every
// service in the cluster, and the per-service guard only wants its own
// surface. Pass "" to keep all entries.
//
// Prefix matching is on the already-normalised path (no /api/v1 prefix,
// trimmed slashes, lowercased). Pass the same form your seed Resources use.
func FilterAPISIXByPathPrefix(routes []APISIXRoute, prefix string) []APISIXRoute {
	if prefix == "" {
		return routes
	}
	prefix = strings.ToLower(strings.Trim(prefix, "/"))
	out := make([]APISIXRoute, 0, len(routes))
	for _, r := range routes {
		if strings.HasPrefix(r.Key.Path, prefix) {
			out = append(out, r)
		}
	}
	return out
}
