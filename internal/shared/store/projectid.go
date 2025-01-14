package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Store) GetProjectIDByDomain(ctx context.Context, domain string) (*uuid.UUID, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProjectID, err := q.GetProjectIDByCustomAuthDomain(ctx, &domain)
	if err != nil {
		return nil, fmt.Errorf("get project id by custom auth domain: %w", err)
	}

	return &qProjectID, nil
}
