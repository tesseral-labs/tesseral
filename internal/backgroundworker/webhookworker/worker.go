package webhookworker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
)

type Worker struct {
	Store *store.Store
	river.WorkerDefaults[Args]
}

type Args struct {
	ProjectID string
	EventName string
	Payload   map[string]any
}

func (Args) Kind() string {
	return "webhook"
}

func (w *Worker) Work(ctx context.Context, job *river.Job[Args]) error {
	slog.InfoContext(ctx, "work", "args", job.Args)

	if err := w.Store.SendWebhook(ctx, &store.SendWebhookRequest{
		ProjectID: job.Args.ProjectID,
		EventType: job.Args.EventName,
		Payload:   job.Args.Payload,
	}); err != nil {
		return fmt.Errorf("store: %w", err)
	}

	return nil
}
