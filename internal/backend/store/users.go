package store

import (
	"context"

	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/crypto/bcrypt"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) UpdateUser(ctx context.Context, req *backendv1.UpdateUserRequest) (*backendv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userId, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	organizationId, err := idformat.Organization.Parse(req.User.OrganizationId)
	if err != nil {
		return nil, err
	}

	updates := queries.UpdateUserParams{
		ID: userId,
	}

	// Conditionally update organizationID
	if req.User.OrganizationId != "" {
		updates.OrganizationID = organizationId
	}

	// Conditionally update email addresses
	if req.User.Email != "" {
		updates.Email = req.User.Email
	}

	// Conditionally update login method user IDs
	if req.User.GoogleUserId != "" {
		updates.GoogleUserID = &req.User.GoogleUserId
	}
	if req.User.MicrosoftUserId != "" {
		updates.MicrosoftUserID = &req.User.MicrosoftUserId
	}

	// Conditionall update password
	if req.User.PasswordBcrypt != "" {
		updates.PasswordBcrypt = &req.User.PasswordBcrypt
	}

	user, err := q.UpdateUser(ctx, updates)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseUser(&user), nil
}

func (s *Store) UpdateUserPassword(ctx context.Context, req *backendv1.UpdateUserPasswordRequest) (*backendv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userId, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	// Bcrypt the password before storing
	passwordBcrypt, err := bcrypt.GenerateBcryptHash(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := q.UpdateUserPassword(ctx, queries.UpdateUserPasswordParams{
		ID:             userId,
		PasswordBcrypt: &passwordBcrypt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseUser(&user), nil
}

func parseUser(qUser *queries.User) *backendv1.User {
	return &backendv1.User{
		Id:             idformat.User.Format(qUser.ID),
		OrganizationId: idformat.Organization.Format(qUser.OrganizationID),
	}
}
