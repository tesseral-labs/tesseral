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
