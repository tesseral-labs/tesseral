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
	"github.com/aws/aws-sdk-go-v2/service/kms"
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
	"github.com/openauth/openauth/internal/googleoauth"
	"github.com/openauth/openauth/internal/hexkey"
	intermediateinterceptor "github.com/openauth/openauth/internal/intermediate/authn/interceptor"
	"github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1/intermediatev1connect"
	intermediateservice "github.com/openauth/openauth/internal/intermediate/service"
	intermediatestore "github.com/openauth/openauth/internal/intermediate/store"
	"github.com/openauth/openauth/internal/loadenv"
	"github.com/openauth/openauth/internal/microsoftoauth"
	oauthservice "github.com/openauth/openauth/internal/oauth/service"
	oauthstore "github.com/openauth/openauth/internal/oauth/store"
	"github.com/openauth/openauth/internal/pagetoken"
	"github.com/openauth/openauth/internal/projectid"
	samlservice "github.com/openauth/openauth/internal/saml/service"
	samlstore "github.com/openauth/openauth/internal/saml/store"
	scimservice "github.com/openauth/openauth/internal/scim/service"
	scimstore "github.com/openauth/openauth/internal/scim/store"
	"github.com/openauth/openauth/internal/secretload"
	"github.com/openauth/openauth/internal/slogcorrelation"
	"github.com/openauth/openauth/internal/store/idformat"
	keyManagementService "github.com/openauth/openauth/internal/store/kms"
	"github.com/rs/cors"
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
		DB                                  string `conf:"db"`
		DogfoodProjectID                    string `conf:"dogfood_project_id"`
		IntermediateSessionKMSKeyID         string `conf:"intermediate_session_kms_key_id"`
		KMSEndpoint                         string `conf:"kms_endpoint_resolver_url,noredact"`
		PageEncodingValue                   string `conf:"page-encoding-value"`
		ServeAddr                           string `conf:"serve_addr,noredact"`
		SessionKMSKeyID                     string `conf:"session_kms_key_id"`
		GoogleOAuthClientSecretsKMSKeyID    string `conf:"google_oauth_client_secrets_kms_key_id,noredact"`
		MicrosoftOAuthClientSecretsKMSKeyID string `conf:"microsoft_oauth_client_secrets_kms_key_id,noredact"`
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

	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Errorf("load aws config: %w", err))
	}

	kmsClient := kms.NewFromConfig(awsConfig, func(opts *kms.Options) {
		if config.KMSEndpoint != "" {
			opts.BaseEndpoint = &config.KMSEndpoint
		}
	})

	kms_ := keyManagementService.NewKeyManagementServiceFromConfig(&awsConfig, &config.KMSEndpoint)

	// Register the backend service
	backendStore := backendstore.New(backendstore.NewStoreParams{
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kmsClient,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		GoogleOAuthClientSecretsKMSKeyID:      config.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID:   config.MicrosoftOAuthClientSecretsKMSKeyID,
	})
	backendConnectPath, backendConnectHandler := backendv1connect.NewBackendServiceHandler(
		&backendservice.Service{
			Store: backendStore,
		},
		connect.WithInterceptors(
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
		GoogleOAuthClient:                     &googleoauth.Client{HTTPClient: &http.Client{}},
		MicrosoftOAuthClient:                  &microsoftoauth.Client{HTTPClient: &http.Client{}},
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		GoogleOAuthClientSecretsKMSKeyID:      config.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID:   config.MicrosoftOAuthClientSecretsKMSKeyID,
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

	oauthStore := oauthstore.New(oauthstore.NewStoreParams{
		DB: db,
	})
	oauthService := oauthservice.Service{
		Store: oauthStore,
	}

	samlService := samlservice.Service{
		Store: samlstore.New(samlstore.NewStoreParams{
			DB: db,
		}),
	}

	scimService := scimservice.Service{
		Store: scimstore.New(scimstore.NewStoreParams{
			DB: db,
		}),
	}

	connect := http.NewServeMux()
	connect.Handle(backendConnectPath, backendConnectHandler)
	connect.Handle(frontendConnectPath, frontendConnectHandler)
	connect.Handle(intermediateConnectPath, intermediateConnectHandler)

	// Register health checks
	mux := http.NewServeMux()
	mux.Handle("/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
	}))

	// Register the connect service
	mux.Handle("/internal/connect/", http.StripPrefix("/internal/connect", connect))

	// Register service transcoders
	mux.Handle("/backend/v1/", backendTranscoder)
	mux.Handle("/frontend/v1/", frontendTranscoder)
	mux.Handle("/intermediate/v1/", intermediateTranscoder)

	// Register oauthservice
	mux.Handle("/oauth/", oauthService.Handler())

	// Register samlservice
	mux.Handle("/saml/", samlService.Handler())

	// Register scimservice
	mux.Handle("/scim/", scimService.Handler())

	// These handlers are registered in a FILO order much like
	// a Matryoshka doll

	// Use the projectid.NewHttpHandler to extract the project ID from the request
	serve := projectid.NewHttpHandler(mux)
	// Use the slogcorrelation.NewHandler to add correlation IDs to the request
	serve = slogcorrelation.NewHandler(serve)
	// Add CORS headers
	serve = cors.AllowAll().Handler(serve)

	// Serve the services
	slog.Info("serve")
	if err := http.ListenAndServe(config.ServeAddr, serve); err != nil {
		panic(err)
	}
}
