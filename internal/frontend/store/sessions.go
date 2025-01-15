package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) Whoami(ctx context.Context, req *frontendv1.WhoAmIRequest) (*frontendv1.WhoAmIResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID := authn.UserID(ctx)

	qUser, err := q.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by id: %w", err))
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	qOrganization, err := q.GetOrganizationByID(ctx, qUser.OrganizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewFailedPreconditionError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	return &frontendv1.WhoAmIResponse{
		Id:             idformat.User.Format(qUser.ID),
		Email:          qUser.Email,
		OrganizationId: idformat.Organization.Format(qUser.OrganizationID),
		Organization: &frontendv1.Organization{
			Id:          idformat.Organization.Format(qOrganization.ID),
			DisplayName: qOrganization.DisplayName,
		},
	}, nil
}

func parseSession(qSession *queries.Session) *frontendv1.Session {
	return &frontendv1.Session{
		Id:         idformat.Session.Format(qSession.ID),
		UserId:     idformat.User.Format(qSession.UserID),
		CreateTime: derefTimeOrNil(qSession.CreateTime),
		ExpireTime: derefTimeOrNil(qSession.ExpireTime),
		Revoked:    qSession.Revoked,
	}
}
