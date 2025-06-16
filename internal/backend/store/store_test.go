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
	Store   *store.Store
	console *storetestutil.Console
}

func (s *BackendSuite) SetupSuite() {
	db := storetestutil.NewDB(s.T())
	console := storetestutil.NewConsole(s.T(), db)
	kms := storetestutil.NewKMS(s.T())

	s.console = console
	s.Store = store.New(store.NewStoreParams{
		DB:                        db,
		S3:                        storetestutil.NewS3Client(s.T()),
		KMS:                       kms.Client,
		SessionSigningKeyKmsKeyID: kms.SessionSigningKeyID,
		DogfoodProjectID:          console.DogfoodProjectID,
		ConsoleDomain:             console.ConsoleDomain,
	})
}

func (s *BackendSuite) NewAuthContext(project storetestutil.Project) context.Context {
	return authn.NewDogfoodSessionContext(s.T().Context(), authn.DogfoodSessionContextData{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})
}

func (s *BackendSuite) CreateProject() storetestutil.Project {
	return s.console.CreateProject(s.T())
}

func (s *BackendSuite) CreateOrganization(params storetestutil.OrganizationParams) storetestutil.Organization {
	return s.console.CreateOrganization(s.T(), params)
}

func TestBackendSuite(t *testing.T) {
	suite.Run(t, new(BackendSuite))
}
