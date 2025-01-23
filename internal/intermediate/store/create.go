package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

const intermediateSessionDuration = time.Minute * 15

func (s *Store) CreateIntermediateSession(ctx context.Context, req *intermediatev1.CreateIntermediateSessionRequest) (*intermediatev1.CreateIntermediateSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

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
