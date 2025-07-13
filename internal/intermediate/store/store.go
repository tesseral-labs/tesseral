package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	stripeclient "github.com/stripe/stripe-go/v82/client"
	svix "github.com/svix/svix-webhooks/go"
	auditlogstore "github.com/tesseral-labs/tesseral/internal/auditlog/store"
	"github.com/tesseral-labs/tesseral/internal/githuboauth"
	"github.com/tesseral-labs/tesseral/internal/googleoauth"
	"github.com/tesseral-labs/tesseral/internal/hibp"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/microsoftoauth"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
)

type Store struct {
	consoleDomain                         string
	authAppsRootDomain                    string
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	hibp                                  *hibp.Client
	intermediateSessionSigningKeyKMSKeyID string
	s3                                    *s3.Client
	s3PresignClient                       *s3.PresignClient
	kms                                   *kms.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	ses                                   *sesv2.Client
	sessionSigningKeyKmsKeyID             string
	githubOAuthClientSecretsKMSKeyID      string
	googleOAuthClientSecretsKMSKeyID      string
	microsoftOAuthClientSecretsKMSKeyID   string
	authenticatorAppSecretsKMSKeyID       string
	githubOAuthClient                     *githuboauth.Client
	googleOAuthClient                     *googleoauth.Client
	microsoftOAuthClient                  *microsoftoauth.Client
	userContentBaseUrl                    string
	s3UserContentBucketName               string
	stripeClient                          *stripeclient.API
	svixClient                            *svix.Svix
	auditlogStore                         *auditlogstore.Store
	defaultGoogleOAuthClientID            string
	defaultGoogleOAuthClientSecret        string
	defaultGoogleOAuthRedirectURI         string
	defaultMicrosoftOAuthClientID         string
	defaultMicrosoftOAuthClientSecret     string
	defaultMicrosoftOAuthRedirectURI      string
}

type NewStoreParams struct {
	ConsoleDomain                         string
	AuthAppsRootDomain                    string
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	IntermediateSessionSigningKeyKMSKeyID string
	S3                                    *s3.Client
	KMS                                   *kms.Client
	PageEncoder                           pagetoken.Encoder
	SES                                   *sesv2.Client
	SessionSigningKeyKmsKeyID             string
	GithubOAuthClientSecretsKMSKeyID      string
	GoogleOAuthClientSecretsKMSKeyID      string
	MicrosoftOAuthClientSecretsKMSKeyID   string
	AuthenticatorAppSecretsKMSKeyID       string
	GithubOAuthClient                     *githuboauth.Client
	GoogleOAuthClient                     *googleoauth.Client
	MicrosoftOAuthClient                  *microsoftoauth.Client
	UserContentBaseUrl                    string
	S3UserContentBucketName               string
	StripeClient                          *stripeclient.API
	SvixClient                            *svix.Svix
	AuditlogStore                         *auditlogstore.Store
	DefaultGoogleOAuthClientID            string
	DefaultGoogleOAuthClientSecret        string
	DefaultGoogleOAuthRedirectURI         string
	DefaultMicrosoftOAuthClientID         string
	DefaultMicrosoftOAuthClientSecret     string
	DefaultMicrosoftOAuthRedirectURI      string
}

func New(p NewStoreParams) *Store {
	store := &Store{
		consoleDomain:      p.ConsoleDomain,
		authAppsRootDomain: p.AuthAppsRootDomain,
		db:                 p.DB,
		dogfoodProjectID:   p.DogfoodProjectID,
		hibp: &hibp.Client{
			HTTPClient: http.DefaultClient,
		},
		intermediateSessionSigningKeyKMSKeyID: p.IntermediateSessionSigningKeyKMSKeyID,
		s3:                                    p.S3,
		s3PresignClient:                       s3.NewPresignClient(p.S3),
		kms:                                   p.KMS,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		ses:                                   p.SES,
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		githubOAuthClient:                     p.GithubOAuthClient,
		googleOAuthClient:                     p.GoogleOAuthClient,
		microsoftOAuthClient:                  p.MicrosoftOAuthClient,
		githubOAuthClientSecretsKMSKeyID:      p.GithubOAuthClientSecretsKMSKeyID,
		googleOAuthClientSecretsKMSKeyID:      p.GoogleOAuthClientSecretsKMSKeyID,
		microsoftOAuthClientSecretsKMSKeyID:   p.MicrosoftOAuthClientSecretsKMSKeyID,
		authenticatorAppSecretsKMSKeyID:       p.AuthenticatorAppSecretsKMSKeyID,
		userContentBaseUrl:                    p.UserContentBaseUrl,
		s3UserContentBucketName:               p.S3UserContentBucketName,
		stripeClient:                          p.StripeClient,
		svixClient:                            p.SvixClient,
		auditlogStore:                         p.AuditlogStore,
		defaultGoogleOAuthClientID:            p.DefaultGoogleOAuthClientID,
		defaultGoogleOAuthClientSecret:        p.DefaultGoogleOAuthClientSecret,
		defaultGoogleOAuthRedirectURI:         p.DefaultGoogleOAuthRedirectURI,
		defaultMicrosoftOAuthClientID:         p.DefaultMicrosoftOAuthClientID,
		defaultMicrosoftOAuthClientSecret:     p.DefaultMicrosoftOAuthClientSecret,
		defaultMicrosoftOAuthRedirectURI:      p.DefaultMicrosoftOAuthRedirectURI,
	}

	return store
}

func (s *Store) tx(ctx context.Context) (tx pgx.Tx, q *queries.Queries, commit func() error, rollback func() error, err error) {
	tx, err = s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("begin tx: %w", err)
	}

	commit = func() error { return tx.Commit(ctx) }
	rollback = func() error { return tx.Rollback(ctx) }
	return tx, queries.New(tx), commit, rollback, nil
}

func derefOrEmpty[T any](t *T) T {
	var z T
	if t == nil {
		return z
	}
	return *t
}

func refOrNil[T comparable](t T) *T {
	var z T
	if t == z {
		return nil
	}
	return &t
}
