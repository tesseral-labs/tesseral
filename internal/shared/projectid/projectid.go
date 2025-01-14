package projectid

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/shared/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errProjectIDNotFound = fmt.Errorf("project ID not found")

type ProjectIDSniffer struct {
	authAppsRootDomain string
	store              *store.Store
}

func NewProjectIDSniffer(authAppsRootDomain string, store *store.Store) *ProjectIDSniffer {
	return &ProjectIDSniffer{
		authAppsRootDomain: authAppsRootDomain,
		store:              store,
	}
}

func (p *ProjectIDSniffer) GetProjectIDFromDomain(domain string) (*uuid.UUID, error) {
	ctx := context.Background()

	projectSubdomainRegexp := regexp.MustCompile(fmt.Sprintf(`([a-zA-Z0-9_-]+)\.%s$`, regexp.QuoteMeta(p.authAppsRootDomain)))

	var projectID *uuid.UUID
	matches := projectSubdomainRegexp.FindStringSubmatch(domain)
	if len(matches) > 1 && strings.HasPrefix(matches[len(matches)-1], "project_") {
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
		foundProjectID, err := p.store.GetProjectIDByDomain(ctx, domain)
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
