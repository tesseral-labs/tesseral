package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateAPIKey_ApiKeysEnabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	res, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.ApiKey)
	require.Equal(t, orgID, res.ApiKey.OrganizationId)
	require.Equal(t, "key1", res.ApiKey.DisplayName)
	require.NotEmpty(t, res.ApiKey.Id)
	require.NotEmpty(t, res.ApiKey.CreateTime)
	require.NotEmpty(t, res.ApiKey.UpdateTime)
	require.NotEmpty(t, res.ApiKey.SecretToken)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_API_KEY, "tesseral.api_keys.create")
}

func TestCreateAPIKey_ApiKeysDisabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(false),
	})

	_, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodePermissionDenied, connectErr.Code())
}

func TestGetAPIKey_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	getResp, err := u.Store.GetAPIKey(ctx, &backendv1.GetAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, getResp.ApiKey)
	require.Equal(t, keyID, getResp.ApiKey.Id)
}

func TestGetAPIKey_DoesNotExist(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetAPIKey(ctx, &backendv1.GetAPIKeyRequest{
		Id: idformat.APIKey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateAPIKey_UpdatesFields(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	updateResp, err := u.Store.UpdateAPIKey(ctx, &backendv1.UpdateAPIKeyRequest{
		Id: keyID,
		ApiKey: &backendv1.APIKey{
			DisplayName: "key2",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "key2", updateResp.ApiKey.DisplayName)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_API_KEY, "tesseral.api_keys.update")
}

func TestRevokeAPIKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	_, err = u.Store.RevokeAPIKey(ctx, &backendv1.RevokeAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_API_KEY, "tesseral.api_keys.revoke")

	getResp, err := u.Store.GetAPIKey(ctx, &backendv1.GetAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, getResp.ApiKey.Revoked)

	_, err = u.Store.RevokeAPIKey(ctx, &backendv1.RevokeAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	getResp2, err := u.Store.GetAPIKey(ctx, &backendv1.GetAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, getResp2.ApiKey.Revoked)
}

func TestDeleteAPIKey_RemovesKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	_, err = u.Store.RevokeAPIKey(ctx, &backendv1.RevokeAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	_, err = u.Store.DeleteAPIKey(ctx, &backendv1.DeleteAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_API_KEY, "tesseral.api_keys.delete")

	_, err = u.Store.GetAPIKey(ctx, &backendv1.GetAPIKeyRequest{Id: keyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListAPIKeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
			ApiKey: &backendv1.APIKey{
				OrganizationId: orgID,
				DisplayName:    "test-key",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.ApiKey.Id)
	}

	listResp, err := u.Store.ListAPIKeys(ctx, &backendv1.ListAPIKeysRequest{OrganizationId: orgID})
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

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
			ApiKey: &backendv1.APIKey{
				OrganizationId: orgID,
				DisplayName:    "test-key",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.ApiKey.Id)
	}

	resp1, err := u.Store.ListAPIKeys(ctx, &backendv1.ListAPIKeysRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.ApiKeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListAPIKeys(ctx, &backendv1.ListAPIKeysRequest{OrganizationId: orgID, PageToken: resp1.NextPageToken})
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

func TestAuthenticateAPIKey_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ApiKey.Id

	authResp, err := u.Store.AuthenticateAPIKey(ctx, &backendv1.AuthenticateAPIKeyRequest{
		SecretToken: createResp.ApiKey.SecretToken,
	})
	require.NoError(t, err)
	require.NotNil(t, authResp)
	require.Equal(t, keyID, authResp.ApiKeyId)
	require.Equal(t, orgID, authResp.OrganizationId)
}

func TestAuthenticateAPIKey_InvalidSecretToken(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.AuthenticateAPIKey(ctx, &backendv1.AuthenticateAPIKeyRequest{
		SecretToken: idformat.APIKey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}
