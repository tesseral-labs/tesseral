package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestGetProjectUISettings_Exists(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	resp, err := u.Store.GetProjectUISettings(ctx, &backendv1.GetProjectUISettingsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp.ProjectUiSettings)
	require.NotEmpty(t, resp.ProjectUiSettings.Id)
	require.Equal(t, u.ProjectID, resp.ProjectUiSettings.ProjectId)
	require.NotEmpty(t, resp.ProjectUiSettings.CreateTime)
	require.NotEmpty(t, resp.ProjectUiSettings.UpdateTime)
}

func TestUpdateProjectUISettings_UpdatesFields(t *testing.T) {
	t.Parallel()
	ctx, u := newTestUtil(t)

	getResp, err := u.Store.GetProjectUISettings(ctx, &backendv1.GetProjectUISettingsRequest{})
	require.NoError(t, err)
	oldColor := getResp.ProjectUiSettings.PrimaryColor

	updateResp, err := u.Store.UpdateProjectUISettings(ctx, &backendv1.UpdateProjectUISettingsRequest{
		PrimaryColor:            refOrNil("#123456"),
		DetectDarkModeEnabled:   refOrNil(true),
		DarkModePrimaryColor:    refOrNil("#654321"),
		AutoCreateOrganizations: refOrNil(true),
		LogInLayout:             "side_by_side",
	})
	require.NoError(t, err)
	require.NotEqual(t, oldColor, updateResp.PrimaryColor)
	require.Equal(t, "#123456", updateResp.PrimaryColor)
	require.True(t, updateResp.DetectDarkModeEnabled)
	require.Equal(t, "#654321", updateResp.DarkModePrimaryColor)
	require.True(t, updateResp.AutoCreateOrganizations)
	require.Equal(t, "side_by_side", updateResp.LogInLayout)

	getResp2, err := u.Store.GetProjectUISettings(ctx, &backendv1.GetProjectUISettingsRequest{})
	require.NoError(t, err)
	require.NotNil(t, getResp2.ProjectUiSettings)
	require.Equal(t, "#123456", getResp2.ProjectUiSettings.PrimaryColor)
	require.True(t, getResp2.ProjectUiSettings.DetectDarkModeEnabled)
	require.Equal(t, "#654321", getResp2.ProjectUiSettings.DarkModePrimaryColor)
	require.True(t, getResp2.ProjectUiSettings.AutoCreateOrganizations)
	require.Equal(t, "side_by_side", getResp2.ProjectUiSettings.LogInLayout)
}
