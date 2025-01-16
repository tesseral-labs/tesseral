package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
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

func (s *Store) getURLForFileKey(fileKey string) string {
	if fileKey == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", s.userContentUrl, fileKey)
}

func (s *Store) parseProjectUISettings(pus queries.ProjectUiSetting) *intermediatev1.ProjectUISettings {
	return &intermediatev1.ProjectUISettings{
		DarkModeLogoUrl:       s.getURLForFileKey(derefOrEmpty(pus.DarkModeLogoFileKey)),
		DetectDarkModeEnabled: pus.DetectDarkModeEnabled,
		DarkModePrimaryColor:  derefOrEmpty(pus.DarkModePrimaryColor),
		FaviconUrl:            s.getURLForFileKey(derefOrEmpty(pus.FaviconFileKey)),
		LogoUrl:               s.getURLForFileKey(derefOrEmpty(pus.LogoFileKey)),
		PrimaryColor:          derefOrEmpty(pus.PrimaryColor),
	}
}
