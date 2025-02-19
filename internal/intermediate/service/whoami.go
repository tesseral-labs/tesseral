package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) Whoami(ctx context.Context, req *connect.Request[intermediatev1.WhoamiRequest]) (*connect.Response[intermediatev1.WhoamiResponse], error) {
	res, err := s.Store.Whoami(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
