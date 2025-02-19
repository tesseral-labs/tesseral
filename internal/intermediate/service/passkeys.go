package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) GetPasskeyOptions(ctx context.Context, req *connect.Request[intermediatev1.GetPasskeyOptionsRequest]) (*connect.Response[intermediatev1.GetPasskeyOptionsResponse], error) {
	res, err := s.Store.GetPasskeyOptions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) RegisterPasskey(ctx context.Context, req *connect.Request[intermediatev1.RegisterPasskeyRequest]) (*connect.Response[intermediatev1.RegisterPasskeyResponse], error) {
	res, err := s.Store.RegisterPasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) IssuePasskeyChallenge(ctx context.Context, req *connect.Request[intermediatev1.IssuePasskeyChallengeRequest]) (*connect.Response[intermediatev1.IssuePasskeyChallengeResponse], error) {
	res, err := s.Store.IssuePasskeyChallenge(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) VerifyPasskey(ctx context.Context, req *connect.Request[intermediatev1.VerifyPasskeyRequest]) (*connect.Response[intermediatev1.VerifyPasskeyResponse], error) {
	res, err := s.Store.VerifyPasskey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
