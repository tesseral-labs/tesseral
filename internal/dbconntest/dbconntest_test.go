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

	// Assert the `users` table exists after migrations
	res, err := pool.Query(ctx, "SELECT * FROM users")
	res.Close()
	require.NoError(t, err, "failed to query users table")
}
