package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type WebhookArgs struct {
	ProjectID    string                 `json:"project_id"`
	EventName    string                 `json:"event_name"`
	EventPayload map[string]interface{} `json:"event_payload"`
}

func (WebhookArgs) Kind() string { return "webhook" }

func (s *Store) SendWebhook(ctx context.Context, args WebhookArgs) error {
	projectID, err := idformat.Project.Parse(args.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	qProjectWebhookSettings, err := s.q.GetProjectWebhookSettings(ctx, projectID)
	if err != nil {
		// We want to ignore this error if the project does not have webhook settings
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("get project by id: %w", err)
	}

	if qProjectWebhookSettings.DirectWebhookUrl != nil {
		// TODO: Handle direct webhooks if they exist
		return fmt.Errorf("direct webhooks are not yet implemented")
	}

	if qProjectWebhookSettings.AppID != nil && *qProjectWebhookSettings.AppID != "" {
		// If the project has an app ID, we can send the webhook via Svix
		if _, err := s.svix.Message.Create(ctx, *qProjectWebhookSettings.AppID, models.MessageIn{
			EventType: args.EventName,
			Payload:   args.EventPayload,
		}, nil); err != nil {
			return fmt.Errorf("send webhook via svix: %w", err)
		}
	}

	return nil
}
