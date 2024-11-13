package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/cyrusaf/ctxlog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth-dev/openauth/internal/hexkey"
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

	config := struct {
		DB 														string `conf:"db"`
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

	awsSDKConfig, err := awsconfig.LoadDefaultConfig(context.Background())
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

	mux := http.NewServeMux()
	mux.Handle("/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
	}))
	
	slog.Info("serve")
	if err := http.ListenAndServe("", slogcorrelation.NewHandler(mux)); err != nil {
		panic(err)
	}
}