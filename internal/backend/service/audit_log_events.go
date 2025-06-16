package service

import (
	"context"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) CreateAuditLogEvent(ctx context.Context, req *connect.Request[backendv1.CreateAuditLogEventRequest]) (*connect.Response[backendv1.CreateAuditLogEventResponse], error) {
	res, err := s.Store.CreateCustomAuditLogEvent(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}
