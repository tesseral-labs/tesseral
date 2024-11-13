package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth-dev/openauth/internal/pagetoken"
)

type Store struct {
	db									*pgxpool.Pool
	pageEncoder					pagetoken.Encoder
}

type NewStoreParams struct {
	DB								*pgxpool.Pool
	PageEncoder				pagetoken.Encoder
}

func New(p NewStoreParams) *Store {
	return &Store{
		db: p.DB,
		pageEncoder: p.PageEncoder,
	}
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