package store

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetestutil"
)

var (
	deps *storetestutil.StoreDependencies
)

func TestMain(m *testing.M) {
	storeDeps, cleanup := storetestutil.NewStoreDependencies()
	defer cleanup()

	deps = storeDeps
	exitCode := m.Run()
	cleanup()

	os.Exit(exitCode)
}

type tester struct {
	Store   *Store
	console *storetestutil.Console
	Project storetestutil.Project
}

func Init(t *testing.T) (context.Context, *tester) {
	store := New(NewStoreParams{
		DB:                        deps.DB,
		S3:                        deps.S3,
		KMS:                       deps.KMS.Client,
		SessionSigningKeyKmsKeyID: deps.KMS.SessionSigningKeyID,
		DogfoodProjectID:          deps.Console.DogfoodProjectID,
		ConsoleDomain:             deps.Console.ConsoleDomain,
	})
	project := deps.Console.NewProject(t)
	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	return ctx, &tester{
		Store:   store,
		console: deps.Console,
		Project: project,
	}
}

func (s *tester) NewOrganization(t *testing.T, organization *backendv1.Organization) storetestutil.Organization {
	return s.console.NewOrganization(t, storetestutil.OrganizationParams{
		Project:      s.Project,
		Organization: organization,
	})
}
