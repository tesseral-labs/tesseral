package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) GetMicrosoftOAuthRedirectURL(ctx context.Context, req *connect.Request[intermediatev1.GetMicrosoftOAuthRedirectURLRequest]) (*connect.Response[intermediatev1.GetMicrosoftOAuthRedirectURLResponse], error) {
	res, err := s.Store.GetMicrosoftOAuthRedirectURL(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "intermediateAccessToken", res.IntermediateSessionToken))

	return connectResponse, nil
}

func (s *Service) RedeemMicrosoftOAuthCode(ctx context.Context, req *connect.Request[intermediatev1.RedeemMicrosoftOAuthCodeRequest]) (*connect.Response[intermediatev1.RedeemMicrosoftOAuthCodeResponse], error) {
	res, err := s.Store.RedeemMicrosoftOAuthCode(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
