package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cyrusaf/ctxlog"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	backendv1 "github.com/openauth/openauth/internal/gen/backend/v1"
	"github.com/openauth/openauth/internal/hexkey"
	"github.com/openauth/openauth/internal/loadenv"
	"github.com/openauth/openauth/internal/pagetoken"
	"github.com/openauth/openauth/internal/secretload"
	"github.com/openauth/openauth/internal/store"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/ssoready/conf"
)

func main() {
	slog.SetDefault(slog.New(ctxlog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))))

	if err := secretload.Load(context.Background()); err != nil {
		panic(fmt.Errorf("load secrets: %w", err))
	}

	// Attempts to load environment variables from a .env file
	loadenv.LoadEnv()

	config := struct {
		DB                string `conf:"db"`
		DogfoodProjectID  string `conf:"dogfood_project_id"`
		PageEncodingValue string `conf:"page-encoding-value"`
		ServeAddr         string `conf:"serve_addr,noredact"`
	}{
		PageEncodingValue: "0000000000000000000000000000000000000000000000000000000000000000",
	}

	conf.Load(&config)
	slog.Info("config", "config", conf.Redact(config))

	// TODO: Set up Sentry apps and error handling

	db, err := pgxpool.New(context.Background(), config.DB)
	if err != nil {
		panic(err)
	}

	pageEncodingValue, err := hexkey.New(config.PageEncodingValue)
	if err != nil {
		panic(fmt.Errorf("parse page encoding secret: %w", err))
	}

	dogfoodProjectID, err := idformat.Project.Parse(config.DogfoodProjectID)
	if err != nil {
		panic(fmt.Errorf("parse dogfood project id: %w", err))
	}
	uuidDogfoodProjectID := uuid.UUID(dogfoodProjectID[:])

	store_ := store.New(store.NewStoreParams{
		DB:               db,
		DogfoodProjectID: &uuidDogfoodProjectID,
		PageEncoder:      pagetoken.Encoder{Secret: pageEncodingValue},
	})

	fmt.Println(store_.CreateProjectAPIKey(context.Background(), &backendv1.CreateProjectAPIKeyRequest{
		ProjectApiKey: &backendv1.ProjectAPIKey{
			ProjectId: config.DogfoodProjectID,
		},
	}))
}
