package service

import (
	"github.com/tesseral-labs/tesseral/internal/common/accesstoken"
	"github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1/frontendv1connect"
	"github.com/tesseral-labs/tesseral/internal/frontend/store"
)

type Service struct {
	Store             *store.Store
	AccessTokenIssuer *accesstoken.Issuer
	frontendv1connect.UnimplementedFrontendServiceHandler
}
