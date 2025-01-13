package interceptor

import (
	"errors"
	"net/http"

	"github.com/openauth/openauth/internal/scim/projectid"
	"github.com/openauth/openauth/internal/scim/store"
)

var ErrProjectIDRequired = errors.New("Project ID is required")

func New(store *store.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostHeader := r.Header.Get("Host")
		projectID, err := store.GetProjectIDByDomain(r.Context(), hostHeader)
		if err != nil {
			http.Error(w, ErrProjectIDRequired.Error(), http.StatusBadRequest)
			return
		}

		ctx := projectid.NewContext(r.Context(), *projectID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
