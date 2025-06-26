package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateOrganization_Success(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	resp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
		Organization: &backendv1.Organization{
			DisplayName: "org1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Organization)
	require.Equal(t, "org1", resp.Organization.DisplayName)
	require.NotEmpty(t, resp.Organization.Id)
	require.NotEmpty(t, resp.Organization.CreateTime)
	require.NotEmpty(t, resp.Organization.UpdateTime)
}

func TestGetOrganization_Exists(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
		Organization: &backendv1.Organization{
			DisplayName: "org1",
		},
	})
	require.NoError(t, err)
	orgID := createResp.Organization.Id

	getResp, err := u.Store.GetOrganization(ctx, &backendv1.GetOrganizationRequest{Id: orgID})
	require.NoError(t, err)
	require.NotNil(t, getResp.Organization)
	require.Equal(t, orgID, getResp.Organization.Id)
}

func TestGetOrganization_DoesNotExist(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	_, err := u.Store.GetOrganization(ctx, &backendv1.GetOrganizationRequest{
		Id: idformat.Organization.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateOrganization_UpdatesFields(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
		Organization: &backendv1.Organization{
			DisplayName: "org1",
		},
	})
	require.NoError(t, err)
	orgID := createResp.Organization.Id

	updateResp, err := u.Store.UpdateOrganization(ctx, &backendv1.UpdateOrganizationRequest{
		Id: orgID,
		Organization: &backendv1.Organization{
			DisplayName: "org2",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "org2", updateResp.Organization.DisplayName)
}

func TestListOrganizations_ReturnsAll(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
			Organization: &backendv1.Organization{
				DisplayName: "test-org",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.Organization.Id)
	}

	listResp, err := u.Store.ListOrganizations(ctx, &backendv1.ListOrganizationsRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Organizations, 3)

	var respIds []string
	for _, org := range listResp.Organizations {
		respIds = append(respIds, org.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListOrganizations_Pagination(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
			Organization: &backendv1.Organization{
				DisplayName: "test-org",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.Organization.Id)
	}

	resp1, err := u.Store.ListOrganizations(ctx, &backendv1.ListOrganizationsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.Organizations, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListOrganizations(ctx, &backendv1.ListOrganizationsRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.Organizations, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, k := range resp1.Organizations {
		allIDs = append(allIDs, k.Id)
	}
	for _, k := range resp2.Organizations {
		allIDs = append(allIDs, k.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}

func TestDeleteOrganization(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	createResp, err := u.Store.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
		Organization: &backendv1.Organization{
			DisplayName: "org1",
		},
	})
	require.NoError(t, err)
	orgID := createResp.Organization.Id

	_, err = u.Store.DeleteOrganization(ctx, &backendv1.DeleteOrganizationRequest{Id: orgID})
	require.NoError(t, err)

	_, err = u.Store.GetOrganization(ctx, &backendv1.GetOrganizationRequest{Id: orgID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}
