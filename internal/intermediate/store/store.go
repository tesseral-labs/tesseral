package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth/openauth/internal/googleoauth"
	"github.com/openauth/openauth/internal/hibp"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/microsoftoauth"
	"github.com/openauth/openauth/internal/pagetoken"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	hibp                                  *hibp.Client
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *kms.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	ses                                   *sesv2.Client
	sessionSigningKeyKmsKeyID             string
	googleOAuthClientSecretsKMSKeyID      string
	microsoftOAuthClientSecretsKMSKeyID   string
	googleOAuthClient                     *googleoauth.Client
	microsoftOAuthClient                  *microsoftoauth.Client
	userContentBaseUrl                    string
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *kms.Client
	PageEncoder                           pagetoken.Encoder
	SES                                   *sesv2.Client
	SessionSigningKeyKmsKeyID             string
	GoogleOAuthClientSecretsKMSKeyID      string
	MicrosoftOAuthClientSecretsKMSKeyID   string
	GoogleOAuthClient                     *googleoauth.Client
	MicrosoftOAuthClient                  *microsoftoauth.Client
	UserContentBaseUrl                    string
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:               p.DB,
		dogfoodProjectID: p.DogfoodProjectID,
		hibp: &hibp.Client{
			HTTPClient: http.DefaultClient,
		},
		intermediateSessionSigningKeyKMSKeyID: p.IntermediateSessionSigningKeyKMSKeyID,
		kms:                                   p.KMS,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		ses:                                   p.SES,
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		googleOAuthClient:                     p.GoogleOAuthClient,
		microsoftOAuthClient:                  p.MicrosoftOAuthClient,
		googleOAuthClientSecretsKMSKeyID:      p.GoogleOAuthClientSecretsKMSKeyID,
		microsoftOAuthClientSecretsKMSKeyID:   p.MicrosoftOAuthClientSecretsKMSKeyID,
		userContentBaseUrl:                    p.UserContentBaseUrl,
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
