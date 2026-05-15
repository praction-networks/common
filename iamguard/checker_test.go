package iamguard

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/common/logger"
)

// TestMain initialises the shared logger so LogReport calls inside the test
// suite do not panic with "Logger not initialized". The logger is a process-
// wide singleton; once initialised here every subsequent test in this binary
// shares the same configuration.
func TestMain(m *testing.M) {
	_ = logger.InitializeLogger(logger.LoggerConfig{LogLevel: "info"})
	os.Exit(m.Run())
}

// TestNormalisePath exercises the path normalisation that joins chi / APISIX
// / seed templates onto a single canonical key.
func TestNormalisePath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		// chi templates
		{"/api/v1/auth/users/{id}", "auth/users/{*}"},
		{"/api/v1/auth/users/{id}/biometric/{factorId}", "auth/users/{*}/biometric/{*}"},

		// APISIX wildcards
		{"/api/v1/auth/users/*", "auth/users/{*}"},
		{"/api/v1/auth/users/*/biometric/*", "auth/users/{*}/biometric/{*}"},

		// seed Resource strings (no /api/v1 prefix)
		{"auth/users/{id}", "auth/users/{*}"},
		{"tenant/{tenantId}/bindings/sms", "tenant/{*}/bindings/sms"},

		// mixed casing + extra slashes
		{"/api/v1/Auth/Me/", "auth/me"},
		{"/api/v2/foo", "foo"},
	}
	for _, tc := range cases {
		got := NormalisePath(tc.in)
		if got != tc.want {
			t.Errorf("NormalisePath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCheckDetectsAllSixDriftKinds wires up synthetic chi routes, seed
// policies, and APISIX yamls so that every DriftKind has exactly one entry
// in the resulting report. Verifies the full Check() data flow end-to-end.
func TestCheckDetectsAllSixDriftKinds(t *testing.T) {
	// 1. APISIX yaml fixture — written to a temp dir.
	tmp := t.TempDir()
	yaml := `service: testsvc
routes:
  # Both chi + seed match this one — should NOT appear as drift.
  - name: aligned-get
    uri: /api/v1/svc/users
    methods: [GET]
    public: false

  # chi missing → APISIXWithoutChi (and APISIXWithoutPolicy when no seed).
  - name: only-in-apisix
    uri: /api/v1/svc/widgets
    methods: [GET]
    public: false

  # public flag set — should NOT generate APISIXWithoutPolicy even with no seed.
  - name: public-login
    uri: /api/v1/svc/login
    methods: [POST]
    public: true
`
	if err := os.WriteFile(filepath.Join(tmp, "testsvc.yaml"), []byte(yaml), 0o600); err != nil {
		t.Fatalf("write yaml fixture: %v", err)
	}

	// 2. Chi router fixture.
	r := chi.NewRouter()
	noop := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	// aligned: matches both seed + APISIX.
	r.Get("/api/v1/svc/users", noop)
	// chi-only: drift = ChiWithoutPolicy + ChiWithoutAPISIX.
	r.Get("/api/v1/svc/orphan", noop)

	// 3. Seed policies fixture.
	seed := []SeedPolicy{
		{Service: "testsvc", Resource: "svc/users", Action: "GET", PermissionKey: "svc.users.view"},
		// policy-only: drift = PolicyWithoutChi + PolicyWithoutAPISIX.
		{Service: "testsvc", Resource: "svc/ghost", Action: "GET", PermissionKey: "svc.ghost.view"},
	}

	// 4. Run the guard.
	report, err := Check(Config{
		Service:         "testsvc",
		Router:          r,
		SeedPolicies:    seed,
		APISIXRoutesDir: tmp,
		// no public or identity-scoped overrides
	})
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	// 5. Expectations per kind.
	// APISIXWithoutChi counts BOTH `svc/widgets` (no chi) AND `svc/login`
	// (public, but every APISIX route still needs a chi upstream). The
	// public flag only suppresses APISIXWithoutPolicy, not the chi check.
	want := map[DriftKind]int{
		ChiWithoutPolicy:    1, // svc/orphan
		PolicyWithoutChi:    1, // svc/ghost
		ChiWithoutAPISIX:    1, // svc/orphan
		APISIXWithoutChi:    2, // svc/widgets + svc/login
		APISIXWithoutPolicy: 1, // svc/widgets (login is public, excluded here)
		PolicyWithoutAPISIX: 1, // svc/ghost
	}
	for kind, n := range want {
		if got := report.Counts[kind]; got != n {
			t.Errorf("Counts[%s] = %d, want %d", kind, got, n)
		}
	}

	totalWant := 0
	for _, n := range want {
		totalWant += n
	}
	if got := len(report.Drift); got != totalWant {
		t.Errorf("total drift entries = %d, want %d; entries: %+v", got, totalWant, report.Drift)
	}

	// 6. Emit logs through the real LogReport path so the test output shows
	// what an operator would see at boot. Captured by `go test -v`.
	LogReport(report)
}

// TestIdentityScopedExemption verifies that paths listed in
// Config.IdentityScopedPaths are excluded from `*_without_policy` checks but
// still subject to chi↔APISIX comparison.
func TestIdentityScopedExemption(t *testing.T) {
	tmp := t.TempDir()
	yaml := `service: testsvc
routes:
  - name: me-bootstrap
    uri: /api/v1/svc/me
    methods: [GET]
    public: false
`
	if err := os.WriteFile(filepath.Join(tmp, "testsvc.yaml"), []byte(yaml), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	r := chi.NewRouter()
	noop := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	r.Get("/api/v1/svc/me", noop)

	report, err := Check(Config{
		Service:      "testsvc",
		Router:       r,
		SeedPolicies: nil, // intentionally no policy for svc/me
		APISIXRoutesDir: tmp,
		IdentityScopedPaths: map[string]struct{}{
			"GET svc/me": {},
		},
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	// /svc/me has no seed policy, but identity-scoped exemption should keep
	// it out of ChiWithoutPolicy AND APISIXWithoutPolicy.
	if report.Counts[ChiWithoutPolicy] != 0 {
		t.Errorf("ChiWithoutPolicy should be 0 for identity-scoped path, got %d", report.Counts[ChiWithoutPolicy])
	}
	if report.Counts[APISIXWithoutPolicy] != 0 {
		t.Errorf("APISIXWithoutPolicy should be 0 for identity-scoped path, got %d", report.Counts[APISIXWithoutPolicy])
	}
	if report.HasDrift() {
		t.Errorf("expected zero drift for fully-aligned identity-scoped route, got %+v", report.Drift)
	}

	LogReport(report)
}

// TestPublicPathExemption verifies that PublicPaths drops chi routes before
// they are even considered by the drift logic.
func TestPublicPathExemption(t *testing.T) {
	r := chi.NewRouter()
	noop := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	r.Post("/api/v1/svc/login", noop)

	report, err := Check(Config{
		Service:      "testsvc",
		Router:       r,
		SeedPolicies: nil,
		// no APISIX dir → only chi↔seed comparison
		PublicPaths: map[string]struct{}{
			"POST svc/login": {},
		},
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if report.HasDrift() {
		t.Errorf("public path should not surface drift, got %+v", report.Drift)
	}
}

// TestCheckMissingConfig surfaces validation errors so misconfigured boots
// fail loudly rather than silently skipping the guard.
func TestCheckMissingConfig(t *testing.T) {
	if _, err := Check(Config{}); err == nil {
		t.Error("Check with empty Config should error on missing Service")
	}
	if _, err := Check(Config{Service: "x"}); err == nil {
		t.Error("Check without Router should error")
	}
}
