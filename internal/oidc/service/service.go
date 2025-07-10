package service

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/tesseral-labs/tesseral/internal/common/accesstoken"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/emailaddr"
	"github.com/tesseral-labs/tesseral/internal/oidc/authn"
	"github.com/tesseral-labs/tesseral/internal/oidc/store"
)

type Service struct {
	AccessTokenIssuer *accesstoken.Issuer
	Store             *store.Store
	Cookier           *cookies.Cookier
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /api/oidc/v1/{oidcConnectionID}/init", withErr(s.authorize))
	mux.Handle("GET /api/oidc/v1/{oidcConnectionID}/callback", withErr(s.exchange))

	return mux
}

func (s *Service) authorize(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	oidcConnectionID := r.PathValue("oidcConnectionID")

	oidcConnectionData, err := s.Store.GetOIDCConnectionInitData(ctx, oidcConnectionID)
	if err != nil {
		return fmt.Errorf("get OIDC connection init data: %w", err)
	}

	http.Redirect(w, r, oidcConnectionData.AuthorizationURL, http.StatusFound)
	return nil
}

func (s *Service) exchange(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	oidcConnectionID := r.PathValue("oidcConnectionID")

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		error := r.URL.Query().Get("error")
		if error != "" {
			return fmt.Errorf("OIDC exchange error: %s", error)
		}
		return fmt.Errorf("missing code or state in OIDC callback")
	}

	intermediateSession := authn.IntermediateSession(ctx)
	expectedState := intermediateSession.OidcState
	if expectedState == nil {
		return fmt.Errorf("OIDC session state is not set in the intermediate session")
	}
	if *expectedState != state {
		return fmt.Errorf("OIDC session state mismatch: expected %s, got %s", *expectedState, state)
	}

	oidcSessionData, err := s.Store.ExchangeOIDCCode(ctx, oidcConnectionID, code)
	if err != nil {
		return fmt.Errorf("exchange OIDC code: %w", err)
	}

	email := oidcSessionData.Email
	domain, err := emailaddr.Parse(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if !slices.Contains(oidcSessionData.OrganizationDomains, domain) {
		http.Error(w, "bad domain", http.StatusBadRequest)
		return nil
	}

	http.Redirect(w, r, oidcSessionData.RedirectURL, http.StatusFound)
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
