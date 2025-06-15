package interceptor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	"github.com/tesseral-labs/tesseral/internal/backend/store"
	"github.com/tesseral-labs/tesseral/internal/ujwt"
)

var errAuthorizationHeaderRequired = errors.New("authorization header is required")

func New(s *store.Store, dogfoodProjectID string) connect.UnaryInterceptorFunc {
	cookieName := fmt.Sprintf("tesseral_%s_access_token", dogfoodProjectID)

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			authorizationHeader := req.Header().Get("Authorization")
			if authorizationHeader != "" {
				// look for backend API in Authorization header
				authorization := req.Header().Get("Authorization")
				if authorization == "" {
					return nil, connect.NewError(connect.CodeUnauthenticated, errAuthorizationHeaderRequired)
				}

				secretValue, ok := strings.CutPrefix(authorization, "Bearer ")
				if !ok {
					return nil, connect.NewError(connect.CodeUnauthenticated, nil)
				}

				res, err := s.AuthenticateBackendAPIKey(ctx, secretValue)
				if err != nil {
					if errors.Is(err, store.ErrBadBackendAPIKey) {
						return nil, connect.NewError(connect.CodeUnauthenticated, err)
					}
					return nil, fmt.Errorf("authenticate project api key: %w", err)
				}

				ctx = authn.NewBackendAPIKeyContext(ctx, &authn.BackendAPIKeyContextData{
					BackendAPIKeyID: res.BackendAPIKeyID,
					ProjectID:       res.ProjectID,
				})
			} else {
				// look for access token in cookie
				var accessToken string
				for _, h := range req.Header().Values("Cookie") {
					cookies, err := http.ParseCookie(h)
					if err != nil {
						return nil, fmt.Errorf("parse cookie: %w", err)
					}

					for _, c := range cookies {
						if c.Name != cookieName {
							continue
						}
						accessToken = c.Value
					}
				}

				if accessToken == "" {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no access token found in cookies"))
				}

				sessionCtxData, err := authenticateAccessToken(ctx, s, dogfoodProjectID, accessToken)
				if err != nil {
					return nil, fmt.Errorf("authenticate access token: %w", err)
				}

				ctx = authn.NewDogfoodSessionContext(ctx, *sessionCtxData)
			}

			return next(ctx, req)
		}
	}
}

func authenticateAccessToken(ctx context.Context, s *store.Store, dogfoodProjectID, accessToken string) (*authn.DogfoodSessionContextData, error) {
	// our customers do this logic using our SDK, but we can't use that
	// ourselves here; fetch the public key indicated by accessToken and then
	// authenticate using that public key

	sessionPublicKeys, err := s.GetSessionPublicKeysByProjectID(ctx, dogfoodProjectID)
	if err != nil {
		return nil, fmt.Errorf("get dogfood session public keys: %w", err)
	}

	kid, err := ujwt.KeyID(accessToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	var pub *ecdsa.PublicKey
	for _, sessionPublicKey := range sessionPublicKeys {
		if sessionPublicKey.ID == kid {
			pub = sessionPublicKey.PublicKey
		}
	}

	aud := fmt.Sprintf("https://%s.tesseral.app", strings.ReplaceAll(dogfoodProjectID, "_", "-"))
	var claims map[string]interface{}
	if err := ujwt.Claims(pub, aud, time.Now(), &claims, accessToken); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	userID := claims["user"].(map[string]any)["id"].(string)
	sessionID := claims["session"].(map[string]any)["id"].(string)
	organizationID := claims["organization"].(map[string]any)["id"].(string)

	projectID, err := s.GetProjectIDOrganizationBacks(ctx, organizationID)
	if err != nil {
		panic(fmt.Errorf("get project id organization backs: %w", err))
	}

	return &authn.DogfoodSessionContextData{
		UserID:    userID,
		SessionID: sessionID,
		ProjectID: projectID,
	}, nil
}
