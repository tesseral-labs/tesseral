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

func TestCreateAPIKeyRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		ApiKeysEnabled:     refOrNil(true),
		CustomRolesEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "key1",
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

	roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1", "test.action.2"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, roleResp.Role)
	roleID := roleResp.Role.Id

	resp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &frontendv1.CreateAPIKeyRoleAssignmentRequest{
		ApiKeyRoleAssignment: &frontendv1.APIKeyRoleAssignment{
			ApiKeyId: apiKeyID,
			RoleId:   roleID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.ApiKeyRoleAssignment)
	require.Equal(t, apiKeyID, resp.ApiKeyRoleAssignment.ApiKeyId)
	require.Equal(t, roleID, resp.ApiKeyRoleAssignment.RoleId)
}

func TestCreateAPIKeyRoleAssignment_InvalidIDs(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &frontendv1.CreateAPIKeyRoleAssignmentRequest{
		ApiKeyRoleAssignment: &frontendv1.APIKeyRoleAssignment{
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

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		ApiKeysEnabled:     refOrNil(true),
		CustomRolesEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, roleResp.Role)
	roleID := roleResp.Role.Id

	createResp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &frontendv1.CreateAPIKeyRoleAssignmentRequest{
		ApiKeyRoleAssignment: &frontendv1.APIKeyRoleAssignment{
			ApiKeyId: apiKeyID,
			RoleId:   roleID,
		},
	})
	require.NoError(t, err)
	assignmentID := createResp.ApiKeyRoleAssignment.Id

	_, err = u.Store.DeleteAPIKeyRoleAssignment(ctx, &frontendv1.DeleteAPIKeyRoleAssignmentRequest{Id: assignmentID})
	require.NoError(t, err)
}

func TestDeleteAPIKeyRoleAssignment_InvalidID(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.DeleteAPIKeyRoleAssignment(ctx, &frontendv1.DeleteAPIKeyRoleAssignmentRequest{Id: "invalid-id"})
	require.Error(t, err)
}

func TestListAPIKeyRoleAssignments_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		ApiKeysEnabled:     refOrNil(true),
		CustomRolesEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	var assignmentIDs []string
	for range 3 {
		roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
			Role: &frontendv1.Role{
				DisplayName: "test-role",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, roleResp.Role)
		roleID := roleResp.Role.Id

		createResp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &frontendv1.CreateAPIKeyRoleAssignmentRequest{
			ApiKeyRoleAssignment: &frontendv1.APIKeyRoleAssignment{
				ApiKeyId: apiKeyID,
				RoleId:   roleID,
			},
		})
		require.NoError(t, err)
		assignmentIDs = append(assignmentIDs, createResp.ApiKeyRoleAssignment.Id)
	}

	listResp, err := u.Store.ListAPIKeyRoleAssignments(ctx, &frontendv1.ListAPIKeyRoleAssignmentsRequest{ApiKeyId: apiKeyID})
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

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		ApiKeysEnabled:     refOrNil(true),
		CustomRolesEnabled: refOrNil(true),
	})

	apiKeyResp, err := u.Store.CreateAPIKey(ctx, &frontendv1.CreateAPIKeyRequest{
		ApiKey: &frontendv1.APIKey{
			DisplayName: "key1",
		},
	})
	require.NoError(t, err)
	apiKeyID := apiKeyResp.ApiKey.Id

	var createdIDs []string
	for range 15 {
		roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
			Role: &frontendv1.Role{
				DisplayName: "test-role",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, roleResp.Role)
		roleID := roleResp.Role.Id

		createResp, err := u.Store.CreateAPIKeyRoleAssignment(ctx, &frontendv1.CreateAPIKeyRoleAssignmentRequest{
			ApiKeyRoleAssignment: &frontendv1.APIKeyRoleAssignment{
				ApiKeyId: apiKeyID,
				RoleId:   roleID,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, createResp.ApiKeyRoleAssignment.Id)
	}

	resp1, err := u.Store.ListAPIKeyRoleAssignments(ctx, &frontendv1.ListAPIKeyRoleAssignmentsRequest{ApiKeyId: apiKeyID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.ApiKeyRoleAssignments, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListAPIKeyRoleAssignments(ctx, &frontendv1.ListAPIKeyRoleAssignmentsRequest{ApiKeyId: apiKeyID, PageToken: resp1.NextPageToken})
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
