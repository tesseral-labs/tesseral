package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetSessionSigningKeyPublicKey(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	var sessionSigningKeyID uuid.UUID
	err = u.Environment.DB.QueryRow(ctx, `
		SELECT id
		FROM session_signing_keys
		WHERE project_id = $1
		LIMIT 1
	`, projectUUID).Scan(&sessionSigningKeyID)
	require.NoError(t, err)

	publicKey, err := u.Store.GetSessionSigningKeyPublicKey(ctx, idformat.SessionSigningKey.Format(sessionSigningKeyID))
	require.NoError(t, err)
	require.NotNil(t, publicKey)
}
