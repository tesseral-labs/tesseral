package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateAPIKey_ApiKeysEnabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	res, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.ApiKey)
	require.NotEmpty(t, res.ApiKey.Id)
	require.Equal(t, "Test Key", res.ApiKey.DisplayName)
	require.NotEmpty(t, res.ApiKey.SecretTokenSuffix)
}

func TestCreateAPIKey_ProjectApiKeysDisabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
	UPDATE projects
		SET api_keys_enabled = false
	WHERE id = $1::uuid
	`,
		projectUUID,
	)
	require.NoError(t, err)

	_, err = u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodePermissionDenied, connectErr.Code())
}

func TestCreateAPIKey_OrganizationApiKeysDisabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(false),
	})

	_, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodePermissionDenied, connectErr.Code())
}

func TestGetAPIKey_Exists(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	getResp, err := u.Store.GetAPIKey(ctx, &frontendv1.GetAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, getResp.ApiKey)
	require.Equal(t, keyID, getResp.ApiKey.Id)
}

func TestGetAPIKey_DoesNotExist(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	_, err := u.Store.GetAPIKey(ctx, &frontendv1.GetAPIKeyRequest{Id: "nonexistent-id"})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestUpdateAPIKey_UpdatesFields(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	updateResp, err := u.Store.UpdateAPIKey(ctx, &frontendv1.UpdateAPIKeyRequest{
		Id: keyID,
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Updated Key",
		},
	})
	require.NoError(t, err)
	updated := updateResp.ApiKey
	require.Equal(t, "Updated Key", updated.DisplayName)
}

func TestRevokeAPIKey_RevokesKey(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	_, err = u.Store.RevokeAPIKey(ctx, &frontendv1.RevokeAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	getResp, err := u.Store.GetAPIKey(ctx, &frontendv1.GetAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, getResp.ApiKey.Revoked)
}

func TestDeleteAPIKey_RemovesKey(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "Test Key",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	_, err = u.Store.RevokeAPIKey(ctx, &frontendv1.RevokeAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	_, err = u.Store.DeleteAPIKey(ctx, &frontendv1.DeleteAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	_, err = u.Store.GetAPIKey(ctx, &frontendv1.GetAPIKeyRequest{Id: keyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListAPIKeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
			ApiKey: &frontendv1.APIKey{
				DisplayName: "Test Key",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.ApiKey.Id)
	}

	listResp, err := u.Store.ListAPIKeys(ctx, &frontendv1.ListAPIKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.ApiKeys, 3)

	var respIds []string
	for _, key := range listResp.ApiKeys {
		respIds = append(respIds, key.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListAPIKeys_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:    "Test Organization",
		ApiKeysEnabled: refOrNil(true),
	})

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
			ApiKey: &frontendv1.APIKey{
				DisplayName: "Test Key",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.ApiKey.Id)
	}

	resp1, err := u.Store.ListAPIKeys(ctx, &frontendv1.ListAPIKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.ApiKeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListAPIKeys(ctx, &frontendv1.ListAPIKeysRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.ApiKeys, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.ApiKeys {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.ApiKeys {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
