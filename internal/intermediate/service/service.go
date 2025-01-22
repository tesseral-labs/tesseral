package service

import (
	"github.com/openauth/openauth/internal/common/accesstoken"
	"github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1/intermediatev1connect"
	"github.com/openauth/openauth/internal/intermediate/store"
)

type Service struct {
	Store             *store.Store
	AccessTokenIssuer *accesstoken.Issuer
	intermediatev1connect.UnimplementedIntermediateServiceHandler
}
