package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/hibp"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	hibp                                  *hibp.Client
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *kms.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	sessionSigningKeyKmsKeyID             string
	authenticatorAppSecretsKMSKeyID       string
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *kms.Client
	PageEncoder                           pagetoken.Encoder
	SessionSigningKeyKmsKeyID             string
	AuthenticatorAppSecretsKMSKeyID       string
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
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		authenticatorAppSecretsKMSKeyID:       p.AuthenticatorAppSecretsKMSKeyID,
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
