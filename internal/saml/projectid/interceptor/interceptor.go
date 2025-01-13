package interceptor

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/saml/projectid"
	"github.com/openauth/openauth/internal/saml/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

var ErrProjectIDRequired = errors.New("project ID is required")

func New(store *store.Store, authAppsRootDomain string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Move project ID logic to a central location to service all authn interceptors that need it

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
			foundProjectID, err := store.GetProjectIDByDomain(r.Context(), host)
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

		ctx := projectid.NewContext(r.Context(), *projectID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
