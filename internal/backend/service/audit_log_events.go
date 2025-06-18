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

func (s *Service) ConsoleListAuditLogEvents(ctx context.Context, req *connect.Request[backendv1.ConsoleListAuditLogEventsRequest]) (*connect.Response[backendv1.ConsoleListAuditLogEventsResponse], error) {
	res, err := s.Store.ConsoleListCustomAuditLogEvents(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ConsoleListAuditLogEventNames(ctx context.Context, req *connect.Request[backendv1.ConsoleListAuditLogEventNamesRequest]) (*connect.Response[backendv1.ConsoleListAuditLogEventNamesResponse], error) {
	res, err := s.Store.ConsoleListAuditLogEventNames(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}
