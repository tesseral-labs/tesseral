package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) GetMicrosoftOAuthRedirectURL(ctx context.Context, req *connect.Request[intermediatev1.GetMicrosoftOAuthRedirectURLRequest]) (*connect.Response[intermediatev1.GetMicrosoftOAuthRedirectURLResponse], error) {
	res, err := s.Store.GetMicrosoftOAuthRedirectURL(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) RedeemMicrosoftOAuthCode(ctx context.Context, req *connect.Request[intermediatev1.RedeemMicrosoftOAuthCodeRequest]) (*connect.Response[intermediatev1.RedeemMicrosoftOAuthCodeResponse], error) {
	res, err := s.Store.RedeemMicrosoftOAuthCode(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
