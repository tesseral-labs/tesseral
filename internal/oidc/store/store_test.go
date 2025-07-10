package store

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/oidc/authn"
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

func newTestUtil(t *testing.T) (context.Context, *testUtil) {
	store := New(NewStoreParams{
		DB:                        environment.DB,
		KMS:                       environment.KMS.Client,
		OIDCClientSecretsKMSKeyID: environment.KMS.OIDCClientSecretsKMSKeyID,
		OIDCClient:                &oidcclient.Client{HTTPClient: http.DefaultClient},
	})

	projectID, _ := environment.NewProject(t)
	projectUUID, err := idformat.Project.Parse(projectID)
	require.NoError(t, err)

	secretToken := environment.NewIntermediateSession(t, projectID)
	qIntermediateSession, err := store.AuthenticateIntermediateSession(t.Context(), projectUUID, secretToken)
	require.NoError(t, err)

	ctx := authn.NewContext(t.Context(), qIntermediateSession, projectUUID)

	return ctx, &testUtil{
		Store:       store,
		Environment: environment,
		ProjectID:   projectID,
	}
}
