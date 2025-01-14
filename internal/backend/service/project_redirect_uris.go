package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) CreateProjectRedirectURI(ctx context.Context, req *connect.Request[backendv1.CreateProjectRedirectURIRequest]) (*connect.Response[backendv1.CreateProjectRedirectURIResponse], error) {
	res, err := s.Store.CreateProjectRedirectURI(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteProjectRedirectURI(ctx context.Context, req *connect.Request[backendv1.DeleteProjectRedirectURIRequest]) (*connect.Response[backendv1.DeleteProjectRedirectURIResponse], error) {
	res, err := s.Store.DeleteProjectRedirectURI(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetProjectRedirectURI(ctx context.Context, req *connect.Request[backendv1.GetProjectRedirectURIRequest]) (*connect.Response[backendv1.GetProjectRedirectURIResponse], error) {
	res, err := s.Store.GetProjectRedirectURI(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListProjectRedirectURIs(ctx context.Context, req *connect.Request[backendv1.ListProjectRedirectURIsRequest]) (*connect.Response[backendv1.ListProjectRedirectURIsResponse], error) {
	res, err := s.Store.ListProjectRedirectURIs(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateProjectRedirectURI(ctx context.Context, req *connect.Request[backendv1.UpdateProjectRedirectURIRequest]) (*connect.Response[backendv1.UpdateProjectRedirectURIResponse], error) {
	res, err := s.Store.UpdateProjectRedirectURI(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
