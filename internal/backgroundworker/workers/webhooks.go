package workers

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
)

type WebhookWorker struct {
	river.WorkerDefaults[store.WebhookArgs]
	Store *store.Store
}

func (w *WebhookWorker) Work(ctx context.Context, job *river.Job[store.WebhookArgs]) error {
	if err := w.Store.SendWebhook(ctx, job.Args); err != nil {
		return err
	}

	return nil
}
