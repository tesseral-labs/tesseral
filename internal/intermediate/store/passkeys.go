package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
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

	qUser, err := q.GetOrganizationUserByEmail(ctx, queries.GetOrganizationUserByEmailParams{
		OrganizationID: *qIntermediateSession.OrganizationID,
		Email:          *qIntermediateSession.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization user by email: %w", err)
	}

	return &intermediatev1.GetPasskeyOptionsResponse{
		RpId:            *qProject.AuthDomain,
		UserId:          idformat.User.Format(qUser.ID),
		UserDisplayName: qUser.Email,
	}, nil
}
