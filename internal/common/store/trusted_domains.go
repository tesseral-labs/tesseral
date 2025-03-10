package store

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

func (s *Store) GetProjectTrustedOrigins(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	domains, err := s.q.GetProjectTrustedDomains(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	var origins []string
	for _, domain := range domains {
		// Trusted domains are all trusted origins, where the origin scheme is
		// https. As a special case, if an origin has the hostname "localhost",
		// then its HTTP (not S) variant is also a trusted origin.
		origin := url.URL{Scheme: "https", Host: domain}
		if origin.Hostname() == "localhost" {
			originHTTP := url.URL{Scheme: "http", Host: domain}
			origins = append(origins, originHTTP.String())
		}

		origins = append(origins, origin.String())
	}

	return origins, nil
}
