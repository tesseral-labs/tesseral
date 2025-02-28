package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type GetPublishableKeyConfigurationResponse struct {
	ProjectID   string `json:"projectId"`
	VaultDomain string `json:"vaultDomain"`
}

func (s *Store) GetPublishableKeyConfiguration(ctx context.Context, publishableKey string) (*GetPublishableKeyConfigurationResponse, error) {
	id, err := idformat.PublishableKey.Parse(publishableKey)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid publishable key id", fmt.Errorf("parse publishable key id: %w", err))
	}

	qConfig, err := s.q.GetPublishableKeyConfiguration(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("publishable key not found", fmt.Errorf("get publishable key configuration: %w", err))
		}
		return nil, fmt.Errorf("get publishable key configuration: %w", err)
	}

	return &GetPublishableKeyConfigurationResponse{
		ProjectID:   idformat.Project.Format(qConfig.ProjectID),
		VaultDomain: qConfig.VaultDomain,
	}, nil
}
