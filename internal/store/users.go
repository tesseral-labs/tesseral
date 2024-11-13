package store

import (
	"context"

	"github.com/google/uuid"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
	"golang.org/x/crypto/bcrypt"
)

func (s * Store) CreateUser(ctx context.Context, req *openauthv1.User) (*openauthv1.User, error) {
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
		ID: uuid.New(),
		OrganizationID: organizationId,
		UnverifiedEmail: &req.UnverifiedEmail,
		VerifiedEmail: &req.VerifiedEmail,
		GoogleUserID: &req.GoogleUserId,
		MicrosoftUserID: &req.MicrosoftUserId,
		PasswordBcrypt: &req.PasswordBcrypt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformUser(&createdUser), nil
}

func (s *Store) CreateUnverifiedUser(ctx context.Context, req *openauthv1.CreateUnverifiedUserRequest) (*openauthv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	createdUser, err := q.CreateUnverifiedUser(ctx, queries.CreateUnverifiedUserParams{
		ID: uuid.New(),
		OrganizationID: organizationId,
		UnverifiedEmail: &req.UnverifiedEmail,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformUser(&createdUser), nil
}

func (s *Store) CreateGoogleUser(ctx context.Context, req *openauthv1.CreateGoogleUserRequest) (*openauthv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	createdUser, err := q.CreateGoogleUser(ctx, queries.CreateGoogleUserParams{
		ID: uuid.New(),
		OrganizationID: organizationId,
		GoogleUserID: &req.GoogleUserId,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformUser(&createdUser), nil
}

func (s *Store) CreateMicrosoftUser(ctx context.Context, req *openauthv1.CreateMicrosoftUserRequest) (*openauthv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	createdUser, err := q.CreateMicrosoftUser(ctx, queries.CreateMicrosoftUserParams{
		ID: uuid.New(),
		OrganizationID: organizationId,
		MicrosoftUserID: &req.MicrosoftUserId,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformUser(&createdUser), nil
}

func (s *Store) GetUser(ctx context.Context, req *openauthv1.ResourceIdRequest) (*openauthv1.User, error) {
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

	return transformUser(&user), nil
}

// TODO: Ensure that this function can only be called via a backend service reuqest
func (s * Store) ListUsers(ctx context.Context, req *backendv1.ListUsersRequest) (*backendv1.ListUsersResponse, error) {
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
		Limit: 				int32(limit + 1),
	})
	if err != nil {
		return nil, err
	}

	users := []*openauthv1.User{}
	for _, userRecord := range userRecords {
		users = append(users, transformUser(&userRecord))
	}

	var nextPageToken string
	if len(users) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(userRecords[limit].ID)
		users = users[:limit]
	}


	return &backendv1.ListUsersResponse{
		Users: users,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) UpdateUser(ctx context.Context, req *openauthv1.User) (*openauthv1.User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userId, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	organizationId, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	updates := queries.UpdateUserParams{
		ID: userId,
	}

	// Conditionally update organizationID
	if req.OrganizationId != "" {
		updates.OrganizationID = organizationId
	}

	// Conditionally update email addresses
	if req.UnverifiedEmail != "" {
		updates.UnverifiedEmail = &req.UnverifiedEmail
	}
	if req.VerifiedEmail != "" {
		updates.VerifiedEmail = &req.VerifiedEmail
	}

	// Conditionally update login method user IDs
	if req.GoogleUserId != "" {
		updates.GoogleUserID = &req.GoogleUserId
	}
	if req.MicrosoftUserId != "" {
		updates.MicrosoftUserID = &req.MicrosoftUserId
	}

	// Conditionall update password
	if req.PasswordBcrypt != "" {
		updates.PasswordBcrypt = &req.PasswordBcrypt
	}

	user, err := q.UpdateUser(ctx, updates)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformUser(&user), nil
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
	passwordBcryptBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	passwordBcrypt := string(passwordBcryptBytes)

	user, err := q.UpdateUserPassword(ctx, queries.UpdateUserPasswordParams{
		ID: userId,
		PasswordBcrypt: &passwordBcrypt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformUser(&user), nil
}


func transformUser(user *queries.User) *openauthv1.User {
	return &openauthv1.User{
		Id: user.ID.String(),
		OrganizationId: user.OrganizationID.String(),
		UnverifiedEmail: *user.UnverifiedEmail,
		VerifiedEmail: *user.VerifiedEmail,
		GoogleUserId: *user.GoogleUserID,
		MicrosoftUserId: *user.MicrosoftUserID,
		PasswordBcrypt: *user.PasswordBcrypt,
	}
}