package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateBackendAPIKey_Entitled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	res, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.BackendApiKey)
	require.Equal(t, "key1", res.BackendApiKey.DisplayName)
	require.NotEmpty(t, res.BackendApiKey.Id)
	require.NotEmpty(t, res.BackendApiKey.CreateTime)
	require.NotEmpty(t, res.BackendApiKey.UpdateTime)
	require.NotEmpty(t, res.BackendApiKey.SecretToken)
}

func TestCreateBackendAPIKey_NotEntitled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	_, err = u.Environment.DB.Exec(ctx, `
	UPDATE projects
	SET entitled_backend_api_keys = false
	WHERE id = $1;
	`,
		uuid.UUID(projectID).String(),
	)
	require.NoError(t, err)

	_, err = u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key1",
		},
	})
	require.Error(t, err)
}

func TestGetBackendAPIKey_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.BackendApiKey.Id

	getResp, err := u.Store.GetBackendAPIKey(ctx, &backendv1.GetBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, getResp.BackendApiKey)
	require.Equal(t, keyID, getResp.BackendApiKey.Id)
}

func TestGetBackendAPIKey_DoesNotExist(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetBackendAPIKey(ctx, &backendv1.GetBackendAPIKeyRequest{
		Id: idformat.BackendAPIKey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateBackendAPIKey_UpdatesFields(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.BackendApiKey.Id

	updateResp, err := u.Store.UpdateBackendAPIKey(ctx, &backendv1.UpdateBackendAPIKeyRequest{
		Id: keyID,
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key2",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "key2", updateResp.BackendApiKey.DisplayName)
}

func TestRevokeBackendAPIKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.BackendApiKey.Id

	_, err = u.Store.RevokeBackendAPIKey(ctx, &backendv1.RevokeBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	getResp, err := u.Store.GetBackendAPIKey(ctx, &backendv1.GetBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, getResp.BackendApiKey.Revoked)

	_, err = u.Store.RevokeBackendAPIKey(ctx, &backendv1.RevokeBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	getResp2, err := u.Store.GetBackendAPIKey(ctx, &backendv1.GetBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, getResp2.BackendApiKey.Revoked)
}

func TestDeleteBackendAPIKey_RemovesKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
		BackendApiKey: &backendv1.BackendAPIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.BackendApiKey.Id

	// Revoke before delete
	_, err = u.Store.RevokeBackendAPIKey(ctx, &backendv1.RevokeBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	_, err = u.Store.DeleteBackendAPIKey(ctx, &backendv1.DeleteBackendAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	_, err = u.Store.GetBackendAPIKey(ctx, &backendv1.GetBackendAPIKeyRequest{Id: keyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListBackendAPIKeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
			BackendApiKey: &backendv1.BackendAPIKey{
				DisplayName: "test-key",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.BackendApiKey.Id)
	}

	listResp, err := u.Store.ListBackendAPIKeys(ctx, &backendv1.ListBackendAPIKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.BackendApiKeys, 3)

	var respIds []string
	for _, key := range listResp.BackendApiKeys {
		respIds = append(respIds, key.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListBackendAPIKeys_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateBackendAPIKey(ctx, &backendv1.CreateBackendAPIKeyRequest{
			BackendApiKey: &backendv1.BackendAPIKey{
				DisplayName: "test-key",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.BackendApiKey.Id)
	}

	resp1, err := u.Store.ListBackendAPIKeys(ctx, &backendv1.ListBackendAPIKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.BackendApiKeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListBackendAPIKeys(ctx, &backendv1.ListBackendAPIKeysRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.BackendApiKeys, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.BackendApiKeys {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.BackendApiKeys {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
