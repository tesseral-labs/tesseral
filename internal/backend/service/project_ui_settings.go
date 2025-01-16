package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) UpdateProjectUISettings(ctx context.Context, req *connect.Request[backendv1.UpdateProjectUISettingsRequest]) (*connect.Response[backendv1.UpdateProjectUISettingsResponse], error) {
	res, err := s.Store.UpdateProjectUISettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to update project UI settings: %w", err)
	}

	return connect.NewResponse(res), nil
}
