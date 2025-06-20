package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateSCIMAPIKey_SCIMEnabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	res, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.ScimApiKey)
	require.Equal(t, orgID, res.ScimApiKey.OrganizationId)
	require.Equal(t, "key1", res.ScimApiKey.DisplayName)
	require.NotEmpty(t, res.ScimApiKey.Id)
	require.NotEmpty(t, res.ScimApiKey.CreateTime)
	require.NotEmpty(t, res.ScimApiKey.UpdateTime)
	require.NotEmpty(t, res.ScimApiKey.SecretToken)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.create")
}

func TestCreateSCIMAPIKey_SCIMDisabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(false),
	})

	_, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestGetSCIMAPIKey_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id

	getResp, err := u.Store.GetSCIMAPIKey(ctx, &backendv1.GetSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, getResp.ScimApiKey)
	require.Equal(t, keyID, getResp.ScimApiKey.Id)
}

func TestGetSCIMAPIKey_DoesNotExist(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetSCIMAPIKey(ctx, &backendv1.GetSCIMAPIKeyRequest{
		Id: idformat.SCIMAPIKey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateSCIMAPIKey_UpdatesFields(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.create")

	updateResp, err := u.Store.UpdateSCIMAPIKey(ctx, &backendv1.UpdateSCIMAPIKeyRequest{
		Id: keyID,
		ScimApiKey: &backendv1.SCIMAPIKey{
			DisplayName: "key2",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "key2", updateResp.ScimApiKey.DisplayName)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.update")
}

func TestRevokeSCIMAPIKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.create")

	revokeResp, err := u.Store.RevokeSCIMAPIKey(ctx, &backendv1.RevokeSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, revokeResp.ScimApiKey)
	require.True(t, revokeResp.ScimApiKey.Revoked)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.revoke")

	getResp, err := u.Store.GetSCIMAPIKey(ctx, &backendv1.GetSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.True(t, getResp.ScimApiKey.Revoked)

	revokeResp2, err := u.Store.RevokeSCIMAPIKey(ctx, &backendv1.RevokeSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	require.NotNil(t, revokeResp2.ScimApiKey)
	require.True(t, revokeResp2.ScimApiKey.Revoked)
}

func TestDeleteSCIMAPIKey_RemovesKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	createResp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
		ScimApiKey: &backendv1.SCIMAPIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	keyID := createResp.ScimApiKey.Id
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.create")

	_, err = u.Store.RevokeSCIMAPIKey(ctx, &backendv1.RevokeSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.revoke")

	_, err = u.Store.DeleteSCIMAPIKey(ctx, &backendv1.DeleteSCIMAPIKeyRequest{Id: keyID})
	require.NoError(t, err)

	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_SCIM_API_KEY, "tesseral.scim_api_keys.delete")

	_, err = u.Store.GetSCIMAPIKey(ctx, &backendv1.GetSCIMAPIKeyRequest{Id: keyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListSCIMAPIKeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
			ScimApiKey: &backendv1.SCIMAPIKey{
				OrganizationId: orgID,
				DisplayName:    "test-key",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.ScimApiKey.Id)
	}

	listResp, err := u.Store.ListSCIMAPIKeys(ctx, &backendv1.ListSCIMAPIKeysRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.ScimApiKeys, 3)

	var respIds []string
	for _, key := range listResp.ScimApiKeys {
		respIds = append(respIds, key.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListSCIMAPIKeys_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
		ScimEnabled: refOrNil(true),
	})

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateSCIMAPIKey(ctx, &backendv1.CreateSCIMAPIKeyRequest{
			ScimApiKey: &backendv1.SCIMAPIKey{
				OrganizationId: orgID,
				DisplayName:    "test-key",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.ScimApiKey.Id)
	}

	resp1, err := u.Store.ListSCIMAPIKeys(ctx, &backendv1.ListSCIMAPIKeysRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.ScimApiKeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListSCIMAPIKeys(ctx, &backendv1.ListSCIMAPIKeysRequest{OrganizationId: orgID, PageToken: resp1.NextPageToken})
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
