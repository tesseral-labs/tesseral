package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store/queries"
)

type Store struct {
	DB                      *pgxpool.Pool
	Svix                    *svix.Svix
	DirectWebhookHTTPClient *http.Client
}

func (s *Store) q() *queries.Queries {
	return queries.New(s.DB)
}

func (s *Store) tx(ctx context.Context) (tx pgx.Tx, q *queries.Queries, commit func() error, rollback func() error, err error) {
	tx, err = s.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("begin tx: %w", err)
	}

	commit = func() error { return tx.Commit(ctx) }
	rollback = func() error { return tx.Rollback(ctx) }
	return tx, queries.New(tx), commit, rollback, nil
}
