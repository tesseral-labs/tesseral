package store

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateUserInvite_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	resp, err := u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
		UserInvite: &frontendv1.UserInvite{
			Email: "invitee@example.com",
			Owner: false,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.UserInvite)
	require.Equal(t, "invitee@example.com", resp.UserInvite.Email)
}

func TestCreateUserInvite_SelfServeUserCreationDisabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	_, err = u.Environment.DB.Exec(t.Context(), `update project_ui_settings set self_serve_create_users = false where project_id = $1`, projectUUID)
	require.NoError(t, err)

	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	_, err = u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
		UserInvite: &frontendv1.UserInvite{
			Email: "invitee@example.com",
			Owner: false,
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodePermissionDenied, connectErr.Code())
}

func TestCreateUserInvite_AlreadyExists(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})
	organizationID := idformat.Organization.Format(authn.OrganizationID(ctx))

	u.Environment.NewUser(t, organizationID, &backendv1.User{
		Email: "invitee@example.com",
	})

	_, err := u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
		UserInvite: &frontendv1.UserInvite{
			Email: "invitee@example.com",
			Owner: false,
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestGetUserInvite_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	createResp, err := u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
		UserInvite: &frontendv1.UserInvite{
			Email: "invitee@example.com",
			Owner: false,
		},
	})
	require.NoError(t, err)
	inviteID := createResp.UserInvite.Id

	getResp, err := u.Store.GetUserInvite(ctx, &frontendv1.GetUserInviteRequest{Id: inviteID})
	require.NoError(t, err)
	require.NotNil(t, getResp.UserInvite)
	require.Equal(t, inviteID, getResp.UserInvite.Id)
}

func TestGetUserInvite_NotFound(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	_, err := u.Store.GetUserInvite(ctx, &frontendv1.GetUserInviteRequest{
		Id: idformat.UserInvite.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestDeleteUserInvite_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	createResp, err := u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
		UserInvite: &frontendv1.UserInvite{
			Email: "invitee@example.com",
			Owner: false,
		},
	})
	require.NoError(t, err)
	inviteID := createResp.UserInvite.Id

	_, err = u.Store.DeleteUserInvite(ctx, &frontendv1.DeleteUserInviteRequest{Id: inviteID})
	require.NoError(t, err)

	_, err = u.Store.GetUserInvite(ctx, &frontendv1.GetUserInviteRequest{Id: inviteID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListUserInvites_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	var ids []string
	for i := range 3 {
		resp, err := u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
			UserInvite: &frontendv1.UserInvite{
				Email: fmt.Sprintf("invitee-%d@example.com", i),
				Owner: false,
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.UserInvite.Id)
	}

	listResp, err := u.Store.ListUserInvites(ctx, &frontendv1.ListUserInvitesRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.UserInvites, 3)

	var respIds []string
	for _, u := range listResp.UserInvites {
		respIds = append(respIds, u.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListUserInvites_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	var createdIDs []string
	for i := range 15 {
		resp, err := u.Store.CreateUserInvite(ctx, &frontendv1.CreateUserInviteRequest{
			UserInvite: &frontendv1.UserInvite{
				Email: fmt.Sprintf("invitee-%d@example.com", i),
				Owner: false,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.UserInvite.Id)
	}

	resp1, err := u.Store.ListUserInvites(ctx, &frontendv1.ListUserInvitesRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.UserInvites, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUserInvites(ctx, &frontendv1.ListUserInvitesRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.UserInvites, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, u := range resp1.UserInvites {
		allIDs = append(allIDs, u.Id)
	}
	for _, u := range resp2.UserInvites {
		allIDs = append(allIDs, u.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
