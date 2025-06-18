package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
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
	projectID, projectUserID := environment.NewProject(t)
	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: projectID,
		UserID:    projectUserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	return ctx, &testUtil{
		Store:       store,
		Environment: environment,
		ProjectID:   projectID,
	}
}

func (u *testUtil) NewOrganization(t *testing.T, organization *backendv1.Organization) string {
	return u.Environment.NewOrganization(t, u.ProjectID, organization)
}
