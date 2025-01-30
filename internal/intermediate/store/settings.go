package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetSettings(ctx context.Context, req *intermediatev1.GetSettingsRequest) (*intermediatev1.GetSettingsResponse, error) {
	projectID := authn.ProjectID(ctx)

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project ui settings: %w", err)
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}
		return nil, fmt.Errorf("get project: %w", err)
	}

	qProjectUISettings, err := q.GetProjectUISettings(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project ui settings: %w", err)
	}

	return &intermediatev1.GetSettingsResponse{
		Settings: s.parseSettings(qProject, qProjectUISettings),
	}, nil
}

func (s *Store) parseSettings(qProject queries.Project, qProjectUISettings queries.ProjectUiSetting) *intermediatev1.Settings {
	projectID := idformat.Project.Format(qProject.ID)

	return &intermediatev1.Settings{
		Id:                        idformat.ProjectUISettings.Format(qProjectUISettings.ID),
		ProjectId:                 projectID,
		DarkModeLogoUrl:           fmt.Sprintf("%s/dark_mode_logos_v1/%s/dark_mode_logo", s.userContentBaseUrl, projectID),
		DetectDarkModeEnabled:     qProjectUISettings.DetectDarkModeEnabled,
		DarkModePrimaryColor:      derefOrEmpty(qProjectUISettings.DarkModePrimaryColor),
		FaviconUrl:                fmt.Sprintf("%s/faviconss_v1/%s/favicon", s.userContentBaseUrl, projectID),
		LogoUrl:                   fmt.Sprintf("%s/logos_v1/%s/logo", s.userContentBaseUrl, projectID),
		PrimaryColor:              derefOrEmpty(qProjectUISettings.PrimaryColor),
		LogInWithGoogleEnabled:    qProject.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: qProject.LogInWithMicrosoftEnabled,
		LogInWithPasswordEnabled:  qProject.LogInWithPasswordEnabled,
	}
}
