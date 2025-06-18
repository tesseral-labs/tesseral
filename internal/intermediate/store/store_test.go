package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetesting"
)

var (
	environment *storetesting.Environment
)

func TestMain(m *testing.M) {
	testEnvironment, cleanup := storetesting.NewEnvironment()
	defer cleanup()

	environment = testEnvironment
	m.Run()
}

type testUtil struct {
	Store       *Store
	Environment *storetesting.Environment
	ProjectID   string
}

func newTestUtil(t *testing.T) (context.Context, *testUtil) {
	store := New(NewStoreParams{
		DB:                        environment.DB,
		S3:                        environment.S3,
		KMS:                       environment.KMS.Client,
		SessionSigningKeyKmsKeyID: environment.KMS.SessionSigningKeyID,
		DogfoodProjectID:          environment.DogfoodProjectID,
		ConsoleDomain:             environment.ConsoleDomain,
	})
	projectID, _ := environment.NewProject(t)

	intermediateSession := &intermediatev1.IntermediateSession{
		Id:        idformat.IntermediateSession.Format(uuid.New()),
		ProjectId: projectID,
	}
	ctx := authn.NewContext(t.Context(), intermediateSession, projectID)

	// Create a true intermediate session so it's in the database
	resp, err := store.CreateIntermediateSession(ctx, &intermediatev1.CreateIntermediateSessionRequest{})
	if err != nil {
		t.Fatalf("failed to create intermediate session: %v", err)
	}
	secretToken := resp.IntermediateSessionSecretToken
	intermediateSession, err = store.AuthenticateIntermediateSession(ctx, projectID, secretToken)
	if err != nil {
		t.Fatalf("failed to create intermediate session: %v", err)
	}

	intermediateSessionUUID, err := idformat.IntermediateSession.Parse(intermediateSession.Id)
	if err != nil {
		t.Fatalf("failed to parse intermediate session id: %v", err)
	}

	email := fmt.Sprintf("%s@%s", uuid.New(), environment.ConsoleDomain)
	_, err = environment.DB.Exec(t.Context(), `
UPDATE intermediate_sessions
SET 
	email_verification_challenge_completed = TRUE,
	email = $2
WHERE id = $1::uuid;
`,
		uuid.UUID(intermediateSessionUUID).String(),
		email)
	if err != nil {
		t.Fatalf("failed to update intermediate session email verified: %v", err)
	}
	intermediateSession.Email = email
	intermediateSession.EmailVerified = true

	// Recreate the context with the new intermediate session
	ctx = authn.NewContext(t.Context(), intermediateSession, projectID)

	return ctx, &testUtil{
		Store:       store,
		Environment: environment,
		ProjectID:   projectID,
	}
}

func (u *testUtil) NewOrganization(t *testing.T, organization *backendv1.Organization) string {
	return u.Environment.NewOrganization(t, u.ProjectID, organization)
}
