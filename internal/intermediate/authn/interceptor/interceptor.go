package intermediateinterceptor

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/intermediate/authn"
	"github.com/openauth/openauth/internal/intermediate/store"
)

var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")

var skipRPCs = []string{
	"/openauth.intermediate.v1.IntermediateService/SignInWithEmail",
	"/openauth.intermediate.v1.IntermediateService/GetGoogleOAuthRedirectURL",
}

func New(s *store.Store) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Check if authentication should be skipped
			for _, rpc := range skipRPCs {
				if req.Spec().Procedure == rpc {
					return next(ctx, req)
				}
			}

			// Enforce authentication if not skipping
			authorization := req.Header().Get("Authorization")
			if authorization == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, ErrAuthorizationHeaderRequired)
			}

			secretValue, ok := strings.CutPrefix(authorization, "Bearer ")
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			}

			intermediateSession, err := s.GetIntermediateSessionByToken(ctx, secretValue)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = authn.NewContext(ctx, intermediateSession)
			return next(ctx, req)
		}
	}
}
