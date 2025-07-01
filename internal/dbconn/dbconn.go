package dbconn

import (
	"context"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN string    `conf:"dsn"`
	IAM IAMConfig `conf:"iam,noredact"`
}

type IAMConfig struct {
	Region string `conf:"region,noredact"`
	Host   string `conf:"host,noredact"`
	Port   int    `conf:"port,noredact"`
	User   string `conf:"user,noredact"`
	DBName string `conf:"dbname,noredact"`
}

func Open(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	poolConfig, err := newPoolConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("new pool config: %w", err)
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

func newPoolConfig(ctx context.Context, config Config) (*pgxpool.Config, error) {
	if config.DSN != "" {
		poolConfig, err := pgxpool.ParseConfig(config.DSN)
		if err != nil {
			return nil, fmt.Errorf("parse dsn: %w", err)
		}
		return poolConfig, nil
	}

	password, err := iamDBPassword(ctx, config.IAM)
	if err != nil {
		return nil, fmt.Errorf("get iam db password: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", config.IAM.Host, config.IAM.Port, config.IAM.User, password, config.IAM.DBName)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	poolConfig.BeforeConnect = func(ctx context.Context, connConfig *pgx.ConnConfig) error {
		password, err := iamDBPassword(ctx, config.IAM)
		if err != nil {
			return fmt.Errorf("get iam db password: %w", err)
		}
		connConfig.Password = password
		return nil
	}

	return poolConfig, nil
}

func iamDBPassword(ctx context.Context, c IAMConfig) (string, error) {
	awscfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("load default aws config: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%d", c.Host, c.Port)
	authToken, err := auth.BuildAuthToken(ctx, endpoint, c.Region, c.User, awscfg.Credentials)
	if err != nil {
		return "", fmt.Errorf("build rds auth token: %w", err)
	}

	return authToken, nil
}
