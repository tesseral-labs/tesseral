package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetRBACPolicy_NoActions(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.GetRBACPolicy(ctx, &frontendv1.GetRBACPolicyRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Empty(t, resp.RbacPolicy.GetActions())
}

func TestGetRBACPolicy(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:        "test",
		ApiKeysEnabled:     refOrNil(true),
		CustomRolesEnabled: refOrNil(true),
	})

	projectID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)
	_, err = u.Environment.DB.Exec(t.Context(), `
INSERT INTO actions (id, project_id, name, description)
  VALUES (gen_random_uuid(), $1::uuid, $2, $2),
  		 (gen_random_uuid(), $1::uuid, $3, $3);
`,
		uuid.UUID(projectID).String(),
		"test.action.1",
		"test.action.2",
	)
	require.NoError(t, err)

	resp, err := u.Store.GetRBACPolicy(ctx, &frontendv1.GetRBACPolicyRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.RbacPolicy.GetActions())

	var actions []string
	for _, action := range resp.RbacPolicy.GetActions() {
		actions = append(actions, action.Name)
	}
	require.ElementsMatch(t, []string{"test.action.1", "test.action.2"}, actions)
}
