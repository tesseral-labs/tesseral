package service

import (
	"fmt"
	"net/http"
	"net/url"

	"connectrpc.com/connect"
)

func (s *Service) googleOAuthCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing code parameter"))
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing state parameter"))
	}

	vaultDomain, err := s.Store.GetVaultDomainByGoogleOAuthState(ctx, state)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("get vault domain: %w", err))
	}

	redirectURL, err := url.Parse(fmt.Sprintf("https://%s/google-oauth-callback", vaultDomain))
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("parse redirect URL: %w", err))
	}

	query := redirectURL.Query()
	query.Set("code", code)
	query.Set("state", state)
	redirectURL.RawQuery = query.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	return nil
}
