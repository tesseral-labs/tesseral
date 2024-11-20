package backendservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s *BackendService) CreateProject(
	ctx context.Context, 
	req *connect.Request[backendv1.CreateProjectRequest],
) (*connect.Response[openauthv1.Project], error) {
	res, err := s.Store.CreateProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *BackendService) GetProject(
	ctx context.Context, 
	req *connect.Request[backendv1.GetProjectRequest],
) (*connect.Response[openauthv1.Project], error) {
	res, err := s.Store.GetProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *BackendService) ListProjects(
	ctx context.Context, 
	req *connect.Request[backendv1.ListProjectsRequest],
) (*connect.Response[backendv1.ListProjectsResponse], error) {
	res, err := s.Store.ListProjects(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *BackendService) UpdateProject(
	ctx context.Context, 
	req *connect.Request[backendv1.UpdateProjectRequest],
) (*connect.Response[openauthv1.Project], error) {
	res, err := s.Store.UpdateProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}