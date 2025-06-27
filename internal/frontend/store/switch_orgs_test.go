package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func TestListSwitchableOrganizations_Solo(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Org Solo"})

	resp, err := u.Store.ListSwitchableOrganizations(ctx, &frontendv1.ListSwitchableOrganizationsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.GreaterOrEqual(t, len(resp.SwitchableOrganizations), 1)
}

func TestListSwitchableOrganizations_Multiple(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{DisplayName: "Org One"})

	var userEmail string
	err := u.Environment.DB.QueryRow(t.Context(), `
		SELECT email
		FROM users
		WHERE id = $1::uuid
	`, authn.UserID(ctx)).Scan(&userEmail)
	require.NoError(t, err)

	org2ID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "Org Two",
	})
	_ = u.Environment.NewUser(t, org2ID, &backendv1.User{
		Email: userEmail,
	})

	resp, err := u.Store.ListSwitchableOrganizations(ctx, &frontendv1.ListSwitchableOrganizationsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.GreaterOrEqual(t, len(resp.SwitchableOrganizations), 2)

	var names []string
	for _, o := range resp.SwitchableOrganizations {
		names = append(names, o.DisplayName)
	}
	require.ElementsMatch(t, []string{"Org One", "Org Two"}, names)
}
