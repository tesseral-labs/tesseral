package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateUserRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	orgID := idformat.Organization.Format(authn.OrganizationID(ctx))
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "user1@example.com",
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		projectID,
		"test.action.1",
	)
	require.NoError(t, err)

	roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1"},
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	resp, err := u.Store.CreateUserRoleAssignment(ctx, &frontendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &frontendv1.UserRoleAssignment{
			UserId: userID,
			RoleId: roleID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.UserRoleAssignment)
	require.Equal(t, userID, resp.UserRoleAssignment.UserId)
	require.Equal(t, roleID, resp.UserRoleAssignment.RoleId)
}

func TestCreateUserRoleAssignment_InvalidIDs(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.CreateUserRoleAssignment(ctx, &frontendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &frontendv1.UserRoleAssignment{
			UserId: "invalid-id",
			RoleId: "invalid-id",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestDeleteUserRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	orgID := idformat.Organization.Format(authn.OrganizationID(ctx))
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "user1@example.com",
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		projectID,
		"test.action.1",
	)
	require.NoError(t, err)

	roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1"},
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	createResp, err := u.Store.CreateUserRoleAssignment(ctx, &frontendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &frontendv1.UserRoleAssignment{
			UserId: userID,
			RoleId: roleID,
		},
	})
	require.NoError(t, err)
	assignmentID := createResp.UserRoleAssignment.Id

	_, err = u.Store.DeleteUserRoleAssignment(ctx, &frontendv1.DeleteUserRoleAssignmentRequest{Id: assignmentID})
	require.NoError(t, err)
}

func TestDeleteUserRoleAssignment_InvalidID(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.DeleteUserRoleAssignment(ctx, &frontendv1.DeleteUserRoleAssignmentRequest{Id: "invalid-id"})
	require.Error(t, err)
}

func TestListUserRoleAssignments_ByUser_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	orgID := idformat.Organization.Format(authn.OrganizationID(ctx))
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "user1@example.com",
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		projectID,
		"test.action.1",
	)
	require.NoError(t, err)

	var assignmentIDs []string
	for range 3 {
		roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
			Role: &frontendv1.Role{
				DisplayName: "role",
				Actions:     []string{"test.action.1"},
			},
		})
		require.NoError(t, err)
		roleID := roleResp.Role.Id

		createResp, err := u.Store.CreateUserRoleAssignment(ctx, &frontendv1.CreateUserRoleAssignmentRequest{
			UserRoleAssignment: &frontendv1.UserRoleAssignment{
				UserId: userID,
				RoleId: roleID,
			},
		})
		require.NoError(t, err)
		assignmentIDs = append(assignmentIDs, createResp.UserRoleAssignment.Id)
	}

	listResp, err := u.Store.ListUserRoleAssignments(ctx, &frontendv1.ListUserRoleAssignmentsRequest{UserId: userID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.UserRoleAssignments, 3)

	var respIDs []string
	for _, a := range listResp.UserRoleAssignments {
		respIDs = append(respIDs, a.Id)
	}
	require.ElementsMatch(t, assignmentIDs, respIDs)
}

func TestListUserRoleAssignments_ByUser_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	orgID := idformat.Organization.Format(authn.OrganizationID(ctx))
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "user1@example.com",
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		projectID,
		"test.action.1",
	)
	require.NoError(t, err)

	var createdIDs []string
	for range 15 {
		roleResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
			Role: &frontendv1.Role{
				DisplayName: "role",
				Actions:     []string{"test.action.1"},
			},
		})
		require.NoError(t, err)
		roleID := roleResp.Role.Id

		createResp, err := u.Store.CreateUserRoleAssignment(ctx, &frontendv1.CreateUserRoleAssignmentRequest{
			UserRoleAssignment: &frontendv1.UserRoleAssignment{
				UserId: userID,
				RoleId: roleID,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, createResp.UserRoleAssignment.Id)
	}

	resp1, err := u.Store.ListUserRoleAssignments(ctx, &frontendv1.ListUserRoleAssignmentsRequest{UserId: userID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.UserRoleAssignments, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUserRoleAssignments(ctx, &frontendv1.ListUserRoleAssignmentsRequest{UserId: userID, PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.UserRoleAssignments, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, a := range resp1.UserRoleAssignments {
		allIDs = append(allIDs, a.Id)
	}
	for _, a := range resp2.UserRoleAssignments {
		allIDs = append(allIDs, a.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
