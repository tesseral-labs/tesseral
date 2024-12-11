package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/openauth/openauth/internal/scim/authn/middleware"
	"github.com/openauth/openauth/internal/scim/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /scim/v1/Users", withErr(s.listUsers))
	mux.Handle("GET /scim/v1/Users/{userID}", withErr(s.getUser))
	mux.Handle("POST /scim/v1/Users", withErr(s.createUser))
	mux.Handle("PUT /scim/v1/Users/{userID}", withErr(s.updateUser))

	return middleware.New(s.Store, mux)
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

	user, err := s.Store.CreateUser(ctx, &reqUser)
	if err != nil {
		var errBadDomain *store.BadEmailDomainError
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

	reqUser.ID = r.PathValue("userID")

	user, err := s.Store.UpdateUser(ctx, &reqUser)
	if err != nil {
		var errBadDomain *store.BadEmailDomainError
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

func withErr(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			panic(err)
		}
	})
}
