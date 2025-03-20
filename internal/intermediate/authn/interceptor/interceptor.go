package intermediateinterceptor

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/common/projectid"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")

var skipRPCs = []string{
	"/tesseral.intermediate.v1.IntermediateService/CreateIntermediateSession",
	"/tesseral.intermediate.v1.IntermediateService/GetSettings",
	"/tesseral.intermediate.v1.IntermediateService/ListSAMLOrganizations",
	"/tesseral.intermediate.v1.IntermediateService/RedeemUserImpersonationToken",
	"/tesseral.intermediate.v1.IntermediateService/ExchangeSessionForIntermediateSession",
	"/tesseral.intermediate.v1.IntermediateService/ExchangeRelayedSessionTokenForSession",
}

func New(s *store.Store, p *projectid.Sniffer, authAppsRootDomain string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			projectID, err := p.GetProjectID(req.Header().Get("X-Tesseral-Host"))
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
			secretValue, err := cookies.GetIntermediateAccessToken(*projectID, req)
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
