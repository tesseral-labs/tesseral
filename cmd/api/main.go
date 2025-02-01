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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/cyrusaf/ctxlog"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	backendinterceptor "github.com/openauth/openauth/internal/backend/authn/interceptor"
	"github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1/backendv1connect"
	backendservice "github.com/openauth/openauth/internal/backend/service"
	backendstore "github.com/openauth/openauth/internal/backend/store"
	"github.com/openauth/openauth/internal/common/accesstoken"
	"github.com/openauth/openauth/internal/common/projectid"
	commonstore "github.com/openauth/openauth/internal/common/store"
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
	samlinterceptor "github.com/openauth/openauth/internal/saml/authn/interceptor"
	samlservice "github.com/openauth/openauth/internal/saml/service"
	samlstore "github.com/openauth/openauth/internal/saml/store"
	scimservice "github.com/openauth/openauth/internal/scim/service"
	scimstore "github.com/openauth/openauth/internal/scim/store"
	"github.com/openauth/openauth/internal/secretload"
	"github.com/openauth/openauth/internal/slogcorrelation"
	"github.com/openauth/openauth/internal/store/idformat"
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
		Host                                string `conf:"host"`
		AuthAppsRootDomain                  string `conf:"auth_apps_root_domain"`
		DB                                  string `conf:"db"`
		DogfoodAuthDomain                   string `conf:"dogfood_auth_domain"`
		DogfoodProjectID                    string `conf:"dogfood_project_id"`
		IntermediateSessionKMSKeyID         string `conf:"intermediate_session_kms_key_id"`
		KMSEndpoint                         string `conf:"kms_endpoint_resolver_url,noredact"`
		PageEncodingValue                   string `conf:"page-encoding-value"`
		S3UserContentBucketName             string `conf:"s3_user_content_bucket_name,noredact"`
		S3Endpoint                          string `conf:"s3_endpoint_resolver_url,noredact"`
		SESEndpoint                         string `conf:"ses_endpoint_resolver_url,noredact"`
		ServeAddr                           string `conf:"serve_addr,noredact"`
		SessionKMSKeyID                     string `conf:"session_kms_key_id"`
		GoogleOAuthClientSecretsKMSKeyID    string `conf:"google_oauth_client_secrets_kms_key_id,noredact"`
		MicrosoftOAuthClientSecretsKMSKeyID string `conf:"microsoft_oauth_client_secrets_kms_key_id,noredact"`
		AuthenticatorAppSecretsKMSKeyID     string `conf:"authenticator_app_secrets_kms_key_id,noredact"`
		UserContentBaseUrl                  string `conf:"user_content_base_url"`
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

	kms_ := kms.NewFromConfig(awsConfig, func(o *kms.Options) {
		if config.KMSEndpoint != "" {
			o.BaseEndpoint = &config.KMSEndpoint
		}
	})

	s3_ := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		if config.S3Endpoint != "" {
			o.BaseEndpoint = &config.S3Endpoint
			o.UsePathStyle = true
		}
	})

	ses_ := sesv2.NewFromConfig(awsConfig, func(o *sesv2.Options) {
		if config.SESEndpoint != "" {
			o.BaseEndpoint = &config.SESEndpoint
		}
	})

	commonStore := commonstore.New(commonstore.NewStoreParams{
		AppAuthRootDomain:         config.AuthAppsRootDomain,
		DB:                        db,
		KMS:                       kms_,
		SessionSigningKeyKMSKeyID: config.SessionKMSKeyID,
	})

	// Register the backend service
	backendStore := backendstore.New(backendstore.NewStoreParams{
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		S3:                                    s3_,
		S3UserContentBucketName:               config.S3UserContentBucketName,
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		GoogleOAuthClientSecretsKMSKeyID:      config.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID:   config.MicrosoftOAuthClientSecretsKMSKeyID,
		UserContentBaseUrl:                    config.UserContentBaseUrl,
	})
	backendConnectPath, backendConnectHandler := backendv1connect.NewBackendServiceHandler(
		&backendservice.Service{
			Store: backendStore,
		},
		connect.WithInterceptors(
			backendinterceptor.New(backendStore, config.Host, config.DogfoodProjectID, config.DogfoodAuthDomain),
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
		&frontendservice.Service{
			Store:             frontendStore,
			AccessTokenIssuer: accesstoken.NewIssuer(commonStore),
		},
		connect.WithInterceptors(
			frontendinterceptor.New(frontendStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), config.AuthAppsRootDomain),
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
		SES:                                   ses_,
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		GoogleOAuthClientSecretsKMSKeyID:      config.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID:   config.MicrosoftOAuthClientSecretsKMSKeyID,
		AuthenticatorAppSecretsKMSKeyID:       config.AuthenticatorAppSecretsKMSKeyID,
		UserContentBaseUrl:                    config.UserContentBaseUrl,
	})
	intermediateConnectPath, intermediateConnectHandler := intermediatev1connect.NewIntermediateServiceHandler(
		&intermediateservice.Service{
			Store:             intermediateStore,
			AccessTokenIssuer: accesstoken.NewIssuer(commonStore),
		},
		connect.WithInterceptors(
			intermediateinterceptor.New(intermediateStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), config.AuthAppsRootDomain),
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

	samlStore := samlstore.New(samlstore.NewStoreParams{
		DB: db,
	})
	samlService := samlservice.Service{
		Store: samlStore,
	}
	samlServiceHandler := samlService.Handler()
	samlServiceHandler = samlinterceptor.New(projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), samlServiceHandler)

	scimStore := scimstore.New(scimstore.NewStoreParams{
		DB: db,
	})
	scimService := scimservice.Service{
		Store: scimStore,
	}
	scimServiceHandler := scimService.Handler(projectid.NewSniffer(config.AuthAppsRootDomain, commonStore))

	connectMux := http.NewServeMux()
	connectMux.Handle(backendConnectPath, backendConnectHandler)
	connectMux.Handle(frontendConnectPath, frontendConnectHandler)
	connectMux.Handle(intermediateConnectPath, intermediateConnectHandler)

	// Register health checks
	mux := http.NewServeMux()
	mux.Handle("/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
	}))

	// Register the connect service
	mux.Handle("/internal/connect/", http.StripPrefix("/internal/connect", connectMux))

	// Register service transcoders
	mux.Handle("/backend/v1/", backendTranscoder)
	mux.Handle("/frontend/v1/", frontendTranscoder)
	mux.Handle("/intermediate/v1/", intermediateTranscoder)

	// Register oauthservice
	mux.Handle("/oauth/", oauthService.Handler())

	// Register samlservice
	mux.Handle("/saml/", samlServiceHandler)

	// Register scimservice
	mux.Handle("/scim/", scimServiceHandler)

	// These handlers are registered in a FILO order much like
	// a Matryoshka doll

	// Use the slogcorrelation.NewHandler to add correlation IDs to the request
	serve := slogcorrelation.NewHandler(mux)
	// Add CORS headers
	serve = cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"*"},
	}).Handler(serve)

	// Serve the services
	slog.Info("serve")
	if err := http.ListenAndServe(config.ServeAddr, serve); err != nil {
		panic(err)
	}
}
