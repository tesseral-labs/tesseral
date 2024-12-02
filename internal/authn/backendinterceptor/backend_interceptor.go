package backendinterceptor

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/authn"
	"github.com/openauth/openauth/internal/store"
	"github.com/openauth/openauth/internal/ujwt"
)

var errAuthorizationHeaderRequired = errors.New("authorization header is required")

var skipRPCs = []string{
	"/frontend.v1.Frontend/SignInWithEmail",
}

func New(s *store.Store, dogfoodProjectID string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			for _, rpc := range skipRPCs {
				if req.Spec().Procedure == rpc {
					return next(ctx, req)
				}
			}

			// Get the authorization header
			authorization := req.Header().Get("Authorization")
			if authorization == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errAuthorizationHeaderRequired)
			}

			secretValue, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			}

			if strings.HasPrefix(secretValue, "openauth_secret_key_") {
				projectAPIKey, err := s.AuthenticateProjectAPIKey(ctx, secretValue)
				if err != nil {
					if errors.Is(err, store.ErrBadProjectAPIKey) {
						return nil, connect.NewError(connect.CodeUnauthenticated, err)
					}
					return nil, fmt.Errorf("authenticate project api key: %w", err)
				}

				ctx = authn.NewProjectAPIKeyContext(ctx, projectAPIKey)
			} else {
				// it's an openauth access token for the dogfood project, no
				// remaining valid options
				if err := authenticateAccessToken(ctx, s, dogfoodProjectID, secretValue); err != nil {
					return nil, err
				}
			}

			return next(ctx, req)
		}
	}
}

func authenticateAccessToken(ctx context.Context, s *store.Store, dogfoodProjectID, accessToken string) error {
	// our customers do this logic using our SDK, but we can't use that
	// ourselves here; fetch the public key indicated by accessToken and then
	// authenticate using that public key

	sessionPublicKeys, err := s.GetSessionPublicKeysByProjectID(ctx, dogfoodProjectID)
	if err != nil {
		return fmt.Errorf("get dogfood session public keys: %w", err)
	}

	kid, err := ujwt.KeyID(accessToken)
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	var sessionPublicKeyJWK map[string]any
	for _, sessionPublicKey := range sessionPublicKeys {
		jwk := sessionPublicKey.PublicKeyJwk.AsMap()
		if jwk["kid"] == kid {
			sessionPublicKeyJWK = jwk
		}
	}

	if sessionPublicKeyJWK == nil {
		return connect.NewError(connect.CodeUnauthenticated, ujwt.ErrBadJWT)
	}

	x, err := base64.RawURLEncoding.DecodeString(sessionPublicKeyJWK["x"].(string))
	if err != nil {
		panic(fmt.Errorf("parse jwk: %w", err))
	}

	y, err := base64.RawURLEncoding.DecodeString(sessionPublicKeyJWK["y"].(string))
	if err != nil {
		panic(fmt.Errorf("parse jwk: %w", err))
	}

	pubX := new(big.Int)
	pubY := new(big.Int)
	pubX.SetBytes(x)
	pubY.SetBytes(y)

	pub := ecdsa.PublicKey{Curve: elliptic.P256(), X: pubX, Y: pubY}

	var claims map[string]interface{}
	if err := ujwt.Claims(&pub, "TODO", time.Now(), &claims, accessToken); err != nil {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	fmt.Println(claims)
	return nil
}
