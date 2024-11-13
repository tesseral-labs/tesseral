package store

import (
	"context"

	"connectrpc.com/connect"
)

func (s *Store) WhoAmI(ctx context.Context, req *connect.Request[any]) (connect.Response[any], error) {
	return connect.Response[any]{}, nil
}