package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestGetOrganizationGoogleHostedDomains_Empty(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	resp, err := u.Store.GetOrganizationGoogleHostedDomains(ctx, &backendv1.GetOrganizationGoogleHostedDomainsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.OrganizationGoogleHostedDomains)
	require.Equal(t, orgID, resp.OrganizationGoogleHostedDomains.OrganizationId)
	require.Empty(t, resp.OrganizationGoogleHostedDomains.GoogleHostedDomains)
}

func TestUpdateOrganizationGoogleHostedDomains_AddAndGet(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	domains := []string{"example.com", "test.org"}
	updateResp, err := u.Store.UpdateOrganizationGoogleHostedDomains(ctx, &backendv1.UpdateOrganizationGoogleHostedDomainsRequest{
		OrganizationId: orgID,
		OrganizationGoogleHostedDomains: &backendv1.OrganizationGoogleHostedDomains{
			GoogleHostedDomains: domains,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, domains, updateResp.OrganizationGoogleHostedDomains.GoogleHostedDomains)

	getResp, err := u.Store.GetOrganizationGoogleHostedDomains(ctx, &backendv1.GetOrganizationGoogleHostedDomainsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, domains, getResp.OrganizationGoogleHostedDomains.GoogleHostedDomains)
}

func TestUpdateOrganizationGoogleHostedDomains_ReplaceDomains(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	orgID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName: "test",
	})

	initial := []string{"a.com", "b.com"}
	_, err := u.Store.UpdateOrganizationGoogleHostedDomains(ctx, &backendv1.UpdateOrganizationGoogleHostedDomainsRequest{
		OrganizationId: orgID,
		OrganizationGoogleHostedDomains: &backendv1.OrganizationGoogleHostedDomains{
			GoogleHostedDomains: initial,
		},
	})
	require.NoError(t, err)

	newDomains := []string{"c.com"}
	updateResp, err := u.Store.UpdateOrganizationGoogleHostedDomains(ctx, &backendv1.UpdateOrganizationGoogleHostedDomainsRequest{
		OrganizationId: orgID,
		OrganizationGoogleHostedDomains: &backendv1.OrganizationGoogleHostedDomains{
			GoogleHostedDomains: newDomains,
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newDomains, updateResp.OrganizationGoogleHostedDomains.GoogleHostedDomains)

	getResp, err := u.Store.GetOrganizationGoogleHostedDomains(ctx, &backendv1.GetOrganizationGoogleHostedDomainsRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, newDomains, getResp.OrganizationGoogleHostedDomains.GoogleHostedDomains)
}
