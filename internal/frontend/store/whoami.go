package store

import (
	"context"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Store) Whoami(ctx context.Context, req *frontendv1.WhoamiRequest) (*frontendv1.WhoamiResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUser, err := q.GetUserByID(ctx, authn.UserID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &frontendv1.WhoamiResponse{
		User: parseUser(qUser),
	}, nil
}
