package interceptor

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/frontend/authn"
	"github.com/openauth/openauth/internal/frontend/store"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/ujwt"
)

var errInvalidProjectID = fmt.Errorf("invalid project ID")

var skipRPCs = []string{
	"/openauth.frontend.v1.FrontendService/GetAccessToken",
}

func New(s *store.Store, authAppsRootDomain string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// TODO: Move project ID logic to a central location to service all authn interceptors that need it

			// --- Start Project ID sniffing

			projectSubdomainRegexp := regexp.MustCompile(fmt.Sprintf(`([a-zA-Z0-9_-]+)\.%s$`, regexp.QuoteMeta(authAppsRootDomain)))
			host := req.Header().Get("Host")

			var projectID *uuid.UUID
			matches := projectSubdomainRegexp.FindStringSubmatch(host)
			if len(matches) > 1 && strings.HasPrefix(matches[len(matches)-1], "project_") {
				// parse the project ID from the host subdomain
				parsedProjectID, err := idformat.Project.Parse(matches[len(matches)-1])
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

			// Ensure the projectID is always present
			ctx = authn.NewContext(ctx, authn.ContextData{
				ProjectID: requestProjectID,
			})

			// --- Start authentication

			for _, rpc := range skipRPCs {
				if req.Spec().Procedure == rpc {
					return next(ctx, req)
				}
			}

			// get the access token from the cookie to enforce authentication
			accessToken, err := cookies.GetCookie(ctx, req, "accessToken", *projectID)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// determine the session signing key for this access token
			kid, err := ujwt.KeyID(accessToken)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// get the public key for this key; the store will check to make
			// sure it's actually a session signing key for the current project
			publicKey, err := s.GetSessionSigningKeyPublicKey(ctx, kid)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			var claims map[string]interface{}
			if err := ujwt.Claims(publicKey, "TODO", time.Now(), &claims, accessToken); err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = authn.NewContext(ctx, authn.ContextData{
				SessionID:      claims["session"].(map[string]any)["id"].(string),
				UserID:         claims["user"].(map[string]any)["id"].(string),
				OrganizationID: claims["organization"].(map[string]any)["id"].(string),
				ProjectID:      requestProjectID,
			})

			return next(ctx, req)
		}
	}
}
