package intermediateinterceptor

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	"github.com/openauth/openauth/internal/intermediate/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errInvalidProjectID = fmt.Errorf("invalid project ID")
var ErrAuthorizationHeaderRequired = errors.New("authorization header is required")

var skipRPCs = []string{
	"/openauth.intermediate.v1.IntermediateService/SignInWithEmail",
	"/openauth.intermediate.v1.IntermediateService/GetGoogleOAuthRedirectURL",
}

func New(s *store.Store, authAppsRootDomain string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// --- Start Project ID sniffing

			projectSubdomainRegexp := regexp.MustCompile(fmt.Sprintf(`([a-zA-Z0-9_-]+)\.%s$`, regexp.QuoteMeta(authAppsRootDomain)))
			host := req.Header().Get("Host")

			var projectID *uuid.UUID
			matches := projectSubdomainRegexp.FindStringSubmatch(host)
			if len(matches) > 1 {
				// parse the project ID from the host subdomain
				parsedProjectID, err := idformat.Project.Parse(matches[0])
				if err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}

				// convert the parsed project ID to a UUID
				projectIDUUID := uuid.UUID(parsedProjectID)
				projectID = &projectIDUUID
			} else {
				// get the project ID by the custom domain
				foundProjectID, err := s.GetProjectIDByDomain(ctx, host)
				if err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}

				projectID = foundProjectID
			}

			if projectID == nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, errInvalidProjectID)
			}
			requestProjectID := idformat.Project.Format(*projectID)

			// --- Start authentication

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

			intermediateSession, err := s.GetIntermediateSessionByToken(ctx, secretValue)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = authn.NewContext(ctx, intermediateSession, requestProjectID)
			return next(ctx, req)
		}
	}
}
