package storetesting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TesNewDB(t *testing.T) {
	t.Parallel()

	pool, cleanup := newDB()
	t.Cleanup(cleanup)

	ctx := context.Background()
	require.NoError(t, pool.Ping(ctx), "failed to ping the database")

	// Assert the `users` table exists after migrations
	res, err := pool.Query(ctx, "SELECT * FROM users")
	res.Close()
	require.NoError(t, err, "failed to query users table")
}
