package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) GetProject(ctx context.Context, req *connect.Request[backendv1.GetProjectRequest]) (*connect.Response[backendv1.GetProjectResponse], error) {
	res, err := s.Store.GetProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateProject(ctx context.Context, req *connect.Request[backendv1.UpdateProjectRequest]) (*connect.Response[backendv1.UpdateProjectResponse], error) {
	res, err := s.Store.UpdateProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) LockoutProject(ctx context.Context, req *connect.Request[backendv1.LockoutProjectRequest]) (*connect.Response[backendv1.LockoutProjectResponse], error) {
	res, err := s.Store.LockoutProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UnlockProject(ctx context.Context, req *connect.Request[backendv1.UnlockProjectRequest]) (*connect.Response[backendv1.UnlockProjectResponse], error) {
	res, err := s.Store.UnlockProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
