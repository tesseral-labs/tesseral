package frontendinterceptor

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/openauth-dev/openauth/internal/authn"
	"github.com/openauth-dev/openauth/internal/store"
)

var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")
var ErrInvalidSessionToken = errors.New("invalid session token")

var skipRPCs = []string{
	"/frontend.v1.Frontend/SignInWithEmail",
}

func New(s *store.Store) connect.UnaryInterceptorFunc {
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
				return nil, connect.NewError(connect.CodeUnauthenticated, ErrAuthorizationHeaderRequired)
			}

			secretValue, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			}

			// Attempt to parse the session token
			sessionJWT, err := s.ParseSessionJWT(ctx, secretValue)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidSessionToken)
			}

			// TODO: Add checks to ensure the intermediate session token is valid

			ctx = authn.NewContext(ctx, authn.ContextData{
				Session: sessionJWT,
			})

			return next(ctx, req)
		}
	}
}
