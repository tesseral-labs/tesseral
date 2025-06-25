package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func TestCreateSAMLConnection_SAMLEnabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})

	res, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
		SamlConnection: &frontendv1.SAMLConnection{
			IdpRedirectUrl: "https://idp.example.com/saml/redirect",
			IdpEntityId:    "https://idp.example.com/saml/idp",
			Primary:        refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.SamlConnection)
	require.NotEmpty(t, res.SamlConnection.SpAcsUrl)
	require.NotEmpty(t, res.SamlConnection.SpEntityId)
	require.Equal(t, "https://idp.example.com/saml/redirect", res.SamlConnection.IdpRedirectUrl)
	require.Equal(t, "https://idp.example.com/saml/idp", res.SamlConnection.IdpEntityId)
	require.NotEmpty(t, res.SamlConnection.CreateTime)
	require.NotEmpty(t, res.SamlConnection.UpdateTime)
	require.True(t, res.SamlConnection.GetPrimary())
}

func TestCreateSAMLConnection_SAMLDisabled(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(false),
	})

	_, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
		SamlConnection: &frontendv1.SAMLConnection{
			IdpRedirectUrl: "https://idp.example.com/saml/redirect",
			IdpEntityId:    "https://idp.example.com/saml/idp",
		},
	})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeFailedPrecondition, connectErr.Code())
}

func TestGetSAMLConnection_Exists(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})
	createResp, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
		SamlConnection: &frontendv1.SAMLConnection{
			IdpRedirectUrl: "https://idp.example.com/saml/redirect",
			IdpEntityId:    "https://idp.example.com/saml/idp",
		},
	})
	require.NoError(t, err)
	connID := createResp.SamlConnection.Id

	getResp, err := u.Store.GetSAMLConnection(ctx, &frontendv1.GetSAMLConnectionRequest{Id: connID})
	require.NoError(t, err)
	require.NotNil(t, getResp.SamlConnection)
	require.Equal(t, connID, getResp.SamlConnection.Id)
}

func TestGetSAMLConnection_DoesNotExist(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})

	_, err := u.Store.GetSAMLConnection(ctx, &frontendv1.GetSAMLConnectionRequest{Id: "nonexistent-id"})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestUpdateSAMLConnection_UpdatesFields(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})
	createResp, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
		SamlConnection: &frontendv1.SAMLConnection{
			IdpRedirectUrl: "https://idp.example.com/saml/redirect",
			IdpEntityId:    "https://idp.example.com/saml/idp",
		},
	})
	require.NoError(t, err)
	connID := createResp.SamlConnection.Id

	updateResp, err := u.Store.UpdateSAMLConnection(ctx, &frontendv1.UpdateSAMLConnectionRequest{
		Id: connID,
		SamlConnection: &frontendv1.SAMLConnection{
			IdpRedirectUrl: "https://idp.example.com/saml/redirect2",
			IdpEntityId:    "https://idp.example.com/saml/idp2",
			Primary:        refOrNil(true),
		},
	})
	require.NoError(t, err)
	updated := updateResp.SamlConnection
	require.Equal(t, "https://idp.example.com/saml/redirect2", updated.IdpRedirectUrl)
	require.Equal(t, "https://idp.example.com/saml/idp2", updated.IdpEntityId)
	require.True(t, updated.GetPrimary())
}

func TestDeleteSAMLConnection_RemovesConnection(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})
	createResp, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
		SamlConnection: &frontendv1.SAMLConnection{
			IdpRedirectUrl: "https://idp.example.com/saml/redirect",
			IdpEntityId:    "https://idp.example.com/saml/idp",
		},
	})
	require.NoError(t, err)
	connID := createResp.SamlConnection.Id

	_, err = u.Store.DeleteSAMLConnection(ctx, &frontendv1.DeleteSAMLConnectionRequest{Id: connID})
	require.NoError(t, err)

	_, err = u.Store.GetSAMLConnection(ctx, &frontendv1.GetSAMLConnectionRequest{Id: connID})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestListSAMLConnections_ReturnsAll(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})

	var ids []string
	for range 3 {
		resp, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
			SamlConnection: &frontendv1.SAMLConnection{
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
			},
		})
		require.NoError(t, err)
		ids = append(ids, resp.SamlConnection.Id)
	}

	listResp, err := u.Store.ListSAMLConnections(ctx, &frontendv1.ListSAMLConnectionsRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.SamlConnections, 3)

	var respIds []string
	for _, conn := range listResp.SamlConnections {
		respIds = append(respIds, conn.Id)
	}
	require.ElementsMatch(t, ids, respIds)
}

func TestListSAMLConnections_Pagination(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:   "Test Organization",
		LogInWithSaml: refOrNil(true),
	})

	var createdIDs []string
	for range 15 {
		resp, err := u.Store.CreateSAMLConnection(ctx, &frontendv1.CreateSAMLConnectionRequest{
			SamlConnection: &frontendv1.SAMLConnection{
				IdpRedirectUrl: "https://idp.example.com/saml/redirect",
				IdpEntityId:    "https://idp.example.com/saml/idp",
			},
		})
		require.NoError(t, err)
		createdIDs = append(createdIDs, resp.SamlConnection.Id)
	}

	resp1, err := u.Store.ListSAMLConnections(ctx, &frontendv1.ListSAMLConnectionsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Len(t, resp1.SamlConnections, 10)
	require.NotEmpty(t, resp1.NextPageToken)

	resp2, err := u.Store.ListSAMLConnections(ctx, &frontendv1.ListSAMLConnectionsRequest{PageToken: resp1.NextPageToken})
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Len(t, resp2.SamlConnections, 5)
	require.Empty(t, resp2.NextPageToken)

	var allIDs []string
	for _, c := range resp1.SamlConnections {
		allIDs = append(allIDs, c.Id)
	}
	for _, c := range resp2.SamlConnections {
		allIDs = append(allIDs, c.Id)
	}
	require.ElementsMatch(t, createdIDs, allIDs)
}
