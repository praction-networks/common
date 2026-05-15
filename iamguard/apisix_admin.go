package iamguard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FetchAPISIXRoutesFromAdmin queries the APISIX admin API and returns one
// APISIXRoute per (uri, method) pair.
//
// Designed for runtime use inside a service pod where mounting the source
// `K8S/apisix/routes/*.yaml` is undesirable (per-environment overrides,
// customer-hosted deployments that should not see the platform team's raw
// yamls). The admin API is the deployed source of truth — what auth-service
// learns here is what users actually hit at the gateway, not what the
// monorepo claims.
//
// adminURL is the APISIX admin endpoint, e.g. https://apisix-admin.i9network.com.
// apiKey is the value to send as the `X-API-KEY` header (required by APISIX
// admin auth). Use the empty string to skip the header (only useful in tests).
//
// The function is intentionally tolerant of unexpected JSON fields — APISIX
// adds keys across versions and we only need uri/methods/name to compute
// drift. Missing or unparseable entries are skipped with no error so a
// single rogue route does not blank out the rest of the dataset.
func FetchAPISIXRoutesFromAdmin(ctx context.Context, adminURL, apiKey string) ([]APISIXRoute, error) {
	if adminURL == "" {
		return nil, fmt.Errorf("apisix admin URL is required")
	}
	base, err := url.Parse(adminURL)
	if err != nil {
		return nil, fmt.Errorf("parse admin URL %q: %w", adminURL, err)
	}
	base.Path = strings.TrimSuffix(base.Path, "/") + "/apisix/admin/routes"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if apiKey != "" {
		req.Header.Set("X-API-KEY", apiKey)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call %s: %w", base.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("apisix admin returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return parseAdminRouteList(body)
}

// apisixAdminListEnvelope is the APISIX admin GET /routes response shape.
// APISIX wraps each route in `{value: {...}}`. We only decode the bits the
// guard uses — unknown fields are ignored.
type apisixAdminListEnvelope struct {
	List []struct {
		Value apisixAdminRouteValue `json:"value"`
	} `json:"list"`
	// Some APISIX versions return the list directly under "node.nodes" (etcd
	// v2 compatibility shape). parseAdminRouteList handles both.
	Node struct {
		Nodes []struct {
			Value json.RawMessage `json:"value"`
		} `json:"nodes"`
	} `json:"node"`
}

type apisixAdminRouteValue struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	URI      string                 `json:"uri"`
	URIs     []string               `json:"uris"`
	Methods  []string               `json:"methods"`
	Status   int                    `json:"status"`
	Labels   map[string]string      `json:"labels"`
	Plugins  map[string]interface{} `json:"plugins"`
	Upstream map[string]interface{} `json:"upstream"`
}

// parseAdminRouteList accepts both the v3-style {list:[{value:...}]} and the
// legacy {node:{nodes:[{value:"<stringified-json>"}]}} envelope APISIX ships
// depending on storage backend. Returns one APISIXRoute per (uri, method).
func parseAdminRouteList(body []byte) ([]APISIXRoute, error) {
	var env apisixAdminListEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("decode admin response: %w", err)
	}

	values := make([]apisixAdminRouteValue, 0, len(env.List)+len(env.Node.Nodes))
	for _, item := range env.List {
		values = append(values, item.Value)
	}
	for _, n := range env.Node.Nodes {
		var v apisixAdminRouteValue
		// Older APISIX wraps the value as a JSON-encoded string; try both.
		if err := json.Unmarshal(n.Value, &v); err == nil {
			values = append(values, v)
			continue
		}
		var raw string
		if err := json.Unmarshal(n.Value, &raw); err == nil {
			if err := json.Unmarshal([]byte(raw), &v); err == nil {
				values = append(values, v)
			}
		}
	}

	out := make([]APISIXRoute, 0, len(values))
	for _, v := range values {
		// Public = no auth gating, which APISIX expresses via labels or by
		// the absence of a forward-auth plugin. We use the same heuristic as
		// the yaml walker: a label `public=true` opts the route out of the
		// `*_without_policy` drift checks.
		public := v.Labels != nil && (strings.EqualFold(v.Labels["public"], "true") || strings.EqualFold(v.Labels["auth"], "public"))

		// Map service from the labels block — apply-routes.py stamps every
		// route with `labels.service=<svc-name>` so we can filter the same
		// way the yaml walker does.
		service := ""
		if v.Labels != nil {
			service = v.Labels["service"]
		}

		uris := v.URIs
		if v.URI != "" {
			uris = append(uris, v.URI)
		}
		methods := v.Methods
		if len(methods) == 0 {
			methods = []string{"GET"}
		}
		for _, uri := range uris {
			for _, m := range methods {
				out = append(out, APISIXRoute{
					Key: RouteKey{
						Path:   NormalisePath(uri),
						Method: NormaliseMethod(m),
					},
					Service: service,
					Name:    v.Name,
					Public:  public,
				})
			}
		}
	}

	return out, nil
}
