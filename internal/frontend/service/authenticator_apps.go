package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) GetAuthenticatorAppOptions(ctx context.Context, req *connect.Request[frontendv1.GetAuthenticatorAppOptionsRequest]) (*connect.Response[frontendv1.GetAuthenticatorAppOptionsResponse], error) {
	res, err := s.Store.GetAuthenticatorAppOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RegisterAuthenticatorApp(ctx context.Context, req *connect.Request[frontendv1.RegisterAuthenticatorAppRequest]) (*connect.Response[frontendv1.RegisterAuthenticatorAppResponse], error) {
	res, err := s.Store.RegisterAuthenticatorApp(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
