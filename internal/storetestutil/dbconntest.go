package storetestutil

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Matches docker-compose.yaml
const imageName = "postgres:15.8"

func NewDB(t *testing.T) *pgxpool.Pool {
	container, err := postgres.Run(t.Context(), imageName,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies())
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("run postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	// Migrate the database schema
	pgx := &pgx.Postgres{}
	db, err := pgx.Open(dsn)
	if err != nil {
		t.Fatalf("open pgx connection: %v", err)
	}

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current file path")
	}

	migrationsDir := filepath.Join(currentFile, "../../../cmd/openauthctl/migrations")
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "pgx", db)
	if err != nil {
		t.Fatalf("create migrate instance: %v", err)
	}
	err = m.Up()
	if err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	err = db.Close()
	if err != nil {
		t.Fatalf("close pgx connection: %v", err)
	}

	// Create a pgx pool for use in tests
	pool, err := pgxpool.New(t.Context(), dsn)
	if err != nil {
		t.Fatalf("create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
