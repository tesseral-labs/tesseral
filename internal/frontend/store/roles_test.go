package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateRole_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2),
         (gen_random_uuid(), $1::uuid, $3, $3);
`,
		projectID,
		"test.action.1",
		"test.action.2",
	)
	require.NoError(t, err)

	resp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1", "test.action.2"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Role)
	require.Equal(t, "role1", resp.Role.DisplayName)
	require.ElementsMatch(t, []string{"test.action.1", "test.action.2"}, resp.Role.Actions)
}

func TestCreateRole_InvalidAction(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	_, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"nonexistent.action"},
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestGetRole_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
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

	createResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1"},
		},
	})
	require.NoError(t, err)
	roleID := createResp.Role.Id

	getResp, err := u.Store.GetRole(ctx, &frontendv1.GetRoleRequest{Id: roleID})
	require.NoError(t, err)
	require.NotNil(t, getResp.Role)
	require.Equal(t, roleID, getResp.Role.Id)
}

func TestGetRole_InvalidID(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.GetRole(ctx, &frontendv1.GetRoleRequest{Id: "invalid-id"})
	require.Error(t, err)
}

func TestUpdateRole_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2),
         (gen_random_uuid(), $1::uuid, $3, $3);
`,
		projectID,
		"test.action.1",
		"test.action.2",
	)
	require.NoError(t, err)

	createResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1"},
		},
	})
	require.NoError(t, err)
	roleID := createResp.Role.Id

	updateResp, err := u.Store.UpdateRole(ctx, &frontendv1.UpdateRoleRequest{
		Id: roleID,
		Role: &frontendv1.Role{
			DisplayName: "role1-updated",
			Actions:     []string{"test.action.2"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "role1-updated", updateResp.Role.DisplayName)
	require.ElementsMatch(t, []string{"test.action.2"}, updateResp.Role.Actions)
}

func TestDeleteRole_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
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

	createResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
		Role: &frontendv1.Role{
			DisplayName: "role1",
			Actions:     []string{"test.action.1"},
		},
	})
	require.NoError(t, err)
	roleID := createResp.Role.Id

	_, err = u.Store.DeleteRole(ctx, &frontendv1.DeleteRoleRequest{Id: roleID})
	require.NoError(t, err)
}

func TestListRoles_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
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

	var roleIDs []string
	for range 3 {
		createResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
			Role: &frontendv1.Role{
				DisplayName: "role",
				Actions:     []string{"test.action.1"},
			},
		})
		require.NoError(t, err)
		roleIDs = append(roleIDs, createResp.Role.Id)
	}

	listResp, err := u.Store.ListRoles(ctx, &frontendv1.ListRolesRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Roles, 3)

	var respIDs []string
	for _, r := range listResp.Roles {
		respIDs = append(respIDs, r.Id)
	}
	require.ElementsMatch(t, roleIDs, respIDs)
}

func TestListRoles_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		CustomRolesEnabled: refOrNil(true),
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
		createResp, err := u.Store.CreateRole(ctx, &frontendv1.CreateRoleRequest{
			Role: &frontendv1.Role{
				DisplayName: "role",
				Actions:     []string{"test.action.1"},
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, createResp.Role.Id)
	}

	resp1, err := u.Store.ListRoles(ctx, &frontendv1.ListRolesRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.Roles, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListRoles(ctx, &frontendv1.ListRolesRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.Roles, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, r := range resp1.Roles {
		allIDs = append(allIDs, r.Id)
	}
	for _, r := range resp2.Roles {
		allIDs = append(allIDs, r.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
