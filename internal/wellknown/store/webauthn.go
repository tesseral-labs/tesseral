package store

import (
	"context"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/wellknown/authn"
)

func (s *Store) GetWebauthnOrigins(ctx context.Context) ([]string, error) {
	qProjectTrustedDomains, err := s.q.GetProjectTrustedDomains(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	var origins []string
	for _, qProjectTrustedDomain := range qProjectTrustedDomains {
		origins = append(origins, fmt.Sprintf("https://%s", qProjectTrustedDomain.Domain))
	}
	return origins, nil
}
