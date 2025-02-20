package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) Refresh(ctx context.Context, req *connect.Request[frontendv1.RefreshRequest]) (*connect.Response[frontendv1.RefreshResponse], error) {
	refreshToken, _ := cookies.GetRefreshToken(authn.ProjectID(ctx), req)
	if refreshToken != "" {
		req.Msg.RefreshToken = refreshToken
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectRes := connect.NewResponse(&frontendv1.RefreshResponse{
		AccessToken: accessToken,
	})
	connectRes.Header().Add("Set-Cookie", cookies.NewAccessToken(authn.ProjectID(ctx), accessToken))

	return connectRes, nil
}
