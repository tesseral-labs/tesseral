package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openauth/openauth/internal/scim/authn/middleware"
	"github.com/openauth/openauth/internal/scim/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /scim/v1/Users", withErr(s.listUsers))

	return middleware.New(s.Store, mux)
}

func (s *Service) listUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	res, err := s.Store.ListUsers(ctx, &store.ListUsersRequest{})
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
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
