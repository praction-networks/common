package iamguard

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ChiRoute is the in-memory record the guard keeps for each registered chi
// endpoint after normalisation.
type ChiRoute struct {
	Key         RouteKey
	HandlerName string
}

// WalkChiRoutes traverses every route registered on the given chi.Router and
// returns one ChiRoute per (method, path), normalised via NormalisePath. The
// optional skip predicate filters out routes that intentionally bypass Casbin
// (login, /metrics, etc.); pass nil to keep every route.
func WalkChiRoutes(r chi.Router, skip func(RouteKey) bool) ([]ChiRoute, error) {
	out := make([]ChiRoute, 0, 256)

	err := chi.Walk(r, func(method, route string, handler http.Handler, _ ...func(http.Handler) http.Handler) error {
		key := RouteKey{
			Path:   NormalisePath(route),
			Method: NormaliseMethod(method),
		}
		if skip != nil && skip(key) {
			return nil
		}
		out = append(out, ChiRoute{
			Key:         key,
			HandlerName: handlerName(handler),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// handlerName best-effort extracts the registered function name so a drift
// entry can point the operator at the chi handler that has no policy.
// Returns "" when the handler is not a recognisable function value (e.g.
// wrapped middleware chains).
func handlerName(h http.Handler) string {
	if h == nil {
		return ""
	}
	v := reflect.ValueOf(h)
	if !v.IsValid() || v.Kind() != reflect.Func {
		return ""
	}
	fn := runtime.FuncForPC(v.Pointer())
	if fn == nil {
		return ""
	}
	return shortFuncName(fn.Name())
}

// shortFuncName trims package paths to a readable receiver.method form.
func shortFuncName(full string) string {
	if i := strings.LastIndex(full, "/"); i >= 0 {
		full = full[i+1:]
	}
	full = strings.TrimSuffix(full, "-fm")
	return full
}

// DefaultPublicPaths returns the keys of routes every service should treat as
// public regardless of which service it serves. Services add their own public
// flows on top (login, OTP, password-reset, etc.) via Config.PublicPaths.
//
// Each entry is "METHOD path" using the canonical RouteKey form.
func DefaultPublicPaths() map[string]struct{} {
	return map[string]struct{}{
		"GET metrics": {},
		"GET healthz": {},
		"GET health":  {},
	}
}

// SkipFunc builds a chi.Walk skip predicate from the union of the defaults
// and the caller's extra entries. Pass nil for extra to use defaults only.
func SkipFunc(extra map[string]struct{}) func(RouteKey) bool {
	merged := DefaultPublicPaths()
	for k := range extra {
		merged[k] = struct{}{}
	}
	return func(k RouteKey) bool {
		_, ok := merged[k.String()]
		return ok
	}
}
