package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
)

type WebhookArgs struct {
	EventType    string                 `json:"event_type"`
	EventPayload map[string]interface{} `json:"event_payload"`
}

func (s *Store) SendWebhook(ctx context.Context, projectID uuid.UUID, args WebhookArgs) error {

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
		message, err := s.svix.Message.Create(ctx, *qProjectWebhookSettings.AppID, models.MessageIn{
			EventType: args.EventType,
			Payload:   args.EventPayload,
		}, nil)
		if err != nil {
			return fmt.Errorf("send webhook via svix: %w", err)
		}

		slog.InfoContext(ctx, "svix_message_created", "message_id", message.Id, "event_type", message.EventType, "event_payload", args.EventPayload)
	}

	return nil
}
