package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Store) GetProjectTrustedOrigins(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	domains, err := s.q.GetProjectTrustedDomains(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	return domains, nil
}
