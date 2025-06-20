package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetProjectIDOrganizationBacks(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	var backingOrgID uuid.UUID
	err = u.Environment.DB.QueryRow(ctx, `
	SELECT organization_id FROM projects
	WHERE id = $1::uuid;
	`,
		uuid.UUID(projectUUID).String(),
	).Scan(&backingOrgID)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, backingOrgID)

	projectID, err := u.Store.GetProjectIDOrganizationBacks(ctx, idformat.Organization.Format(backingOrgID))
	require.NoError(t, err)
	require.Equal(t, u.ProjectID, projectID)
}
