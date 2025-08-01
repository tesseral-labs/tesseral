package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/kms"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
	"github.com/tesseral-labs/tesseral/internal/store/queries"
)

type Store struct {
	db                   *pgxpool.Pool
	consoleProjectID     *uuid.UUID
	sessionSigningKeyKMS *kms.KMS
	pageEncoder          pagetoken.Encoder
	q                    *queries.Queries
}

type NewStoreParams struct {
	DB                   *pgxpool.Pool
	ConsoleProjectID     *uuid.UUID
	SessionSigningKeyKMS *kms.KMS
	PageEncoder          pagetoken.Encoder
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:                   p.DB,
		consoleProjectID:     p.ConsoleProjectID,
		sessionSigningKeyKMS: p.SessionSigningKeyKMS,
		pageEncoder:          p.PageEncoder,
		q:                    queries.New(p.DB),
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
