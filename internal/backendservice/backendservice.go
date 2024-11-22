package backendservice

import (
	"github.com/openauth-dev/openauth/internal/gen/backend/v1/backendv1connect"
	"github.com/openauth-dev/openauth/internal/store"
)

type BackendService struct {
	Store *store.Store
	backendv1connect.UnimplementedBackendServiceHandler
}
