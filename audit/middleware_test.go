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

func TestExtractResourceAndID(t *testing.T) {
	cases := []struct {
		path       string
		wantRes    string
		wantID     string
	}{
		// cuid2 (32 chars, lowercase alphanumeric)
		{"/api/v1/olts/n4ybdyowbpo5nkpjzcf3n7vf79w5z1nz", "olt", "n4ybdyowbpo5nkpjzcf3n7vf79w5z1nz"},
		{"/api/v1/tenant-users/abc123def456ghi789jkl012mno34567", "tenant-user", "abc123def456ghi789jkl012mno34567"},

		// UUID
		{"/api/v1/audit-logs/f483bddf-3221-4e98-8507-980f2016fc75", "audit-log", "f483bddf-3221-4e98-8507-980f2016fc75"},

		// Numeric
		{"/api/v1/invoices/42", "invoice", "42"},

		// Nested: resource then id then sub-resource — we keep top-level resource + id
		{"/api/v1/olts/n4ybdyowbpo5nkpjzcf3n7vf79w5z1nz/onts", "olt", "n4ybdyowbpo5nkpjzcf3n7vf79w5z1nz"},
		{"/api/v1/tenants/ryyggweeb13wccx45uaopwvt25g0oq7e/verify", "tenant", "ryyggweeb13wccx45uaopwvt25g0oq7e"},

		// No ID — bare collection
		{"/api/v1/olts", "olt", ""},
		{"/api/v1/subscribers", "subscriber", ""},

		// Action verb sitting where ID would be — must not be misread as ID
		{"/api/v1/olts/lookup", "olt", ""},
		{"/api/v1/tenant-users/lookup", "tenant-user", ""},
	}

	for _, tc := range cases {
		gotRes, gotID := extractResourceAndID(tc.path)
		if gotRes != tc.wantRes || gotID != tc.wantID {
			t.Errorf("extractResourceAndID(%q) = (%q, %q), want (%q, %q)",
				tc.path, gotRes, gotID, tc.wantRes, tc.wantID)
		}
	}
}
