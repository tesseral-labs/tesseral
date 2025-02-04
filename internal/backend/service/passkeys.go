package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListPasskeys(ctx context.Context, req *connect.Request[backendv1.ListPasskeysRequest]) (*connect.Response[backendv1.ListPasskeysResponse], error) {
	res, err := s.Store.ListPasskeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetPasskey(ctx context.Context, req *connect.Request[backendv1.GetPasskeyRequest]) (*connect.Response[backendv1.GetPasskeyResponse], error) {
	res, err := s.Store.GetPasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) DeletePasskey(ctx context.Context, req *connect.Request[backendv1.DeletePasskeyRequest]) (*connect.Response[backendv1.DeletePasskeyResponse], error) {
	res, err := s.Store.DeletePasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
