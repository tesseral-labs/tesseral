package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cyrusaf/ctxlog"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/google/uuid"
	"github.com/ssoready/conf"
	stripeclient "github.com/stripe/stripe-go/v82/client"
	svix "github.com/svix/svix-webhooks/go"
	backendinterceptor "github.com/tesseral-labs/tesseral/internal/backend/authn/interceptor"
	"github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1/backendv1connect"
	backendservice "github.com/tesseral-labs/tesseral/internal/backend/service"
	backendstore "github.com/tesseral-labs/tesseral/internal/backend/store"
	"github.com/tesseral-labs/tesseral/internal/cloudflaredoh"
	"github.com/tesseral-labs/tesseral/internal/common/accesstoken"
	"github.com/tesseral-labs/tesseral/internal/common/corstrusteddomains"
	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	commonstore "github.com/tesseral-labs/tesseral/internal/common/store"
	configapiservice "github.com/tesseral-labs/tesseral/internal/configapi/service"
	configapistore "github.com/tesseral-labs/tesseral/internal/configapi/store"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/dbconn"
	frontendinterceptor "github.com/tesseral-labs/tesseral/internal/frontend/authn/interceptor"
	"github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1/frontendv1connect"
	frontendservice "github.com/tesseral-labs/tesseral/internal/frontend/service"
	frontendstore "github.com/tesseral-labs/tesseral/internal/frontend/store"
	"github.com/tesseral-labs/tesseral/internal/googleoauth"
	"github.com/tesseral-labs/tesseral/internal/hexkey"
	"github.com/tesseral-labs/tesseral/internal/httplambda"
	"github.com/tesseral-labs/tesseral/internal/httplog"
	intermediateinterceptor "github.com/tesseral-labs/tesseral/internal/intermediate/authn/interceptor"
	"github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1/intermediatev1connect"
	intermediateservice "github.com/tesseral-labs/tesseral/internal/intermediate/service"
	intermediatestore "github.com/tesseral-labs/tesseral/internal/intermediate/store"
	"github.com/tesseral-labs/tesseral/internal/loadenv"
	"github.com/tesseral-labs/tesseral/internal/microsoftoauth"
	"github.com/tesseral-labs/tesseral/internal/opaqueinternalerror"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
	samlinterceptor "github.com/tesseral-labs/tesseral/internal/saml/authn/interceptor"
	samlservice "github.com/tesseral-labs/tesseral/internal/saml/service"
	samlstore "github.com/tesseral-labs/tesseral/internal/saml/store"
	scimservice "github.com/tesseral-labs/tesseral/internal/scim/service"
	scimstore "github.com/tesseral-labs/tesseral/internal/scim/store"
	"github.com/tesseral-labs/tesseral/internal/secretload"
	"github.com/tesseral-labs/tesseral/internal/slogcorrelation"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	wellknownservice "github.com/tesseral-labs/tesseral/internal/wellknown/service"
	wellknownstore "github.com/tesseral-labs/tesseral/internal/wellknown/store"
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

	slog.SetDefault(slog.New(ctxlog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))))

	if err := secretload.Load(context.Background()); err != nil {
		panic(fmt.Errorf("load secrets: %w", err))
	}

	// Attempts to load environment variables from a .env file
	loadenv.LoadEnv()

	config := struct {
		RunAsLambda                         bool          `conf:"run_as_lambda,noredact"`
		ConsoleDomain                       string        `conf:"console_domain,noredact"`
		AuthAppsRootDomain                  string        `conf:"auth_apps_root_domain,noredact"`
		TesseralDNSVaultCNAMEValue          string        `conf:"tesseral_dns_vault_cname_value,noredact"`
		SESSPFMXRecordValue                 string        `conf:"ses_spf_mx_record_value,noredact"`
		DB                                  dbconn.Config `conf:"db,noredact"`
		CloudflareAPIToken                  string        `conf:"cloudflare_api_token"`
		DogfoodProjectID                    string        `conf:"dogfood_project_id,noredact"`
		IntermediateSessionKMSKeyID         string        `conf:"intermediate_session_kms_key_id,noredact"`
		KMSEndpoint                         string        `conf:"kms_endpoint_resolver_url,noredact"`
		PageEncodingValue                   string        `conf:"page-encoding-value"`
		S3UserContentBucketName             string        `conf:"s3_user_content_bucket_name,noredact"`
		S3Endpoint                          string        `conf:"s3_endpoint_resolver_url,noredact"`
		SESEndpoint                         string        `conf:"ses_endpoint_resolver_url,noredact"`
		ServeAddr                           string        `conf:"serve_addr,noredact"`
		SessionKMSKeyID                     string        `conf:"session_kms_key_id,noredact"`
		GoogleOAuthClientSecretsKMSKeyID    string        `conf:"google_oauth_client_secrets_kms_key_id,noredact"`
		MicrosoftOAuthClientSecretsKMSKeyID string        `conf:"microsoft_oauth_client_secrets_kms_key_id,noredact"`
		AuthenticatorAppSecretsKMSKeyID     string        `conf:"authenticator_app_secrets_kms_key_id,noredact"`
		UserContentBaseUrl                  string        `conf:"user_content_base_url,redact"`
		TesseralDNSCloudflareZoneID         string        `conf:"tesseral_dns_cloudflare_zone_id,noredact"`
		StripeAPIKey                        string        `conf:"stripe_api_key"`
		StripePriceIDGrowthTier             string        `conf:"stripe_price_id_growth_tier,noredact"`
		SvixApiKey                          string        `conf:"svix_api_key,noredact"`
	}{
		PageEncodingValue: "0000000000000000000000000000000000000000000000000000000000000000",
	}

	conf.Load(&config)
	slog.Info("config", "config", conf.Redact(config))

	// TODO: Set up Sentry apps and error handling

	db, err := dbconn.Open(context.Background(), config.DB)
	if err != nil {
		panic(fmt.Errorf("open database: %w", err))
	}
	defer db.Close()

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

	svixClient, err := svix.New(config.SvixApiKey, nil)
	if err != nil {
		panic(fmt.Errorf("create svix client: %w", err))
	}

	stripeClient := stripeclient.New(config.StripeAPIKey, nil)

	commonStore := commonstore.New(commonstore.NewStoreParams{
		AppAuthRootDomain:         config.AuthAppsRootDomain,
		DB:                        db,
		KMS:                       kms_,
		SessionSigningKeyKMSKeyID: config.SessionKMSKeyID,
	})

	cookier := cookies.Cookier{Store: commonStore}

	// Register the backend service
	backendStore := backendstore.New(backendstore.NewStoreParams{
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		ConsoleDomain:                         config.ConsoleDomain,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		SES:                                   ses_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		S3:                                    s3_,
		S3UserContentBucketName:               config.S3UserContentBucketName,
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		GoogleOAuthClientSecretsKMSKeyID:      config.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID:   config.MicrosoftOAuthClientSecretsKMSKeyID,
		UserContentBaseUrl:                    config.UserContentBaseUrl,
		AuthAppsRootDomain:                    config.AuthAppsRootDomain,
		TesseralDNSVaultCNAMEValue:            config.TesseralDNSVaultCNAMEValue,
		SESSPFMXRecordValue:                   config.SESSPFMXRecordValue,
		TesseralDNSCloudflareZoneID:           config.TesseralDNSCloudflareZoneID,
		Cloudflare:                            cloudflare.NewClient(option.WithAPIToken(config.CloudflareAPIToken)),
		CloudflareDOH:                         &cloudflaredoh.Client{HTTPClient: &http.Client{}},
		Stripe:                                stripeClient,
		StripePriceIDGrowthTier:               config.StripePriceIDGrowthTier,
		SvixClient:                            svixClient,
	})
	backendConnectPath, backendConnectHandler := backendv1connect.NewBackendServiceHandler(
		&backendservice.Service{
			Store: backendStore,
		},
		connect.WithInterceptors(
			opaqueinternalerror.NewInterceptor(),
			httplog.NewInterceptor(),
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
		ConsoleDomain:                         config.ConsoleDomain,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		SES:                                   ses_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		AuthenticatorAppSecretsKMSKeyID:       config.AuthenticatorAppSecretsKMSKeyID,
	})
	frontendConnectPath, frontendConnectHandler := frontendv1connect.NewFrontendServiceHandler(
		&frontendservice.Service{
			Store:             frontendStore,
			AccessTokenIssuer: accesstoken.NewIssuer(commonStore),
			Cookier:           &cookier,
		},
		connect.WithInterceptors(
			opaqueinternalerror.NewInterceptor(),
			httplog.NewInterceptor(),
			frontendinterceptor.New(frontendStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), &cookier),
		),
	)
	frontend := vanguard.NewService(frontendConnectPath, frontendConnectHandler)
	frontendTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{frontend})
	if err != nil {
		panic(err)
	}

	// Register the intermediate service
	intermediateStore := intermediatestore.New(intermediatestore.NewStoreParams{
		ConsoleDomain:                         config.ConsoleDomain,
		AuthAppsRootDomain:                    config.AuthAppsRootDomain,
		DB:                                    db,
		DogfoodProjectID:                      &uuidDogfoodProjectID,
		IntermediateSessionSigningKeyKMSKeyID: config.IntermediateSessionKMSKeyID,
		KMS:                                   kms_,
		PageEncoder:                           pagetoken.Encoder{Secret: pageEncodingValue},
		GoogleOAuthClient:                     &googleoauth.Client{HTTPClient: &http.Client{}},
		MicrosoftOAuthClient:                  &microsoftoauth.Client{HTTPClient: &http.Client{}},
		S3:                                    s3_,
		SES:                                   ses_,
		SessionSigningKeyKmsKeyID:             config.SessionKMSKeyID,
		GoogleOAuthClientSecretsKMSKeyID:      config.GoogleOAuthClientSecretsKMSKeyID,
		MicrosoftOAuthClientSecretsKMSKeyID:   config.MicrosoftOAuthClientSecretsKMSKeyID,
		AuthenticatorAppSecretsKMSKeyID:       config.AuthenticatorAppSecretsKMSKeyID,
		UserContentBaseUrl:                    config.UserContentBaseUrl,
		S3UserContentBucketName:               config.S3UserContentBucketName,
		StripeClient:                          stripeClient,
		SvixClient:                            svixClient,
	})
	intermediateConnectPath, intermediateConnectHandler := intermediatev1connect.NewIntermediateServiceHandler(
		&intermediateservice.Service{
			Store:             intermediateStore,
			AccessTokenIssuer: accesstoken.NewIssuer(commonStore),
			Cookier:           &cookier,
		},
		connect.WithInterceptors(
			opaqueinternalerror.NewInterceptor(),
			httplog.NewInterceptor(),
			intermediateinterceptor.New(intermediateStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), &cookier),
		),
	)
	intermediate := vanguard.NewService(intermediateConnectPath, intermediateConnectHandler)
	intermediateTranscoder, err := vanguard.NewTranscoder([]*vanguard.Service{intermediate})
	if err != nil {
		panic(err)
	}

	samlStore := samlstore.New(samlstore.NewStoreParams{
		DB: db,
	})
	samlService := samlservice.Service{
		Store:             samlStore,
		AccessTokenIssuer: accesstoken.NewIssuer(commonStore),
		Cookier:           &cookier,
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

	wellknownStore := wellknownstore.New(wellknownstore.NewStoreParams{
		DB: db,
	})
	wellknownService := wellknownservice.Service{
		Store: wellknownStore,
	}
	wellknownServiceHandler := wellknownService.Handler(projectid.NewSniffer(config.AuthAppsRootDomain, commonStore))

	configapiStore := configapistore.New(configapistore.NewStoreParams{
		DB: db,
	})
	configapiService := configapiservice.Service{
		Store: configapiStore,
	}
	configapiServiceHandler := configapiService.Handler()

	connectMux := http.NewServeMux()
	connectMux.Handle(backendConnectPath, backendConnectHandler)
	connectMux.Handle(frontendConnectPath, frontendConnectHandler)
	connectMux.Handle(intermediateConnectPath, intermediateConnectHandler)

	// Register health checks
	mux := http.NewServeMux()
	mux.Handle("/api/internal/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "health")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	mux.Handle("/api/internal/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("deliberate panic")
	}))

	// Register the connect service
	mux.Handle("/api/internal/connect/", corstrusteddomains.Handler(commonStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), http.StripPrefix("/api/internal/connect", connectMux)))

	// Register service transcoders
	mux.Handle("/api/backend/v1/", http.StripPrefix("/api/backend", backendTranscoder))
	mux.Handle("/api/frontend/v1/", corstrusteddomains.Handler(commonStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), http.StripPrefix("/api", frontendTranscoder)))
	mux.Handle("/api/intermediate/v1/", corstrusteddomains.Handler(commonStore, projectid.NewSniffer(config.AuthAppsRootDomain, commonStore), http.StripPrefix("/api", intermediateTranscoder)))

	// Register samlservice
	mux.Handle("/api/saml/", samlServiceHandler)

	// Register scimservice
	mux.Handle("/api/scim/", scimServiceHandler)

	// Register wellknownservice
	mux.Handle("/.well-known/", wellknownServiceHandler)

	// Register configapiservice
	mux.Handle("/api/config-api/", http.StripPrefix("/api/config-api", configapiServiceHandler))

	// These handlers are registered in a FILO order much like
	// a Matryoshka doll

	// wrap all http requests with sentry
	serve := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle(mux)

	// add correlation IDs to logs
	serve = slogcorrelation.NewHandler(serve)

	slog.Info("serve")
	if config.RunAsLambda {
		lambda.Start(httplambda.Handler(serve))
	} else {
		if err := http.ListenAndServe(config.ServeAddr, serve); err != nil {
			panic(err)
		}
	}
}
