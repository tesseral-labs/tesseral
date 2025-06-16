package dbconntest

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Matches docker-compose.yaml
const imageName = "postgres:15.8"

func Open(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	container, err := postgres.Run(ctx, imageName,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies())
	if err != nil {
		t.Fatalf("run postgres container: %v", err)
	}
	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
