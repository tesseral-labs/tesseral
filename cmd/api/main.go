package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/cyrusaf/ctxlog"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	backendinterceptor "github.com/openauth/openauth/internal/backend/authn/interceptor"
	"github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1/backendv1connect"
	backendservice "github.com/openauth/openauth/internal/backend/service"
	backendstore "github.com/openauth/openauth/internal/backend/store"
	frontendinterceptor "github.com/openauth/openauth/internal/frontend/authn/interceptor"
	"github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1/frontendv1connect"
	frontendservice "github.com/openauth/openauth/internal/frontend/service"
	frontendstore "github.com/openauth/openauth/internal/frontend/store"
	"github.com/openauth/openauth/internal/hexkey"
	intermediateinterceptor "github.com/openauth/openauth/internal/intermediate/authn/interceptor"
	"github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1/intermediatev1connect"
	intermediateservice "github.com/openauth/openauth/internal/intermediate/service"
	intermediatestore "github.com/openauth/openauth/internal/intermediate/store"
	"github.com/openauth/openauth/internal/loadenv"
	"github.com/openauth/openauth/internal/oauthservice"
	"github.com/openauth/openauth/internal/pagetoken"
	"github.com/openauth/openauth/internal/secretload"
	"github.com/openauth/openauth/internal/slogcorrelation"
	"github.com/openauth/openauth/internal/store"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/kms"
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
		DB                          string `conf:"db"`
		DogfoodProjectID            string `conf:"dogfood_project_id"`
		IntermediateSessionKMSKeyID string `conf:"intermediate_session_kms_key_id"`
		KMSEndpoint                 string `conf:"kms_endpoint_resolver_url,noredact"`
		PageEncodingValue           string `conf:"page-encoding-value"`
		ServeAddr                   string `conf:"serve_addr,noredact"`
		SessionKMSKeyID             string `conf:"session_kms_key_id"`
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
		panic(fmt.Errorf("parse dogfood project ID: %w", err))
	}
	uuidDogfoodProjectID := uuid.UUID(dogfoodProjectID[:])

	awsConf, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("load aws config: %w", err))
	}

	kms_ := kms.NewKeyManagementServiceFromConfig(&awsConf, &config.KMSEndpoint)

	// Register the backend service
	backendStore := backendstore.New(backendstore.NewStoreParams{
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
	})
	backendConnectPath, backendConnectHandler := backendv1connect.NewBackendServiceHandler(
		&backendservice.Service{
			Store: backendStore,
		},
		connect.WithInterceptors(
			// We may want to use separate auth interceptors for backend and frontend services
			backendinterceptor.New(backendStore, config.DogfoodProjectID),
		),
	)
	backend := vanguard.NewService(backendConnectPath, backendConnectHandler)
	backendTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{backend})
	if err != nil {
		panic(err)
	}

	// Register the frontend service
	frontendStore := frontendstore.New(frontendstore.NewStoreParams{
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
	})
	frontendConnectPath, frontendConnectHandler := frontendv1connect.NewFrontendServiceHandler(
		&frontendservice.FrontendService{
			Store: frontendStore,
		},
		connect.WithInterceptors(
			// We may want to use separate auth interceptors for backend and frontend services
			frontendinterceptor.New(frontendStore),
		),
	)
	frontend := vanguard.NewService(frontendConnectPath, frontendConnectHandler)
	frontendTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{frontend})
	if err != nil {
		panic(err)
	}

	// Register the intermediate service
	intermediateStore := intermediatestore.New(intermediatestore.NewStoreParams{
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
	})
	intermediateConnectPath, intermediateConnectHandler := intermediatev1connect.NewIntermediateServiceHandler(
		&intermediateservice.Service{
			Store: intermediateStore,
		},
		connect.WithInterceptors(
			intermediateinterceptor.New(intermediateStore),
		),
	)
	intermediate := vanguard.NewService(intermediateConnectPath, intermediateConnectHandler)
	intermediateTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{intermediate})
	if err != nil {
		panic(err)
	}

	oauthService := oauthservice.Service{
		Store: store.New(store.NewStoreParams{
			DB:                                    db,
			DogfoodProjectID:                      &uuidDogfoodProjectID,
			IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
			KMS:                                   kms_,
			PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
			SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		}),
	}

	// Register health checks
	mux := http.NewServeMux()
	mux.Handle("/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
	}))

	// Register service transcoders
	mux.Handle("/backend/v1/", backendTranscoder)
	mux.Handle("/frontend/v1/", frontendTranscoder)
	mux.Handle("/intermediate/v1/", intermediateTranscoder)

	// Register oauthservice
	mux.Handle("/oauth/", oauthService.Handler())

	// Serve the services
	slog.Info("serve")
	if err := http.ListenAndServe(config.ServeAddr, slogcorrelation.NewHandler(mux)); err != nil {
		panic(err)
	}
}
