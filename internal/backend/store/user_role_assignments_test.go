package store

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateUserRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})

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

	roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "test-role",
			Actions:        []string{"test.action.1", "test.action.2"},
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	resp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &backendv1.UserRoleAssignment{
			UserId: userID,
			RoleId: roleID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.UserRoleAssignment)
	require.Equal(t, userID, resp.UserRoleAssignment.UserId)
	require.Equal(t, roleID, resp.UserRoleAssignment.RoleId)
}

func TestCreateUserRoleAssignment_NotFound(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &backendv1.UserRoleAssignment{
			UserId: idformat.User.Format(uuid.New()),
			RoleId: idformat.Role.Format(uuid.New()),
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestDeleteUserRoleAssignment_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})

	roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "test-role",
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	createResp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &backendv1.UserRoleAssignment{
			UserId: userID,
			RoleId: roleID,
		},
	})
	require.NoError(t, err)
	assignmentID := createResp.UserRoleAssignment.Id

	_, err = u.Store.DeleteUserRoleAssignment(ctx, &backendv1.DeleteUserRoleAssignmentRequest{Id: assignmentID})
	require.NoError(t, err)
}

func TestGetUserRoleAssignment_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})

	roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "test-role",
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	createResp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
		UserRoleAssignment: &backendv1.UserRoleAssignment{
			UserId: userID,
			RoleId: roleID,
		},
	})
	require.NoError(t, err)
	assignmentID := createResp.UserRoleAssignment.Id

	getResp, err := u.Store.GetUserRoleAssignment(ctx, &backendv1.GetUserRoleAssignmentRequest{Id: assignmentID})
	require.NoError(t, err)
	require.NotNil(t, getResp.UserRoleAssignment)
	require.Equal(t, assignmentID, getResp.UserRoleAssignment.Id)
}

func TestGetUserRoleAssignment_NotFound(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetUserRoleAssignment(ctx, &backendv1.GetUserRoleAssignmentRequest{
		Id: idformat.UserRoleAssignment.Format(uuid.New()),
	})
	require.Error(t, err)
}

func TestListUserRoleAssignments_ByUser_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})

	var assignmentIDs []string
	for range 3 {
		roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: orgID,
				DisplayName:    "test-role",
			},
		})
		require.NoError(t, err)
		roleID := roleResp.Role.Id
		createResp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
			UserRoleAssignment: &backendv1.UserRoleAssignment{
				UserId: userID,
				RoleId: roleID,
			},
		})
		require.NoError(t, err)
		assignmentIDs = append(assignmentIDs, createResp.UserRoleAssignment.Id)
	}

	listResp, err := u.Store.ListUserRoleAssignments(ctx, &backendv1.ListUserRoleAssignmentsRequest{UserId: userID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.UserRoleAssignments, 3)

	var respIDs []string
	for _, a := range listResp.UserRoleAssignments {
		respIDs = append(respIDs, a.Id)
	}
	require.ElementsMatch(t, assignmentIDs, respIDs)
}

func TestListUserRoleAssignments_ByRole_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "test-role",
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	var assignmentIDs []string
	for i := range 3 {
		userID := u.Environment.NewUser(t, orgID, &backendv1.User{
			Email: fmt.Sprintf("user%d@example.com", i),
		})
		createResp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
			UserRoleAssignment: &backendv1.UserRoleAssignment{
				UserId: userID,
				RoleId: roleID,
			},
		})
		require.NoError(t, err)
		assignmentIDs = append(assignmentIDs, createResp.UserRoleAssignment.Id)
	}

	listResp, err := u.Store.ListUserRoleAssignments(ctx, &backendv1.ListUserRoleAssignmentsRequest{RoleId: roleID})
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

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})
	userID := u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})

	var createdIDs []string
	for range 15 {
		roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: orgID,
				DisplayName:    "test-role",
			},
		})
		require.NoError(t, err)
		roleID := roleResp.Role.Id

		createResp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
			UserRoleAssignment: &backendv1.UserRoleAssignment{
				UserId: userID,
				RoleId: roleID,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, createResp.UserRoleAssignment.Id)
	}

	resp1, err := u.Store.ListUserRoleAssignments(ctx, &backendv1.ListUserRoleAssignmentsRequest{UserId: userID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.UserRoleAssignments, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUserRoleAssignments(ctx, &backendv1.ListUserRoleAssignmentsRequest{UserId: userID, PageToken: resp1.NextPageToken})
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

func TestListUserRoleAssignments_ByRole_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	roleResp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "roleY",
		},
	})
	require.NoError(t, err)
	roleID := roleResp.Role.Id

	var createdIDs []string
	for i := range 15 {
		userID := u.Environment.NewUser(t, orgID, &backendv1.User{
			Email: fmt.Sprintf("user%d@example.com", i),
		})
		createResp, err := u.Store.CreateUserRoleAssignment(ctx, &backendv1.CreateUserRoleAssignmentRequest{
			UserRoleAssignment: &backendv1.UserRoleAssignment{
				UserId: userID,
				RoleId: roleID,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, createResp.UserRoleAssignment.Id)
	}

	resp1, err := u.Store.ListUserRoleAssignments(ctx, &backendv1.ListUserRoleAssignmentsRequest{RoleId: roleID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.UserRoleAssignments, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUserRoleAssignments(ctx, &backendv1.ListUserRoleAssignmentsRequest{RoleId: roleID, PageToken: resp1.NextPageToken})
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
