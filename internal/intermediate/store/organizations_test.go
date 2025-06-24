package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/emailaddr"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestListSAMLOrganizations_SAMLEnabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	res, err := u.Store.ListSAMLOrganizations(ctx, &intermediatev1.ListSAMLOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	require.NoError(t, err)
	require.Empty(t, res.Organizations)
}

func TestListSAMLOrganizations_SAMLDisabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	_, err = u.Environment.DB.Exec(t.Context(), `
	UPDATE projects
	SET log_in_with_saml = false
	WHERE id = $1::uuid;
`,
		uuid.UUID(projectUUID).String())
	require.NoError(t, err)

	_, err = u.Store.ListSAMLOrganizations(ctx, &intermediatev1.ListSAMLOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestListSAMLOrganizations_ForActiveDomain(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	organizationID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	domain, err := emailaddr.Parse(authn.IntermediateSession(ctx).Email)
	require.NoError(t, err)

	// Create the organization domain
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO organization_domains (id, organization_id, domain)
VALUES (gen_random_uuid(), $1::uuid, $2);
`,
		uuid.UUID(organizationUUID).String(),
		domain)
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

	res, err := u.Store.ListSAMLOrganizations(ctx, &intermediatev1.ListSAMLOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	require.NoError(t, err)
	require.Len(t, res.Organizations, 1)
	require.Equal(t, organizationID, res.Organizations[0].Id)
	require.Equal(t, idformat.SAMLConnection.Format(samlConnectionID), res.Organizations[0].PrimarySamlConnectionId)
}

func TestListSAMLOrganizations_ForDifferentDomain(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	organizationID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	// Create the organization domain
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO organization_domains (id, organization_id, domain)
VALUES (gen_random_uuid(), $1::uuid, 'example.com');
`,
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

	res, err := u.Store.ListSAMLOrganizations(ctx, &intermediatev1.ListSAMLOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	require.NoError(t, err)
	require.Empty(t, res.Organizations)
}

func TestListOIDCOrganizations_OIDCEnabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	res, err := u.Store.ListOIDCOrganizations(ctx, &intermediatev1.ListOIDCOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	require.NoError(t, err)
	require.Empty(t, res.Organizations)
}

func TestListOIDCOrganizations_OIDCDisabled(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	_, err = u.Environment.DB.Exec(t.Context(), `
	UPDATE projects
	SET log_in_with_oidc = false
	WHERE id = $1::uuid;
`,
		uuid.UUID(projectUUID).String())
	require.NoError(t, err)

	_, err = u.Store.ListOIDCOrganizations(ctx, &intermediatev1.ListOIDCOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestListOIDCOrganizations_ForActiveDomain(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	organizationID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:   "Test OIDC Organization",
		LogInWithOidc: refOrNil(true),
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	domain, err := emailaddr.Parse(authn.IntermediateSession(ctx).Email)
	require.NoError(t, err)

	// Create the organization domain
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO organization_domains (id, organization_id, domain)
VALUES (gen_random_uuid(), $1::uuid, $2);
`,
		uuid.UUID(organizationUUID).String(),
		domain)
	require.NoError(t, err)

	// Create an OIDC connection for the organization
	oidcConnectionID := uuid.New()
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO oidc_connections (id, organization_id, is_primary, configuration_url, issuer, client_id)
VALUES ($1::uuid, $2::uuid, true, 'https://issuer.example.com/.well-known/openid-configuration', 'https://issuer.example.com', 'client-id');
`,
		oidcConnectionID.String(),
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	res, err := u.Store.ListOIDCOrganizations(ctx, &intermediatev1.ListOIDCOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	require.NoError(t, err)
	require.Len(t, res.Organizations, 1)
	require.Equal(t, organizationID, res.Organizations[0].Id)
	require.Equal(t, idformat.OIDCConnection.Format(oidcConnectionID), res.Organizations[0].PrimaryOidcConnectionId)
}

func TestListOIDCOrganizations_ForDifferentDomain(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	organizationID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName:   "Test OIDC Organization",
		LogInWithOidc: refOrNil(true),
	})
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	// Create the organization domain
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO organization_domains (id, organization_id, domain)
VALUES (gen_random_uuid(), $1::uuid, 'otherdomain.com');
`,
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	// Create an OIDC connection for the organization
	oidcConnectionID := uuid.New()
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO oidc_connections (id, organization_id, is_primary, configuration_url, issuer, client_id)
VALUES ($1::uuid, $2::uuid, true, 'https://issuer.example.com/.well-known/openid-configuration', 'https://issuer.example.com', 'client-id');
`,
		oidcConnectionID.String(),
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	res, err := u.Store.ListOIDCOrganizations(ctx, &intermediatev1.ListOIDCOrganizationsRequest{
		Email: authn.IntermediateSession(ctx).Email,
	})
	require.NoError(t, err)
	require.Empty(t, res.Organizations)
}
