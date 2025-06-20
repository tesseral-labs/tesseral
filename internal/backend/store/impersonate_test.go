package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestCreateUserImpersonationToken(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	userID := u.NewUser(t, orgID, "test@example.com")

	res, err := u.Store.CreateUserImpersonationToken(ctx, &backendv1.CreateUserImpersonationTokenRequest{
		UserImpersonationToken: &backendv1.UserImpersonationToken{
			ImpersonatedId: userID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.UserImpersonationToken)
	require.NotEmpty(t, res.UserImpersonationToken.Id)
	require.True(t, res.UserImpersonationToken.CreateTime.AsTime().Before(time.Now()))
	require.True(t, res.UserImpersonationToken.ExpireTime.AsTime().After(time.Now()))
	require.Equal(t, userID, res.UserImpersonationToken.ImpersonatedId)
	require.NotEmpty(t, authn.GetContextData(ctx).DogfoodSession.UserID, res.UserImpersonationToken.ImpersonatorId)
	require.NotEmpty(t, res.UserImpersonationToken.SecretToken)
}
