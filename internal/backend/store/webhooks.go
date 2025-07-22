package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/svix/svix-webhooks/go/models"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) createProjectWebhookSettings(ctx context.Context, q *queries.Queries, qProject queries.Project) (*backendv1.ProjectWebhookSettings, error) {
	// Create a new Svix application
	svixApplication, err := s.svixClient.Application.Create(ctx, models.ApplicationIn{
		Name: qProject.DisplayName,
	}, nil)
	if err != nil {
		return nil, err
	}

	qWebhook, err := q.CreateProjectWebhookSettings(ctx, queries.CreateProjectWebhookSettingsParams{
		ID:        uuid.New(),
		ProjectID: qProject.ID,
		AppID:     refOrNil(svixApplication.Id),
	})
	if err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}

	return parseProjectWebhookSettings(qWebhook), nil
}

func parseProjectWebhookSettings(qWebhook queries.ProjectWebhookSetting) *backendv1.ProjectWebhookSettings {
	return &backendv1.ProjectWebhookSettings{
		Id:         idformat.ProjectWebhookSettings.Format(qWebhook.ID),
		AppId:      derefOrEmpty(qWebhook.AppID),
		CreateTime: timestamppb.New(*qWebhook.CreateTime),
		UpdateTime: timestamppb.New(*qWebhook.UpdateTime),
	}

}
