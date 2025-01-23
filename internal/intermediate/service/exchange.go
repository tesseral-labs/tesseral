package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) ExchangeIntermediateSessionForSession(ctx context.Context, req *connect.Request[intermediatev1.ExchangeIntermediateSessionForSessionRequest]) (*connect.Response[intermediatev1.ExchangeIntermediateSessionForSessionResponse], error) {
	res, err := s.Store.ExchangeIntermediateSessionForSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, res.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	res.AccessToken = accessToken

	connectRes := connect.NewResponse(res)
	connectRes.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "refreshToken", res.RefreshToken, authn.ProjectID(ctx)))
	connectRes.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "accessToken", res.AccessToken, authn.ProjectID(ctx)))
	return connectRes, nil
}

func (s *Service) ExchangeIntermediateSessionForNewOrganizationSession(ctx context.Context, req *connect.Request[intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionRequest]) (*connect.Response[intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse], error) {
	res, err := s.Store.ExchangeIntermediateSessionForNewOrganizationSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, res.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	res.AccessToken = accessToken

	connectRes := connect.NewResponse(res)
	connectRes.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "refreshToken", res.RefreshToken, authn.ProjectID(ctx)))
	connectRes.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "accessToken", res.AccessToken, authn.ProjectID(ctx)))
	return connectRes, nil
}
