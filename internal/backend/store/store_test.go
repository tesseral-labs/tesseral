package store

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	commonstore "github.com/tesseral-labs/tesseral/internal/common/store"
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
	Common      *commonstore.Store
	Environment *storetesting.Environment
	ProjectID   string
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
		OIDCClientSecretsKMSKeyID:           environment.KMS.OIDCClientSecretsKMSKeyID,
		DogfoodProjectID:                    environment.DogfoodProjectID,
		ConsoleDomain:                       environment.ConsoleDomain,
		AuthAppsRootDomain:                  environment.AuthAppsRootDomain,
		OIDCClient:                          &oidcclient.Client{HTTPClient: http.DefaultClient},
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
	}
}
