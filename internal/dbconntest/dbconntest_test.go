package dbconntest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	pool := Open(t)
	require.NoError(t, pool.Ping(context.Background()), "failed to ping the database")
}
