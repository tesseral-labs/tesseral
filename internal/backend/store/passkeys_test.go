package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (u *testUtil) newPasskey(t *testing.T, userID string) string {
	passkeyID := uuid.New()
	userUUID, err := idformat.User.Parse(userID)
	require.NoError(t, err)

	_, err = u.Environment.DB.Exec(t.Context(), `
	INSERT INTO passkeys (id, user_id, credential_id, public_key, aaguid, rp_id)
	VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6)
	`,
		passkeyID.String(),
		uuid.UUID(userUUID).String(),
		[]byte{},
		[]byte{},
		"",
		"example.com",
	)
	require.NoError(t, err)

	return idformat.Passkey.Format(passkeyID)
}

func TestGetPasskey_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	userID := u.NewUser(t, orgID, "test@example.com")
	passkeyID := u.newPasskey(t, userID)

	resp, err := u.Store.GetPasskey(ctx, &backendv1.GetPasskeyRequest{Id: passkeyID})
	require.NoError(t, err)
	require.NotNil(t, resp.Passkey)
	require.Equal(t, passkeyID, resp.Passkey.Id)
}

func TestGetPasskey_DoesNotExist(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetPasskey(ctx, &backendv1.GetPasskeyRequest{
		Id: idformat.Passkey.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdatePasskey_Disable(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	userID := u.NewUser(t, orgID, "test@example.com")
	passkeyID := u.newPasskey(t, userID)

	resp, err := u.Store.UpdatePasskey(ctx, &backendv1.UpdatePasskeyRequest{
		Id: passkeyID,
		Passkey: &backendv1.Passkey{
			Disabled: refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Passkey)
	require.True(t, resp.Passkey.GetDisabled())
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_PASSKEY, "tesseral.passkeys.update")
}

func TestListPasskeys_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	userID := u.NewUser(t, orgID, "test@example.com")

	var ids []string
	for range 3 {
		id := u.newPasskey(t, userID)
		ids = append(ids, id)
	}

	resp, err := u.Store.ListPasskeys(ctx, &backendv1.ListPasskeysRequest{UserId: userID})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Passkeys, 3)

	var respIds []string
	for _, pk := range resp.Passkeys {
		respIds = append(respIds, pk.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListPasskeys_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	userID := u.NewUser(t, orgID, "test@example.com")

	var createdIDs []string
	for range 15 {
		id := u.newPasskey(t, userID)
		createdIDs = append(createdIDs, id)
	}

	resp1, err := u.Store.ListPasskeys(ctx, &backendv1.ListPasskeysRequest{UserId: userID})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.Passkeys, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListPasskeys(ctx, &backendv1.ListPasskeysRequest{UserId: userID, PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.Passkeys, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.Passkeys {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.Passkeys {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}

func TestDeletePasskey_RemovesKey(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	userID := u.NewUser(t, orgID, "test@example.com")
	passkeyID := u.newPasskey(t, userID)

	_, err := u.Store.DeletePasskey(ctx, &backendv1.DeletePasskeyRequest{Id: passkeyID})
	require.NoError(t, err)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_PASSKEY, "tesseral.passkeys.delete")

	_, err = u.Store.GetPasskey(ctx, &backendv1.GetPasskeyRequest{Id: passkeyID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}
