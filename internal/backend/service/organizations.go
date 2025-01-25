package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
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

func (s *Service) LockoutOrganization(ctx context.Context, req *connect.Request[backendv1.LockoutOrganizationRequest]) (*connect.Response[backendv1.LockoutOrganizationResponse], error) {
	res, err := s.Store.LockoutOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UnlockOrganization(ctx context.Context, req *connect.Request[backendv1.UnlockOrganizationRequest]) (*connect.Response[backendv1.UnlockOrganizationResponse], error) {
	res, err := s.Store.UnlockOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
