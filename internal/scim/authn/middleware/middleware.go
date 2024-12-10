package middleware

import (
	"net/http"
	"strings"

	"github.com/openauth/openauth/internal/scim/authn"
	"github.com/openauth/openauth/internal/scim/store"
)

func New(s *store.Store, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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
		})
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
