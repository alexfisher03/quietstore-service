package db

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migFS embed.FS

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := migFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			return err
		}

		sql := string(b)
		for _, s := range strings.Split(sql, ";") {
			q := strings.TrimSpace(s)
			if q == "" {
				continue
			}
			if _, err := pool.Exec(ctx, q); err != nil {
				return fmt.Errorf("migration %s failed: %w", e.Name(), err)
			}
		}
	}
	return nil
}
