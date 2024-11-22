package backendservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s *BackendService) UpdateUser(
	ctx context.Context,
	req *connect.Request[backendv1.UpdateUserRequest],
) (*connect.Response[openauthv1.User], error) {
	res, err := s.Store.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *BackendService) UpdateUserPassword(
	ctx context.Context,
	req *connect.Request[backendv1.UpdateUserPasswordRequest],
) (*connect.Response[openauthv1.User], error) {
	res, err := s.Store.UpdateUserPassword(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
