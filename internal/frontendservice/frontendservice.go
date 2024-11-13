package frontendservice

import (
	"github.com/openauth-dev/openauth/internal/gen/frontend/v1/frontendv1connect"
	"github.com/openauth-dev/openauth/internal/store"
)

type FrontendService struct {
  Store *store.Store
	frontendv1connect.UnimplementedFrontendServiceHandler
}