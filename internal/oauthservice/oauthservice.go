package oauthservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openauth/openauth/internal/store"
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
	ctx := r.Context()
	projectID := r.PathValue("projectID")

	sessionPublicKeys, err := s.Store.GetSessionPublicKeysByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get session public key: %w", err)
	}

	var jwksKeys []any
	for _, key := range sessionPublicKeys {
		jwksKeys = append(jwksKeys, key.PublicKeyJwk)
	}
	res := map[string]any{"keys": jwksKeys}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
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
