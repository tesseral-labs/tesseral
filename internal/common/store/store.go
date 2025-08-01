package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"github.com/tesseral-labs/tesseral/internal/kms"
)

type Store struct {
	appAuthRootDomain     string
	db                    *pgxpool.Pool
	sessionSigningKeysKMS *kms.KMS
	q                     *queries.Queries
}

type NewStoreParams struct {
	AppAuthRootDomain     string
	DB                    *pgxpool.Pool
	SessionSigningKeysKMS *kms.KMS
}

func New(p NewStoreParams) *Store {
	store := &Store{
		appAuthRootDomain:     p.AppAuthRootDomain,
		db:                    p.DB,
		sessionSigningKeysKMS: p.SessionSigningKeysKMS,
		q:                     queries.New(p.DB),
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
