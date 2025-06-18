package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
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
		SessionSigningKeyKmsKeyID: environment.KMS.SessionSigningKeyID,
		DogfoodProjectID:          environment.DogfoodProjectID,
		ConsoleDomain:             environment.ConsoleDomain,
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

	userID := uuid.New()
	email := fmt.Sprintf("%s@%s", userID.String(), u.Environment.ConsoleDomain)

	organizationUUID, err := idformat.Organization.Parse(organizationID)
	require.NoError(t, err)

	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
VALUES ($1::uuid, $2, crypt('password', gen_salt('bf', 14)), $3, true);
`,
		userID.String(),
		email,
		uuid.UUID(organizationUUID).String())
	require.NoError(t, err)

	ctx := authn.NewContext(t.Context(), authn.ContextData{
		ProjectID:      u.ProjectID,
		OrganizationID: organizationID,
		UserID:         idformat.User.Format(userID),
		SessionID:      idformat.Session.Format(uuid.New()),
	})

	return ctx
}
