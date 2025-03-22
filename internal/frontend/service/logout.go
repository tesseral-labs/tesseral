package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) Logout(ctx context.Context, req *connect.Request[frontendv1.LogoutRequest]) (*connect.Response[frontendv1.LogoutResponse], error) {
	res, err := s.Store.Logout(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectRes := connect.NewResponse(res)

	accessTokenCookie, err := s.Cookier.ExpiredAccessToken(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("create expired access token cookie: %w", err)
	}

	intermediateAccessTokenCookie, err := s.Cookier.ExpiredIntermediateAccessToken(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("create expired intermediate access token cookie: %w", err)
	}

	refreshTokenCookie, err := s.Cookier.ExpiredRefreshToken(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("create expired refresh token cookie: %w", err)
	}

	connectRes.Header().Add("Set-Cookie", accessTokenCookie)
	connectRes.Header().Add("Set-Cookie", intermediateAccessTokenCookie)
	connectRes.Header().Add("Set-Cookie", refreshTokenCookie)

	return connectRes, nil
}
