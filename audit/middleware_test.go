package audit

import "testing"

func TestExtractResource(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		// Standard plural plain
		{"/api/v1/subscribers", "subscriber"},
		{"/api/v1/subscribers/abc-123", "subscriber"},
		{"/api/v1/plans", "plan"},
		{"/api/v1/tenants/t-001/users", "tenant"},

		// Words that trip up naive TrimSuffix("s")
		{"/api/v1/radius", "radius"},   // not "radiu"
		{"/api/v1/series", "series"},   // not "serie"
		{"/api/v1/status", "status"},   // not "statu"

		// Hyphenated paths
		{"/api/v1/tenant-users", "tenant-user"},
		{"/api/v1/onts", "ont"},
		{"/api/v1/olts/o-001", "olt"},

		// Platform resources (regression fix)
		{"/api/v1/products", "product"},
		{"/api/v1/devices", "device"},
		{"/api/v1/subscriptions", "subscription"},

		// Unknown path — fall back to raw segment, never empty
		{"/api/v1/something-new", "something-new"},

		// Health / metrics — middleware short-circuits before this is called,
		// but defensive default still useful
		{"/", "unknown"},
		{"", "unknown"},
	}

	for _, tc := range cases {
		got := extractResource(tc.path)
		if got != tc.want {
			t.Errorf("extractResource(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}
