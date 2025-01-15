package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/shared/apierror"
)

func (s *Store) GetProjectIDByDomain(ctx context.Context, domain string) (*uuid.UUID, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProjectID, err := q.GetProjectIDByCustomAuthDomain(ctx, &domain)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project id not found", fmt.Errorf("project id not found"))
		}

		return nil, fmt.Errorf("get project id by custom auth domain: %w", err)
	}

	return &qProjectID, nil
}
