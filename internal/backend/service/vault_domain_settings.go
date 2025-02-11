package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) GetVaultDomainSettings(ctx context.Context, req *connect.Request[backendv1.GetVaultDomainSettingsRequest]) (*connect.Response[backendv1.GetVaultDomainSettingsResponse], error) {
	res, err := s.Store.GetVaultDomainSettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
