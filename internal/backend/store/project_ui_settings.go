package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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

	logoURL := fmt.Sprintf("%s/vault-ui-settings-v1/%s/logo", s.userContentBaseUrl, idformat.Project.Format(projectID))
	logoKey := fmt.Sprintf("vault-ui-settings-v1/%s/logo", idformat.Project.Format(projectID))
	logoExists, err := s.getUserContentFileExists(ctx, logoKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check if logo file exists: %w", err)
	}
	if !logoExists {
		logoURL = ""
	}

	darkModeLogoURL := fmt.Sprintf("%s/vault-ui-settings-v1/%s/logo-dark", s.userContentBaseUrl, idformat.Project.Format(projectID))
	darkModeLogoKey := fmt.Sprintf("vault-ui-settings-v1/%s/logo-dark", idformat.Project.Format(projectID))
	darkModeLogoExists, err := s.getUserContentFileExists(ctx, darkModeLogoKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check if dark mode logo file exists: %w", err)
	}
	if !darkModeLogoExists {
		darkModeLogoURL = ""
	}

	return &backendv1.GetProjectUISettingsResponse{
		ProjectUiSettings: &backendv1.ProjectUISettings{
			PrimaryColor:          derefOrEmpty(qProjectUISettings.PrimaryColor),
			DetectDarkModeEnabled: qProjectUISettings.DetectDarkModeEnabled,
			DarkModePrimaryColor:  derefOrEmpty(qProjectUISettings.DarkModePrimaryColor),
			LogInLayout:           string(qProjectUISettings.LogInLayout),
			LogoUrl:               logoURL,
			DarkModeLogoUrl:       darkModeLogoURL,
			CreateTime:            timestamppb.New(*qProjectUISettings.CreateTime),
			UpdateTime:            timestamppb.New(*qProjectUISettings.UpdateTime),
		},
	}, nil
}

func (s *Store) UpdateProjectUISettings(ctx context.Context, req *backendv1.UpdateProjectUISettingsRequest) (*backendv1.UpdateProjectUISettingsResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer rollback()

	qProjectUISettings, err := q.GetProjectUISettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project ui settings not found", fmt.Errorf("failed to get project ui settings: %w", err))
		}

		return nil, fmt.Errorf("failed to get project ui settings: %w", err)
	}

	updates := queries.UpdateProjectUISettingsParams{
		ID:        qProjectUISettings.ID,
		ProjectID: authn.ProjectID(ctx),
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

	if req.LogInLayout != "" {
		updates.LogInLayout = queries.LogInLayout(req.LogInLayout)
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
		ProjectId:             idformat.Project.Format(authn.ProjectID(ctx)),
		CreateTime:            timestamppb.New(*qUpdatedProjectUISettings.CreateTime),
		UpdateTime:            timestamppb.New(*qUpdatedProjectUISettings.UpdateTime),
		DarkModePrimaryColor:  derefOrEmpty(qUpdatedProjectUISettings.DarkModePrimaryColor),
		DetectDarkModeEnabled: qUpdatedProjectUISettings.DetectDarkModeEnabled,
		PrimaryColor:          derefOrEmpty(qUpdatedProjectUISettings.PrimaryColor),
	}

	// generate a presigned URL for the dark mode logo file
	darkModeLogoPresignedUploadUrl, err := s.getPresignedUrlForFile(ctx, fmt.Sprintf("vault-ui-settings-v1/%s/logo-dark", idformat.Project.Format(authn.ProjectID(ctx))), req.DarkModeLogoContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for dark mode logo file: %w", err)
	}
	res.DarkModeLogoPresignedUploadUrl = darkModeLogoPresignedUploadUrl

	// generate a presigned URL for the logo file
	logoPresignedUploadUrl, err := s.getPresignedUrlForFile(ctx, fmt.Sprintf("vault-ui-settings-v1/%s/logo", idformat.Project.Format(authn.ProjectID(ctx))), req.LogoContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to get presigned URL for logo file: %w", err)
	}
	res.LogoPresignedUploadUrl = logoPresignedUploadUrl

	return res, nil
}

func (s *Store) getPresignedUrlForFile(ctx context.Context, fileKey string, contentType *string) (string, error) {
	putObjectParams := &s3.PutObjectInput{
		Bucket: aws.String(s.s3UserContentBucketName),
		Key:    aws.String(fileKey),
	}

	if contentType != nil {
		putObjectParams.ContentType = aws.String(*contentType)
		putObjectParams.Metadata = map[string]string{
			"x-amz-meta-trigger": "true", // anything to force header signing path
		}
	}

	req, err := s.s3PresignClient.PresignPutObject(ctx, putObjectParams, func(opts *s3.PresignOptions) {
		opts.Expires = time.Minute // set expiry to one minute
	})

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return req.URL, nil
}

func (s *Store) getUserContentFileExists(ctx context.Context, key string) (bool, error) {
	if _, err := s.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.s3UserContentBucketName,
		Key:    &key,
	}); err != nil {
		var notFoundErr *types.NotFound
		if errors.As(err, &notFoundErr) {
			return false, nil
		}

		// Return other errors
		return false, fmt.Errorf("failed to check if user content file exists: %w", err)
	}

	return true, nil
}
