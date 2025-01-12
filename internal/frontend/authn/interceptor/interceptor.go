package interceptor

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/frontend/authn"
	"github.com/openauth/openauth/internal/frontend/store"
	"github.com/openauth/openauth/internal/ujwt"
)

var skipRPCs = []string{
	"/openauth.frontend.v1.FrontendService/GetAccessToken",
}

func New(s *store.Store) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			for _, rpc := range skipRPCs {
				if req.Spec().Procedure == rpc {
					return next(ctx, req)
				}
			}

			// get the access token from the cookie to enforce authentication
			accessToken, err := cookies.GetCookie(ctx, req, "accessToken")
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// determine the session signing key for this access token
			kid, err := ujwt.KeyID(accessToken)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// get the public key for this key; the store will check to make
			// sure it's actually a session signing key for the current project
			publicKey, err := s.GetSessionSigningKeyPublicKey(ctx, kid)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

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
