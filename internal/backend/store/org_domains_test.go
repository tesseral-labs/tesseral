package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestGetOrganizationDomains_Empty(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.GetOrganizationDomains(ctx, &backendv1.GetOrganizationDomainsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.OrganizationDomains)
	require.Equal(t, orgID, resp.OrganizationDomains.OrganizationId)
	require.Empty(t, resp.OrganizationDomains.Domains)
}

func TestUpdateOrganizationDomains_AddAndGet(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
	})

	domains := []string{"example.com", "test.org"}
	updateResp, err := u.Store.UpdateOrganizationDomains(ctx, &backendv1.UpdateOrganizationDomainsRequest{
		OrganizationId: orgID,
		OrganizationDomains: &backendv1.OrganizationDomains{
			Domains: domains,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, domains, updateResp.OrganizationDomains.Domains)
	u.EnsureAuditLogEvent(t, backendv1.AuditLogEventResourceType_AUDIT_LOG_EVENT_RESOURCE_TYPE_ORGANIZATION, "tesseral.organizations.update_domains")

	getResp, err := u.Store.GetOrganizationDomains(ctx, &backendv1.GetOrganizationDomainsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, domains, getResp.OrganizationDomains.Domains)
}

func TestUpdateOrganizationDomains_ReplaceDomains(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)
	orgID := u.NewOrganization(t, &backendv1.Organization{
		DisplayName: "test",
	})

	initial := []string{"a.com", "b.com"}
	_, err := u.Store.UpdateOrganizationDomains(ctx, &backendv1.UpdateOrganizationDomainsRequest{
		OrganizationId: orgID,
		OrganizationDomains: &backendv1.OrganizationDomains{
			Domains: initial,
		},
	})
	require.NoError(t, err)

	newDomains := []string{"c.com"}
	updateResp, err := u.Store.UpdateOrganizationDomains(ctx, &backendv1.UpdateOrganizationDomainsRequest{
		OrganizationId: orgID,
		OrganizationDomains: &backendv1.OrganizationDomains{
			Domains: newDomains,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newDomains, updateResp.OrganizationDomains.Domains)

	getResp, err := u.Store.GetOrganizationDomains(ctx, &backendv1.GetOrganizationDomainsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newDomains, getResp.OrganizationDomains.Domains)
}
