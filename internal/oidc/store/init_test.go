package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetOIDCConnectionInitData(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	// Create a new organization with OIDC enabled
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	oidcConnectionUUID := uuid.New()
	oidcConnectionID := idformat.OIDCConnection.Format(oidcConnectionUUID)

	_, err = u.Environment.DB.Exec(t.Context(), `
	INSERT INTO oidc_connections (id, organization_id, configuration_url, client_id, is_primary)
	VALUES ($1::uuid, $2::uuid, 'https://accounts.google.com/.well-known/openid-configuration', 'client-id', true);
	`, oidcConnectionUUID, organizationUUID)
	require.NoError(t, err)

	// Call the GetOIDCConnectionInitData method
	res, err := u.Store.GetOIDCConnectionInitData(ctx, oidcConnectionID)
	require.NoError(t, err)

	// Validate the response
	require.NotNil(t, res)
	require.Contains(t, res.AuthorizationURL, "https://accounts.google.com/o/oauth2/v2/auth")
	require.NotEmpty(t, res.State)
}
