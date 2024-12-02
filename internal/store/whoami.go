package store

import (
	"context"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/gen/frontend/v1"
)

func (s *Store) WhoAmI(ctx context.Context, req *connect.Request[frontendv1.WhoAmIRequest]) (connect.Response[frontendv1.WhoAmIResponse], error) {
	return connect.Response[frontendv1.WhoAmIResponse]{}, nil
}
