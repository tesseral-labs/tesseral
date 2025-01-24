package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth/openauth/internal/backend/authn"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/pagetoken"
	"github.com/openauth/openauth/internal/store/idformat"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *kms.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	s3PresignClient                       *s3.PresignClient
	s3UserContentBucketName               string
	sessionSigningKeyKmsKeyID             string
	googleOAuthClientSecretsKMSKeyID      string
	microsoftOAuthClientSecretsKMSKeyID   string
	userContentBaseUrl                    string
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *kms.Client
	PageEncoder                           pagetoken.Encoder
	S3                                    *s3.Client
	S3UserContentBucketName               string
	SessionSigningKeyKmsKeyID             string
	GoogleOAuthClientSecretsKMSKeyID      string
	MicrosoftOAuthClientSecretsKMSKeyID   string
	UserContentBaseUrl                    string
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:                                    p.DB,
		dogfoodProjectID:                      p.DogfoodProjectID,
		intermediateSessionSigningKeyKMSKeyID: p.IntermediateSessionSigningKeyKMSKeyID,
		kms:                                   p.KMS,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		s3PresignClient:                       s3.NewPresignClient(p.S3),
		s3UserContentBucketName:               p.S3UserContentBucketName,
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		googleOAuthClientSecretsKMSKeyID:      p.GoogleOAuthClientSecretsKMSKeyID,
		microsoftOAuthClientSecretsKMSKeyID:   p.MicrosoftOAuthClientSecretsKMSKeyID,
		userContentBaseUrl:                    p.UserContentBaseUrl,
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

// validateIsDogfoodSession returns an error if the caller isn't a dogfood
// session.
//
// The intention of this method is to allow endpoints to prevent themselves from
// being called by project API keys.
func validateIsDogfoodSession(ctx context.Context) error {
	data := authn.GetContextData(ctx)
	if data.DogfoodSession == nil {
		return apierror.NewUnauthenticatedError("this endpoint cannot be invoked by project API keys", fmt.Errorf("non-dogfood session request"))
	}
	return nil
}

// validateIsDogfoodSessionOfProjectOwner returns an error if the caller isn't a
// dogfood session (i.e. a session in the app) or of an owner (i.e. an owner of
// their respective project).
//
// The intention of this method is to allow endpoints to prevent themselves from
// being called by project API keys or by non-owners.
//
// In production, this method validates that you're logged into app.tesseral.com
// and that you're an owner of your Tesseral project. No other types of callers
// are allowed.
func (s *Store) validateIsDogfoodSessionOfProjectOwner(ctx context.Context) error {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return fmt.Errorf("validate is dogfood session: %w", err)
	}

	data := authn.GetContextData(ctx)
	userID, err := idformat.User.Parse(data.DogfoodSession.UserID)
	if err != nil {
		panic(fmt.Errorf("parse user id: %w", err))
	}

	qUser, err := s.q.GetUser(ctx, queries.GetUserParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        userID,
	})
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	if !qUser.IsOwner {
		return apierror.NewPermissionDeniedError("user must be an owner of the organization", fmt.Errorf("user is not an owner"))
	}
	return nil
}
