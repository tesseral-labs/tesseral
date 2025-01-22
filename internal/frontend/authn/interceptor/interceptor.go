package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/common/projectid"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/frontend/authn"
	"github.com/openauth/openauth/internal/frontend/store"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/ujwt"
)

var skipRPCs = []string{
	"/openauth.frontend.v1.FrontendService/GetAccessToken",
}

func New(s *store.Store, p *projectid.Sniffer, authAppsRootDomain string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			projectID, err := p.GetProjectID(req.Header().Get("Host"))
			if err != nil {
				return nil, connect.NewError(connect.CodeNotFound, err)
			}
			requestProjectID := idformat.Project.Format(*projectID)

			// Ensure the projectID is always present
			ctx = authn.NewContext(ctx, authn.ContextData{
				ProjectID: requestProjectID,
			})

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

			aud := fmt.Sprintf("https://%s.tesseral.app", strings.ReplaceAll(requestProjectID, "_", "-"))

			var claims map[string]interface{}
			if err := ujwt.Claims(publicKey, aud, time.Now(), &claims, accessToken); err != nil {
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
