package service

import (
	"github.com/openauth/openauth/internal/common/accesstoken"
	"github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1/frontendv1connect"
	"github.com/openauth/openauth/internal/frontend/store"
)

type Service struct {
	Store             *store.Store
	AccessTokenIssuer *accesstoken.Issuer
	frontendv1connect.UnimplementedFrontendServiceHandler
}
