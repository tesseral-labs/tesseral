package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) ListOIDCConnections(ctx context.Context, req *connect.Request[frontendv1.ListOIDCConnectionsRequest]) (*connect.Response[frontendv1.ListOIDCConnectionsResponse], error) {
	res, err := s.Store.ListOIDCConnections(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetOIDCConnection(ctx context.Context, req *connect.Request[frontendv1.GetOIDCConnectionRequest]) (*connect.Response[frontendv1.GetOIDCConnectionResponse], error) {
	res, err := s.Store.GetOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateOIDCConnection(ctx context.Context, req *connect.Request[frontendv1.CreateOIDCConnectionRequest]) (*connect.Response[frontendv1.CreateOIDCConnectionResponse], error) {
	res, err := s.Store.CreateOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOIDCConnection(ctx context.Context, req *connect.Request[frontendv1.UpdateOIDCConnectionRequest]) (*connect.Response[frontendv1.UpdateOIDCConnectionResponse], error) {
	res, err := s.Store.UpdateOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteOIDCConnection(ctx context.Context, req *connect.Request[frontendv1.DeleteOIDCConnectionRequest]) (*connect.Response[frontendv1.DeleteOIDCConnectionResponse], error) {
	res, err := s.Store.DeleteOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
