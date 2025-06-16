package store

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/dbconntest"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

var (
	// company1
	projectID      = idformat.Project.Format(uuid.MustParse("7abd6d2e-c314-456e-b9c5-bdbb62f0345f"))
	organizationID = idformat.Organization.Format(uuid.MustParse("3b1a04a1-0803-47af-bfdd-831349e2aac6"))
	userID         = idformat.User.Format(uuid.MustParse("125edb51-a832-445f-b45b-cba6acc0fb75"))
	sessionID      = idformat.Session.Format(uuid.Must(uuid.NewRandom()))
)

func TestCreateSAMLConnection(t *testing.T) {
	pool := dbconntest.Open(t)

	store := New(NewStoreParams{
		DB: pool,
		S3: s3.NewFromConfig(*aws.NewConfig()),
	})

	ctx := context.Background()
	ctx = authn.NewDogfoodSessionContext(ctx, authn.DogfoodSessionContextData{
		ProjectID: projectID,
		UserID:    userID,
		SessionID: sessionID,
	})

	_, err := store.CreateSAMLConnection(ctx, &backendv1.CreateSAMLConnectionRequest{
		SamlConnection: &backendv1.SAMLConnection{
			SpAcsUrl:       "https://example.com/saml/acs",
			SpEntityId:     "https://example.com/saml/sp",
			IdpRedirectUrl: "https://idp.example.com/saml/redirect",
			IdpEntityId:    "https://idp.example.com/saml/idp",
			OrganizationId: organizationID,
		},
	})
	require.NoError(t, err, "failed to create SAML connection")
}
