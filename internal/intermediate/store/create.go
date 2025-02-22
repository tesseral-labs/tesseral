package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

const intermediateSessionDuration = time.Minute * 15

func (s *Store) CreateIntermediateSession(ctx context.Context, req *intermediatev1.CreateIntermediateSessionRequest) (*intermediatev1.CreateIntermediateSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("get project by id: %w", fmt.Errorf("project not found: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	expireTime := time.Now().Add(intermediateSessionDuration)

	secretToken := uuid.New()
	secretTokenSHA256 := sha256.Sum256(secretToken[:])
	if _, err := q.CreateIntermediateSession(ctx, queries.CreateIntermediateSessionParams{
		ID:                uuid.Must(uuid.NewV7()),
		ProjectID:         authn.ProjectID(ctx),
		ExpireTime:        &expireTime,
		SecretTokenSha256: secretTokenSHA256[:],
	}); err != nil {
		return nil, fmt.Errorf("create intermediate session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.CreateIntermediateSessionResponse{
		IntermediateSessionSecretToken: idformat.IntermediateSessionSecretToken.Format(secretToken),
	}, nil
}
