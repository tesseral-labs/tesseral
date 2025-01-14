package projectid

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/shared/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errProjectIDNotFound = fmt.Errorf("project ID not found")

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

	projectSubdomainRegexp := regexp.MustCompile(fmt.Sprintf(`([a-zA-Z0-9_-]+)\.%s$`, regexp.QuoteMeta(p.authAppsRootDomain)))

	var projectID *uuid.UUID
	matches := projectSubdomainRegexp.FindStringSubmatch(hostname)
	if len(matches) > 1 {
		// parse the project ID from the host subdomain
		parsedProjectID, err := idformat.Project.Parse(matches[len(matches)-1])
		if err != nil {
			return nil, err
		}

		// convert the parsed project ID to a UUID
		projectIDUUID := uuid.UUID(parsedProjectID)
		projectID = &projectIDUUID
	} else {
		// get the project ID by the custom domain
		foundProjectID, err := p.store.GetProjectIDByDomain(ctx, hostname)
		if err != nil {
			return nil, err
		}

		projectID = foundProjectID
	}

	if projectID == nil {
		return nil, errProjectIDNotFound
	}

	return projectID, nil
}
