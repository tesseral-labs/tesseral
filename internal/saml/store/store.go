package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	auditlogstore "github.com/tesseral-labs/tesseral/internal/auditlog/store"
	"github.com/tesseral-labs/tesseral/internal/saml/store/queries"
)

type Store struct {
	db            *pgxpool.Pool
	q             *queries.Queries
	auditlogStore *auditlogstore.Store
}

type NewStoreParams struct {
	DB            *pgxpool.Pool
	AuditlogStore *auditlogstore.Store
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:            p.DB,
		q:             queries.New(p.DB),
		auditlogStore: p.AuditlogStore,
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
