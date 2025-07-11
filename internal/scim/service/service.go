package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"

	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	"github.com/tesseral-labs/tesseral/internal/scim/authn/authnmiddleware"
	"github.com/tesseral-labs/tesseral/internal/scim/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler(p *projectid.Sniffer) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /api/scim/v1/Users", withErr(s.listUsers))
	mux.Handle("GET /api/scim/v1/Users/{userID}", withErr(s.getUser))
	mux.Handle("POST /api/scim/v1/Users", withErr(s.createUser))
	mux.Handle("PUT /api/scim/v1/Users/{userID}", withErr(s.updateUser))
	mux.Handle("PATCH /api/scim/v1/Users/{userID}", withErr(s.patchUser))
	mux.Handle("DELETE /api/scim/v1/Users/{userID}", withErr(s.deleteUser))

	return logHTTP(authnmiddleware.New(s.Store, p, mux))
}

var (
	filterEmailPat = regexp.MustCompile(`userName eq "(.*)"`)
)

func (s *Service) listUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var filterUserName string
	if r.URL.Query().Has("filter") {
		match := filterEmailPat.FindStringSubmatch(r.URL.Query().Get("filter"))
		if match == nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		// scimvalidator.microsoft.com sends url-encoded values; harmless to
		// "normal" emails to url-parse them
		filterUserName, _ = url.QueryUnescape(match[1])
	}

	var count int
	if r.URL.Query().Has("count") {
		count, _ = strconv.Atoi(r.URL.Query().Get("count"))
	}

	var startIndex int
	if r.URL.Query().Has("startIndex") {
		startIndex, _ = strconv.Atoi(r.URL.Query().Get("startIndex"))
	}

	res, err := s.Store.ListUsers(ctx, &store.ListUsersRequest{
		Count:      count,
		StartIndex: startIndex,
		UserName:   filterUserName,
	})
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

func (s *Service) getUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	user, err := s.Store.GetUser(ctx, r.PathValue("userID"))
	if err != nil {
		var scimError *store.SCIMError
		if errors.As(err, &scimError) {
			w.Header().Set("Content-Type", "application/scim+json")
			w.WriteHeader(scimError.Status)
			if err := json.NewEncoder(w).Encode(scimError); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			return nil
		}

		return fmt.Errorf("store: %w", err)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("read body: %s", err), http.StatusBadRequest)
		return nil
	}

	var reqUser store.User
	if err := json.Unmarshal(body, &reqUser); err != nil {
		http.Error(w, fmt.Sprintf("unmarshal body: %s", err), http.StatusBadRequest)
		return nil
	}

	user, err := s.Store.CreateUser(ctx, reqUser)
	if err != nil {
		var errBadDomain *store.SCIMError
		if errors.As(err, &errBadDomain) {
			w.Header().Set("Content-Type", "application/scim+json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errBadDomain); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			return nil
		}

		return fmt.Errorf("store: %w", err)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

func (s *Service) updateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("read body: %s", err), http.StatusBadRequest)
		return nil
	}

	var reqUser store.User
	if err := json.Unmarshal(body, &reqUser); err != nil {
		http.Error(w, fmt.Sprintf("unmarshal body: %s", err), http.StatusBadRequest)
		return nil
	}

	user, err := s.Store.UpdateUser(ctx, r.PathValue("userID"), reqUser)
	if err != nil {
		var errBadDomain *store.SCIMError
		if errors.As(err, &errBadDomain) {
			w.Header().Set("Content-Type", "application/scim+json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errBadDomain); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			return nil
		}

		return fmt.Errorf("store: %w", err)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

func (s *Service) patchUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("read body: %s", err), http.StatusBadRequest)
		return nil
	}

	var operations store.PatchOperations
	if err := json.Unmarshal(body, &operations); err != nil {
		http.Error(w, fmt.Sprintf("unmarshal body: %s", err), http.StatusBadRequest)
		return nil
	}

	user, err := s.Store.PatchUser(ctx, r.PathValue("userID"), operations)
	if err != nil {
		var scimError *store.SCIMError
		if errors.As(err, &scimError) {
			w.Header().Set("Content-Type", "application/scim+json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(scimError); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			return nil
		}

		return fmt.Errorf("store: %w", err)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if _, err := s.Store.DeleteUser(ctx, r.PathValue("userID")); err != nil {
		var scimError *store.SCIMError
		if errors.As(err, &scimError) {
			w.Header().Set("Content-Type", "application/scim+json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(scimError); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			return nil
		}

		return fmt.Errorf("store: %w", err)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusNoContent)
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
