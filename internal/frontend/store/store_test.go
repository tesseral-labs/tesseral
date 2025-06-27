package store

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	"github.com/tesseral-labs/tesseral/internal/oidcclient"
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

func newTestUtil(t *testing.T) *testUtil {
	store := New(NewStoreParams{
		DB:                        environment.DB,
		KMS:                       environment.KMS.Client,
		OIDCClientSecretsKMSKeyID: environment.KMS.OIDCClientSecretsKMSKeyID,
		SessionSigningKeyKmsKeyID: environment.KMS.SessionSigningKeyID,
		DogfoodProjectID:          environment.DogfoodProjectID,
		ConsoleDomain:             environment.ConsoleDomain,
		OIDCClient:                &oidcclient.Client{HTTPClient: http.DefaultClient},
	})
	projectID, _ := environment.NewProject(t)

	return &testUtil{
		Store:       store,
		Environment: environment,
		ProjectID:   projectID,
	}
}

func (u *testUtil) NewOrganizationContext(t *testing.T, organization *backendv1.Organization) context.Context {
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, organization)
	userID := u.Environment.NewUser(t, organizationID, &backendv1.User{
		Owner: refOrNil(true),
	})

	ctx := authn.NewContext(t.Context(), authn.ContextData{
		ProjectID:      u.ProjectID,
		OrganizationID: organizationID,
		UserID:         userID,
		SessionID:      idformat.Session.Format(uuid.New()),
	})

	return ctx
}
