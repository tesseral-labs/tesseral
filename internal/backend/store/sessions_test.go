package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetSession_Success(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test-org",
	})
	userID := u.NewUser(t, orgID, "test@example.com")
	sessionID, _ := u.Environment.NewSession(t, userID)

	resp, err := u.Store.GetSession(ctx, &backendv1.GetSessionRequest{Id: sessionID})
	require.NoError(t, err)
	require.NotNil(t, resp.Session)
	require.Equal(t, sessionID, resp.Session.Id)
}

func TestGetSession_NotFound(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	_, err := u.Store.GetSession(ctx, &backendv1.GetSessionRequest{
		Id: idformat.Session.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListSessions(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test-org",
	})
	userID := u.NewUser(t, orgID, "test@example.com")

	var sessionIDs []string
	for range 12 {
		sessionID, _ := u.Environment.NewSession(t, userID)
		sessionIDs = append(sessionIDs, sessionID)
	}

	resp1, err := u.Store.ListSessions(ctx, &backendv1.ListSessionsRequest{UserId: userID})
	require.NoError(t, err)
	require.Len(t, resp1.Sessions, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListSessions(ctx, &backendv1.ListSessionsRequest{UserId: userID, PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.Len(t, resp2.Sessions, 2)

	var gotIDs []string
	for _, s := range append(resp1.Sessions, resp2.Sessions...) {
		gotIDs = append(gotIDs, s.Id)
	}
	require.ElementsMatch(t, sessionIDs, gotIDs)
}

func TestListSessions_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.ListSessions(ctx, &backendv1.ListSessionsRequest{
		UserId: idformat.User.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}
