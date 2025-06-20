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

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{DisplayName: "test"})

	resp, err := u.Store.CreateUser(ctx, &backendv1.CreateUserRequest{
		User: &backendv1.User{
			OrganizationId: orgID,
			Email:          "test@example.com",
			Owner:          refOrNil(false),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.User)
	require.Equal(t, orgID, resp.User.OrganizationId)
	require.Equal(t, "test@example.com", resp.User.Email)
	require.NotEmpty(t, resp.User.Id)
	require.NotEmpty(t, resp.User.CreateTime)
	require.NotEmpty(t, resp.User.UpdateTime)
	require.False(t, resp.User.GetOwner())
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_USER, "tesseral.users.create")
}

func TestCreateUser_OrgNotFound(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.CreateUser(ctx, &backendv1.CreateUserRequest{
		User: &backendv1.User{
			OrganizationId: idformat.Organization.Format(uuid.New()),
			Email:          "test@example.com",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestGetUser_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{DisplayName: "test"})
	userID := u.NewUser(t, orgID, "test@example.com")

	getResp, err := u.Store.GetUser(ctx, &backendv1.GetUserRequest{Id: userID})
	require.NoError(t, err)
	require.NotNil(t, getResp.User)
	require.Equal(t, userID, getResp.User.Id)
	require.Equal(t, orgID, getResp.User.OrganizationId)
	require.Equal(t, "test@example.com", getResp.User.Email)
	require.NotEmpty(t, getResp.User.CreateTime)
	require.NotEmpty(t, getResp.User.UpdateTime)
	require.False(t, getResp.User.GetOwner())
}

func TestGetUser_NotFound(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetUser(ctx, &backendv1.GetUserRequest{Id: idformat.User.Format(uuid.New())})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateUser_UpdatesFields(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{DisplayName: "test"})
	userID := u.NewUser(t, orgID, "test@example.com")

	updateResp, err := u.Store.UpdateUser(ctx, &backendv1.UpdateUserRequest{
		Id: userID,
		User: &backendv1.User{
			Email:             "test-updated@example.com",
			DisplayName:       refOrNil("Updated Name"),
			Owner:             refOrNil(true),
			GoogleUserId:      refOrNil("google-user-id-123"),
			MicrosoftUserId:   refOrNil("microsoft-user-id-456"),
			GithubUserId:      refOrNil("apple-user-id-789"),
			ProfilePictureUrl: refOrNil("https://example.com/profile.jpg"),
		},
	})
	require.NoError(t, err)
	require.Equal(t, "test-updated@example.com", updateResp.User.Email)
	require.Equal(t, "Updated Name", updateResp.User.GetDisplayName())
	require.True(t, updateResp.User.GetOwner())
	require.Equal(t, "google-user-id-123", updateResp.User.GetGoogleUserId())
	require.Equal(t, "microsoft-user-id-456", updateResp.User.GetMicrosoftUserId())
	require.Equal(t, "apple-user-id-789", updateResp.User.GetGithubUserId())
	require.Equal(t, "https://example.com/profile.jpg", updateResp.User.GetProfilePictureUrl())
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_USER, "tesseral.users.update")
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{DisplayName: "test"})
	userID := u.NewUser(t, orgID, "user5@example.com")

	_, err := u.Store.DeleteUser(ctx, &backendv1.DeleteUserRequest{Id: userID})
	require.NoError(t, err)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_USER, "tesseral.users.delete")

	_, err = u.Store.GetUser(ctx, &backendv1.GetUserRequest{Id: userID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListUsers_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{DisplayName: "test"})

	var ids []string
	for i := range 3 {
		email := fmt.Sprintf("user%d@example.com", i)
		id := u.NewUser(t, orgID, email)
		ids = append(ids, id)
	}

	listResp, err := u.Store.ListUsers(ctx, &backendv1.ListUsersRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Users, 3)

	var respIds []string
	for _, user := range listResp.Users {
		respIds = append(respIds, user.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListUsers_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{DisplayName: "test"})

	var createdIDs []string
	for i := range 15 {
		email := fmt.Sprintf("user%d@example.com", i)
		id := u.NewUser(t, orgID, email)
		createdIDs = append(createdIDs, id)
	}

	resp1, err := u.Store.ListUsers(ctx, &backendv1.ListUsersRequest{OrganizationId: orgID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.Users, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListUsers(ctx, &backendv1.ListUsersRequest{OrganizationId: orgID, PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.Users, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, u := range resp1.Users {
		allIDs = append(allIDs, u.Id)
	}
	for _, u := range resp2.Users {
		allIDs = append(allIDs, u.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
