package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) ListSCIMAPIKeys(ctx context.Context, req *connect.Request[frontendv1.ListSCIMAPIKeysRequest]) (*connect.Response[frontendv1.ListSCIMAPIKeysResponse], error) {
	res, err := s.Store.ListSCIMAPIKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetSCIMAPIKey(ctx context.Context, req *connect.Request[frontendv1.GetSCIMAPIKeyRequest]) (*connect.Response[frontendv1.GetSCIMAPIKeyResponse], error) {
	res, err := s.Store.GetSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateSCIMAPIKey(ctx context.Context, req *connect.Request[frontendv1.CreateSCIMAPIKeyRequest]) (*connect.Response[frontendv1.CreateSCIMAPIKeyResponse], error) {
	res, err := s.Store.CreateSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateSCIMAPIKey(ctx context.Context, req *connect.Request[frontendv1.UpdateSCIMAPIKeyRequest]) (*connect.Response[frontendv1.UpdateSCIMAPIKeyResponse], error) {
	res, err := s.Store.UpdateSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteSCIMAPIKey(ctx context.Context, req *connect.Request[frontendv1.DeleteSCIMAPIKeyRequest]) (*connect.Response[frontendv1.DeleteSCIMAPIKeyResponse], error) {
	res, err := s.Store.DeleteSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeSCIMAPIKey(ctx context.Context, req *connect.Request[frontendv1.RevokeSCIMAPIKeyRequest]) (*connect.Response[frontendv1.RevokeSCIMAPIKeyResponse], error) {
	res, err := s.Store.RevokeSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
