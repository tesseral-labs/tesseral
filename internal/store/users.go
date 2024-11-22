package store

import (
	"context"

	"github.com/google/uuid"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

func (s *Store) CreateUser(ctx context.Context, req *openauthv1.User) (*openauthv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	createdUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:              uuid.New(),
		OrganizationID:  organizationId,
		UnverifiedEmail: &req.UnverifiedEmail,
		VerifiedEmail:   &req.VerifiedEmail,
		GoogleUserID:    &req.GoogleUserId,
		MicrosoftUserID: &req.MicrosoftUserId,
		PasswordBcrypt:  &req.PasswordBcrypt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseUser(&createdUser), nil
}

func (s *Store) GetUser(ctx context.Context, req *backendv1.GetUserRequest) (*openauthv1.User, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userId, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := q.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return parseUser(&user), nil
}

// TODO: Ensure that this function can only be called via a backend service reuqest
func (s *Store) ListUsers(ctx context.Context, req *backendv1.ListUsersRequest) (*backendv1.ListUsersResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	userRecords, err := q.ListUsersByOrganization(ctx, queries.ListUsersByOrganizationParams{
		OrganizationID: organizationId,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, err
	}

	users := []*openauthv1.User{}
	for _, userRecord := range userRecords {
		users = append(users, parseUser(&userRecord))
	}

	var nextPageToken string
	if len(users) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(userRecords[limit].ID)
		users = users[:limit]
	}

	return &backendv1.ListUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) UpdateUser(ctx context.Context, req *backendv1.UpdateUserRequest) (*openauthv1.User, error) {
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
	if req.User.UnverifiedEmail != "" {
		updates.UnverifiedEmail = &req.User.UnverifiedEmail
	}
	if req.User.VerifiedEmail != "" {
		updates.VerifiedEmail = &req.User.VerifiedEmail
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

func (s *Store) UpdateUserPassword(ctx context.Context, req *backendv1.UpdateUserPasswordRequest) (*openauthv1.User, error) {
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
	passwordBcrypt, err := generateBcryptHash(req.Password)
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

func parseUser(user *queries.User) *openauthv1.User {
	return &openauthv1.User{
		Id:              user.ID.String(),
		OrganizationId:  user.OrganizationID.String(),
		UnverifiedEmail: *user.UnverifiedEmail,
		VerifiedEmail:   *user.VerifiedEmail,
		GoogleUserId:    *user.GoogleUserID,
		MicrosoftUserId: *user.MicrosoftUserID,
		PasswordBcrypt:  *user.PasswordBcrypt,
	}
}
