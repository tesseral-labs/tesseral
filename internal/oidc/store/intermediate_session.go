package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) AuthenticateIntermediateSession(ctx context.Context, projectUUID uuid.UUID, secretToken string) (*queries.IntermediateSession, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretTokenUUID, err := idformat.IntermediateSessionSecretToken.Parse(secretToken)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid intermediate session secret token", fmt.Errorf("parse intermediate session secret token: %w", err))
	}

	secretTokenSHA := sha256.Sum256(secretTokenUUID[:])
	qIntermediateSession, err := q.GetIntermediateSessionByTokenSHA256AndProjectID(ctx, queries.GetIntermediateSessionByTokenSHA256AndProjectIDParams{
		ProjectID:         projectUUID,
		SecretTokenSha256: secretTokenSHA[:],
	})
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by token sha256 and project id: %w", err)
	}

	return &qIntermediateSession, nil
}
