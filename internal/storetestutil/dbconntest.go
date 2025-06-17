package storetestutil

import (
	"context"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Matches docker-compose.yaml
const imageName = "postgres:15.8"

func NewDB() (*pgxpool.Pool, func()) {
	ctx := context.Background()

	container, err := postgres.Run(ctx, imageName,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies())
	if err != nil {
		_ = testcontainers.TerminateContainer(container)
		log.Panicf("run postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Panicf("get connection string: %v", err)
	}

	// Migrate the database schema
	pgx := &pgx.Postgres{}
	db, err := pgx.Open(dsn)
	if err != nil {
		log.Panicf("open pgx connection: %v", err)
	}

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("failed to get current file path")
	}

	migrationsDir := filepath.Join(currentFile, "../../../cmd/openauthctl/migrations")
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "pgx", db)
	if err != nil {
		log.Panicf("create migrate instance: %v", err)
	}
	err = m.Up()
	if err != nil {
		log.Panicf("run migrations: %v", err)
	}
	err = db.Close()
	if err != nil {
		log.Panicf("close pgx connection: %v", err)
	}

	// Create a pgx pool for use in tests
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Panicf("create pgx pool: %v", err)
	}

	return pool, func() {
		pool.Close()
		_ = testcontainers.TerminateContainer(container)
	}
}
