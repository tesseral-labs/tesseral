package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/ssoready/conf"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/workers"
	"github.com/tesseral-labs/tesseral/internal/dbconn"
	"github.com/tesseral-labs/tesseral/internal/loadenv"
)

func main() {
	loadenv.LoadEnv()

	config := struct {
		DB         dbconn.Config `conf:"db,noredact"`
		SvixApiKey string        `conf:"svix_api_key,noredact"`
	}{}

	conf.Load(&config)

	slog.Info("config", "config", conf.Redact(config))

	db, err := dbconn.Open(context.Background(), config.DB)
	if err != nil {
		panic(fmt.Errorf("open database: %w", err))
	}
	defer db.Close()

	svixClient, err := svix.New(config.SvixApiKey, nil)
	if err != nil {
		panic(fmt.Errorf("create svix client: %w", err))
	}

	riverWorkers := river.NewWorkers()

	if err := river.AddWorkerSafely(riverWorkers, &workers.WebhookWorker{
		Store: store.New(store.NewStoreParams{
			DB:   db,
			Svix: svixClient,
		}),
	}); err != nil {
		panic(err)
	}

	// TODO: Add key rotation worker
}
