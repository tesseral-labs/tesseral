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

func TestCreateUserInvite_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
		UserInvite: &backendv1.UserInvite{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          false,
		},
		SendEmail: false,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.UserInvite)
	require.Equal(t, orgID, resp.UserInvite.OrganizationId)
	require.Equal(t, "test@example.com", resp.UserInvite.Email)
	require.NotEmpty(t, resp.UserInvite.Id)
	require.NotEmpty(t, resp.UserInvite.CreateTime)
	require.NotEmpty(t, resp.UserInvite.UpdateTime)
}

func TestCreateUserInvite_DuplicateEmail(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
		UserInvite: &backendv1.UserInvite{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          false,
		},
		SendEmail: false,
	})
	require.NoError(t, err)

	// Second invite with same email
	_, err = u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
		UserInvite: &backendv1.UserInvite{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          false,
		},
		SendEmail: false,
	})
	require.NoError(t, err)
}

func TestCreateUserInvite_ExistingUser(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})
	_ = u.Environment.NewUser(t, orgID, &backendv1.User{
		Email: "test@example.com",
	})

	_, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
		UserInvite: &backendv1.UserInvite{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          false,
		},
		SendEmail: false,
	})
	require.Error(t, err)

}

func TestGetUserInvite_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	createResp, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
		UserInvite: &backendv1.UserInvite{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          false,
		},
		SendEmail: false,
	})
	require.NoError(t, err)
	inviteID := createResp.UserInvite.Id

	getResp, err := u.Store.GetUserInvite(ctx, &backendv1.GetUserInviteRequest{Id: inviteID})
	require.NoError(t, err)
	require.NotNil(t, getResp.UserInvite)
	require.Equal(t, inviteID, getResp.UserInvite.Id)
	require.Equal(t, orgID, getResp.UserInvite.OrganizationId)
	require.Equal(t, "test@example.com", getResp.UserInvite.Email)
	require.NotEmpty(t, getResp.UserInvite.CreateTime)
	require.NotEmpty(t, getResp.UserInvite.UpdateTime)
}

func TestGetUserInvite_NotFound(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetUserInvite(ctx, &backendv1.GetUserInviteRequest{
		Id: idformat.UserInvite.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestDeleteUserInvite_Deletes(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	createResp, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
		UserInvite: &backendv1.UserInvite{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          false,
		},
		SendEmail: false,
	})
	require.NoError(t, err)
	inviteID := createResp.UserInvite.Id

	_, err = u.Store.DeleteUserInvite(ctx, &backendv1.DeleteUserInviteRequest{Id: inviteID})
	require.NoError(t, err)

	_, err = u.Store.GetUserInvite(ctx, &backendv1.GetUserInviteRequest{Id: inviteID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListUserInvites_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	var ids []string
	for i := range 3 {
		email := fmt.Sprintf("user%d@example.com", i)
		resp, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
			UserInvite: &backendv1.UserInvite{
				OrganizationId: orgID,
				Email:          email,
				Owner:          false,
			},
			SendEmail: false,
		})
		require.NoError(t, err)
		ids = append(ids, resp.UserInvite.Id)
	}

	listResp, err := u.Store.ListUserInvites(ctx, &backendv1.ListUserInvitesRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.UserInvites, 3)

	var respIds []string
	for _, invite := range listResp.UserInvites {
		respIds = append(respIds, invite.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListUserInvites_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	var createdIDs []string
	for i := range 15 {
		email := fmt.Sprintf("user%d@example.com", i)
		resp, err := u.Store.CreateUserInvite(ctx, &backendv1.CreateUserInviteRequest{
			UserInvite: &backendv1.UserInvite{
				OrganizationId: orgID,
				Email:          email,
				Owner:          false,
			},
			SendEmail: false,
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.UserInvite.Id)
	}

	resp1, err := u.Store.ListUserInvites(ctx, &backendv1.ListUserInvitesRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.UserInvites, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUserInvites(ctx, &backendv1.ListUserInvitesRequest{OrganizationId: orgID, PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.UserInvites, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.UserInvites {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.UserInvites {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
