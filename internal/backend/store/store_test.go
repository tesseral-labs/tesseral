package store_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	"github.com/tesseral-labs/tesseral/internal/backend/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetestutil"
)

var (
	storeT *tester
)

type tester struct {
	*store.Store
	console *storetestutil.Console
}

func TestMain(m *testing.M) {
	tester, cleanup := newTester()
	defer cleanup()

	storeT = tester
	exitCode := m.Run()
	cleanup()

	os.Exit(exitCode)
}

func newTester() (*tester, func()) {
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

	return &tester{
		Store:   store,
		console: console,
	}, cleanup
}

func (s *tester) Init(t *testing.T) (context.Context, storetestutil.Project) {
	t.Helper()

	project := s.console.NewProject(t)
	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	return ctx, project
}

func (s *tester) NewOrganization(t *testing.T, params storetestutil.OrganizationParams) storetestutil.Organization {
	return s.console.NewOrganization(t, params)
}
