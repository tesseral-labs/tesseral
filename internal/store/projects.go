package store

import (
	"context"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/store/queries"
)

type UpdateProjectRequest struct {
	ProjectID       string
	VaultDomain     string
	LogoURL         string
	DarkModeLogoURL string
}

func (s *Store) UpdateProject(ctx context.Context, req *UpdateProjectRequest) error {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	if req.VaultDomain != "" {
		if _, err := q.UpdateProject(ctx, queries.UpdateProjectParams{
			ID:          projectID,
			VaultDomain: req.VaultDomain,
		}); err != nil {
			return fmt.Errorf("update project: %w", err)
		}
	}

	if req.LogoURL != "" {
		if _, err := q.UpdateProjectLogoURL(ctx, queries.UpdateProjectLogoURLParams{
			ProjectID: projectID,
			LogoUrl:   &req.LogoURL,
		}); err != nil {
			return fmt.Errorf("update project logo url: %w", err)
		}
	}

	if req.DarkModeLogoURL != "" {
		if _, err := q.UpdateProjectDarkModeLogoURL(ctx, queries.UpdateProjectDarkModeLogoURLParams{
			ProjectID:       projectID,
			DarkModeLogoUrl: &req.DarkModeLogoURL,
		}); err != nil {
			return fmt.Errorf("update project dark mode logo url: %w", err)
		}
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
