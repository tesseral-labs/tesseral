package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestAuthenticateBackendAPIKey_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "test-key",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp.BackendApiKey)

	authResp, err := u.Store.AuthenticateBackendAPIKey(ctx, createResp.BackendApiKey.SecretToken)
	require.NoError(t, err)
	require.NotNil(t, authResp)
	require.Equal(t, createResp.BackendApiKey.Id, authResp.BackendAPIKeyID)
	require.Equal(t, u.ProjectID, authResp.ProjectID)
}

func TestAuthenticateBackendAPIKey_InvalidToken(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.AuthenticateBackendAPIKey(ctx, idformat.BackendAPIKey.Format(uuid.New()))
	require.Error(t, err)
}
