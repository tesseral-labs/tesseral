package backendinterceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/authn"
	"github.com/openauth/openauth/internal/store"
)

var errAuthorizationHeaderRequired = errors.New("authorization header is required")

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
			}

			return next(ctx, req)
		}
	}
}
