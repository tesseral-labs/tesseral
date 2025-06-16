package consoletest

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/dbconntest"
)

func TestCreate(t *testing.T) {
	pool := dbconntest.Open(t)
	store := New(t, pool)

	project := store.CreateProject(t)
	require.NotEmpty(t, project.ProjectID, "project ID should not be empty")

	organization := store.CreateOrganization(t, OrganizationParams{
		Project: project,
		Organization: &backendv1.Organization{
			DisplayName: "test",
		},
	})
	require.NotEmpty(t, organization.OrganizationID, "organization ID should not be empty")
	require.NotEmpty(t, organization.UserID, "user ID should not be empty")
	require.Equal(t, project.ProjectID, organization.ProjectID, "project ID in organization should match project ID")
}
