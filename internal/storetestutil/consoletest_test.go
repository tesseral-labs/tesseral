package storetestutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestNewConsole(t *testing.T) {
	t.Parallel()

	pool, cleanupDB := NewDB()
	t.Cleanup(cleanupDB)

	kms, cleanupKMS := NewKMS()
	t.Cleanup(cleanupKMS)

	console := NewConsole(pool, kms)

	project := console.NewProject(t)
	require.NotEmpty(t, project.ProjectID, "project ID should not be empty")

	organization := console.NewOrganization(t, OrganizationParams{
		Project: project,
		Organization: &backendv1.Organization{
			DisplayName: "test",
		},
	})
	require.NotEmpty(t, organization.OrganizationID, "organization ID should not be empty")
	require.NotEmpty(t, organization.UserID, "user ID should not be empty")
	require.Equal(t, project.ProjectID, organization.ProjectID, "project ID in organization should match project ID")
}
