package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) GetAuthenticatorAppOptions(ctx context.Context, req *connect.Request[intermediatev1.GetAuthenticatorAppOptionsRequest]) (*connect.Response[intermediatev1.GetAuthenticatorAppOptionsResponse], error) {
	res, err := s.Store.GetAuthenticatorAppOptions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) RegisterAuthenticatorApp(ctx context.Context, req *connect.Request[intermediatev1.RegisterAuthenticatorAppRequest]) (*connect.Response[intermediatev1.RegisterAuthenticatorAppResponse], error) {
	res, err := s.Store.RegisterAuthenticatorApp(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) VerifyAuthenticatorApp(ctx context.Context, req *connect.Request[intermediatev1.VerifyAuthenticatorAppRequest]) (*connect.Response[intermediatev1.VerifyAuthenticatorAppResponse], error) {
	res, err := s.Store.VerifyAuthenticatorApp(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
