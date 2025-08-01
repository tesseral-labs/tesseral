package store

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	svix "github.com/svix/svix-webhooks/go"
	auditlogstore "github.com/tesseral-labs/tesseral/internal/auditlog/store"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/hibp"
	"github.com/tesseral-labs/tesseral/internal/kms"
	"github.com/tesseral-labs/tesseral/internal/oidcclient"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	db                         *pgxpool.Pool
	consoleProjectID           *uuid.UUID
	consoleDomain              string
	hibp                       *hibp.Client
	oidcClientSecretsKMS       *kms.KMS
	authenticatorAppSecretsKMS *kms.KMS
	ses                        *sesv2.Client
	pageEncoder                pagetoken.Encoder
	q                          *queries.Queries
	svixClient                 *svix.Svix
	auditlogStore              *auditlogstore.Store
	oidc                       *oidcclient.Client
}

type NewStoreParams struct {
	DB                         *pgxpool.Pool
	ConsoleProjectID           *uuid.UUID
	ConsoleDomain              string
	OIDCClientSecretsKMS       *kms.KMS
	AuthenticatorAppSecretsKMS *kms.KMS
	SES                        *sesv2.Client
	PageEncoder                pagetoken.Encoder
	SvixClient                 *svix.Svix
	AuditlogStore              *auditlogstore.Store
	OIDCClient                 *oidcclient.Client
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:               p.DB,
		consoleProjectID: p.ConsoleProjectID,
		consoleDomain:    p.ConsoleDomain,
		hibp: &hibp.Client{
			HTTPClient: http.DefaultClient,
		},
		oidcClientSecretsKMS:       p.OIDCClientSecretsKMS,
		authenticatorAppSecretsKMS: p.AuthenticatorAppSecretsKMS,
		ses:                        p.SES,
		pageEncoder:                p.PageEncoder,
		q:                          queries.New(p.DB),
		svixClient:                 p.SvixClient,
		auditlogStore:              p.AuditlogStore,
		oidc:                       p.OIDCClient,
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

func refOrNil[T comparable](t T) *T {
	var z T
	if t == z {
		return nil
	}
	return &t
}
