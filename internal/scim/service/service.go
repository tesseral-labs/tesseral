package service

import (
	"net/http"

	"github.com/openauth/openauth/internal/scim/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /scim/v1/{organizationID}/Users", withErr(s.listUsers))

	return mux
}

func (s *Service) listUsers(w http.ResponseWriter, r *http.Request) error {
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
