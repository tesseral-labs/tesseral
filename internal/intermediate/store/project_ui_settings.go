package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetProjectUISettings(ctx context.Context, req *intermediatev1.GetProjectUISettingsRequest) (*intermediatev1.GetProjectUISettingsResponse, error) {
	projectID := authn.ProjectID(ctx)

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project ui settings: %w", err)
	}
	defer rollback()

	qProjectUISettings, err := q.GetProjectUISettings(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project ui settings: %w", err)
	}

	return &intermediatev1.GetProjectUISettingsResponse{
		ProjectUiSettings: s.parseProjectUISettings(qProjectUISettings),
	}, nil
}

func (s *Store) parseProjectUISettings(pus queries.ProjectUiSetting) *intermediatev1.ProjectUISettings {
	projectID := idformat.Project.Format(pus.ProjectID)

	return &intermediatev1.ProjectUISettings{
		Id:                    idformat.ProjectUISettings.Format(pus.ID),
		ProjectId:             projectID,
		DarkModeLogoUrl:       fmt.Sprintf("%s/dark_mode_logos_v1/%s/dark_mode_logo", s.userContentUrl, projectID),
		DetectDarkModeEnabled: pus.DetectDarkModeEnabled,
		DarkModePrimaryColor:  derefOrEmpty(pus.DarkModePrimaryColor),
		FaviconUrl:            fmt.Sprintf("%s/faviconss_v1/%s/favicon", s.userContentUrl, projectID),
		LogoUrl:               fmt.Sprintf("%s/logos_v1/%s/logo", s.userContentUrl, projectID),
		PrimaryColor:          derefOrEmpty(pus.PrimaryColor),
	}
}
