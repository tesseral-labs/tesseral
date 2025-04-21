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
	kms                                   *kms.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	ses                                   *sesv2.Client
	sessionSigningKeyKmsKeyID             string
	googleOAuthClientSecretsKMSKeyID      string
	microsoftOAuthClientSecretsKMSKeyID   string
	authenticatorAppSecretsKMSKeyID       string
	googleOAuthClient                     *googleoauth.Client
	microsoftOAuthClient                  *microsoftoauth.Client
	userContentBaseUrl                    string
	s3UserContentBucketName               string
	stripeClient                          *stripeclient.API
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
	GoogleOAuthClientSecretsKMSKeyID      string
	MicrosoftOAuthClientSecretsKMSKeyID   string
	AuthenticatorAppSecretsKMSKeyID       string
	GoogleOAuthClient                     *googleoauth.Client
	MicrosoftOAuthClient                  *microsoftoauth.Client
	UserContentBaseUrl                    string
	S3UserContentBucketName               string
	StripeClient                          *stripeclient.API
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
		kms:                                   p.KMS,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		ses:                                   p.SES,
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		googleOAuthClient:                     p.GoogleOAuthClient,
		microsoftOAuthClient:                  p.MicrosoftOAuthClient,
		googleOAuthClientSecretsKMSKeyID:      p.GoogleOAuthClientSecretsKMSKeyID,
		microsoftOAuthClientSecretsKMSKeyID:   p.MicrosoftOAuthClientSecretsKMSKeyID,
		authenticatorAppSecretsKMSKeyID:       p.AuthenticatorAppSecretsKMSKeyID,
		userContentBaseUrl:                    p.UserContentBaseUrl,
		s3UserContentBucketName:               p.S3UserContentBucketName,
		stripeClient:                          p.StripeClient,
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
