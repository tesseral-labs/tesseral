package service

import (
	"github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1/frontendv1connect"
	"github.com/openauth/openauth/internal/frontend/store"
)

type Service struct {
	Store *store.Store
	frontendv1connect.UnimplementedFrontendServiceHandler
}
