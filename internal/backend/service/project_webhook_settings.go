package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetProjectWebhookManagementURL(ctx context.Context, req *connect.Request[backendv1.GetProjectWebhookManagementURLRequest]) (*connect.Response[backendv1.GetProjectWebhookManagementURLResponse], error) {
	res, err := s.Store.GetProjectWebhookManagementURL(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("get project webhook management url: %w", err)
	}

	return connect.NewResponse(res), nil
}
