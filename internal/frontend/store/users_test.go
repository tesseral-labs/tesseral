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

func TestListUsers_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})
	organizationID := idformat.Organization.Format(authn.OrganizationID(ctx))

	createdIDs := []string{
		idformat.User.Format(authn.UserID(ctx)), // Include self user
	}
	for i := range 3 {
		userID := u.Environment.NewUser(t, organizationID, &backendv1.User{
			Email: fmt.Sprintf("user-%d@example.com", i),
		})
		createdIDs = append(createdIDs, userID)
	}

	resp, err := u.Store.ListUsers(ctx, &frontendv1.ListUsersRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Users, 3+1) // +1 for the self user
	require.Empty(t, resp.NextPageToken)

	var respIDs []string
	for _, user := range resp.Users {
		respIDs = append(respIDs, user.Id)
	}
	require.ElementsMatch(t, createdIDs, respIDs)
}

func TestListUsers_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})
	organizationID := idformat.Organization.Format(authn.OrganizationID(ctx))

	createdIDs := []string{
		idformat.User.Format(authn.UserID(ctx)), // Include self user
	}
	for i := range 15 {
		userID := u.Environment.NewUser(t, organizationID, &backendv1.User{
			Email: fmt.Sprintf("user-%d@example.com", i),
		})
		createdIDs = append(createdIDs, userID)
	}

	resp1, err := u.Store.ListUsers(ctx, &frontendv1.ListUsersRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.Users, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUsers(ctx, &frontendv1.ListUsersRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.Users, 5+1) // +1 for the self user
	require.Empty(t, resp2.NextPageToken)

	var respIDs []string
	for _, user := range resp1.Users {
		respIDs = append(respIDs, user.Id)
	}
	for _, user := range resp2.Users {
		respIDs = append(respIDs, user.Id)
	}
	require.ElementsMatch(t, createdIDs, respIDs)
}

func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})
	organizationID := idformat.Organization.Format(authn.OrganizationID(ctx))

	email := "test-123@example.com"
	userID := u.Environment.NewUser(t, organizationID, &backendv1.User{
		Email:           email,
		Owner:           refOrNil(true),
		GoogleUserId:    refOrNil("google-123"),
		MicrosoftUserId: refOrNil("microsoft-123"),
		GithubUserId:    refOrNil("github-123"),
		DisplayName:     refOrNil("Test User"),
	})

	resp, err := u.Store.GetUser(ctx, &frontendv1.GetUserRequest{Id: userID})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, userID, resp.User.Id)
	require.Equal(t, email, resp.User.Email)
	require.True(t, resp.User.GetOwner())
	require.Equal(t, "Test User", resp.User.GetDisplayName())
	require.Equal(t, "google-123", resp.User.GetGoogleUserId())
	require.Equal(t, "microsoft-123", resp.User.GetMicrosoftUserId())
	require.Equal(t, "github-123", resp.User.GetGithubUserId())
}

func TestGetUser_NotFound(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})

	_, err := u.Store.GetUser(ctx, &frontendv1.GetUserRequest{
		Id: idformat.User.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateUser_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})
	organizationID := idformat.Organization.Format(authn.OrganizationID(ctx))

	userID := u.Environment.NewUser(t, organizationID, &backendv1.User{
		Email: "test-123@example.com",
	})

	resp, err := u.Store.UpdateUser(ctx, &frontendv1.UpdateUserRequest{
		Id: userID,
		User: &frontendv1.User{
			Owner:       refOrNil(true),
			DisplayName: refOrNil("Updated Name"),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.User)
	require.True(t, resp.User.GetOwner())
	require.Equal(t, "Updated Name", resp.User.GetDisplayName())

	unchanged, err := u.Store.UpdateUser(ctx, &frontendv1.UpdateUserRequest{
		Id:   userID,
		User: &frontendv1.User{},
	})
	require.NoError(t, err)
	require.NotNil(t, unchanged.User)
	require.True(t, unchanged.User.GetOwner())
	require.Equal(t, "Updated Name", unchanged.User.GetDisplayName())
}

func TestDeleteUser_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})
	organizationID := idformat.Organization.Format(authn.OrganizationID(ctx))

	userID := u.Environment.NewUser(t, organizationID, &backendv1.User{
		Email: "test-123@example.com",
	})

	_, err := u.Store.DeleteUser(ctx, &frontendv1.DeleteUserRequest{Id: userID})
	require.NoError(t, err)

	_, err = u.Store.GetUser(ctx, &frontendv1.GetUserRequest{Id: userID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestDeleteUser_CannotDeleteSelf(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Test Org"})

	_, err := u.Store.DeleteUser(ctx, &frontendv1.DeleteUserRequest{
		Id: idformat.User.Format(authn.UserID(ctx)),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}
