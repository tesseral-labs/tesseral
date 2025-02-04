package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/bcryptcost"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) SetUserPassword(ctx context.Context, req *frontendv1.SetPasswordRequest) (*frontendv1.SetPasswordResponse, error) {
	// Check if the password is compromised.
	pwned, err := s.hibp.Pwned(ctx, req.Password)
	if err != nil {
		return nil, fmt.Errorf("check password against HIBP: %w", err)
	}
	if pwned {
		return nil, apierror.NewFailedPreconditionError("password is compromised", errors.New("password is compromised"))
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	passwordBcryptBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptcost.Cost)
	if err != nil {
		return nil, apierror.NewFailedPreconditionError("could not generate password hash", fmt.Errorf("generate bcrypt hash: %w", err))
	}

	passwordBcrypt := string(passwordBcryptBytes)
	if _, err = q.SetPassword(ctx, queries.SetPasswordParams{
		ID:             authn.UserID(ctx),
		PasswordBcrypt: &passwordBcrypt,
	}); err != nil {
		return nil, fmt.Errorf("set password: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.SetPasswordResponse{}, nil
}

func (s *Store) ListUsers(ctx context.Context, req *frontendv1.ListUsersRequest) (*frontendv1.ListUsersResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qUsers, err := q.ListUsers(ctx, queries.ListUsersParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	var users []*frontendv1.User
	for _, qUser := range qUsers {
		users = append(users, parseUser(qUser))
	}

	var nextPageToken string
	if len(users) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUsers[limit].ID)
		users = users[:limit]
	}

	return &frontendv1.ListUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetUser(ctx context.Context, req *frontendv1.GetUserRequest) (*frontendv1.GetUserResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	qUser, err := q.GetUser(ctx, queries.GetUserParams{
		ID:             userID,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by id: %w", err))
		}

		return nil, fmt.Errorf("get user: %w", err)
	}

	return &frontendv1.GetUserResponse{User: parseUser(qUser)}, nil
}

func (s *Store) UpdateUser(ctx context.Context, req *frontendv1.UpdateUserRequest) (*frontendv1.UpdateUserResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// Only owners can call UpdateUser.
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	userID, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	// Fetch the existing user details. Also acts as authz check.
	qUser, err := q.GetUser(ctx, queries.GetUserParams{
		ID:             userID,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by id: %w", err))
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	updates := queries.UpdateUserParams{
		ID:      userID,
		IsOwner: qUser.IsOwner,
	}

	if req.User.Owner != nil {
		updates.IsOwner = *req.User.Owner
	}

	// Perform the update.
	qUpdatedUser, err := q.UpdateUser(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	// Commit the transaction.
	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.UpdateUserResponse{User: parseUser(qUpdatedUser)}, nil
}

func parseUser(qUser queries.User) *frontendv1.User {
	return &frontendv1.User{
		Id:                  idformat.User.Format(qUser.ID),
		CreateTime:          timestamppb.New(*qUser.CreateTime),
		UpdateTime:          timestamppb.New(*qUser.UpdateTime),
		Email:               qUser.Email,
		Owner:               &qUser.IsOwner,
		GoogleUserId:        derefOrEmpty(qUser.GoogleUserID),
		MicrosoftUserId:     derefOrEmpty(qUser.MicrosoftUserID),
		HasAuthenticatorApp: qUser.AuthenticatorAppSecretCiphertext != nil,
	}
}
