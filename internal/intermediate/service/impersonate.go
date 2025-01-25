package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) RedeemUserImpersonationToken(ctx context.Context, req *connect.Request[intermediatev1.RedeemUserImpersonationTokenRequest]) (*connect.Response[intermediatev1.RedeemUserImpersonationTokenResponse], error) {
	res, err := s.Store.RedeemUserImpersonationToken(ctx, req.Msg)
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
