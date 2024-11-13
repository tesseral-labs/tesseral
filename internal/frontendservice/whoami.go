package frontendservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
)

func (s *FrontendService) WhoAmI(ctx context.Context, req *connect.Request[any]) (*connect.Response[any], error) {
	res, err := s.Store.WhoAmI(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return &res, nil
}