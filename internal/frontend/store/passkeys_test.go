package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetPasskeyOptions(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	resp, err := u.Store.GetPasskeyOptions(ctx, &frontendv1.GetPasskeyOptionsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.RpId)
	require.NotEmpty(t, resp.RpName)
	require.NotEmpty(t, resp.UserId)
	require.NotEmpty(t, resp.UserDisplayName)
}

func TestListMyPasskeys_Empty(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	resp, err := u.Store.ListMyPasskeys(ctx, &frontendv1.ListMyPasskeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Empty(t, resp.Passkeys)
}

func TestListMyPasskeys_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	// Create 3 passkeys
	for range 15 {
		_, err := u.Environment.DB.Exec(ctx, `
		INSERT INTO passkeys (id, user_id, credential_id, public_key, aaguid, rp_id)
    		VALUES (gen_random_uuid(), $1, ''::bytea, ''::bytea, '', '')
		`,
			authn.UserID(ctx))
		require.NoError(t, err)
	}

	resp, err := u.Store.ListMyPasskeys(ctx, &frontendv1.ListMyPasskeysRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Passkeys, 10)
	require.NotEmpty(t, resp.NextPageToken)

	resp, err = u.Store.ListMyPasskeys(ctx, &frontendv1.ListMyPasskeysRequest{
		PageToken: resp.NextPageToken,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Passkeys, 5)
	require.Empty(t, resp.NextPageToken)
}

func TestDeleteMyPasskey_NotFound(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Org",
	})

	_, err := u.Store.DeleteMyPasskey(ctx, &frontendv1.DeleteMyPasskeyRequest{
		Id: idformat.Passkey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}
