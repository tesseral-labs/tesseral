package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetProjectUISettings(ctx context.Context, req *backendv1.GetProjectUISettingsRequest) (*backendv1.GetProjectUISettingsResponse, error) {
	projectID := authn.ProjectID(ctx)

	qProjectUISettings, err := s.q.GetProjectUISettings(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project ui settings not found", fmt.Errorf("failed to get project ui settings: %w", err))
		}

		return nil, fmt.Errorf("failed to get project ui settings: %w", err)
	}

	return &backendv1.GetProjectUISettingsResponse{
		ProjectUiSettings: s.parseProjectUISettings(qProjectUISettings),
	}, nil
}

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

	qUpdatedProjectUISettings, err := q.UpdateProjectUISettings(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update project ui settings: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	res := &backendv1.UpdateProjectUISettingsResponse{
		Id:                    idformat.ProjectUISettings.Format(qUpdatedProjectUISettings.ID),
		ProjectId:             idformat.Project.Format(projectID),
		CreateTime:            timestamppb.New(*qUpdatedProjectUISettings.CreateTime),
		UpdateTime:            timestamppb.New(*qUpdatedProjectUISettings.UpdateTime),
		DarkModePrimaryColor:  derefOrEmpty(qUpdatedProjectUISettings.DarkModePrimaryColor),
		DetectDarkModeEnabled: qUpdatedProjectUISettings.DetectDarkModeEnabled,
		PrimaryColor:          derefOrEmpty(qUpdatedProjectUISettings.PrimaryColor),
	}

	// generate a presigned URL for the dark mode logo file
	darkModeLogoPresignedUploadUrl, err := s.getPresignedUrlForFile(ctx, fmt.Sprintf("dark_mode_logos_v1/%s/dark_mode_logo", projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for dark mode logo file: %w", err)
	}
	res.DarkModeLogoPresignedUploadUrl = darkModeLogoPresignedUploadUrl

	// generate a presigned URL for the favicon file
	faviconPresignedUploadUrl, err := s.getPresignedUrlForFile(ctx, fmt.Sprintf("favicons_v1/%s/favicon", projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for favicon file: %w", err)
	}
	res.FaviconPresignedUploadUrl = faviconPresignedUploadUrl

	// generate a presigned URL for the logo file
	logoPresignedUploadUrl, err := s.getPresignedUrlForFile(ctx, fmt.Sprintf("logos_v1/%s/logo", projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for logo file: %w", err)
	}
	res.LogoPresignedUploadUrl = logoPresignedUploadUrl

	return res, nil
}

func (s *Store) getPresignedUrlForFile(ctx context.Context, fileKey string) (string, error) {
	req, err := s.s3PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.s3UserContentBucketName),
		Key:    aws.String(fileKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Minute // set expiry to one minute
	})

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return req.URL, nil
}

func (s *Store) parseProjectUISettings(pus queries.ProjectUiSetting) *backendv1.ProjectUISettings {
	projectID := idformat.Project.Format(pus.ProjectID)

	return &backendv1.ProjectUISettings{
		PrimaryColor:          derefOrEmpty(pus.PrimaryColor),
		DetectDarkModeEnabled: pus.DetectDarkModeEnabled,
		DarkModePrimaryColor:  derefOrEmpty(pus.DarkModePrimaryColor),
		LogoUrl:               fmt.Sprintf("%s/logos_v1/%s/logo", s.userContentUrl, projectID),
		FaviconUrl:            fmt.Sprintf("%s/favicons_v1/%s/favicon", s.userContentUrl, projectID),
		DarkModeLogoUrl:       fmt.Sprintf("%s/dark_mode_logos_v1/%s/dark_mode_logo", s.userContentUrl, projectID),
		CreateTime:            timestamppb.New(*pus.CreateTime),
		UpdateTime:            timestamppb.New(*pus.UpdateTime),
	}
}
