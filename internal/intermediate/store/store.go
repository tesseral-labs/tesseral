package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth/openauth/internal/googleoauth"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/microsoftoauth"
	"github.com/openauth/openauth/internal/pagetoken"
	keyManagementService "github.com/openauth/openauth/internal/store/kms"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *keyManagementService.KeyManagementService
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	sessionSigningKeyKmsKeyID             string
	googleOAuthClientSecretsKMSKeyID      string
	microsoftOAuthClientSecretsKMSKeyID   string
	googleOAuthClient                     *googleoauth.Client
	microsoftOAuthClient                  *microsoftoauth.Client
	uiUrl                                 string
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *keyManagementService.KeyManagementService
	PageEncoder                           pagetoken.Encoder
	SessionSigningKeyKmsKeyID             string
	GoogleOAuthClientSecretsKMSKeyID      string
	MicrosoftOAuthClientSecretsKMSKeyID   string
	GoogleOAuthClient                     *googleoauth.Client
	MicrosoftOAuthClient                  *microsoftoauth.Client
	UIUrl                                 string
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
		googleOAuthClient:                     p.GoogleOAuthClient,
		microsoftOAuthClient:                  p.MicrosoftOAuthClient,
		googleOAuthClientSecretsKMSKeyID:      p.GoogleOAuthClientSecretsKMSKeyID,
		microsoftOAuthClientSecretsKMSKeyID:   p.MicrosoftOAuthClientSecretsKMSKeyID,
		uiUrl:                                 p.UIUrl,
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
