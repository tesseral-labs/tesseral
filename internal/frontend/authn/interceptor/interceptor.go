package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/frontend/authn"
	"github.com/openauth/openauth/internal/frontend/store"
	"github.com/openauth/openauth/internal/ujwt"
)

func New(s *store.Store) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			fmt.Println("authn interceptor", req)

			// Get the authorization header
			authorization := req.Header().Get("Authorization")
			if authorization == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
			}

			accessToken, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid authorization header"))
			}

			// determine the session signing key for this access token
			kid, err := ujwt.KeyID(accessToken)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			fmt.Println("before pub key")
			// get the public key for this key; the store will check to make
			// sure it's actually a session signing key for the current project
			publicKey, err := s.GetSessionSigningKeyPublicKey(ctx, kid)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			fmt.Println("after pub key")

			var claims map[string]interface{}
			if err := ujwt.Claims(publicKey, "TODO", time.Now(), &claims, accessToken); err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = authn.NewContext(ctx, authn.ContextData{
				SessionID:      claims["session"].(map[string]any)["id"].(string),
				UserID:         claims["user"].(map[string]any)["id"].(string),
				OrganizationID: claims["organization"].(map[string]any)["id"].(string),
			})

			return next(ctx, req)
		}
	}
}
