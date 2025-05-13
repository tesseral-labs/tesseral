package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetRBACPolicy(ctx context.Context, req *connect.Request[backendv1.GetRBACPolicyRequest]) (*connect.Response[backendv1.GetRBACPolicyResponse], error) {
	res, err := s.Store.GetRBACPolicy(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateRBACPolicy(ctx context.Context, req *connect.Request[backendv1.UpdateRBACPolicyRequest]) (*connect.Response[backendv1.UpdateRBACPolicyResponse], error) {
	res, err := s.Store.UpdateRBACPolicy(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
