package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/wellknown/authn"
)

func (s *Store) GetWebauthnOrigins(ctx context.Context) ([]string, error) {
	qProject, err := s.q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	qProject.AuthDomain
}
