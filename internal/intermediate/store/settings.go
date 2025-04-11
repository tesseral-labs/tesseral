package store

import (
	"context"
	"errors"
	"fmt"

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
			LogInWithMicrosoft:         qProject.LogInWithMicrosoft,
			LogInWithPassword:          qProject.LogInWithPassword,
			LogInWithSaml:              qProject.LogInWithSaml,
			RedirectUri:                qProject.RedirectUri,
			AfterLoginRedirectUri:      qProject.AfterLoginRedirectUri,
			AfterSignupRedirectUri:     qProject.AfterSignupRedirectUri,
		},
	}, nil
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
