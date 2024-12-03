package intermediateinterceptor

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/store"
)

var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")
var ErrInvalidSessionToken = errors.New("invalid session token")

var skipRPCs = []string{
	"/intermediate.v1.IntermediateService/SignInWithEmail",
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
			// authorization := req.Header().Get("Authorization")
			// if authorization == "" {
			// 	return nil, connect.NewError(connect.CodeUnauthenticated, ErrAuthorizationHeaderRequired)
			// }

			// secretValue, ok := strings.CutPrefix(authorization, "Bearer ")
			// if !ok {
			// 	return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			// }

			// intermediateSessionKID, err := ujwt.KeyID(secretValue)
			// if err != nil {
			// 	return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidSessionToken)
			// }

			// // Get the intermediate session signing key
			// intermediateSessionSigningKey, err := s.GetIntermediateSessionSigningKeyByID(ctx, intermediateSessionKID)
			// if err != nil {
			// 	return nil, connect.NewError(connect.CodeUnauthenticated, err)
			// }

			// intermediateSessionClaims := &intermediatev1.IntermediateSessionClaims{}

			// // Attempt to parse the intermediate session token claims
			// err = ujwt.Claims(intermediateSessionSigningKey.PublicKey, "aud1", time.Unix(2, 0), intermediateSessionClaims, secretValue)
			// if err != nil {
			// 	return nil, connect.NewError(connect.CodeUnauthenticated, err)
			// }

			// // TODO: Add checks to ensure the intermediate session token is valid

			// ctx = authn.NewContext(ctx, authn.ContextData{})

			return next(ctx, req)
		}
	}
}
