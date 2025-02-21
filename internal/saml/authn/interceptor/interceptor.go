package interceptor

import (
	"net/http"

	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	"github.com/tesseral-labs/tesseral/internal/saml/authn"
)

func New(p *projectid.Sniffer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID, err := p.GetProjectID(r.Header.Get("X-Tesseral-Host"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := authn.NewContext(r.Context(), *projectID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
