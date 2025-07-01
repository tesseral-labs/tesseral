package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func TestCreateOIDCConnection_OIDCEnabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})

	res, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			Primary:          refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.OidcConnection)
	require.NotEmpty(t, res.OidcConnection.RedirectUri)
	require.Equal(t, "https://accounts.google.com/.well-known/openid-configuration", res.OidcConnection.ConfigurationUrl)
	require.Equal(t, "client-id", res.OidcConnection.ClientId)
	require.True(t, res.OidcConnection.GetPrimary())
}

func TestCreateOIDCConnection_OIDCDisabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(false),
	})

	_, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestGetOIDCConnection_Exists(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})
	createResp, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
		},
	})
	require.NoError(t, err)
	connID := createResp.OidcConnection.Id

	getResp, err := u.Store.GetOIDCConnection(ctx, &frontendv1.GetOIDCConnectionRequest{Id: connID})
	require.NoError(t, err)
	require.NotNil(t, getResp.OidcConnection)
	require.Equal(t, connID, getResp.OidcConnection.Id)
}

func TestGetOIDCConnection_DoesNotExist(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})

	_, err := u.Store.GetOIDCConnection(ctx, &frontendv1.GetOIDCConnectionRequest{Id: "nonexistent-id"})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestUpdateOIDCConnection_UpdatesFields(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})
	createResp, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
		},
	})
	require.NoError(t, err)
	connID := createResp.OidcConnection.Id

	updateResp, err := u.Store.UpdateOIDCConnection(ctx, &frontendv1.UpdateOIDCConnectionRequest{
		Id: connID,
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
			ClientId:         "client-id-2",
			ClientSecret:     "client-secret-2",
			Primary:          refOrNil(true),
		},
	})
	require.NoError(t, err)
	updated := updateResp.OidcConnection
	require.Equal(t, "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration", updated.ConfigurationUrl)
	require.Equal(t, "client-id-2", updated.ClientId)
	require.True(t, updated.GetPrimary())
}

func TestUpdateOIDCConnection_SetPrimary(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})

	original, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
			Primary:          refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.True(t, original.OidcConnection.GetPrimary())

	new, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
			ClientId:         "client-id-2",
			ClientSecret:     "client-secret-2",
			Primary:          refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.True(t, new.OidcConnection.GetPrimary())

	getResp, err := u.Store.GetOIDCConnection(ctx, &frontendv1.GetOIDCConnectionRequest{Id: original.OidcConnection.Id})
	require.NoError(t, err)
	require.NotNil(t, getResp.OidcConnection)
	require.False(t, getResp.OidcConnection.GetPrimary(), "original connection should no longer be primary")
}

func TestDeleteOIDCConnection_RemovesConnection(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})
	createResp, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
		OidcConnection: &frontendv1.OIDCConnection{
			ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
			ClientId:         "client-id",
			ClientSecret:     "client-secret",
		},
	})
	require.NoError(t, err)
	connID := createResp.OidcConnection.Id

	_, err = u.Store.DeleteOIDCConnection(ctx, &frontendv1.DeleteOIDCConnectionRequest{Id: connID})
	require.NoError(t, err)

	_, err = u.Store.GetOIDCConnection(ctx, &frontendv1.GetOIDCConnectionRequest{Id: connID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListOIDCConnections_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
			OidcConnection: &frontendv1.OIDCConnection{
				ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
				ClientId:         "client-id",
				ClientSecret:     "client-secret",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.OidcConnection.Id)
	}

	listResp, err := u.Store.ListOIDCConnections(ctx, &frontendv1.ListOIDCConnectionsRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.OidcConnections, 3)

	var respIds []string
	for _, conn := range listResp.OidcConnections {
		respIds = append(respIds, conn.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListOIDCConnections_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithOidc: refOrNil(true),
	})

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateOIDCConnection(ctx, &frontendv1.CreateOIDCConnectionRequest{
			OidcConnection: &frontendv1.OIDCConnection{
				ConfigurationUrl: "https://accounts.google.com/.well-known/openid-configuration",
				ClientId:         "client-id",
				ClientSecret:     "client-secret",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.OidcConnection.Id)
	}

	resp1, err := u.Store.ListOIDCConnections(ctx, &frontendv1.ListOIDCConnectionsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.OidcConnections, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListOIDCConnections(ctx, &frontendv1.ListOIDCConnectionsRequest{PageToken: resp1.NextPageToken})
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
