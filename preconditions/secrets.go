package preconditions

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// EnvVarsSet returns a check that every named env var is non-empty.
// Catches missing secret references in deployment manifests — without
// this, services often boot with zero-value JWT keys / empty PG
// passwords and only fail at first use.
func EnvVarsSet(name string, vars ...string) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     "verify the deployment manifest sets these via secretKeyRef or value",
		Check: func(_ context.Context) error {
			var missing []string
			for _, v := range vars {
				if strings.TrimSpace(os.Getenv(v)) == "" {
					missing = append(missing, v)
				}
			}
			if len(missing) > 0 {
				return fmt.Errorf("env vars unset/empty: %s", strings.Join(missing, ", "))
			}
			return nil
		},
	}
}

// NonEmptyString returns a check that a config-derived string is
// non-empty. Use for values that load from yaml/configmap rather than
// raw env vars (e.g. cfg.JWTEnv.PublicKey).
func NonEmptyString(name string, value string, hint string) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     hint,
		Check: func(_ context.Context) error {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("required value is empty")
			}
			return nil
		},
	}
}
