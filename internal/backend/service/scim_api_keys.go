package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ListSCIMAPIKeys(ctx context.Context, req *connect.Request[backendv1.ListSCIMAPIKeysRequest]) (*connect.Response[backendv1.ListSCIMAPIKeysResponse], error) {
	res, err := s.Store.ListSCIMAPIKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetSCIMAPIKey(ctx context.Context, req *connect.Request[backendv1.GetSCIMAPIKeyRequest]) (*connect.Response[backendv1.GetSCIMAPIKeyResponse], error) {
	res, err := s.Store.GetSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateSCIMAPIKey(ctx context.Context, req *connect.Request[backendv1.CreateSCIMAPIKeyRequest]) (*connect.Response[backendv1.CreateSCIMAPIKeyResponse], error) {
	res, err := s.Store.CreateSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateSCIMAPIKey(ctx context.Context, req *connect.Request[backendv1.UpdateSCIMAPIKeyRequest]) (*connect.Response[backendv1.UpdateSCIMAPIKeyResponse], error) {
	res, err := s.Store.UpdateSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteSCIMAPIKey(ctx context.Context, req *connect.Request[backendv1.DeleteSCIMAPIKeyRequest]) (*connect.Response[backendv1.DeleteSCIMAPIKeyResponse], error) {
	res, err := s.Store.DeleteSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeSCIMAPIKey(ctx context.Context, req *connect.Request[backendv1.RevokeSCIMAPIKeyRequest]) (*connect.Response[backendv1.RevokeSCIMAPIKeyResponse], error) {
	res, err := s.Store.RevokeSCIMAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
