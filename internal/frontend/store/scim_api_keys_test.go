package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateSCIMAPIKey_SCIMEnabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})

	resp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.ScimApiKey)
	require.Equal(t, "scim-key-1", resp.ScimApiKey.DisplayName)
}

func TestCreateSCIMAPIKey_SCIMDisabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(false),
	})

	_, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestGetSCIMAPIKey_Exists(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})
	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id

	getResp, err := u.Store.GetSCIMAPIKey(ctx, &frontendv1.GetSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, getResp.ScimApiKey)
	require.Equal(t, keyID, getResp.ScimApiKey.Id)
}

func TestGetSCIMAPIKey_DoesNotExist(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})

	_, err := u.Store.GetSCIMAPIKey(ctx, &frontendv1.GetSCIMAPIKeyRequest{
		Id: idformat.SCIMAPIKey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateSCIMAPIKey_UpdatesFields(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})
	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id

	updateResp, err := u.Store.UpdateSCIMAPIKey(ctx, &frontendv1.UpdateSCIMAPIKeyRequest{
		Id: keyID,
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-2",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "scim-key-2", updateResp.ScimApiKey.DisplayName)
}

func TestRevokeSCIMAPIKey_SetsRevoked(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})
	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id

	revokeResp, err := u.Store.RevokeSCIMAPIKey(ctx, &frontendv1.RevokeSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, revokeResp.ScimApiKey.Revoked)
}

func TestDeleteSCIMAPIKey_Revoked(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})
	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id

	_, err = u.Store.RevokeSCIMAPIKey(ctx, &frontendv1.RevokeSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	_, err = u.Store.DeleteSCIMAPIKey(ctx, &frontendv1.DeleteSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
}

func TestDeleteSCIMAPIKey_NotRevoked(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})
	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &frontendv1.SCIMAPIKey{
			DisplayName: "scim-key-1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id

	_, err = u.Store.DeleteSCIMAPIKey(ctx, &frontendv1.DeleteSCIMAPIKeyRequest{Id: keyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestListSCIMAPIKeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})

	var ids []string
	for i := 0; i < 3; i++ {
		resp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
			ScimApiKey: &frontendv1.SCIMAPIKey{
				DisplayName: "scim-key",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.ScimApiKey.Id)
	}

	listResp, err := u.Store.ListSCIMAPIKeys(ctx, &frontendv1.ListSCIMAPIKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.ScimApiKeys, 3)

	var respIds []string
	for _, k := range listResp.ScimApiKeys {
		respIds = append(respIds, k.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListSCIMAPIKeys_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
		ScimEnabled: refOrNil(true),
	})

	var createdIDs []string
	for i := 0; i < 15; i++ {
		resp, err := u.Store.CreateSCIMAPIKey(ctx, &frontendv1.CreateSCIMAPIKeyRequest{
			ScimApiKey: &frontendv1.SCIMAPIKey{
				DisplayName: "scim-key",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.ScimApiKey.Id)
	}

	resp1, err := u.Store.ListSCIMAPIKeys(ctx, &frontendv1.ListSCIMAPIKeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.ScimApiKeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListSCIMAPIKeys(ctx, &frontendv1.ListSCIMAPIKeysRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.ScimApiKeys, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.ScimApiKeys {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.ScimApiKeys {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
