package iamguard

// WalkSeedPolicies normalises a slice of SeedPolicy entries into RouteKey-
// comparable records, filtered to the policies that target servicePrefix
// (pass "" to keep every policy).
//
// Policies with empty Action or Action == "*" are skipped: wildcard entries
// represent role-level catch-alls (e.g. SuperAdmin's "manage everything"
// grant) and do not correspond to a single HTTP verb.
func WalkSeedPolicies(policies []SeedPolicy, servicePrefix string) []seedRoute {
	out := make([]seedRoute, 0, len(policies))
	for _, p := range policies {
		if servicePrefix != "" && p.Service != servicePrefix {
			continue
		}
		if p.Action == "" || p.Action == "*" {
			continue
		}
		out = append(out, seedRoute{
			Key: RouteKey{
				Path:   NormalisePath(p.Resource),
				Method: NormaliseMethod(p.Action),
			},
			PermissionKey: p.PermissionKey,
			Service:       p.Service,
			Name:          p.Name,
		})
	}
	return out
}

// seedRoute is the post-normalisation record used internally by the checker.
type seedRoute struct {
	Key           RouteKey
	PermissionKey string
	Service       string
	Name          string
}
