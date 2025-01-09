package intermediateinterceptor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	"github.com/openauth/openauth/internal/intermediate/store"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
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

			projectID := projectid.ProjectID(ctx)
			cookieName := fmt.Sprintf("tesseral:%s:intermediateAccessToken", idformat.Project.Format(projectID))

			// Enforce authentication if not skipping
			secretValue, err := cookies.GetCookie(ctx, req, cookieName)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			slog.InfoContext(ctx, "authn.interceptor", "intermediateAccessToken", secretValue)

			intermediateSession, err := s.GetIntermediateSessionByToken(ctx, secretValue)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = authn.NewContext(ctx, intermediateSession)
			return next(ctx, req)
		}
	}
}
