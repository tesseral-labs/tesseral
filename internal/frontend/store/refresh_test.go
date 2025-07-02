package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateRefreshAuditLogEvent_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})

	userID := authn.UserID(ctx)
	_, refreshToken := u.Environment.NewSession(t, idformat.User.Format(userID))

	accessToken, err := u.Common.IssueAccessToken(ctx, authn.ProjectID(ctx), refreshToken)
	require.NoError(t, err)

	err = u.Store.CreateRefreshAuditLogEvent(ctx, accessToken)
	require.NoError(t, err)
}

func TestCreateRefreshAuditLogEvent_InvalidToken(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})

	err := u.Store.CreateRefreshAuditLogEvent(ctx, "header.body.signature")
	require.Error(t, err)
}
