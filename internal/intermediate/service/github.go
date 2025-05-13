package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) GetGithubOAuthRedirectURL(ctx context.Context, req *connect.Request[intermediatev1.GetGithubOAuthRedirectURLRequest]) (*connect.Response[intermediatev1.GetGithubOAuthRedirectURLResponse], error) {
	res, err := s.Store.GetGithubOAuthRedirectURL(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) RedeemGithubOAuthCode(ctx context.Context, req *connect.Request[intermediatev1.RedeemGithubOAuthCodeRequest]) (*connect.Response[intermediatev1.RedeemGithubOAuthCodeResponse], error) {
	res, err := s.Store.RedeemGithubOAuthCode(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
