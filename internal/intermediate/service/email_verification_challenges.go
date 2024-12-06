package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) VerifyEmailChallenge(
	ctx context.Context,
	req *connect.Request[intermediatev1.VerifyEmailChallengeRequest],
) (*connect.Response[intermediatev1.VerifyEmailChallengeResponse], error) {
	res, err := s.Store.CompleteEmailVerificationChallenge(ctx, req.Msg.Challenge)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
