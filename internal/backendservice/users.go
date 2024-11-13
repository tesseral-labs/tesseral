package backendservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s *BackendService) CreateUser(
	ctx context.Context, 
	req *connect.Request[openauthv1.User],
) (*connect.Response[openauthv1.User], error) {
		res, err := s.Store.CreateUser(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) CreateUnverifiedUser(
	ctx context.Context, 
	req *connect.Request[openauthv1.CreateUnverifiedUserRequest],
) (*connect.Response[openauthv1.User], error) {
		res, err := s.Store.CreateUnverifiedUser(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) CreateGoogleUser(
	ctx context.Context, 
	req *connect.Request[openauthv1.CreateGoogleUserRequest],
) (*connect.Response[openauthv1.User], error) {
		res, err := s.Store.CreateGoogleUser(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) CreateMicrosoftUser(
	ctx context.Context, 
	req *connect.Request[openauthv1.CreateMicrosoftUserRequest],
) (*connect.Response[openauthv1.User], error) {
		res, err := s.Store.CreateMicrosoftUser(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) UpdateUser(
	ctx context.Context, 
	req *connect.Request[openauthv1.User],
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