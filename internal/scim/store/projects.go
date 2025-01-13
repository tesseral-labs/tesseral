package store

import (
	"context"

	"github.com/google/uuid"
)

func (s *Store) GetProjectIDByDomain(ctx context.Context, domain string) (*uuid.UUID, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := q.GetProjectIDByCustomDomain(ctx, &domain)
	if err != nil {
		return nil, err
	}

	return &projectID, nil
}
