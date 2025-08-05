package workers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type BackgroundWorkerArgs struct {
	ProjectID      string            `json:"project_id"`
	EventName      string            `json:"event_name"`
	WebhookPayload store.WebhookArgs `json:"webhook_payload"`
}

func (args BackgroundWorkerArgs) Kind() string {
	return "background_worker"
}

type BackgroundWorker struct {
	river.WorkerDefaults[BackgroundWorkerArgs]
	Store *store.Store
}

func NewBackgroundWorker(store *store.Store) *BackgroundWorker {
	return &BackgroundWorker{
		Store: store,
	}
}

func (w *BackgroundWorker) Work(ctx context.Context, job *river.Job[BackgroundWorkerArgs]) error {
	slog.InfoContext(ctx, "work", "args", job.Args)

	projectID, err := idformat.Project.Parse(job.Args.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	switch job.Args.EventName {
	case "send_webhook":
		if err := w.Store.SendWebhook(ctx, projectID, job.Args.WebhookPayload); err != nil {
			return fmt.Errorf("send webhook: %w", err)
		}
	// TODO: Handle email events
	case "send_email":
		return fmt.Errorf("email jobs not implemented")
	default:
		return fmt.Errorf("unknown event name: %s", job.Args.EventName)
	}

	return nil
}
