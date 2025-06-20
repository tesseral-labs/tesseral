package store

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	commonstore "github.com/tesseral-labs/tesseral/internal/common/store"
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
	Common      *commonstore.Store
	Environment *storetesting.Environment
	ProjectID   string
	ctx         context.Context
}

func newTestUtil(t *testing.T) (context.Context, *testUtil) {
	store := New(NewStoreParams{
		DB:                                  environment.DB,
		S3:                                  environment.S3.Client,
		S3UserContentBucketName:             environment.S3.UserContentBucketName,
		KMS:                                 environment.KMS.Client,
		SessionSigningKeyKmsKeyID:           environment.KMS.SessionSigningKeyID,
		GoogleOAuthClientSecretsKMSKeyID:    environment.KMS.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID: environment.KMS.MicrosoftOAuthClientSecretsKMSKeyID,
		GithubOAuthClientSecretsKMSKeyID:    environment.KMS.GithubOAuthClientSecretsKMSKeyID,
		DogfoodProjectID:                    environment.DogfoodProjectID,
		ConsoleDomain:                       environment.ConsoleDomain,
		AuthAppsRootDomain:                  environment.AuthAppsRootDomain,
	})
	commonStore := commonstore.New(commonstore.NewStoreParams{
		AppAuthRootDomain:         environment.ConsoleDomain,
		DB:                        environment.DB,
		KMS:                       environment.KMS.Client,
		SessionSigningKeyKMSKeyID: environment.KMS.SessionSigningKeyID,
	})

	projectID, projectUserID := environment.NewProject(t)
	ctx := authn.NewDogfoodSessionContext(t.Context(), authn.DogfoodSessionContextData{
		ProjectID: projectID,
		UserID:    projectUserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	return ctx, &testUtil{
		Store:       store,
		Common:      commonStore,
		Environment: environment,
		ProjectID:   projectID,
		ctx:         ctx,
	}
}

func (u *testUtil) NewOrganization(t *testing.T, organization *backendv1.Organization) string {
	return u.Environment.NewOrganization(t, u.ProjectID, organization)
}

func (u *testUtil) CreateActions(t *testing.T, names ...string) {
	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	for _, name := range names {
		_, err := u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2);
`,
			uuid.UUID(projectID).String(),
			name,
		)
		require.NoError(t, err)
	}
}

func (u *testUtil) EnsureAuditLogEvent(t *testing.T, resourceType backendv1.AuditLogEventResourceType, eventType string) {
	resp, err := u.Store.ConsoleListAuditLogEventNames(u.ctx, &backendv1.ConsoleListAuditLogEventNamesRequest{
		ResourceType: resourceType,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.EventNames, eventType)
}

func (u *testUtil) NewUser(t *testing.T, organizationID string, email string) string {
	return u.Environment.NewUser(t, organizationID, &backendv1.User{
		Email: email,
	})
}
