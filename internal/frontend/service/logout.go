package service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/backend/authn"
	"github.com/openauth/openauth/internal/cookies"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) Logout(ctx context.Context, req *connect.Request[frontendv1.LogoutRequest]) (*connect.Response[frontendv1.LogoutResponse], error) {
	connectRes := connect.NewResponse(&frontendv1.LogoutResponse{})

	connectRes.Header().Add("Set-Cookie", cookies.ExpiredAccessToken(authn.ProjectID(ctx)))
	connectRes.Header().Add("Set-Cookie", cookies.ExpiredIntermediateAccessToken(authn.ProjectID(ctx)))
	connectRes.Header().Add("Set-Cookie", cookies.ExpiredRefreshToken(authn.ProjectID(ctx)))

	return connectRes, nil
}
