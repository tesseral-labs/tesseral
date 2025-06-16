package store_test

import (
	"connectrpc.com/connect"
	"github.com/google/uuid"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetestutil"
)

func (s *BackendSuite) TestCreateSAMLConnection() {
	s.T().Parallel()

	project := s.CreateProject()

	s.Run("SAML enabled", func() {
		ctx := s.NewAuthContext(project)
		require := s.Require()

		organization := s.CreateOrganization(storetestutil.OrganizationParams{
			Project: project,
			Organization: &backendv1.Organization{
				DisplayName:   "test",
				LogInWithSaml: refOrNil(true),
			},
		})
		_, err := s.Store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
			SamlConnection: &backendv1.SAMLConnection{
				SpAcsUrl:       "https://example.com/saml/acs",
				SpEntityId:     "https://example.com/saml/sp",
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
				OrganizationId: organization.OrganizationID,
			},
		})
		require.NoError(err, "failed to create SAML connection")
	})

	s.Run("SAML disabled", func() {
		ctx := s.NewAuthContext(project)
		require := s.Require()

		organization := s.CreateOrganization(storetestutil.OrganizationParams{
			Project: project,
			Organization: &backendv1.Organization{
				DisplayName:   "test",
				LogInWithSaml: refOrNil(false), // SAML not enabled
			},
		})
		_, err := s.Store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
			SamlConnection: &backendv1.SAMLConnection{
				SpAcsUrl:       "https://example.com/saml/acs",
				SpEntityId:     "https://example.com/saml/sp",
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
				OrganizationId: organization.OrganizationID,
			},
		})

		var connectErr *connect.Error
		require.ErrorAs(err, &connectErr)
		require.Equal(connect.CodeFailedPrecondition, connectErr.Code(), "expected error when creating SAML connection for organization without SAML enabled")
	})
}

func (s *BackendSuite) TestGetSAMLConnection() {
	s.T().Parallel()

	project := s.CreateProject()

	s.Run("SAML connection exists", func() {
		ctx := s.NewAuthContext(project)
		require := s.Require()

		organization := s.CreateOrganization(storetestutil.OrganizationParams{
			Project: project,
			Organization: &backendv1.Organization{
				DisplayName:   "test",
				LogInWithSaml: refOrNil(true),
			},
		})
		samlConnection, err := s.Store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
			SamlConnection: &backendv1.SAMLConnection{
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
				OrganizationId: organization.OrganizationID,
			},
		})
		require.NoError(err, "failed to create SAML connection")

		res, err := s.Store.GetSAMLConnection(ctx, &backendv1.GetSAMLConnectionRequest{
			Id: samlConnection.SamlConnection.Id,
		})
		require.NoError(err, "failed to get SAML connection")
		require.NotNil(res.SamlConnection, "expected SAML connection to be returned")
		require.Equal(samlConnection.SamlConnection.Id, res.SamlConnection.Id, "expected SAML connection ID to match")
		require.NotEmpty(res.SamlConnection.SpAcsUrl, "expected SAML connection SP ACS URL to be set")
		require.NotEmpty(res.SamlConnection.SpEntityId, "expected SAML connection SP Entity ID to be set")
		require.Equal("https://idp.example.com/saml/redirect", res.SamlConnection.IdpRedirectUrl, "expected SAML connection IdP Redirect URL to match")
		require.Equal("https://idp.example.com/saml/idp", res.SamlConnection.IdpEntityId, "expected SAML connection IdP Entity ID to match")
		require.Equal(organization.OrganizationID, res.SamlConnection.OrganizationId, "expected SAML connection Organization ID to match")
		require.NotEmpty(res.SamlConnection.CreateTime, "expected SAML connection CreatedAt to be set")
		require.NotEmpty(res.SamlConnection.UpdateTime, "expected SAML connection UpdatedAt to be set")
	})

	s.Run("SAML connection does not exist", func() {
		ctx := s.NewAuthContext(project)
		require := s.Require()

		res, err := s.Store.GetSAMLConnection(ctx, &backendv1.GetSAMLConnectionRequest{
			Id: idformat.SAMLConnection.Format(uuid.New()),
		})

		var connectErr *connect.Error
		require.ErrorAs(err, &connectErr)
		require.Equal(connect.CodeNotFound, connectErr.Code(), "expected error when getting non-existent SAML connection")
		require.Nil(res, "expected no SAML connection to be returned")
	})
}

func refOrNil[T comparable](t T) *T {
	var z T
	if t == z {
		return nil
	}
	return &t
}
