package store

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	svix "github.com/svix/svix-webhooks/go"
	common "github.com/tesseral-labs/tesseral/internal/common/store"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/hibp"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	consoleDomain                         string
	hibp                                  *hibp.Client
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *kms.Client
	ses                                   *sesv2.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	sessionSigningKeyKmsKeyID             string
	authenticatorAppSecretsKMSKeyID       string
	svixClient                            *svix.Svix
	common                                *common.Store
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	ConsoleDomain                         string
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *kms.Client
	SES                                   *sesv2.Client
	PageEncoder                           pagetoken.Encoder
	SessionSigningKeyKmsKeyID             string
	AuthenticatorAppSecretsKMSKeyID       string
	SvixClient                            *svix.Svix
	CommonStore                           *common.Store
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:               p.DB,
		dogfoodProjectID: p.DogfoodProjectID,
		consoleDomain:    p.ConsoleDomain,
		hibp: &hibp.Client{
			HTTPClient: http.DefaultClient,
		},
		intermediateSessionSigningKeyKMSKeyID: p.IntermediateSessionSigningKeyKMSKeyID,
		kms:                                   p.KMS,
		ses:                                   p.SES,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		authenticatorAppSecretsKMSKeyID:       p.AuthenticatorAppSecretsKMSKeyID,
		svixClient:                            p.SvixClient,
		common:                                p.CommonStore,
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

func timestampOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}
