package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestWhoami(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	var expectedEmail string
	err := u.Environment.DB.QueryRow(t.Context(), `
		SELECT email
		FROM users
		WHERE id = $1::uuid
	`, authn.UserID(ctx)).Scan(&expectedEmail)
	require.NoError(t, err)

	resp, err := u.Store.Whoami(ctx, &frontendv1.WhoamiRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp.User)
	require.Equal(t, idformat.User.Format(authn.UserID(ctx)), resp.User.Id)
	require.Equal(t, expectedEmail, resp.User.Email)
}
