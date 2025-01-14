package authnmiddleware

import (
	"net/http"
	"strings"

	"github.com/openauth/openauth/internal/scim/authn"
	"github.com/openauth/openauth/internal/scim/store"
	"github.com/openauth/openauth/internal/shared/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func New(s *store.Store, p *projectid.Sniffer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Project ID sniffing
		projectID, err := p.GetProjectID(r.Host)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		requestProjectID := idformat.Project.Format(*projectID)
		ctx = authn.NewContext(ctx, nil, requestProjectID)

		// Authentication

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
