package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) UpdateUser(
	ctx context.Context,
	req *connect.Request[backendv1.UpdateUserRequest],
) (*connect.Response[backendv1.User], error) {
	res, err := s.Store.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateUserPassword(
	ctx context.Context,
	req *connect.Request[backendv1.UpdateUserPasswordRequest],
) (*connect.Response[backendv1.User], error) {
	res, err := s.Store.UpdateUserPassword(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
