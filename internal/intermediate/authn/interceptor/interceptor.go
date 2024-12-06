package intermediateinterceptor

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/intermediate/authn"
	"github.com/openauth/openauth/internal/intermediate/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")

var skipRPCs = []string{
	"/openauth.intermediate.v1.IntermediateService/SignInWithEmail",
}

func New(s *store.Store) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			projectIDHeader := req.Header().Get("X-OpenAuth-Project-ID")
			if projectIDHeader == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("project ID header is required"))
			}

			projectID, err := idformat.Project.Parse(projectIDHeader)
			if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}

			slog.Info("New", "projectID", projectID)

			// Check if authentication should be skipped

			for _, rpc := range skipRPCs {
				if req.Spec().Procedure == rpc {
					// Still need to add the project ID to the context even when skipping authentication
					ctx = authn.NewContext(ctx, nil, projectID)
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

			ctx = authn.NewContext(ctx, intermediateSession, projectID)
			return next(ctx, req)
		}
	}
}
