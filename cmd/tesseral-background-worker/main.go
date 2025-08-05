package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cyrusaf/ctxlog"
	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/ssoready/conf"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/workers"
	"github.com/tesseral-labs/tesseral/internal/common/sentryintegration"
	"github.com/tesseral-labs/tesseral/internal/dbconn"
	"github.com/tesseral-labs/tesseral/internal/loadenv"
	"github.com/tesseral-labs/tesseral/internal/secretload"
)

func main() {
	// do direct os.Getenv here so that we don't depend on secretload, conf, or
	// other things that themselves may fail
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         os.Getenv("API_SENTRY_DSN"),
		Environment: os.Getenv("API_SENTRY_ENVIRONMENT"),
	}); err != nil {
		panic(fmt.Errorf("init sentry: %w", err))
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	slogHandler := ctxlog.NewHandler(
		sentryintegration.NewSlogHandler(jsonHandler),
	)
	slog.SetDefault(slog.New(slogHandler))

	if err := secretload.Load(context.Background()); err != nil {
		panic(fmt.Errorf("load secrets: %w", err))
	}

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

	river.AddWorker(riverWorkers, workers.NewBackgroundWorker(store_))

	riverClient, err := river.NewClient(riverpgxv5.New(db), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {
				MaxWorkers: 100,
			},
		},
		Workers: riverWorkers,
	})
	if err != nil {
		panic(fmt.Errorf("create river client: %w", err))
	}

	slog.Info("start")

	if err := riverClient.Start(context.Background()); err != nil {
		panic(fmt.Errorf("start river: %w", err))
	}

	sigintOrTerm := make(chan os.Signal, 1)
	signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("wait_sigint_sigterm")
	<-sigintOrTerm

	slog.Info("stop")
	if err := riverClient.Stop(context.Background()); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return
		}
		panic(fmt.Errorf("stop river: %w", err))
	}
	slog.Info("stopped")
}
