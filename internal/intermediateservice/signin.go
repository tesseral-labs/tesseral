package intermediateservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/gen/intermediate/v1"
)

func (s *IntermediateService) SignInWithEmail(
	ctx context.Context,
	req *connect.Request[intermediatev1.SignInWithEmailRequest],
) (*connect.Response[intermediatev1.SignInWithEmailResponse], error) {
	res, err := s.Store.SignInWithEmail(&ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
