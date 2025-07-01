package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	auditlogstore "github.com/tesseral-labs/tesseral/internal/auditlog/store"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
	"github.com/tesseral-labs/tesseral/internal/oidcclient"
)

type Store struct {
	db                        *pgxpool.Pool
	q                         *queries.Queries
	oidcClientSecretsKMSKeyID string
	kms                       *kms.Client
	oidc                      *oidcclient.Client
	auditlogStore             *auditlogstore.Store
}

type NewStoreParams struct {
	DB                        *pgxpool.Pool
	KMS                       *kms.Client
	OIDCClientSecretsKMSKeyID string
	OIDCClient                *oidcclient.Client
	AuditlogStore             *auditlogstore.Store
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:                        p.DB,
		q:                         queries.New(p.DB),
		oidcClientSecretsKMSKeyID: p.OIDCClientSecretsKMSKeyID,
		kms:                       p.KMS,
		oidc:                      p.OIDCClient,
		auditlogStore:             p.AuditlogStore,
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

func refOrNil[T comparable](t T) *T {
	var z T
	if t == z {
		return nil
	}
	return &t
}
