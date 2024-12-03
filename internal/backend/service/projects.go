package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) CreateProject(
	ctx context.Context,
	req *connect.Request[backendv1.CreateProjectRequest],
) (*connect.Response[backendv1.Project], error) {
	res, err := s.Store.CreateProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetProject(
	ctx context.Context,
	req *connect.Request[backendv1.GetProjectRequest],
) (*connect.Response[backendv1.Project], error) {
	res, err := s.Store.GetProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListProjects(
	ctx context.Context,
	req *connect.Request[backendv1.ListProjectsRequest],
) (*connect.Response[backendv1.ListProjectsResponse], error) {
	res, err := s.Store.ListProjects(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateProject(
	ctx context.Context,
	req *connect.Request[backendv1.UpdateProjectRequest],
) (*connect.Response[backendv1.Project], error) {
	res, err := s.Store.UpdateProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
