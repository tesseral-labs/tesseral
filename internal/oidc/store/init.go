package store

import (
	"context"
	"fmt"
	"net/url"
	"slices"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/oidc/authn"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type OIDCConnectionInitData struct {
	AuthorizationURL string
	State            string
}

func (s *Store) GetOIDCConnectionInitData(ctx context.Context, oidcConnectionID string) (*OIDCConnectionInitData, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	oidcConnectionUUID, err := idformat.OIDCConnection.Parse(oidcConnectionID)
	if err != nil {
		return nil, fmt.Errorf("parse oidc connection id: %w", err)
	}

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        oidcConnectionUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	config, err := s.oidc.GetConfiguration(ctx, qOIDCConnection.ConfigurationUrl)
	if err != nil {
		return nil, fmt.Errorf("get OIDC configuration: %w", err)
	}

	authorizationURL, err := url.Parse(config.AuthorizationEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parse authorization endpoint URL: %w", err)
	}
	query := authorizationURL.Query()
	query.Set("client_id", qOIDCConnection.ClientID)
	query.Set("response_type", "code")
	query.Set("scope", "openid email profile")
	query.Set("redirect_uri", fmt.Sprintf("https://%s/api/oidc/v1/%s/callback", qProject.VaultDomain, oidcConnectionID))

	intermediateSessionID := uuid.New()
	formattedIntermediateSessionID := idformat.OIDCIntermediateSession.Format(intermediateSessionID)
	query.Set("state", formattedIntermediateSessionID)

	// If PKCE is supported, generate code verifier and challenge.
	//
	// Even if it's not required by the OIDC provider, we still generate it to ensure
	// compatibility with clients that may require it.
	var codeVerifier *string
	if len(config.CodeChallengeMethodsSupported) != 0 {
		if !slices.Contains(config.CodeChallengeMethodsSupported, "S256") {
			return nil, fmt.Errorf("OIDC provider does not support S256 code challenge method")
		}
		verifier, codeChallenge, err := s.oidc.GenerateCodeVerifierAndChallenge()
		if err != nil {
			return nil, fmt.Errorf("generate code verifier and challenge: %w", err)
		}
		query.Set("code_challenge", codeChallenge)
		query.Set("code_challenge_method", "S256")
		codeVerifier = &verifier
	}

	authorizationURL.RawQuery = query.Encode()

	if _, err := q.CreateOIDCIntermediateSession(ctx, queries.CreateOIDCIntermediateSessionParams{
		OidcIntermediateSessionID: intermediateSessionID,
		OidcConnectionID:          oidcConnectionUUID,
		CodeVerifier:              codeVerifier,
	}); err != nil {
		return nil, fmt.Errorf("create OIDC session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &OIDCConnectionInitData{
		AuthorizationURL: authorizationURL.String(),
		State:            formattedIntermediateSessionID,
	}, nil
}
