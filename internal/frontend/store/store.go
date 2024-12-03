package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/pagetoken"
	keyManagementService "github.com/openauth/openauth/internal/store/kms"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *keyManagementService.KeyManagementService
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	sessionSigningKeyKmsKeyID             string
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *keyManagementService.KeyManagementService
	PageEncoder                           pagetoken.Encoder
	SessionSigningKeyKmsKeyID             string
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:                                    p.DB,
		dogfoodProjectID:                      p.DogfoodProjectID,
		intermediateSessionSigningKeyKMSKeyID: p.IntermediateSessionSigningKeyKMSKeyID,
		kms:                                   p.KMS,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
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

func derefTimeOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
