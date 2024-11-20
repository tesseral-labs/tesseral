package authinterceptor

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/openauth-dev/openauth/internal/authn"
	"github.com/openauth-dev/openauth/internal/jwt"
	"github.com/openauth-dev/openauth/internal/store"
)

var ErrAuthorizationHeaderRequired = errors.New("Authorization header is required")
var ErrInvalidSessionToken = errors.New("invalid session token")

var skipRPCs = []string{

}

func New(j *jwt.JWT, s *store.Store) connect.UnaryInterceptorFunc {
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

			if strings.HasPrefix(secretValue, "openauth_secret_") {
				// It's an API key
				// TODO: Implement API key authentication
				
			} else {
				// Check whether the session token is an intermediate session token or a session token
				intermediateSessionJWT, intermediateErr := j.ParseIntermediateSessionJWT(ctx, secretValue)
				sessionJWT, sessionErr := j.ParseSessionJWT(ctx, secretValue)
				if intermediateErr != nil && sessionErr != nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidSessionToken)
				}
				
				if intermediateSessionJWT != nil {
					// It's an intermediate session token
					ctx = authn.NewContext(ctx, authn.ContextData{
						IntermediateSession: intermediateSessionJWT,
					})
				}

				if sessionJWT != nil {
					// It's a session token
					ctx = authn.NewContext(ctx, authn.ContextData{
						Session: sessionJWT,
					})
				}

				return next(ctx, req)
			} 

			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
	}	
}