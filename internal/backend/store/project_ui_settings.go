package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) UpdateProjectUISettings(ctx context.Context, req *backendv1.UpdateProjectUISettingsRequest) (*backendv1.UpdateProjectUISettingsResponse, error) {
	projectID := authn.ProjectID(ctx)

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer rollback()

	qProjectUISettings, err := q.GetProjectUISettings(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project ui settings not found", fmt.Errorf("failed to get project ui settings: %w", err))
		}

		return nil, fmt.Errorf("failed to get project ui settings: %w", err)
	}

	updates := queries.UpdateProjectUISettingsParams{
		ID:        qProjectUISettings.ID,
		ProjectID: projectID,
	}

	if req.PrimaryColor != nil {
		updates.PrimaryColor = req.PrimaryColor
	}

	if req.DetectDarkModeEnabled != qProjectUISettings.DetectDarkModeEnabled {
		updates.DetectDarkModeEnabled = req.DetectDarkModeEnabled
	}

	if req.DarkModePrimaryColor != nil {
		updates.DarkModePrimaryColor = req.DarkModePrimaryColor
	}

	// Process image uploads
	if req.Logo != nil {
		logoFileKey, err := s.processImageUpload(ctx, "logo", req.Logo)
		if err != nil {
			return nil, fmt.Errorf("process image upload: %w", err)
		}

		updates.LogoFileKey = &logoFileKey
	}

	if req.Favicon != nil {
		faviconFileKey, err := s.processImageUpload(ctx, "favicon", req.Favicon)
		if err != nil {
			return nil, fmt.Errorf("process image upload: %w", err)
		}

		updates.FaviconFileKey = &faviconFileKey
	}

	if req.DarkModeLogo != nil {
		darkModeLogoFileKey, err := s.processImageUpload(ctx, "dark_mode_logo", req.DarkModeLogo)
		if err != nil {
			return nil, fmt.Errorf("process image upload: %w", err)
		}

		updates.DarkModeLogoFileKey = &darkModeLogoFileKey
	}

	qUpdatedProjectUISettings, err := q.UpdateProjectUISettings(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update project ui settings: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectUISettingsResponse{
		ProjectUiSettings: s.parseProjectUISettings(qUpdatedProjectUISettings),
	}, nil
}

func (s *Store) processImageUpload(ctx context.Context, imageType string, req *backendv1.ImageUploadRequest) (string, error) {
	fileKey, err := getFileKeyForImageType(authn.ProjectID(ctx), imageType, req.MimeType)
	if err != nil {
		return "", fmt.Errorf("get filename for image type: %w", err)
	}

	err = s.uploadToS3(ctx, req, fileKey)
	if err != nil {
		return "", fmt.Errorf("upload file to S3: %w", err)
	}

	return fileKey, nil
}

func (s *Store) parseProjectUISettings(pus queries.ProjectUiSetting) *backendv1.ProjectUISettings {
	return &backendv1.ProjectUISettings{
		PrimaryColor:          *pus.PrimaryColor,
		DetectDarkModeEnabled: pus.DetectDarkModeEnabled,
		DarkModePrimaryColor:  *pus.DarkModePrimaryColor,
		LogoUrl:               s.getURLForFileKey(*pus.LogoFileKey),
		FaviconUrl:            s.getURLForFileKey(*pus.FaviconFileKey),
		DarkModeLogoUrl:       s.getURLForFileKey(*pus.DarkModeLogoFileKey),
		CreateTime:            timestamppb.New(*pus.CreateTime),
		UpdateTime:            timestamppb.New(*pus.UpdateTime),
	}
}
