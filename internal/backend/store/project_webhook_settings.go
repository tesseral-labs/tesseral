package store

import (
	"context"
	"fmt"

	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Store) GetProjectWebhookManagementURL(ctx context.Context, req *backendv1.GetProjectWebhookManagementURLRequest) (*backendv1.GetProjectWebhookManagementURLResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	qProjectWebhookSettings, err := q.GetProjectWebhookSettings(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project webhook settings: %w", err)
	}

	dashboard, err := s.svixClient.Authentication.AppPortalAccess(ctx, qProjectWebhookSettings.AppID, models.AppPortalAccessIn{}, nil)
	if err != nil {
		return nil, fmt.Errorf("get app portal access: %w", err)
	}

	return &backendv1.GetProjectWebhookManagementURLResponse{
		Url: dashboard.Url,
	}, nil
}
