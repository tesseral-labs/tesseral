package store

import (
	"context"
	"errors"

	"github.com/openauth/openauth/internal/crypto/bcrypt"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errInvalidUserId = errors.New("invalid user id")

func (s *Store) SetUserPassword(ctx context.Context, req *frontendv1.SetUserPasswordRequest) (*frontendv1.SetUserPasswordResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	sessionUserID := authn.UserID(ctx)
	userID, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	if sessionUserID != userID {
		return nil, errInvalidUserId
	}

	passwordBcrypt, err := bcrypt.GenerateBcryptHash(req.Password)
	if err != nil {
		return nil, err
	}

	_, err = q.SetUserPassword(ctx, queries.SetUserPasswordParams{
		ID:             userID,
		PasswordBcrypt: &passwordBcrypt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &frontendv1.SetUserPasswordResponse{}, nil
}

func parseUser(qUser *queries.User) *frontendv1.User {
	return &frontendv1.User{
		Id:             idformat.User.Format(qUser.ID),
		OrganizationId: idformat.Organization.Format(qUser.OrganizationID),
	}
}
