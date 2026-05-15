// Package iamguard enforces the IAM new-route checklist at service boot.
//
// Every protected route registered by a Go service MUST have:
//  1. a chi handler (so the request can reach the service),
//  2. a seed Casbin policy (so RBAC roles can grant it), and
//  3. an APISIX route entry (so the gateway forwards traffic and runs forward-auth).
//
// The guard walks each of those three sources, normalises path templates to
// a common key, and reports drift in either direction. It does NOT panic —
// the service must still boot in production even when seed lags reality —
// but it logs every drift entry at WARN and increments a Prometheus counter
// so the drift is visible in dashboards and alertable.
//
// Each service consumes this package by calling iamguard.Check at boot with
// the chi router, its seed-policy slice (converted via iamguard.SeedPolicy),
// and the directory holding K8S/apisix/routes yamls. See the auth-service
// startup wiring in `internal/app/app.start.go` for a reference integration.
//
// Public + identity-scoped path lists are intentionally configurable per
// service: every service has different public flows (login vs. provision
// vs. captive-portal) and different self-service identity-scoped endpoints
// (`/auth/me`, `/tenant/switch`, etc.). The defaults exposed by this
// package cover the common cases; services pass extra entries via
// Config.PublicPaths and Config.IdentityScopedPaths.
package iamguard
