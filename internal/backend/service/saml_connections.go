package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListSAMLConnections(ctx context.Context, req *connect.Request[backendv1.ListSAMLConnectionsRequest]) (*connect.Response[backendv1.ListSAMLConnectionsResponse], error) {
	res, err := s.Store.ListSAMLConnections(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetSAMLConnection(ctx context.Context, req *connect.Request[backendv1.GetSAMLConnectionRequest]) (*connect.Response[backendv1.GetSAMLConnectionResponse], error) {
	res, err := s.Store.GetSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateSAMLConnection(ctx context.Context, req *connect.Request[backendv1.CreateSAMLConnectionRequest]) (*connect.Response[backendv1.CreateSAMLConnectionResponse], error) {
	res, err := s.Store.CreateSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateSAMLConnection(ctx context.Context, req *connect.Request[backendv1.UpdateSAMLConnectionRequest]) (*connect.Response[backendv1.UpdateSAMLConnectionResponse], error) {
	res, err := s.Store.UpdateSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteSAMLConnection(ctx context.Context, req *connect.Request[backendv1.DeleteSAMLConnectionRequest]) (*connect.Response[backendv1.DeleteSAMLConnectionResponse], error) {
	res, err := s.Store.DeleteSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
