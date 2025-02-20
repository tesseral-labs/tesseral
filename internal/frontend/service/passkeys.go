package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) DeleteMyPasskey(ctx context.Context, req *connect.Request[frontendv1.DeleteMyPasskeyRequest]) (*connect.Response[frontendv1.DeleteMyPasskeyResponse], error) {
	res, err := s.Store.DeleteMyPasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetPasskeyOptions(ctx context.Context, req *connect.Request[frontendv1.GetPasskeyOptionsRequest]) (*connect.Response[frontendv1.GetPasskeyOptionsResponse], error) {
	res, err := s.Store.GetPasskeyOptions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) ListMyPasskeys(ctx context.Context, req *connect.Request[frontendv1.ListMyPasskeysRequest]) (*connect.Response[frontendv1.ListMyPasskeysResponse], error) {
	res, err := s.Store.ListMyPasskeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RegisterPasskey(ctx context.Context, req *connect.Request[frontendv1.RegisterPasskeyRequest]) (*connect.Response[frontendv1.RegisterPasskeyResponse], error) {
	res, err := s.Store.RegisterPasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
