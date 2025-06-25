package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestCreateOIDCConnection_OIDCEnabled(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	res, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			Issuer:           "https://issuer.example.com",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			OrganizationId:   organizationID,
			Primary:          refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.OidcConnection)
	require.Equal(t, organizationID, res.OidcConnection.OrganizationId)
	require.NotEmpty(t, res.OidcConnection.CreateTime)
	require.NotEmpty(t, res.OidcConnection.UpdateTime)
	require.True(t, res.OidcConnection.GetPrimary())
}

func TestCreateOIDCConnection_OIDCDisabled(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(false),
	})
	_, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			Issuer:           "https://issuer.example.com",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			OrganizationId:   organizationID,
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestGetOIDCConnection_Exists(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	createResp, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			Issuer:           "https://issuer.example.com",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			OrganizationId:   organizationID,
		},
	})
	require.NoError(t, err)
	connID := createResp.OidcConnection.Id
	res, err := u.Store.GetOIDCConnection(ctx, &backendv1.GetOIDCConnectionRequest{Id: connID})
	require.NoError(t, err)
	require.NotNil(t, res.OidcConnection)
	require.Equal(t, connID, res.OidcConnection.Id)
}

func TestGetOIDCConnection_DoesNotExist(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	_, err := u.Store.GetOIDCConnection(ctx, &backendv1.GetOIDCConnectionRequest{
		Id: idformat.OIDCConnection.Format(uuid.New()),
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateOIDCConnection_UpdatesFields(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	createResp, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			Issuer:           "https://issuer.example.com",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			OrganizationId:   organizationID,
		},
	})
	require.NoError(t, err)
	connID := createResp.OidcConnection.Id
	updateResp, err := u.Store.UpdateOIDCConnection(ctx, &backendv1.UpdateOIDCConnectionRequest{
		Id: connID,
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
			Issuer:           "https://issuer2.example.com",
			ClientId:         "new-client-id",
			ClientSecret:     "new-client-secret",
			Primary:          refOrNil(true),
		},
	})
	require.NoError(t, err)
	updated := updateResp.OidcConnection
	require.Equal(t, "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration", updated.ConfigurationUrl)
	require.Equal(t, "https://issuer2.example.com", updated.Issuer)
	require.Equal(t, "new-client-id", updated.ClientId)
	require.True(t, updated.GetPrimary())
}

func TestDeleteOIDCConnection_RemovesConnection(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	createResp, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			Issuer:           "https://issuer.example.com",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			OrganizationId:   organizationID,
		},
	})
	require.NoError(t, err)
	connID := createResp.OidcConnection.Id
	_, err = u.Store.DeleteOIDCConnection(ctx, &backendv1.DeleteOIDCConnectionRequest{Id: connID})
	require.NoError(t, err)
	_, err = u.Store.GetOIDCConnection(ctx, &backendv1.GetOIDCConnectionRequest{Id: connID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestCreateOIDCConnection_InvalidConfigURL(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	_, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
		OidcConnection: &backendv1.OIDCConnection{
			ConfigurationUrl: "not-a-url",
			OrganizationId:   organizationID,
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestListOIDCConnections_ReturnsAllForOrg(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	var ids []string
	for range 3 {
		resp, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
			OidcConnection: &backendv1.OIDCConnection{
				ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
				Issuer:           "https://issuer.example.com",
				ClientId:         "client-id",
				ClientSecret:     "client-secret",
				OrganizationId:   organizationID,
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.OidcConnection.Id)
	}
	listResp, err := u.Store.ListOIDCConnections(ctx, &backendv1.ListOIDCConnectionsRequest{
		OrganizationId: organizationID,
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.OidcConnections, 3)
	respIds := []string{}
	for _, conn := range listResp.OidcConnections {
		respIds = append(respIds, conn.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListOIDCConnections_Pagination(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)
	organizationID := u.Environment.NewOrganization(t, u.ProjectID, &backendv1.Organization{
		DisplayName:   "test",
		LogInWithOidc: refOrNil(true),
	})
	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateOIDCConnection(ctx, &backendv1.CreateOIDCConnectionRequest{
			OidcConnection: &backendv1.OIDCConnection{
				ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
				Issuer:           "https://issuer.example.com",
				ClientId:         "client-id",
				ClientSecret:     "client-secret",
				OrganizationId:   organizationID,
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.OidcConnection.Id)
	}
	resp1, err := u.Store.ListOIDCConnections(ctx, &backendv1.ListOIDCConnectionsRequest{
		OrganizationId: organizationID,
	})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.OidcConnections, 10)
	require.NotEmpty(t, resp1.NextPageToken)
	resp2, err := u.Store.ListOIDCConnections(ctx, &backendv1.ListOIDCConnectionsRequest{
		OrganizationId: organizationID,
		PageToken:      resp1.NextPageToken,
	})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.OidcConnections, 5)
	require.Empty(t, resp2.NextPageToken)
	var allIDs []string
	for _, c := range resp1.OidcConnections {
		allIDs = append(allIDs, c.Id)
	}
	for _, c := range resp2.OidcConnections {
		allIDs = append(allIDs, c.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
