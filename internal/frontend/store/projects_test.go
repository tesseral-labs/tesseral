package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func TestGetProject(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "Test Organization",
	})

	project, err := u.Store.GetProject(ctx, &frontendv1.GetProjectRequest{})
	require.NoError(t, err)
	require.NotNil(t, project)
	require.Equal(t, u.ProjectID, project.Project.Id)
}
