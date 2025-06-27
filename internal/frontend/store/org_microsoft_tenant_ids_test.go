package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func TestGetOrganizationMicrosoftTenantIDs_Empty(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.GetOrganizationMicrosoftTenantIDs(ctx, &frontendv1.GetOrganizationMicrosoftTenantIDsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp.OrganizationMicrosoftTenantIds)
	require.Empty(t, resp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
}

func TestUpdateOrganizationMicrosoftTenantIDs_AddAndGet(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	ids := []string{"tenant1", "tenant2"}
	updateResp, err := u.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, &frontendv1.UpdateOrganizationMicrosoftTenantIDsRequest{
		OrganizationMicrosoftTenantIds: &frontendv1.OrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds: ids,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, ids, updateResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)

	getResp, err := u.Store.GetOrganizationMicrosoftTenantIDs(ctx, &frontendv1.GetOrganizationMicrosoftTenantIDsRequest{})
	require.NoError(t, err)
	require.ElementsMatch(t, ids, getResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
}

func TestUpdateOrganizationMicrosoftTenantIDs_ReplaceIDs(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "test",
	})

	initial := []string{"a", "b"}
	_, err := u.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, &frontendv1.UpdateOrganizationMicrosoftTenantIDsRequest{
		OrganizationMicrosoftTenantIds: &frontendv1.OrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds: initial,
		},
	})
	require.NoError(t, err)

	newIDs := []string{"c"}
	updateResp, err := u.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, &frontendv1.UpdateOrganizationMicrosoftTenantIDsRequest{
		OrganizationMicrosoftTenantIds: &frontendv1.OrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds: newIDs,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newIDs, updateResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)

	getResp, err := u.Store.GetOrganizationMicrosoftTenantIDs(ctx, &frontendv1.GetOrganizationMicrosoftTenantIDsRequest{})
	require.NoError(t, err)
	require.ElementsMatch(t, newIDs, getResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
}
