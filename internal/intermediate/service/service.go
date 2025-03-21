package service

import (
	"github.com/tesseral-labs/tesseral/internal/common/accesstoken"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1/intermediatev1connect"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store"
)

type Service struct {
	Store             *store.Store
	AccessTokenIssuer *accesstoken.Issuer
	Cookier           *cookies.Cookier
	intermediatev1connect.UnimplementedIntermediateServiceHandler
}
