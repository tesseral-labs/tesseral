package corstrusteddomains

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	"github.com/tesseral-labs/tesseral/internal/common/store"
)

func Handler(s *store.Store, p *projectid.Sniffer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		projectID, err := p.GetProjectID(r.Header.Get("X-Tesseral-Host"))
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		trustedOrigins, err := s.GetProjectTrustedOrigins(ctx, *projectID)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		cors.New(cors.Options{
			AllowedOrigins:   trustedOrigins,
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
		}).Handler(h).ServeHTTP(w, r)
	})
}
