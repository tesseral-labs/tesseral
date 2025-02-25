package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetVaultDomainSettings(ctx context.Context, req *connect.Request[backendv1.GetVaultDomainSettingsRequest]) (*connect.Response[backendv1.GetVaultDomainSettingsResponse], error) {
	res, err := s.Store.GetVaultDomainSettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateVaultDomainSettings(ctx context.Context, req *connect.Request[backendv1.UpdateVaultDomainSettingsRequest]) (*connect.Response[backendv1.UpdateVaultDomainSettingsResponse], error) {
	res, err := s.Store.UpdateVaultDomainSettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) EnableCustomVaultDomain(ctx context.Context, req *connect.Request[backendv1.EnableCustomVaultDomainRequest]) (*connect.Response[backendv1.EnableCustomVaultDomainResponse], error) {
	res, err := s.Store.EnableCustomVaultDomain(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
