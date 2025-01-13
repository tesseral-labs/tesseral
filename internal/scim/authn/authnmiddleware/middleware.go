package authnmiddleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/scim/authn"
	"github.com/openauth/openauth/internal/scim/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

func New(s *store.Store, authAppsRootDomain string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// --- Start Project ID sniffing

		projectSubdomainRegexp := regexp.MustCompile(fmt.Sprintf(`([a-zA-Z0-9_-]+)\.%s$`, regexp.QuoteMeta(authAppsRootDomain)))
		host := r.Header.Get("Host")

		var projectID *uuid.UUID
		matches := projectSubdomainRegexp.FindStringSubmatch(host)
		if len(matches) > 1 && strings.HasPrefix(matches[len(matches)-1], "project_") {
			// parse the project ID from the host subdomain
			parsedProjectID, err := idformat.Project.Parse(matches[len(matches)-1])
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			// convert the parsed project ID to a UUID
			projectIDUUID := uuid.UUID(parsedProjectID)
			projectID = &projectIDUUID
		} else {
			// get the project ID by the custom domain
			foundProjectID, err := s.GetProjectIDByDomain(r.Context(), host)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			projectID = foundProjectID
		}

		if projectID == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		requestProjectID := idformat.Project.Format(*projectID)

		ctx = authn.NewContext(ctx, nil, requestProjectID)

		bearerToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if bearerToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		scimAPIKey, err := s.GetSCIMAPIKeyByToken(ctx, bearerToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = authn.NewContext(ctx, &authn.SCIMAPIKey{
			ID:             scimAPIKey.ID,
			OrganizationID: scimAPIKey.OrganizationID,
		}, requestProjectID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
