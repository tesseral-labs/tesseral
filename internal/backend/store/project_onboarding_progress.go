package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) GetProjectOnboardingProgress(ctx context.Context, req *backendv1.GetProjectOnboardingProgressRequest) (*backendv1.GetProjectOnboardingProgressResponse, error) {
	if err := validateIsConsoleSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is console session: %w", err)
	}

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer rollback()

	qProgress, err := q.GetProjectOnboardingProgress(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &backendv1.GetProjectOnboardingProgressResponse{
				Progress: &backendv1.ProjectOnboardingProgress{
					ProjectId: idformat.Project.Format(authn.ProjectID(ctx)),
				},
			}, nil
		}
		return nil, fmt.Errorf("get project onboarding progress: %w", err)
	}

	return &backendv1.GetProjectOnboardingProgressResponse{
		Progress: parseProjectOnboardingProgress(qProgress),
	}, nil
}

func (s *Store) UpdateProjectOnboardingProgress(ctx context.Context, req *backendv1.UpdateProjectOnboardingProgressRequest) (*backendv1.UpdateProjectOnboardingProgressResponse, error) {
	if err := validateIsConsoleSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is console session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer rollback()

	updates := queries.UpsertProjectOnboardingProgressParams{
		ProjectID: authn.ProjectID(ctx),
	}

	if req.Progress.ConfigureAuthenticationTime != nil {
		updates.ConfigureAuthenticationTime = refOrNil(req.Progress.ConfigureAuthenticationTime.AsTime())
	}
	if req.Progress.LogInToVaultTime != nil {
		updates.LogInToVaultTime = refOrNil(req.Progress.LogInToVaultTime.AsTime())
	}
	if req.Progress.ManageOrganizationsTime != nil {
		updates.ManageOrganizationsTime = refOrNil(req.Progress.ManageOrganizationsTime.AsTime())
	}
	if req.Progress.OnboardingSkipped != nil {
		updates.OnboardingSkipped = req.Progress.OnboardingSkipped
	}

	qProgress, err := q.UpsertProjectOnboardingProgress(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update project onboarding progress: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectOnboardingProgressResponse{
		Progress: parseProjectOnboardingProgress(qProgress),
	}, nil
}

func parseProjectOnboardingProgress(qProgress queries.ProjectOnboardingProgress) *backendv1.ProjectOnboardingProgress {
	return &backendv1.ProjectOnboardingProgress{
		ProjectId:                   idformat.Project.Format(qProgress.ProjectID),
		ConfigureAuthenticationTime: timestampOrNil(qProgress.ConfigureAuthenticationTime),
		LogInToVaultTime:            timestampOrNil(qProgress.LogInToVaultTime),
		ManageOrganizationsTime:     timestampOrNil(qProgress.ManageOrganizationsTime),
		OnboardingSkipped:           qProgress.OnboardingSkipped,
		CreateTime:                  timestampOrNil(qProgress.CreateTime),
		UpdateTime:                  timestampOrNil(qProgress.UpdateTime),
	}
}
