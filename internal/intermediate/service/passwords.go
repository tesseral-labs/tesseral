package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) RegisterPassword(ctx context.Context, req *connect.Request[intermediatev1.RegisterPasswordRequest]) (*connect.Response[intermediatev1.RegisterPasswordResponse], error) {
	res, err := s.Store.RegisterPassword(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) VerifyPassword(ctx context.Context, req *connect.Request[intermediatev1.VerifyPasswordRequest]) (*connect.Response[intermediatev1.VerifyPasswordResponse], error) {
	res, err := s.Store.VerifyPassword(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) IssuePasswordResetCode(ctx context.Context, req *connect.Request[intermediatev1.IssuePasswordResetCodeRequest]) (*connect.Response[intermediatev1.IssuePasswordResetCodeResponse], error) {
	res, err := s.Store.IssuePasswordResetCode(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) VerifyPasswordResetCode(ctx context.Context, req *connect.Request[intermediatev1.VerifyPasswordResetCodeRequest]) (*connect.Response[intermediatev1.VerifyPasswordResetCodeResponse], error) {
	res, err := s.Store.VerifyPasswordResetCode(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
