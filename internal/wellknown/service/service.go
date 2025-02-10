package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/openauth/openauth/internal/common/projectid"
	"github.com/openauth/openauth/internal/wellknown/authn/authnmiddleware"
	"github.com/openauth/openauth/internal/wellknown/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler(p *projectid.Sniffer) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /.well-known/webauthn", withErr(s.webauthn))

	return logHTTP(authnmiddleware.New(s.Store, p, mux))
}

func (s *Service) webauthn(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	origins, err := s.Store.GetWebauthnOrigins(ctx)
	if err != nil {
		return fmt.Errorf("get webauthn origins: %w", err)
	}

	body := struct {
		Origins []string `json:"origins"`
	}{
		Origins: origins,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return fmt.Errorf("encode json: %w", err)
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

func logHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(fmt.Errorf("read body: %w", err))
		}

		// log request before ServeHTTP in case that panics
		slog.InfoContext(r.Context(), "http_request", "method", r.Method, "path", r.URL.Path, "request_body", string(body))

		// rewrite the response to be a recorded one, and the request to have the original body
		recorder := httptest.NewRecorder()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		h.ServeHTTP(recorder, r)

		slog.InfoContext(r.Context(), "http_response", "method", r.Method, "path", r.URL.Path, "request_body", string(body), "response_status", recorder.Code, "response_headers", recorder.Header(), "response_body", recorder.Body.String())

		// write out recorded response to w
		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(recorder.Code)
		if _, err := recorder.Body.WriteTo(w); err != nil {
			panic(fmt.Errorf("write body: %w", err))
		}
	})
}
