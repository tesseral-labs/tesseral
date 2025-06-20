package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestRole_CRUD(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{DisplayName: "org"})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2),
  		 (gen_random_uuid(), $1::uuid, $3, $3);
`,
		uuid.UUID(projectID).String(),
		"foo.bar.baz",
		"foo.bar.qux",
	)
	require.NoError(t, err)

	resp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "role1",
			Description:    "desc1",
			Actions:        []string{"foo.bar.baz", "foo.bar.qux"},
		},
	})
	require.NoError(t, err)
	role := resp.Role
	require.NotNil(t, role)
	require.Equal(t, "role1", role.DisplayName)
	require.ElementsMatch(t, []string{"foo.bar.baz", "foo.bar.qux"}, role.Actions)

	getResp, err := u.Store.GetRole(ctx, &backendv1.GetRoleRequest{Id: role.Id})
	require.NoError(t, err)
	require.Equal(t, role.Id, getResp.Role.Id)

	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		uuid.UUID(projectID).String(),
		"foo.bar.new",
	)
	require.NoError(t, err)

	updResp, err := u.Store.UpdateRole(ctx, &backendv1.UpdateRoleRequest{
		Id: role.Id,
		Role: &backendv1.Role{
			DisplayName: "role1-upd",
			Description: "desc1-upd",
			Actions:     []string{"foo.bar.new"},
		},
	})
	require.NoError(t, err)
	upd := updResp.Role
	require.Equal(t, "role1-upd", upd.DisplayName)
	require.Equal(t, "desc1-upd", upd.Description)
	require.ElementsMatch(t, []string{"foo.bar.new"}, upd.Actions)

	_, err = u.Store.DeleteRole(ctx, &backendv1.DeleteRoleRequest{Id: role.Id})
	require.NoError(t, err)

	_, err = u.Store.GetRole(ctx, &backendv1.GetRoleRequest{Id: role.Id})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestRole_Create_InvalidOrgOrAction(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		uuid.UUID(projectID).String(),
		"foo.bar.baz",
	)
	require.NoError(t, err)

	_, err = u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: "invalid-id",
			DisplayName:    "role1",
			Actions:        []string{"foo.bar.baz"},
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())

	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{DisplayName: "org"})
	_, err = u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "role2",
			Actions:        []string{"not.exist.action"},
		},
	})
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestRole_Update_InvalidIDOrAction(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{DisplayName: "org"})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
		uuid.UUID(projectID).String(),
		"foo.bar.baz",
	)
	require.NoError(t, err)

	resp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
		Role: &backendv1.Role{
			OrganizationId: orgID,
			DisplayName:    "role1",
			Actions:        []string{"foo.bar.baz"},
		},
	})
	require.NoError(t, err)

	_, err = u.Store.UpdateRole(ctx, &backendv1.UpdateRoleRequest{
		Id: resp.Role.Id,
		Role: &backendv1.Role{
			Actions: []string{"not.exist.action"},
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestRole_List_ByProjectAndByOrg_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{DisplayName: "org1"})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2),
  		 (gen_random_uuid(), $1::uuid, $3, $3);
`,
		uuid.UUID(projectID).String(),
		"foo.bar.baz",
		"foo.bar.qux",
	)
	require.NoError(t, err)

	var projectRoleIDs []string
	for range 12 {
		resp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: "", // project-level role
				DisplayName:    "test-role",
				Actions:        []string{"foo.bar.baz"},
			},
		})
		require.NoError(t, err)
		projectRoleIDs = append(projectRoleIDs, resp.Role.Id)
	}

	var organizationRoleIDs []string
	for range 12 {
		resp, err := u.Store.CreateRole(ctx, &backendv1.CreateRoleRequest{
			Role: &backendv1.Role{
				OrganizationId: orgID,
				DisplayName:    "test-role",
				Actions:        []string{"foo.bar.baz"},
			},
		})
		require.NoError(t, err)
		organizationRoleIDs = append(organizationRoleIDs, resp.Role.Id)
	}

	// List project roles
	resp1, err := u.Store.ListRoles(ctx, &backendv1.ListRolesRequest{})
	require.NoError(t, err)
	require.Len(t, resp1.Roles, 10)
	require.NotEmpty(t, resp1.NextPageToken)
	resp2, err := u.Store.ListRoles(ctx, &backendv1.ListRolesRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.Len(t, resp2.Roles, 2)
	var ids []string
	for _, r := range append(resp1.Roles, resp2.Roles...) {
		ids = append(ids, r.Id)
	}
	require.ElementsMatch(t, projectRoleIDs, ids)

	// List organization roles
	org1Resp, err := u.Store.ListRoles(ctx, &backendv1.ListRolesRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.Len(t, org1Resp.Roles, 10)
	require.NotEmpty(t, org1Resp.NextPageToken)
	org2Resp, err := u.Store.ListRoles(ctx, &backendv1.ListRolesRequest{
		OrganizationId: orgID,
		PageToken:      org1Resp.NextPageToken,
	})
	require.NoError(t, err)
	require.Len(t, org2Resp.Roles, 2)
	var orgIDs []string
	for _, r := range append(org1Resp.Roles, org2Resp.Roles...) {
		orgIDs = append(orgIDs, r.Id)
	}
	require.ElementsMatch(t, organizationRoleIDs, orgIDs)
}
