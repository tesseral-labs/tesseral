package service

import (
	"github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1/backendv1connect"
	"github.com/openauth/openauth/internal/backend/store"
)

type Service struct {
	Store *store.Store
	backendv1connect.UnimplementedBackendServiceHandler
}
