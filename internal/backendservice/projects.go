package backendservice

import (
	"context"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s *BackendService) CreateProject(
	ctx context.Context, 
	req *connect.Request[openauthv1.CreateProjectRequest],
) (*connect.Response[backendv1.CreateProjectResponse], error) {
	res, err := s.Store.CreateProject(ctx, req.Msg)
	return nil, nil
}

func (s *BackendService) GetProject(
	ctx context.Context, 
	req *connect.Request[openauthv1.ResourceIdRequest],
) (*connect.Response[backendv1.GetProjectResponse], error) {
	return nil, nil
}

func (s *BackendService) ListProjects(
	ctx context.Context, 
	req *connect.Request[backendv1.ListProjectsRequest],
) (*connect.Response[backendv1.ListProjectsResponse], error) {
	return nil, nil
}

func (s *BackendService) UpdateProject(
	ctx context.Context, 
	req *connect.Request[backendv1.UpdateProjectRequest],
) (*connect.Response[backendv1.UpdateProjectResponse], error) {
	return nil, nil
}