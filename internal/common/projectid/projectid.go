package projectid

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/common/store"
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
	projectID, err := p.store.GetProjectIDByDomain(ctx, hostname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project id not found", err)
		}
		return nil, fmt.Errorf("get project id by domain: %w", err)
	}

	return projectID, nil
}
