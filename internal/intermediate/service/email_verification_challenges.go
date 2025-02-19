package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) IssueEmailVerificationChallenge(ctx context.Context, req *connect.Request[intermediatev1.IssueEmailVerificationChallengeRequest]) (*connect.Response[intermediatev1.IssueEmailVerificationChallengeResponse], error) {
	res, err := s.Store.IssueEmailVerificationChallenge(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) VerifyEmailChallenge(ctx context.Context, req *connect.Request[intermediatev1.VerifyEmailChallengeRequest]) (*connect.Response[intermediatev1.VerifyEmailChallengeResponse], error) {
	res, err := s.Store.VerifyEmailChallenge(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
