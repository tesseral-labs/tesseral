package store

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetOrganization_Exists(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName: "org1",
	})

	getResp, err := u.Store.GetOrganization(ctx, &frontendv1.GetOrganizationRequest{})
	require.NoError(t, err)
	require.NotNil(t, getResp.Organization)
	require.Equal(t, "org1", getResp.Organization.DisplayName)
}

func TestGetOrganization_DoesNotExist(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := authn.NewContext(t.Context(), authn.ContextData{
		UserID:         idformat.User.Format(uuid.New()),
		OrganizationID: idformat.Organization.Format(uuid.New()),
		ProjectID:      u.ProjectID,
	})

	_, err := u.Store.GetOrganization(ctx, &frontendv1.GetOrganizationRequest{})
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	require.Equal(t, connect.CodeNotFound, connectErr.Code())
}

func TestUpdateOrganization_UpdatesFields(t *testing.T) {
	t.Parallel()

	u := newTestUtil(t)
	ctx := u.NewOrganizationContext(t, &backendv1.Organization{
		DisplayName:               "org1",
		LogInWithGoogle:           refOrNil(false),
		LogInWithMicrosoft:        refOrNil(false),
		LogInWithGithub:           refOrNil(false),
		LogInWithEmail:            refOrNil(false),
		LogInWithPassword:         refOrNil(false),
		LogInWithSaml:             refOrNil(false),
		LogInWithAuthenticatorApp: refOrNil(false),
		LogInWithPasskey:          refOrNil(false),
		RequireMfa:                refOrNil(false),
	})

	updateResp, err := u.Store.UpdateOrganization(ctx, &frontendv1.UpdateOrganizationRequest{
		Organization: &frontendv1.Organization{
			DisplayName:               "org2",
			LogInWithGoogle:           refOrNil(true),
			LogInWithMicrosoft:        refOrNil(true),
			LogInWithGithub:           refOrNil(true),
			LogInWithEmail:            refOrNil(true),
			LogInWithPassword:         refOrNil(true),
			LogInWithSaml:             refOrNil(true),
			LogInWithAuthenticatorApp: refOrNil(true),
			LogInWithPasskey:          refOrNil(true),
			RequireMfa:                refOrNil(true),
		},
	})
	require.NoError(t, err)
	require.Equal(t, "org2", updateResp.Organization.DisplayName)
	require.Equal(t, true, updateResp.Organization.GetLogInWithGoogle())
	require.Equal(t, true, updateResp.Organization.GetLogInWithMicrosoft())
	require.Equal(t, true, updateResp.Organization.GetLogInWithGithub())
	require.Equal(t, true, updateResp.Organization.GetLogInWithEmail())
	require.Equal(t, true, updateResp.Organization.GetLogInWithPassword())
	require.Equal(t, true, updateResp.Organization.GetLogInWithSaml())
	require.Equal(t, true, updateResp.Organization.GetLogInWithAuthenticatorApp())
	require.Equal(t, true, updateResp.Organization.GetLogInWithPasskey())
	require.Equal(t, true, updateResp.Organization.GetRequireMfa())

	updateUnchangedResp, err := u.Store.UpdateOrganization(ctx, &frontendv1.UpdateOrganizationRequest{
		Organization: &frontendv1.Organization{
			LogInWithGoogle:           nil,
			LogInWithMicrosoft:        nil,
			LogInWithGithub:           nil,
			LogInWithEmail:            nil,
			LogInWithPassword:         nil,
			LogInWithSaml:             nil,
			LogInWithAuthenticatorApp: nil,
			LogInWithPasskey:          nil,
			RequireMfa:                nil,
		},
	})
	require.NoError(t, err)
	require.Equal(t, "org2", updateUnchangedResp.Organization.DisplayName)
	require.Equal(t, true, updateResp.Organization.GetLogInWithGoogle())
	require.Equal(t, true, updateResp.Organization.GetLogInWithMicrosoft())
	require.Equal(t, true, updateResp.Organization.GetLogInWithGithub())
	require.Equal(t, true, updateResp.Organization.GetLogInWithEmail())
	require.Equal(t, true, updateResp.Organization.GetLogInWithPassword())
	require.Equal(t, true, updateResp.Organization.GetLogInWithSaml())
	require.Equal(t, true, updateResp.Organization.GetLogInWithAuthenticatorApp())
	require.Equal(t, true, updateResp.Organization.GetLogInWithPasskey())
	require.Equal(t, true, updateResp.Organization.GetRequireMfa())

}
