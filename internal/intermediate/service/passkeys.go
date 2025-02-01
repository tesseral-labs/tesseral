package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) GetPasskeyOptions(ctx context.Context, req *connect.Request[intermediatev1.GetPasskeyOptionsRequest]) (*connect.Response[intermediatev1.GetPasskeyOptionsResponse], error) {
	res, err := s.Store.GetPasskeyOptions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) RegisterPasskey(ctx context.Context, req *connect.Request[intermediatev1.RegisterPasskeyRequest]) (*connect.Response[intermediatev1.RegisterPasskeyResponse], error) {
	res, err := s.Store.RegisterPasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
