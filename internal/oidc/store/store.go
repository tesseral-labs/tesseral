package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/oidc"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
)

type Store struct {
	db                        *pgxpool.Pool
	q                         *queries.Queries
	oidcClientSecretsKMSKeyID string
	kms                       *kms.Client
	oidc                      *oidc.Client
}

type NewStoreParams struct {
	DB                        *pgxpool.Pool
	KMS                       *kms.Client
	OIDCClientSecretsKMSKeyID string
	OIDCClient                *oidc.Client
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:                        p.DB,
		q:                         queries.New(p.DB),
		oidcClientSecretsKMSKeyID: p.OIDCClientSecretsKMSKeyID,
		kms:                       p.KMS,
		oidc:                      p.OIDCClient,
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
