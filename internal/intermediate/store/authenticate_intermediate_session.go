package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) AuthenticateIntermediateSession(ctx context.Context, projectID, secretToken string) (*intermediatev1.IntermediateSession, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectUUID, err := idformat.Project.Parse(projectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}

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

	emailVerified, err := s.getIntermediateSessionEmailVerified(ctx, q, qIntermediateSession.ID)
	if err != nil {
		return nil, fmt.Errorf("get intermediate session verified: %w", err)
	}

	return parseIntermediateSession(qIntermediateSession, emailVerified), nil
}
