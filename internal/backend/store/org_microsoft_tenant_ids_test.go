package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestGetOrganizationMicrosoftTenantIDs_Empty(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.GetOrganizationMicrosoftTenantIDs(ctx, &backendv1.GetOrganizationMicrosoftTenantIDsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.OrganizationMicrosoftTenantIds)
	require.Equal(t, orgID, resp.OrganizationMicrosoftTenantIds.OrganizationId)
	require.Empty(t, resp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
}

func TestUpdateOrganizationMicrosoftTenantIDs_AddAndGet(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
	})

	tenantIDs := []string{"tenant1", "tenant2"}
	updateResp, err := u.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, &backendv1.UpdateOrganizationMicrosoftTenantIDsRequest{
		OrganizationId: orgID,
		OrganizationMicrosoftTenantIds: &backendv1.OrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds: tenantIDs,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, tenantIDs, updateResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_ORGANIZATION, "tesseral.organizations.update_microsoft_tenant_ids")

	getResp, err := u.Store.GetOrganizationMicrosoftTenantIDs(ctx, &backendv1.GetOrganizationMicrosoftTenantIDsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, tenantIDs, getResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
}

func TestUpdateOrganizationMicrosoftTenantIDs_ReplaceTenantIDs(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
	})

	initial := []string{"a", "b"}
	_, err := u.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, &backendv1.UpdateOrganizationMicrosoftTenantIDsRequest{
		OrganizationId: orgID,
		OrganizationMicrosoftTenantIds: &backendv1.OrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds: initial,
		},
	})
	require.NoError(t, err)

	newTenantIDs := []string{"c"}
	updateResp, err := u.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, &backendv1.UpdateOrganizationMicrosoftTenantIDsRequest{
		OrganizationId: orgID,
		OrganizationMicrosoftTenantIds: &backendv1.OrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds: newTenantIDs,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newTenantIDs, updateResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)

	getResp, err := u.Store.GetOrganizationMicrosoftTenantIDs(ctx, &backendv1.GetOrganizationMicrosoftTenantIDsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newTenantIDs, getResp.OrganizationMicrosoftTenantIds.MicrosoftTenantIds)
}
