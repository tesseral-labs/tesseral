package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) GetGoogleOAuthRedirectURL(ctx context.Context, req *connect.Request[intermediatev1.GetGoogleOAuthRedirectURLRequest]) (*connect.Response[intermediatev1.GetGoogleOAuthRedirectURLResponse], error) {
	res, err := s.Store.GetGoogleOAuthRedirectURL(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "intermediateAccessToken", res.IntermediateSessionToken, authn.ProjectID(ctx)))

	return connectResponse, nil
}

func (s *Service) RedeemGoogleOAuthCode(ctx context.Context, req *connect.Request[intermediatev1.RedeemGoogleOAuthCodeRequest]) (*connect.Response[intermediatev1.RedeemGoogleOAuthCodeResponse], error) {
	res, err := s.Store.RedeemGoogleOAuthCode(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
