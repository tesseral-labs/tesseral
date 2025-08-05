package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type SendWebhookRequest struct {
	ProjectID string
	EventType string
	Payload   map[string]any
}

func (s *Store) SendWebhook(ctx context.Context, req *SendWebhookRequest) error {
	projectID, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	qProjectWebhookSettings, err := s.q().GetProjectWebhookSettings(ctx, projectID)
	if err != nil {
		// We want to ignore this error if the project does not have webhook settings
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("get project by id: %w", err)
	}

	slog.InfoContext(ctx, "project_webhook_settings", "direct_webhook_url", qProjectWebhookSettings.DirectWebhookUrl, "svix_app_id", qProjectWebhookSettings.AppID)

	if qProjectWebhookSettings.DirectWebhookUrl != nil {
		slog.InfoContext(ctx, "handle_direct_webhook", "url", *qProjectWebhookSettings.DirectWebhookUrl)

		body, err := json.Marshal(req.Payload)
		if err != nil {
			return fmt.Errorf("marshal webhook body: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, *qProjectWebhookSettings.DirectWebhookUrl, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("new http request: %w", err)
		}

		httpReq.Header.Set("Content-Type", "application/json")

		httpRes, err := s.DirectWebhookHTTPClient.Do(httpReq)
		if err != nil {
			return fmt.Errorf("send http request: %w", err)
		}

		defer func() { _ = httpRes.Body.Close() }()

		if httpRes.StatusCode != http.StatusOK {
			return fmt.Errorf("bad response status code: %s", httpRes.Status)
		}
	} else if qProjectWebhookSettings.AppID != nil && *qProjectWebhookSettings.AppID != "" {
		slog.InfoContext(ctx, "handle_svix_webhook", "svix_app_id", *qProjectWebhookSettings.AppID)

		// If the project has an app ID, we can send the webhook via Svix
		message, err := s.Svix.Message.Create(ctx, *qProjectWebhookSettings.AppID, models.MessageIn{
			EventType: req.EventType,
			Payload:   req.Payload,
		}, nil)
		if err != nil {
			return fmt.Errorf("send webhook via svix: %w", err)
		}

		slog.InfoContext(ctx, "svix_message_created", "message_id", message.Id)
	}

	return nil
}
