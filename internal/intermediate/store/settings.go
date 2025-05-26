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
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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

	logoKey := fmt.Sprintf("vault-ui-settings-v1/%s/logo", idformat.Project.Format(projectID))
	logoURL := fmt.Sprintf("%s/%s", s.userContentBaseUrl, logoKey)
	logoExists, err := s.getUserContentFileExists(ctx, logoKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check if logo file exists: %w", err)
	}
	if !logoExists {
		logoURL = ""
	} else {
		logoURL, err = s.buildPresignedGetUrlForFile(ctx, logoKey)
		if err != nil {
			return nil, fmt.Errorf("failed to build presigned URL for logo file: %w", err)
		}
	}

	darkModeLogoKey := fmt.Sprintf("vault-ui-settings-v1/%s/logo-dark", idformat.Project.Format(projectID))
	darkModeLogoURL := fmt.Sprintf("%s/%s", s.userContentBaseUrl, darkModeLogoKey)
	darkModeLogoExists, err := s.getUserContentFileExists(ctx, darkModeLogoKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check if dark mode logo file exists: %w", err)
	}
	if !darkModeLogoExists {
		darkModeLogoURL = ""
	} else {
		darkModeLogoURL, err = s.buildPresignedGetUrlForFile(ctx, darkModeLogoKey)
		if err != nil {
			return nil, fmt.Errorf("failed to build presigned URL for dark mode logo file: %w", err)
		}
	}

	return &intermediatev1.GetSettingsResponse{
		Settings: &intermediatev1.Settings{
			Id:                         idformat.ProjectUISettings.Format(qProjectUISettings.ID),
			ProjectId:                  idformat.Project.Format(qProject.ID),
			ProjectDisplayName:         qProject.DisplayName,
			ProjectEmailSendFromDomain: qProject.EmailSendFromDomain,
			LogoUrl:                    logoURL,
			PrimaryColor:               derefOrEmpty(qProjectUISettings.PrimaryColor),
			DetectDarkModeEnabled:      qProjectUISettings.DetectDarkModeEnabled,
			DarkModeLogoUrl:            darkModeLogoURL,
			DarkModePrimaryColor:       derefOrEmpty(qProjectUISettings.DarkModePrimaryColor),
			LogInLayout:                string(qProjectUISettings.LogInLayout),
			LogInWithEmail:             qProject.LogInWithEmail,
			LogInWithGoogle:            qProject.LogInWithGoogle,
			LogInWithGithub:            qProject.LogInWithGithub,
			LogInWithMicrosoft:         qProject.LogInWithMicrosoft,
			LogInWithPassword:          qProject.LogInWithPassword,
			LogInWithSaml:              qProject.LogInWithSaml,
			RedirectUri:                qProject.RedirectUri,
			AfterLoginRedirectUri:      qProject.AfterLoginRedirectUri,
			AfterSignupRedirectUri:     qProject.AfterSignupRedirectUri,
			AutoCreateOrganizations:    qProjectUISettings.AutoCreateOrganizations,
		},
	}, nil
}

func (s *Store) buildPresignedGetUrlForFile(ctx context.Context, fileKey string) (string, error) {
	req, err := s.s3PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.s3UserContentBucketName),
		Key:    aws.String(fileKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Hour * 12 // set expiry to 12 hours
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
