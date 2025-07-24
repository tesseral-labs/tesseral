package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
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

	db, err := dbconn.Open(context.Background(), config.DB)
	if err != nil {
		panic(fmt.Errorf("open database: %w", err))
	}
	defer db.Close()

	svixClient, err := svix.New(config.SvixApiKey, nil)
	if err != nil {
		panic(fmt.Errorf("create svix client: %w", err))
	}

	store_ := store.New(store.NewStoreParams{
		DB:   db,
		Svix: svixClient,
	})

	riverWorkers := river.NewWorkers()

	if err := river.AddWorkerSafely(riverWorkers, workers.NewBackgroundWorker(store_)); err != nil {
		panic(err)
	}

	riverClient, err := river.NewClient(riverpgxv5.New(db), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {
				MaxWorkers: 100,
			},
		},
		Workers: riverWorkers,
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	defer func() {
		// Handle graceful shutdown
		if err := riverClient.Stop(ctx); err != nil {
			slog.Error("failed to close river client", "error", err)
		}
	}()

	if err := riverClient.Start(ctx); err != nil {
		panic(err)
	}
}
