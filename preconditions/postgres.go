package preconditions

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// PostgresConnect returns a check that pings the database. Catches the
// most common boot failure: wrong host, wrong creds, network unreachable.
func PostgresConnect(name string, db *sql.DB) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     "verify POSTGRES_*_HOST/USERNAME/PASSWORD env vars and the cluster is reachable from this pod",
		Check: func(ctx context.Context) error {
			if db == nil {
				return fmt.Errorf("nil *sql.DB")
			}
			return db.PingContext(ctx)
		},
	}
}

// PostgresHasTables returns a check that every named table exists in
// the public schema. Useful when migrations are run out-of-band (init
// container, kubectl exec) — if the operator forgot to apply the
// latest migration, this fails the boot with the exact missing table
// list rather than blowing up at first query.
func PostgresHasTables(name string, db *sql.DB, tables ...string) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     "run pending migrations (kubectl exec deploy/<svc> -- /app/migrate up, or rebuild image to trigger init container)",
		Check: func(ctx context.Context) error {
			if db == nil {
				return fmt.Errorf("nil *sql.DB")
			}
			var missing []string
			for _, t := range tables {
				var exists bool
				if err := db.QueryRowContext(ctx,
					"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name=$1)",
					t).Scan(&exists); err != nil {
					return fmt.Errorf("query for table %q: %w", t, err)
				}
				if !exists {
					missing = append(missing, t)
				}
			}
			if len(missing) > 0 {
				return fmt.Errorf("tables missing in public schema: %s", strings.Join(missing, ", "))
			}
			return nil
		},
	}
}

// PostgresReplicationGrant returns a check that the connected user has
// the REPLICATION attribute. Without it, CreateReplicationSlot fails
// with a confusing pgproto error mid-stream — this catches it at boot
// instead.
func PostgresReplicationGrant(name string, db *sql.DB, expectedUser string) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     fmt.Sprintf("ALTER USER %s WITH REPLICATION; (run as postgres superuser)", expectedUser),
		Check: func(ctx context.Context) error {
			if db == nil {
				return fmt.Errorf("nil *sql.DB")
			}
			var hasReplication bool
			if err := db.QueryRowContext(ctx,
				"SELECT rolreplication FROM pg_roles WHERE rolname=$1",
				expectedUser).Scan(&hasReplication); err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("CDC user %q does not exist on the cluster", expectedUser)
				}
				return fmt.Errorf("query pg_roles: %w", err)
			}
			if !hasReplication {
				return fmt.Errorf("user %q lacks REPLICATION attribute", expectedUser)
			}
			return nil
		},
	}
}

// PostgresWALLevelLogical returns a check that wal_level is set to
// "logical" — required for CDC via pg_logical_replication. wal_level
// is a postmaster-level setting requiring a server restart to change,
// so it's an operator/cluster-config issue, not a runtime one.
func PostgresWALLevelLogical(name string, db *sql.DB) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     "set wal_level=logical in postgresql.conf and restart the cluster (CNPG: spec.postgresql.parameters.wal_level: logical)",
		Check: func(ctx context.Context) error {
			if db == nil {
				return fmt.Errorf("nil *sql.DB")
			}
			var level string
			if err := db.QueryRowContext(ctx, "SHOW wal_level").Scan(&level); err != nil {
				return fmt.Errorf("SHOW wal_level: %w", err)
			}
			if level != "logical" {
				return fmt.Errorf("wal_level is %q, need 'logical' for CDC", level)
			}
			return nil
		},
	}
}

// PostgresPublicationExists returns a check that a logical-replication
// publication exists. Auto-creating publications at runtime is fine in
// dev but masks "wrong publication name in env" mistakes in prod.
func PostgresPublicationExists(name string, db *sql.DB, publicationName string) Check {
	return Check{
		Name:     name,
		Required: true,
		Hint:     fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE <tables>; (or check the publication name env var matches)", publicationName),
		Check: func(ctx context.Context) error {
			if db == nil {
				return fmt.Errorf("nil *sql.DB")
			}
			var exists bool
			if err := db.QueryRowContext(ctx,
				"SELECT EXISTS(SELECT 1 FROM pg_publication WHERE pubname=$1)",
				publicationName).Scan(&exists); err != nil {
				return fmt.Errorf("query pg_publication: %w", err)
			}
			if !exists {
				return fmt.Errorf("publication %q does not exist", publicationName)
			}
			return nil
		},
	}
}
