package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) ListSAMLConnections(ctx context.Context, req *connect.Request[frontendv1.ListSAMLConnectionsRequest]) (*connect.Response[frontendv1.ListSAMLConnectionsResponse], error) {
	res, err := s.Store.ListSAMLConnections(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetSAMLConnection(ctx context.Context, req *connect.Request[frontendv1.GetSAMLConnectionRequest]) (*connect.Response[frontendv1.GetSAMLConnectionResponse], error) {
	res, err := s.Store.GetSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateSAMLConnection(ctx context.Context, req *connect.Request[frontendv1.CreateSAMLConnectionRequest]) (*connect.Response[frontendv1.CreateSAMLConnectionResponse], error) {
	res, err := s.Store.CreateSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateSAMLConnection(ctx context.Context, req *connect.Request[frontendv1.UpdateSAMLConnectionRequest]) (*connect.Response[frontendv1.UpdateSAMLConnectionResponse], error) {
	res, err := s.Store.UpdateSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteSAMLConnection(ctx context.Context, req *connect.Request[frontendv1.DeleteSAMLConnectionRequest]) (*connect.Response[frontendv1.DeleteSAMLConnectionResponse], error) {
	res, err := s.Store.DeleteSAMLConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
