package dbconn

import (
	"context"
	"fmt"
	"strings"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN string `conf:"dsn"`
}

func Open(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName(), otelpgx.WithSpanNameFunc(func(stmt string) string {
		// If stmt is of the sqlc form "-- name: Example :one\n...", extract
		// "Example". Otherwise, leave as-is.
		stmt = strings.TrimPrefix(stmt, "-- name: ")
		if i := strings.IndexByte(stmt, ' '); i != -1 {
			return stmt[:i]
		}
		return stmt
	}))

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool from dsn: %w", err)
	}
	return pool, nil
}
