package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestUpdateMe(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	beforeResp, err := u.Store.GetUser(ctx, &frontendv1.GetUserRequest{
		Id: idformat.User.Format(authn.UserID(ctx)),
	})
	require.NoError(t, err)
	require.NotNil(t, beforeResp)
	require.Empty(t, beforeResp.User.GetDisplayName())
	email := beforeResp.User.GetEmail()
	require.NotEmpty(t, email)

	newDisplayName := "Updated User"
	resp, err := u.Store.UpdateMe(ctx, &frontendv1.UpdateMeRequest{
		User: &frontendv1.User{
			DisplayName: &newDisplayName,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, newDisplayName, resp.User.GetDisplayName())
	require.Equal(t, email, resp.User.GetEmail())

	afterResp, err := u.Store.GetUser(ctx, &frontendv1.GetUserRequest{
		Id: idformat.User.Format(authn.UserID(ctx)),
	})
	require.NoError(t, err)
	require.NotNil(t, afterResp)
	require.Equal(t, newDisplayName, afterResp.User.GetDisplayName())
	require.Equal(t, email, afterResp.User.GetEmail())
}
