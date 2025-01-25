package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListSessions(ctx context.Context, req *connect.Request[backendv1.ListSessionsRequest]) (*connect.Response[backendv1.ListSessionsResponse], error) {
	res, err := s.Store.ListSessions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetSession(ctx context.Context, req *connect.Request[backendv1.GetSessionRequest]) (*connect.Response[backendv1.GetSessionResponse], error) {
	res, err := s.Store.GetSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeAllOrganizationSessions(ctx context.Context, req *connect.Request[backendv1.RevokeAllOrganizationSessionsRequest]) (*connect.Response[backendv1.RevokeAllOrganizationSessionsResponse], error) {
	res, err := s.Store.RevokeAllOrganizationSessions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeAllProjectSessions(ctx context.Context, req *connect.Request[backendv1.RevokeAllProjectSessionsRequest]) (*connect.Response[backendv1.RevokeAllProjectSessionsResponse], error) {
	res, err := s.Store.RevokeAllProjectSessions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
