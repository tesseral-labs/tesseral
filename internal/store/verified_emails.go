package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type VerifiedEmail struct {
	ID string
}

func (s *Store) CreateVerifiedEmail(ctx context.Context) (*VerifiedEmail, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	_, err = q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
		ID: uuid.New(),
	})
	if err != nil {
		return nil, err
	}

	commit()

	return &VerifiedEmail{
		ID: verifiedEmailID,
	}, nil

}
