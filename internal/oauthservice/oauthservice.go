package oauthservice

import (
	"fmt"
	"net/http"

	"github.com/openauth-dev/openauth/internal/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /oauth/v1/{projectID}/jwks", withErr(s.jwks))

	return mux
}

func (s *Service) jwks(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("projectID", r.PathValue("projectID"))
	return nil
}

func withErr(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			panic(err)
		}
	})
}
