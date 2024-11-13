package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	"github.com/cyrusaf/ctxlog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth-dev/openauth/internal/backendservice"
	"github.com/openauth-dev/openauth/internal/frontendservice"
	"github.com/openauth-dev/openauth/internal/gen/backend/v1/backendv1connect"
	"github.com/openauth-dev/openauth/internal/gen/frontend/v1/frontendv1connect"
	"github.com/openauth-dev/openauth/internal/hexkey"
	"github.com/openauth-dev/openauth/internal/loadenv"
	"github.com/openauth-dev/openauth/internal/pagetoken"
	"github.com/openauth-dev/openauth/internal/secretload"
	"github.com/openauth-dev/openauth/internal/slogcorrelation"
	"github.com/openauth-dev/openauth/internal/store"
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
		DB 														string `conf:"db"`
		DogfoodProjectID 							string `conf:"dogfood_project_id"`
		PageEncodingValue            	string `conf:"page-encoding-value"`
		ServeAddr 										string `conf:"serve_addr,noredact"`
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

	store_ := store.New(store.NewStoreParams{
		DB: db,
		PageEncoder: pagetoken.Encoder{Secret: pageEncodingValue},
	})

	backendConnectPath, backendConnectHandler := backendv1connect.NewBackendServiceHandler(
		&backendservice.BackendService{
			Store: store_,
		},
		connect.WithInterceptors(),
	)
	backend := vanguard.NewService(backendConnectPath, backendConnectHandler)
	backendTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{backend})
	if err != nil {
		panic(err)
	}

	frontendConnectPath, frontendConnectHandler := frontendv1connect.NewFrontendServiceHandler(
		&frontendservice.FrontendService{
			Store: store_,
		},
		connect.WithInterceptors(),
	)
	frontend := vanguard.NewService(frontendConnectPath, frontendConnectHandler)
	frontendTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{frontend})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
	}))

	mux.Handle("/backend/v1/", backendTranscoder)
	mux.Handle("/frontend/v1/", frontendTranscoder)
	
	slog.Info("serve")
	if err := http.ListenAndServe(config.ServeAddr, slogcorrelation.NewHandler(mux)); err != nil {
		panic(err)
	}
}
