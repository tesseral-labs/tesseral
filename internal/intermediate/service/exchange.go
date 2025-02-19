package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
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
	connectRes.Header().Add("Set-Cookie", cookies.NewRefreshToken(authn.ProjectID(ctx), res.RefreshToken))
	connectRes.Header().Add("Set-Cookie", cookies.NewAccessToken(authn.ProjectID(ctx), res.AccessToken))
	return connectRes, nil
}
