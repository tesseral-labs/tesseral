package service

import (
	"context"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetProjectOnboardingProgress(ctx context.Context, req *connect.Request[backendv1.GetProjectOnboardingProgressRequest]) (*connect.Response[backendv1.GetProjectOnboardingProgressResponse], error) {
	resp, err := s.Store.GetProjectOnboardingProgress(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}

func (s *Service) UpdateProjectOnboardingProgress(ctx context.Context, req *connect.Request[backendv1.UpdateProjectOnboardingProgressRequest]) (*connect.Response[backendv1.UpdateProjectOnboardingProgressResponse], error) {
	resp, err := s.Store.UpdateProjectOnboardingProgress(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}
