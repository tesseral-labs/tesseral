package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
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

	return &intermediatev1.GetSettingsResponse{
		Settings: s.parseSettings(qProject, qProjectUISettings),
	}, nil
}

func (s *Store) parseSettings(qProject queries.Project, qProjectUISettings queries.ProjectUiSetting) *intermediatev1.Settings {
	return &intermediatev1.Settings{
		Id:                         idformat.ProjectUISettings.Format(qProjectUISettings.ID),
		ProjectId:                  idformat.Project.Format(qProject.ID),
		ProjectDisplayName:         qProject.DisplayName,
		ProjectEmailSendFromDomain: qProject.EmailSendFromDomain,
		LogoUrl:                    fmt.Sprintf("%s/vault-ui-settings-v1/%s/logo", s.userContentBaseUrl, idformat.Project.Format(qProject.ID)),
		PrimaryColor:               derefOrEmpty(qProjectUISettings.PrimaryColor),
		DetectDarkModeEnabled:      qProjectUISettings.DetectDarkModeEnabled,
		DarkModeLogoUrl:            fmt.Sprintf("%s/vault-ui-settings-v1/%s/logo-dark", s.userContentBaseUrl, idformat.Project.Format(qProject.ID)),
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
	}
}
