package intermediateinterceptor

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/common/projectid"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	"github.com/openauth/openauth/internal/intermediate/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")

var skipRPCs = []string{
	"/openauth.intermediate.v1.IntermediateService/CreateIntermediateSession",
	"/openauth.intermediate.v1.IntermediateService/GetProjectUISettings",
	"/openauth.intermediate.v1.IntermediateService/ListSAMLOrganizations",
}

func New(s *store.Store, p *projectid.Sniffer, authAppsRootDomain string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			projectID, err := p.GetProjectID(req.Header().Get("Host"))
			if err != nil {
				return nil, connect.NewError(connect.CodeNotFound, err)
			}
			requestProjectID := idformat.Project.Format(*projectID)

			// Ensure the projectID is always present on the context
			ctx = authn.NewContext(ctx, nil, requestProjectID)

			// Check if authentication should be skipped
			for _, rpc := range skipRPCs {
				if req.Spec().Procedure == rpc {
					return next(ctx, req)
				}
			}

			// Enforce authentication if not skipping
			secretValue, err := cookies.GetCookie(ctx, req, "intermediateAccessToken", *projectID)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			intermediateSession, err := s.AuthenticateIntermediateSession(ctx, requestProjectID, secretValue)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = authn.NewContext(ctx, intermediateSession, requestProjectID)
			return next(ctx, req)
		}
	}
}
