package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestListOrganizations_NoActiveOrganizations(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_ = u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})

	res, err := u.Store.ListOrganizations(ctx, &intermediatev1.ListOrganizationsRequest{})
	require.NoError(t, err)
	require.Empty(t, res.Organizations)
}

func TestListOrganizations_ActiveByEmail(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	organizationID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	email := authn.IntermediateSession(ctx).Email

	// Create a user in the organization with the same email
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
VALUES (gen_random_uuid(), $1, crypt('password', gen_salt('bf', 14)), $2, true);
`,
		email,
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	res, err := u.Store.ListOrganizations(ctx, &intermediatev1.ListOrganizationsRequest{})
	require.NoError(t, err)
	require.Len(t, res.Organizations, 1)
	require.Equal(t, organizationID, res.Organizations[0].Id)
}

func TestListOrganizations_ActiveByEmailWithSamlConnection(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	organizationID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	email := authn.IntermediateSession(ctx).Email

	// Create a user in the organization with the same email
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
VALUES (gen_random_uuid(), $1, crypt('password', gen_salt('bf', 14)), $2, true);
`,
		email,
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	// Create a SAML connection for the organization
	samlConnectionID := uuid.New()
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO saml_connections (id, organization_id, is_primary, idp_redirect_url, idp_x509_certificate, idp_entity_id)
VALUES ($1::uuid, $2::uuid, true, 'https://idp.example.com/saml/redirect', ''::bytea, 'https://idp.example.com/saml/idp');
`,
		samlConnectionID.String(),
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	res, err := u.Store.ListOrganizations(ctx, &intermediatev1.ListOrganizationsRequest{})
	require.NoError(t, err)
	require.Len(t, res.Organizations, 1)
	require.Equal(t, organizationID, res.Organizations[0].Id)
	require.Equal(t, idformat.SAMLConnection.Format(samlConnectionID), res.Organizations[0].PrimarySamlConnectionId)
}
