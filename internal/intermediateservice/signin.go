package intermediateservice

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth-dev/openauth/internal/gen/intermediate/v1"
)

func (s *IntermediateService) SignInWithEmail(
	ctx context.Context,
	req *connect.Request[intermediatev1.SignInWithEmailRequest],
) (*connect.Response[intermediatev1.SignInWithEmailResponse], error) {
	slog.Info("sign in with email", "email", req.Msg.Email)

	res, err := s.Store.SignInWithEmail(&ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	slog.Info("sign in with email", "res", res)

	return connect.NewResponse(res), nil
}
