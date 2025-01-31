package store

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Store) GetPasskeyOptions(ctx context.Context, req *intermediatev1.GetPasskeyOptionsRequest) (*intermediatev1.GetPasskeyOptionsResponse, error) {
	if !authn.IntermediateSession(ctx).EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		panic(fmt.Errorf("read random bytes: %w", err))
	}

	// TODO save this

	return &intermediatev1.GetPasskeyOptionsResponse{
		RpId:            *qProject.AuthDomain,
		RpName:          qProject.DisplayName,
		UserId:          fmt.Sprintf("%s|%s", *qIntermediateSession.Email, *qIntermediateSession.OrganizationID),
		UserDisplayName: *qIntermediateSession.Email,
		Challenge:       challenge,
	}, nil
}
