package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetOrganization(ctx context.Context, req *connect.Request[backendv1.GetOrganizationRequest]) (*connect.Response[backendv1.GetOrganizationResponse], error) {
	res, err := s.Store.GetOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListOrganizations(ctx context.Context, req *connect.Request[backendv1.ListOrganizationsRequest]) (*connect.Response[backendv1.ListOrganizationsResponse], error) {
	res, err := s.Store.ListOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateOrganization(ctx context.Context, req *connect.Request[backendv1.CreateOrganizationRequest]) (*connect.Response[backendv1.CreateOrganizationResponse], error) {
	res, err := s.Store.CreateOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganization(ctx context.Context, req *connect.Request[backendv1.UpdateOrganizationRequest]) (*connect.Response[backendv1.UpdateOrganizationResponse], error) {
	res, err := s.Store.UpdateOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteOrganization(ctx context.Context, req *connect.Request[backendv1.DeleteOrganizationRequest]) (*connect.Response[backendv1.DeleteOrganizationResponse], error) {
	res, err := s.Store.DeleteOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DisableOrganizationLogins(ctx context.Context, req *connect.Request[backendv1.DisableOrganizationLoginsRequest]) (*connect.Response[backendv1.DisableOrganizationLoginsResponse], error) {
	res, err := s.Store.DisableOrganizationLogins(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) EnableOrganizationLogins(ctx context.Context, req *connect.Request[backendv1.EnableOrganizationLoginsRequest]) (*connect.Response[backendv1.EnableOrganizationLoginsResponse], error) {
	res, err := s.Store.EnableOrganizationLogins(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
