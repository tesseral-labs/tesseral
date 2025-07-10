package corstrusteddomains

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	"github.com/tesseral-labs/tesseral/internal/common/store"
	"github.com/tesseral-labs/tesseral/internal/common/trusteddomains"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/tesseral-labs/tesseral/internal/common/corstrusteddomains")

func Handler(s *store.Store, p *projectid.Sniffer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := tracer.Start(ctx, "common/corstrusteddomains/handler")
		defer span.End()

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
			AllowOriginVaryRequestFunc: func(r *http.Request, origin string) (bool, []string) {
				isTrusted, err := trusteddomains.IsTrustedDomain(trustedOrigins, origin)
				if err != nil {
					http.Error(w, "", http.StatusInternalServerError)
					return false, nil
				}

				if isTrusted {
					return true, []string{origin}
				}

				return false, nil
			},
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
