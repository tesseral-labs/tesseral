package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetProjectLogoURLs(ctx context.Context, req *connect.Request[backendv1.GetProjectLogoURLsRequest]) (*connect.Response[backendv1.GetProjectLogoURLsResponse], error) {
	res, err := s.Store.GetProjectLogoURLs(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get project logo URLs: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetProjectUISettings(ctx context.Context, req *connect.Request[backendv1.GetProjectUISettingsRequest]) (*connect.Response[backendv1.GetProjectUISettingsResponse], error) {
	res, err := s.Store.GetProjectUISettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get project UI settings: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateProjectUISettings(ctx context.Context, req *connect.Request[backendv1.UpdateProjectUISettingsRequest]) (*connect.Response[backendv1.UpdateProjectUISettingsResponse], error) {
	res, err := s.Store.UpdateProjectUISettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to update project UI settings: %w", err)
	}

	return connect.NewResponse(res), nil
}
