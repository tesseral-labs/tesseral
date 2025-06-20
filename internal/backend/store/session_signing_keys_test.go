package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSessionPublicKeysByProjectID(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	keys, err := u.Store.GetSessionPublicKeysByProjectID(ctx, u.ProjectID)
	require.NoError(t, err)
	require.NotEmpty(t, keys)
}
