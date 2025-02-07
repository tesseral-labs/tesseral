package authnmiddleware

import (
	"net/http"

	"github.com/openauth/openauth/internal/common/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/wellknown/authn"
	"github.com/openauth/openauth/internal/wellknown/store"
)

func New(s *store.Store, p *projectid.Sniffer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// project id sniffing
		projectID, err := p.GetProjectID(r.Host)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		requestProjectID := idformat.Project.Format(*projectID)
		ctx = authn.NewContext(ctx, requestProjectID)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
