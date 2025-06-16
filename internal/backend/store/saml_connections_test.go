package store

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/consoletest"
	"github.com/tesseral-labs/tesseral/internal/dbconntest"
)

func TestCreateSAMLConnection(t *testing.T) {
	pool := dbconntest.Open(t)
	console := consoletest.New(t, pool)

	store := New(NewStoreParams{
		DB:               pool,
		S3:               s3.NewFromConfig(*aws.NewConfig()),
		DogfoodProjectID: console.DogfoodProjectID,
	})

	project := console.CreateProject(t, consoletest.ProjectParams{
		Name:          "test",
		LoginWithSaml: true,
	})

	ctx := context.Background()
	ctx = authn.NewBackendAPIKeyContext(ctx, &authn.BackendAPIKeyContextData{
		ProjectID:       project.ProjectID,
		BackendAPIKeyID: project.BackendAPIKey,
	})

	{
		organization := console.CreateOrganization(t, consoletest.OrganizationParams{
			ProjectID:     project.ProjectID,
			Name:          "test",
			LoginWithSaml: true,
		})
		_, err := store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
			SamlConnection: &backendv1.SAMLConnection{
				SpAcsUrl:       "https://example.com/saml/acs",
				SpEntityId:     "https://example.com/saml/sp",
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
				OrganizationId: organization.OrganizationID,
			},
		})
		require.NoError(t, err, "failed to create SAML connection")
	}

	{
		organization := console.CreateOrganization(t, consoletest.OrganizationParams{
			ProjectID:     project.ProjectID,
			Name:          "test",
			LoginWithSaml: false,
		})
		_, err := store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
			SamlConnection: &backendv1.SAMLConnection{
				SpAcsUrl:       "https://example.com/saml/acs",
				SpEntityId:     "https://example.com/saml/sp",
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
				OrganizationId: organization.OrganizationID,
			},
		})

		var connectErr *connect.Error
		require.ErrorAs(t, err, &connectErr)
		require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code(), "expected error when creating SAML connection for organization without SAML enabled")
	}
}
