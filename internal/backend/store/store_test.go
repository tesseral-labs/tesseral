package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	"github.com/tesseral-labs/tesseral/internal/backend/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/storetestutil"
)

type BackendSuite struct {
	suite.Suite
	*tester
}

func (s *BackendSuite) SetupSuite() {
	s.tester = newTester(s.T())
}

type tester struct {
	Store *store.Store
	*storetestutil.Console
}

func newTester(t *testing.T) *tester {
	t.Helper()

	db := storetestutil.NewDB(t)
	console := storetestutil.NewConsole(t, db)
	kms := storetestutil.NewKMS(t)
	store := store.New(store.NewStoreParams{
		DB:                        db,
		S3:                        storetestutil.NewS3Client(t),
		KMS:                       kms.Client,
		SessionSigningKeyKmsKeyID: kms.SessionSigningKeyID,
		DogfoodProjectID:          console.DogfoodProjectID,
		ConsoleDomain:             console.ConsoleDomain,
	})

	return &tester{
		Store:   store,
		Console: console,
	}
}

func (tester *tester) NewAuthContext(t *testing.T, project storetestutil.Project) context.Context {
	return authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})
}

func TestBackendSuite(t *testing.T) {
	suite.Run(t, new(BackendSuite))
}
