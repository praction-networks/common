package iamguard

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds the per-service Prometheus collectors. Each service should
// instantiate one set via NewMetrics(service) at package init or first boot.
// Collectors are kept in an internal cache so repeated NewMetrics calls with
// the same service name return the same instance — avoids the
// "duplicate metrics collector registration" panic when boot retries.
type Metrics struct {
	RouteDriftTotal    *prometheus.CounterVec
	RouteDriftLastPass *prometheus.GaugeVec
}

var (
	metricsCache   = map[string]*Metrics{}
	metricsCacheMu sync.Mutex
)

// NewMetrics returns the Metrics bundle for a service, creating + registering
// it on first call and returning the cached instance on subsequent calls.
//
// The metric names are prefixed with the service name to keep dashboards
// scoped: e.g. `auth_iam_route_drift_total{kind=...}` for auth-service,
// `tenant_iam_route_drift_total{kind=...}` for tenant-service. The service
// suffix `-service` is stripped to keep names short.
func NewMetrics(service string) *Metrics {
	prefix := metricPrefix(service)

	metricsCacheMu.Lock()
	defer metricsCacheMu.Unlock()
	if m, ok := metricsCache[prefix]; ok {
		return m
	}

	m := &Metrics{
		RouteDriftTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_iam_route_drift_total",
			Help: "IAM guard route-drift events detected at boot, by drift kind.",
		}, []string{"kind"}),
		RouteDriftLastPass: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: prefix + "_iam_route_drift_last_pass",
			Help: "IAM guard route-drift count from the most recent boot-time pass, by drift kind.",
		}, []string{"kind"}),
	}
	metricsCache[prefix] = m
	return m
}

// metricPrefix converts "auth-service" -> "auth", "olt-manager" -> "olt_manager".
// Prometheus name regex disallows '-' so we substitute '_'.
func metricPrefix(service string) string {
	out := make([]rune, 0, len(service))
	skipSuffix := "-service"
	if len(service) > len(skipSuffix) && service[len(service)-len(skipSuffix):] == skipSuffix {
		service = service[:len(service)-len(skipSuffix)]
	}
	for _, r := range service {
		if r == '-' {
			out = append(out, '_')
			continue
		}
		out = append(out, r)
	}
	if len(out) == 0 {
		return "service"
	}
	return string(out)
}
