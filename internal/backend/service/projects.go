package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
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

func (s *Service) DisableProjectLogins(ctx context.Context, req *connect.Request[backendv1.DisableProjectLoginsRequest]) (*connect.Response[backendv1.DisableProjectLoginsResponse], error) {
	res, err := s.Store.DisableProjectLogins(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) EnableProjectLogins(ctx context.Context, req *connect.Request[backendv1.EnableProjectLoginsRequest]) (*connect.Response[backendv1.EnableProjectLoginsResponse], error) {
	res, err := s.Store.EnableProjectLogins(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
