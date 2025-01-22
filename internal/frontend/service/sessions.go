package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) GetAccessToken(ctx context.Context, req *connect.Request[frontendv1.GetAccessTokenRequest]) (*connect.Response[frontendv1.GetAccessTokenResponse], error) {
	refreshToken, _ := cookies.GetCookie(ctx, req, "refreshToken", authn.ProjectID(ctx))
	if refreshToken != "" {
		req.Msg.RefreshToken = refreshToken
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectRes := connect.NewResponse(&frontendv1.GetAccessTokenResponse{
		AccessToken: accessToken,
	})
	connectRes.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "accessToken", accessToken, authn.ProjectID(ctx)))

	return connectRes, nil
}

func (s *Service) Whoami(ctx context.Context, req *connect.Request[frontendv1.WhoAmIRequest]) (*connect.Response[frontendv1.WhoAmIResponse], error) {
	res, err := s.Store.Whoami(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
