package frontendservice

import (
	"context"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth-dev/openauth/internal/gen/frontend/v1"
)

func (s *FrontendService) WhoAmI(
	ctx context.Context,
	req *connect.Request[frontendv1.WhoAmIRequest],
) (*connect.Response[frontendv1.WhoAmIResponse], error) {
	res := connect.Response[frontendv1.WhoAmIResponse]{
		Msg: &frontendv1.WhoAmIResponse{
			Id:          "123",
			DisplayName: "John Doe",
			Email:       "john.doe@example.com",
		},
	}

	// res, err := s.Store.WhoAmI(ctx, req)
	// if err != nil {
	// 	return nil, fmt.Errorf("store: %w", err)
	// }

	return &res, nil
}
