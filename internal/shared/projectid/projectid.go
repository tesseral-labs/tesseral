package projectid

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/shared/store"
)

type Sniffer struct {
	authAppsRootDomain string
	store              *store.Store
}

func NewSniffer(authAppsRootDomain string, store *store.Store) *Sniffer {
	return &Sniffer{
		authAppsRootDomain: authAppsRootDomain,
		store:              store,
	}
}

func (p *Sniffer) GetProjectID(hostname string) (*uuid.UUID, error) {
	ctx := context.Background()

	// get the project ID by the custom domain
	projectID, err := p.store.GetProjectIDByDomain(ctx, strings.Replace(hostname, "-", "_", 1))
	if err != nil {
		return nil, fmt.Errorf("get project id: %w", err)
	}

	return projectID, nil
}
