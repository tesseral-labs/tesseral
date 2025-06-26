package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateAPIKeyRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2),
  		 (gen_random_uuid(), $1::uuid, $3, $3);
`,
		uuid.UUID(projectID).String(),
		"test.action.1",
		"test.action.2",
	)
	require.NoError(t, err)

	role, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "role1",
			Actions:        []string{"test.action.1", "test.action.2"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, role.Role)
	roleID := role.Role.Id

	resp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &backendv1.CreateAPIKeyRoleAssignmentRequest{
		ApiKeyRoleAssignment: &backendv1.APIKeyRoleAssignment{
			ApiKeyId: apiKeyID,
			RoleId:   roleID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.ApiKeyRoleAssignment)
	require.Equal(t, apiKeyID, resp.ApiKeyRoleAssignment.ApiKeyId)
	require.Equal(t, roleID, resp.ApiKeyRoleAssignment.RoleId)

	apiKey, err := u.Store.AuthenticateAPIKey(ctx, &backendv1.AuthenticateAPIKeyRequest{
		SecretToken: apiKeyResp.ApiKey.SecretToken,
	})
	require.NoError(t, err)
	require.NotNil(t, apiKey)
	require.Equal(t, apiKeyID, apiKey.ApiKeyId)
	require.Equal(t, orgID, apiKey.OrganizationId)
	require.Equal(t, []string{"test.action.1", "test.action.2"}, apiKey.Actions)
}

func TestCreateAPIKeyRoleAssignment_InvalidIDs(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &backendv1.CreateAPIKeyRoleAssignmentRequest{
		ApiKeyRoleAssignment: &backendv1.APIKeyRoleAssignment{
			ApiKeyId: "invalid-id",
			RoleId:   "invalid-id",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestDeleteAPIKeyRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	role, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "role1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, role.Role)
	roleID := role.Role.Id

	createResp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &backendv1.CreateAPIKeyRoleAssignmentRequest{
		ApiKeyRoleAssignment: &backendv1.APIKeyRoleAssignment{
			ApiKeyId: apiKeyID,
			RoleId:   roleID,
		},
	})
	require.NoError(t, err)
	assignmentID := createResp.ApiKeyRoleAssignment.Id

	_, err = u.Store.DeleteAPIKeyRoleAssignment(ctx, &backendv1.DeleteAPIKeyRoleAssignmentRequest{Id: assignmentID})
	require.NoError(t, err)
}

func TestDeleteAPIKeyRoleAssignment_InvalidID(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.DeleteAPIKeyRoleAssignment(ctx, &backendv1.DeleteAPIKeyRoleAssignmentRequest{Id: "invalid-id"})
	require.Error(t, err)
}

func TestListAPIKeyRoleAssignments_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	var assignmentIDs []string
	for range 3 {
		role, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: orgID,
				DisplayName:    "test-role",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, role.Role)
		roleID := role.Role.Id

		createResp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &backendv1.CreateAPIKeyRoleAssignmentRequest{
			ApiKeyRoleAssignment: &backendv1.APIKeyRoleAssignment{
				ApiKeyId: apiKeyID,
				RoleId:   roleID,
			},
		})
		require.NoError(t, err)
		assignmentIDs = append(assignmentIDs, createResp.ApiKeyRoleAssignment.Id)
	}

	listResp, err := u.Store.ListAPIKeyRoleAssignments(ctx, &backendv1.ListAPIKeyRoleAssignmentsRequest{ApiKeyId: apiKeyID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.ApiKeyRoleAssignments, 3)

	var respIDs []string
	for _, a := range listResp.ApiKeyRoleAssignments {
		respIDs = append(respIDs, a.Id)
	}
	require.ElementsMatch(t, assignmentIDs, respIDs)
}

func TestListAPIKeyRoleAssignments_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:    "test",
		ApiKeysEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &backendv1.CreateAPIKeyRequest{
		ApiKey: &backendv1.APIKey{
			OrganizationId: orgID,
			DisplayName:    "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	var createdIDs []string
	for range 15 {
		role, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: orgID,
				DisplayName:    "test-role",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, role.Role)
		roleID := role.Role.Id

		createResp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &backendv1.CreateAPIKeyRoleAssignmentRequest{
			ApiKeyRoleAssignment: &backendv1.APIKeyRoleAssignment{
				ApiKeyId: apiKeyID,
				RoleId:   roleID,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, createResp.ApiKeyRoleAssignment.Id)
	}

	resp1, err := u.Store.ListAPIKeyRoleAssignments(ctx, &backendv1.ListAPIKeyRoleAssignmentsRequest{ApiKeyId: apiKeyID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.ApiKeyRoleAssignments, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListAPIKeyRoleAssignments(ctx, &backendv1.ListAPIKeyRoleAssignmentsRequest{ApiKeyId: apiKeyID, PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.ApiKeyRoleAssignments, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, a := range resp1.ApiKeyRoleAssignments {
		allIDs = append(allIDs, a.Id)
	}
	for _, a := range resp2.ApiKeyRoleAssignments {
		allIDs = append(allIDs, a.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
