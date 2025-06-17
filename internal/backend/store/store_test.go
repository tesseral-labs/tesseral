package store_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetestutil"
)

var (
	deps *packageDeps
)

type packageDeps struct {
	store   *store.Store
	console *storetestutil.Console
}

func TestMain(m *testing.M) {
	packageDeps, cleanup := initPackageDeps()
	defer cleanup()

	deps = packageDeps
	exitCode := m.Run()
	cleanup()

	os.Exit(exitCode)
}

func initPackageDeps() (*packageDeps, func()) {
	db, cleanupDB := storetestutil.NewDB()
	kms, cleanupKms := storetestutil.NewKMS()
	s3, cleanupS3 := storetestutil.NewS3()

	console := storetestutil.NewConsole(db, kms)
	store := store.New(store.NewStoreParams{
		DB:                        db,
		S3:                        s3,
		KMS:                       console.KMS.Client,
		SessionSigningKeyKmsKeyID: console.KMS.SessionSigningKeyID,
		DogfoodProjectID:          console.DogfoodProjectID,
		ConsoleDomain:             console.ConsoleDomain,
	})

	cleanup := func() {
		cleanupS3()
		cleanupKms()
		cleanupDB()
	}

	return &packageDeps{
		store:   store,
		console: console,
	}, cleanup
}

type tester struct {
	Store   *store.Store
	console *storetestutil.Console
	Project storetestutil.Project
}

func Init(t *testing.T) (context.Context, *tester) {
	t.Helper()

	project := deps.console.NewProject(t)
	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	return ctx, &tester{
		Store:   deps.store,
		console: deps.console,
		Project: project,
	}
}

func (s *tester) NewOrganization(t *testing.T, organization *backendv1.Organization) storetestutil.Organization {
	return s.console.NewOrganization(t, storetestutil.OrganizationParams{
		Project:      s.Project,
		Organization: organization,
	})
}
