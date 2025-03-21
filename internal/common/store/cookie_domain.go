package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Store) GetProjectCookieDomain(ctx context.Context, projectID uuid.UUID) (string, error) {
	res, err := s.q.GetProjectCookieDomainByProjectID(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("get project cookie domain: %w", err)
	}

	return res, nil
}
