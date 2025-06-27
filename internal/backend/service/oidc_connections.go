package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ListOIDCConnections(ctx context.Context, req *connect.Request[backendv1.ListOIDCConnectionsRequest]) (*connect.Response[backendv1.ListOIDCConnectionsResponse], error) {
	res, err := s.Store.ListOIDCConnections(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetOIDCConnection(ctx context.Context, req *connect.Request[backendv1.GetOIDCConnectionRequest]) (*connect.Response[backendv1.GetOIDCConnectionResponse], error) {
	res, err := s.Store.GetOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateOIDCConnection(ctx context.Context, req *connect.Request[backendv1.CreateOIDCConnectionRequest]) (*connect.Response[backendv1.CreateOIDCConnectionResponse], error) {
	res, err := s.Store.CreateOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOIDCConnection(ctx context.Context, req *connect.Request[backendv1.UpdateOIDCConnectionRequest]) (*connect.Response[backendv1.UpdateOIDCConnectionResponse], error) {
	res, err := s.Store.UpdateOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteOIDCConnection(ctx context.Context, req *connect.Request[backendv1.DeleteOIDCConnectionRequest]) (*connect.Response[backendv1.DeleteOIDCConnectionResponse], error) {
	res, err := s.Store.DeleteOIDCConnection(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
