package dbconntest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	pool := Open(t)

	ctx := context.Background()
	require.NoError(t, pool.Ping(ctx), "failed to ping the database")

	_, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, name TEXT)")
	require.NoError(t, err, "failed to create test table")

	_, err = pool.Exec(ctx, "INSERT INTO test_table (name) VALUES ('test')")
	require.NoError(t, err, "failed to insert into test table")
}
