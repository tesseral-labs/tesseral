package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/totp"
)

func TestGetAuthenticatorAppOptions_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.OtpauthUri)
}

func TestRegisterAuthenticatorApp_Success(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)

	secret, err := u.Store.getUserAuthenticatorAppChallengeSecret(ctx)
	require.NoError(t, err)

	key := totp.Key{Secret: secret}
	code := key.Gen(time.Now())

	resp, err := u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: code,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.RecoveryCodes)
}

func TestRegisterAuthenticatorApp_InvalidCode(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)

	resp, err := u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: "123456",
	})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestRegisterAuthenticatorApp_ExpiredCode(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)

	secret, err := u.Store.getUserAuthenticatorAppChallengeSecret(ctx)
	require.NoError(t, err)

	key := totp.Key{Secret: secret}
	code := key.Gen(time.Now().Add(-time.Hour))

	resp, err := u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: code,
	})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestRegisterAuthenticatorApp_AlreadyRegistered(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	_, err := u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)

	secret, err := u.Store.getUserAuthenticatorAppChallengeSecret(ctx)
	require.NoError(t, err)

	key := totp.Key{Secret: secret}
	code := key.Gen(time.Now())

	_, err = u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: code,
	})
	require.NoError(t, err)

	_, err = u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)

	secret, err = u.Store.getUserAuthenticatorAppChallengeSecret(ctx)
	require.NoError(t, err)

	key = totp.Key{Secret: secret}
	code = key.Gen(time.Now())

	_, err = u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: code,
	})
	require.NoError(t, err)
}
