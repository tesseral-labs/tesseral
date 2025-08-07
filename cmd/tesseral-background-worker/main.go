package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/cyrusaf/ctxlog"
	"github.com/getsentry/sentry-go"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/riverqueue/rivercontrib/otelriver"
	"github.com/ssoready/conf"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/emailworker"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/webhookworker"
	"github.com/tesseral-labs/tesseral/internal/common/sentryintegration"
	"github.com/tesseral-labs/tesseral/internal/dbconn"
	"github.com/tesseral-labs/tesseral/internal/loadenv"
	"github.com/tesseral-labs/tesseral/internal/multislog"
	"github.com/tesseral-labs/tesseral/internal/secretload"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	// do direct os.Getenv here so that we don't depend on secretload, conf, or
	// other things that themselves may fail
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         os.Getenv("TESSERAL_BACKGROUND_WORKER_SENTRY_DSN"),
		Environment: os.Getenv("TESSERAL_BACKGROUND_WORKER_SENTRY_ENVIRONMENT"),
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
		OTELExportTraces      bool          `conf:"otel_export_traces,noredact"`
		OTLPTraceGRPCInsecure bool          `conf:"otlp_trace_grpc_insecure,noredact"`
		ServeAddr             string        `conf:"serve_addr,noredact"`
		DB                    dbconn.Config `conf:"db,noredact"`
		SvixApiKey            string        `conf:"svix_api_key"`
		SESBaseEndpoint       string        `conf:"ses_base_endpoint,noredact"`
		ConsoleProjectID      string        `conf:"console_project_id,noredact"`
		ConsoleDomain         string        `conf:"console_domain,noredact"`
	}{}

	conf.Load(&config)

	if config.OTELExportTraces {
		var exporterOpts []otlptracegrpc.Option
		if config.OTLPTraceGRPCInsecure {
			exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
		}

		exporter, err := otlptracegrpc.New(context.Background(), exporterOpts...)
		if err != nil {
			panic(fmt.Errorf("create otel trace exporter: %w", err))
		}

		tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))

		defer func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				panic(fmt.Errorf("shutdown tracer provider: %w", err))
			}
		}()

		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.TraceContext{})

		logExporter, err := otlploghttp.New(context.Background())
		if err != nil {
			panic(fmt.Errorf("create otel log exporter: %w", err))
		}

		lp := log.NewLoggerProvider(
			log.WithProcessor(
				log.NewBatchProcessor(logExporter),
			),
		)
		defer func() {
			if err := lp.Shutdown(context.Background()); err != nil {
				panic(fmt.Errorf("shutdown logger provider: %w", err))
			}
		}()

		global.SetLoggerProvider(lp)

		slogHandler := ctxlog.NewHandler(
			sentryintegration.NewSlogHandler(
				multislog.Handler{jsonHandler, otelslog.NewHandler("api")},
			),
		)
		slog.SetDefault(slog.New(slogHandler))
	}

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

	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Errorf("load aws config: %w", err))
	}

	sesClient := sesv2.NewFromConfig(awsConfig, func(o *sesv2.Options) {
		if config.SESBaseEndpoint != "" {
			o.BaseEndpoint = &config.SESBaseEndpoint
		}
	})

	backgroundStore := &store.Store{
		DB:                      db,
		Svix:                    svixClient,
		DirectWebhookHTTPClient: &http.Client{},
		SES:                     sesClient,
		ConsoleProjectID:        config.ConsoleProjectID,
		ConsoleDomain:           config.ConsoleDomain,
	}

	riverWorkers := river.NewWorkers()
	river.AddWorker(riverWorkers, &webhookworker.Worker{
		Store: backgroundStore,
	})
	river.AddWorker(riverWorkers, &emailworker.Worker{
		Store: backgroundStore,
	})

	riverClient, err := river.NewClient(riverpgxv5.New(db), &river.Config{
		Logger: slog.Default(),
		Middleware: []rivertype.Middleware{
			otelriver.NewMiddleware(nil),
		},
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

	go func() {
		sigintOrTerm := make(chan os.Signal, 1)
		signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

		slog.Info("wait_sigint_sigterm")
		<-sigintOrTerm

		slog.Info("stop_river")
		if err := riverClient.Stop(context.Background()); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}
			panic(fmt.Errorf("stop river: %w", err))
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/api/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	slog.Info("serve")
	if err := http.ListenAndServe(config.ServeAddr, mux); err != nil {
		panic(err)
	}
}
