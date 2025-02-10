package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/wellknown/authn"
)

func (s *Store) GetWebauthnOrigins(ctx context.Context) ([]string, error) {
	qProjectTrustedDomains, err := s.q.GetProjectTrustedDomains(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	var origins []string
	for _, qProjectTrustedDomain := range qProjectTrustedDomains {
		origins = append(origins, qProjectTrustedDomain.Domain)
	}
	return origins, nil
}
