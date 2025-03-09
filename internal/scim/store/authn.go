package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/scim/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type SCIMAPIKey struct {
	ID             string
	OrganizationID string
}

func (s *Store) GetSCIMAPIKeyByToken(ctx context.Context, projectID uuid.UUID, token string) (*SCIMAPIKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	tokenUUID, err := idformat.SCIMAPIKeySecretToken.Parse(token)
	if err != nil {
		return nil, fmt.Errorf("parse scim api key token: %w", err)
	}

	tokenSHA := sha256.Sum256(tokenUUID[:])
	qSCIMAPIKey, err := q.GetSCIMAPIKeyByTokenSHA256(ctx, queries.GetSCIMAPIKeyByTokenSHA256Params{
		ProjectID:         projectID,
		SecretTokenSha256: tokenSHA[:],
	})
	if err != nil {
		return nil, fmt.Errorf("get scim api key by token sha256: %w", err)
	}

	return &SCIMAPIKey{
		ID:             idformat.SCIMAPIKey.Format(qSCIMAPIKey.ID),
		OrganizationID: idformat.Organization.Format(qSCIMAPIKey.OrganizationID),
	}, nil
}
