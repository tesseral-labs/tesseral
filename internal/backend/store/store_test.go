package store

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetestutil"
)

var (
	deps *packageDeps
)

type packageDeps struct {
	db      *pgxpool.Pool
	s3      *s3.Client
	kms     *storetestutil.KMS
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

	cleanup := func() {
		cleanupS3()
		cleanupKms()
		cleanupDB()
	}

	return &packageDeps{
		db:      db,
		s3:      s3,
		kms:     kms,
		console: console,
	}, cleanup
}

type tester struct {
	Store   *Store
	console *storetestutil.Console
	Project storetestutil.Project
}

func Init(t *testing.T) (context.Context, *tester) {
	t.Helper()

	store := New(NewStoreParams{
		DB:                        deps.db,
		S3:                        deps.s3,
		KMS:                       deps.kms.Client,
		SessionSigningKeyKmsKeyID: deps.kms.SessionSigningKeyID,
		DogfoodProjectID:          deps.console.DogfoodProjectID,
		ConsoleDomain:             deps.console.ConsoleDomain,
	})
	project := deps.console.NewProject(t)
	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	return ctx, &tester{
		Store:   store,
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
