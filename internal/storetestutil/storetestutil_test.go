package storetestutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestNewEnvironment(t *testing.T) {
	t.Parallel()

	env, cleanup := NewEnvironment()
	defer cleanup()

	projectID, _ := env.NewProject(t)
	require.NotEmpty(t, projectID, "project ID should not be empty")

	organizationID := env.NewOrganization(t, projectID, &backendv1.Organization{
		DisplayName: "test",
	})
	require.NotEmpty(t, organizationID, "organization ID should not be empty")
}
