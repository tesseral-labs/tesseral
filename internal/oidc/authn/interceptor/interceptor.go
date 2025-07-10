package interceptor

import (
	"net/http"

	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/oidc/authn"
	"github.com/tesseral-labs/tesseral/internal/oidc/store"
)

func New(s *store.Store, p *projectid.Sniffer, cookier *cookies.Cookier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID, err := p.GetProjectID(r.Header.Get("X-Tesseral-Host"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Ensure the projectID is always present on the context
		ctx := authn.NewContext(r.Context(), nil, *projectID)

		intermediateSessionToken, _ := cookier.GetIntermediateAccessTokenHTTP(*projectID, r)
		if intermediateSessionToken == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		intermediateSession, err := s.AuthenticateIntermediateSession(ctx, *projectID, intermediateSessionToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx = authn.NewContext(r.Context(), intermediateSession, *projectID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
