package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth-dev/openauth/internal/pagetoken"
	keyManagementService "github.com/openauth-dev/openauth/internal/store/kms"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type Store struct {
	db																				*pgxpool.Pool
	dogfoodProjectID													*uuid.UUID
	intermediateSessionSigningKeyKMSKeyID 		string
	kms 																			*keyManagementService.KeyManagementService
	pageEncoder																pagetoken.Encoder
	q																					*queries.Queries
	sessionSigningKeyKmsKeyID 								string
}

type NewStoreParams struct {
	AwsConfig 																*aws.Config
	DB																				*pgxpool.Pool
	DogfoodProjectID													string
	IntermediateSessionSigningKeyKMSKeyID 		string
	PageEncoder																pagetoken.Encoder
	SessionSigningKeyKmsKeyID 								string
}

func New(p NewStoreParams) *Store {
	dogfoodProjectID := uuid.MustParse(p.DogfoodProjectID)

	store := &Store{
		db: 																		p.DB,
		dogfoodProjectID: 											&dogfoodProjectID,
		intermediateSessionSigningKeyKMSKeyID: 	p.IntermediateSessionSigningKeyKMSKeyID,
		pageEncoder: 														p.PageEncoder,
		q:                    									queries.New(p.DB),
		sessionSigningKeyKmsKeyID: 							p.SessionSigningKeyKmsKeyID,
	}

	if p.AwsConfig != nil {
		store.kms = keyManagementService.NewKeyManagementServiceFromConfig(p.AwsConfig)
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