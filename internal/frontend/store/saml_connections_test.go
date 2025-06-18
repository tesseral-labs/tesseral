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
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res.SamlConnection)
	require.NotEmpty(t, res.SamlConnection.SpAcsUrl)
	require.NotEmpty(t, res.SamlConnection.SpEntityId)
	require.Equal(t, "https://idp.example.com/saml/redirect", res.SamlConnection.IdpRedirectUrl)
	require.Equal(t, "https://idp.example.com/saml/idp", res.SamlConnection.IdpEntityId)
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
