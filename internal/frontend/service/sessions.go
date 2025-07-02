package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) Refresh(ctx context.Context, req *connect.Request[frontendv1.RefreshRequest]) (*connect.Response[frontendv1.RefreshResponse], error) {
	refreshToken, _ := s.Cookier.GetRefreshToken(authn.ProjectID(ctx), req)
	if refreshToken != "" {
		req.Msg.RefreshToken = refreshToken
	}

	if req.Msg.RefreshToken == "" {
		return nil, apierror.NewUnauthenticatedError("no refresh token provided", nil)
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, authn.ProjectID(ctx), req.Msg.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	if err := s.Store.CreateRefreshAuditLogEvent(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("log refresh event: %w", err)
	}

	connectRes := connect.NewResponse(&frontendv1.RefreshResponse{
		AccessToken: accessToken,
	})

	accessTokenCookie, err := s.Cookier.NewAccessToken(ctx, authn.ProjectID(ctx), accessToken)
	if err != nil {
		return nil, fmt.Errorf("create access token cookie: %w", err)
	}

	connectRes.Header().Add("Set-Cookie", accessTokenCookie)

	return connectRes, nil
}
