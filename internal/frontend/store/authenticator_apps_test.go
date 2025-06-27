package store

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
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
	code := genTOTPCode(key, time.Now())

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
	code := genTOTPCode(key, time.Now().Add(-time.Hour))

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
	code := genTOTPCode(key, time.Now())

	_, err = u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: code,
	})
	require.NoError(t, err)

	_, err = u.Store.GetAuthenticatorAppOptions(ctx)
	require.NoError(t, err)

	secret, err = u.Store.getUserAuthenticatorAppChallengeSecret(ctx)
	require.NoError(t, err)

	key = totp.Key{Secret: secret}
	code = genTOTPCode(key, time.Now())

	_, err = u.Store.RegisterAuthenticatorApp(ctx, &frontendv1.RegisterAuthenticatorAppRequest{
		TotpCode: code,
	})
	require.NoError(t, err)
}

func genTOTPCode(key totp.Key, now time.Time) string {
	counter := now.Unix() / 30

	mac := hmac.New(sha1.New, key.Secret)
	_ = binary.Write(mac, binary.BigEndian, counter)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	return fmt.Sprintf("%06d", value%1_000_000)
}
