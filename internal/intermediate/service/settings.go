package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) GetSettings(ctx context.Context, req *connect.Request[intermediatev1.GetSettingsRequest]) (*connect.Response[intermediatev1.GetSettingsResponse], error) {
	res, err := s.Store.GetSettings(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
