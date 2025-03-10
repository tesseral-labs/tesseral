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

	var origins []string
	for _, domain := range domains {
		origins = append(origins, fmt.Sprintf("https://%s", domain))
	}

	return origins, nil
}
