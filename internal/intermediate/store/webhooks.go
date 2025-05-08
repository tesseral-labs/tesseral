package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/svix/svix-webhooks/go/models"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) createSvixApplication(ctx context.Context, displayName string) (*models.ApplicationOut, error) {
	// Create a new Svix application
	app, err := s.svixClient.Application.Create(ctx, models.ApplicationIn{
		Name: displayName,
	}, nil)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *Store) createProjectWebhookSettings(ctx context.Context, q *queries.Queries, qProject queries.Project) (*intermediatev1.ProjectWebhookSettings, error) {
	svixApplication, err := s.createSvixApplication(ctx, qProject.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("create svix application: %w", err)
	}

	qWebhook, err := q.CreateProjectWebhookSettings(ctx, queries.CreateProjectWebhookSettingsParams{
		ID:        uuid.New(),
		ProjectID: qProject.ID,
		AppID:     svixApplication.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}

	return parseProjectWebhookSettings(qWebhook), nil
}

func parseProjectWebhookSettings(qWebhook queries.ProjectWebhookSetting) *intermediatev1.ProjectWebhookSettings {
	return &intermediatev1.ProjectWebhookSettings{
		Id:         idformat.ProjectWebhookSettings.Format(qWebhook.ID),
		AppId:      qWebhook.AppID,
		CreateTime: timestamppb.New(*qWebhook.CreateTime),
		UpdateTime: timestamppb.New(*qWebhook.UpdateTime),
	}

}
