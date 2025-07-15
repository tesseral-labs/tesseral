package store

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func TestGetProject_Exists(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	resp, err := u.Store.GetProject(ctx, &backendv1.GetProjectRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp.Project)
	require.NotEmpty(t, resp.Project.Id)
	require.NotEmpty(t, resp.Project.DisplayName)
	require.NotEmpty(t, resp.Project.CreateTime)
	require.NotEmpty(t, resp.Project.UpdateTime)
	require.True(t, resp.Project.GetLogInWithGoogle())
	require.True(t, resp.Project.GetLogInWithMicrosoft())
	require.True(t, resp.Project.GetLogInWithGithub())
	require.True(t, resp.Project.GetLogInWithEmail())
	require.True(t, resp.Project.GetLogInWithPassword())
	require.True(t, resp.Project.GetLogInWithSaml())
	require.True(t, resp.Project.GetLogInWithOidc())
	require.True(t, resp.Project.GetLogInWithAuthenticatorApp())
	require.True(t, resp.Project.GetLogInWithPasskey())
	require.Empty(t, resp.Project.GoogleOauthClientId)
	require.Empty(t, resp.Project.GoogleOauthClientSecret)
	require.Empty(t, resp.Project.MicrosoftOauthClientId)
	require.Empty(t, resp.Project.MicrosoftOauthClientSecret)
	require.Empty(t, resp.Project.GithubOauthClientId)
	require.Empty(t, resp.Project.GithubOauthClientSecret)
	require.NotEmpty(t, resp.Project.VaultDomain)
	require.True(t, resp.Project.VaultDomainCustom)
	require.NotEmpty(t, resp.Project.TrustedDomains)
	require.NotEmpty(t, resp.Project.CookieDomain)
	require.NotEmpty(t, resp.Project.EmailSendFromDomain)
	require.True(t, resp.Project.GetApiKeysEnabled())
	require.NotEmpty(t, resp.Project.ApiKeySecretTokenPrefix)
	require.True(t, resp.Project.GetAuditLogsEnabled())
}

func TestUpdateProject_UpdatesFields(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	getResp, err := u.Store.GetProject(ctx, &backendv1.GetProjectRequest{})
	require.NoError(t, err)
	oldName := getResp.Project.DisplayName

	updateResp, err := u.Store.UpdateProject(ctx, &backendv1.UpdateProjectRequest{
		Project: &backendv1.Project{
			DisplayName:                "new-project-name",
			GoogleOauthClientId:        refOrNil("new-google-client-id"),
			GoogleOauthClientSecret:    "new-google-client-secret",
			MicrosoftOauthClientId:     refOrNil("new-microsoft-client-id"),
			MicrosoftOauthClientSecret: "new-microsoft-client-secret",
			GithubOauthClientId:        refOrNil("new-github-client-id"),
			GithubOauthClientSecret:    "new-github-client-secret",
			CookieDomain:               "cookies.example.com",
			TrustedDomains: []string{
				"trusted.example.com",
				"another-trusted.example.com",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "new-project-name", updateResp.Project.DisplayName)

	getResp2, err := u.Store.GetProject(ctx, &backendv1.GetProjectRequest{})
	require.NoError(t, err)
	require.Equal(t, "new-project-name", getResp2.Project.DisplayName)
	require.NotEqual(t, oldName, getResp2.Project.DisplayName)
	require.Equal(t, refOrNil("new-google-client-id"), getResp2.Project.GoogleOauthClientId)
	require.Empty(t, getResp2.Project.GoogleOauthClientSecret)
	require.Equal(t, refOrNil("new-microsoft-client-id"), getResp2.Project.MicrosoftOauthClientId)
	require.Empty(t, getResp2.Project.MicrosoftOauthClientSecret)
	require.Equal(t, refOrNil("new-github-client-id"), getResp2.Project.GithubOauthClientId)
	require.Empty(t, getResp2.Project.GithubOauthClientSecret)
	require.Equal(t, "cookies.example.com", getResp2.Project.CookieDomain)
	require.ElementsMatch(t, []string{
		"trusted.example.com",
		"another-trusted.example.com",
		fmt.Sprintf("%s.example.com", strings.ReplaceAll(u.ProjectID, "_", "-")), // Default domain for the project
		fmt.Sprintf("%s.%s", strings.ReplaceAll(u.ProjectID, "_", "-"), u.Environment.AuthAppsRootDomain),
	}, getResp2.Project.TrustedDomains)
}

func TestDisableProjectLogins(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	_, err := u.Store.DisableProjectLogins(ctx, &backendv1.DisableProjectLoginsRequest{})
	require.NoError(t, err)

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	var loginsDisabled bool
	err = u.Environment.DB.QueryRow(ctx, "SELECT logins_disabled FROM projects WHERE id = $1::uuid", uuid.UUID(projectUUID).String()).Scan(&loginsDisabled)
	require.NoError(t, err)
	require.True(t, loginsDisabled)
}

func TestEnableProjectLogins(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	_, err := u.Store.EnableProjectLogins(ctx, &backendv1.EnableProjectLoginsRequest{})
	require.NoError(t, err)

	projectUUID, err := idformat.Project.Parse(u.ProjectID)
	require.NoError(t, err)

	var loginsDisabled bool
	err = u.Environment.DB.QueryRow(ctx, "SELECT logins_disabled FROM projects WHERE id = $1::uuid", uuid.UUID(projectUUID).String()).Scan(&loginsDisabled)
	require.NoError(t, err)
	require.False(t, loginsDisabled)
}
