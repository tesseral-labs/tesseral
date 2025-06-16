package consoletest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/dbconntest"
)

func TestCreate(t *testing.T) {
	pool := dbconntest.Open(t)
	store := New(t, pool)

	project := store.CreateProject(t, ProjectParams{Name: "test"})
	require.NotEmpty(t, project.ProjectID, "project ID should not be empty")
	require.NotEmpty(t, project.BackendAPIKey, "project name should not be empty")

	organization := store.CreateOrganization(t, OrganizationParams{
		ProjectID: project.ProjectID,
		Name:      "test",
	})
	require.NotEmpty(t, organization.OrganizationID, "organization ID should not be empty")
	require.NotEmpty(t, organization.UserID, "user ID should not be empty")
	require.Equal(t, project.ProjectID, organization.ProjectID, "project ID in organization should match project ID")
}
