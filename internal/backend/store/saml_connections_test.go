package store_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store"
	"github.com/tesseral-labs/tesseral/internal/consoletest"
	"github.com/tesseral-labs/tesseral/internal/dbconntest"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateSAMLConnection(t *testing.T) {
	pool := dbconntest.Open(t)
	console := consoletest.New(t, pool)

	store := store.New(store.NewStoreParams{
		DB:               pool,
		S3:               s3.NewFromConfig(*aws.NewConfig()),
		DogfoodProjectID: console.DogfoodProjectID,
	})

	project := console.CreateProject(t)

	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	{
		organization := console.CreateOrganization(t, consoletest.OrganizationParams{
			Project: project,
			Organization: &backendv1.Organization{
				DisplayName:   "test",
				LogInWithSaml: refOrNil(true),
			},
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
			Project: project,
			Organization: &backendv1.Organization{
				DisplayName:   "test",
				LogInWithSaml: refOrNil(false), // SAML not enabled
			},
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

func refOrNil[T comparable](t T) *T {
	var z T
	if t == z {
		return nil
	}
	return &t
}
