package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/bcryptcost"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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
		return nil, apierror.NewPasswordCompromisedError("password is compromised", errors.New("password is compromised"))
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

	user := parseUser(qUpdatedUser)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.update",
		EventDetails: map[string]any{
			"user":         user,
			"previousUser": parseUser(qUser),
		},
		ResourceType: queries.AuditLogEventResourceTypeUser,
		ResourceID:   &qUser.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	// Commit the transaction.
	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// send sync user event
	if err := s.sendSyncUserEvent(ctx, qUpdatedUser); err != nil {
		return nil, fmt.Errorf("send sync user event: %w", err)
	}

	return &frontendv1.UpdateUserResponse{User: user}, nil
}

func (s *Store) DeleteUser(ctx context.Context, req *frontendv1.DeleteUserRequest) (*frontendv1.DeleteUserResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	if userID == authn.UserID(ctx) {
		return nil, apierror.NewFailedPreconditionError("cannot delete self", errors.New("cannot delete self"))
	}

	// Fetch the user to ensure it exists and belongs to the organization.
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

	if err := q.DeleteUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.delete",
		EventDetails: map[string]any{
			"user": parseUser(qUser),
		},
		ResourceType: queries.AuditLogEventResourceTypeUser,
		ResourceID:   &qUser.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// send sync user event
	if err := s.sendSyncUserEvent(ctx, qUser); err != nil {
		return nil, fmt.Errorf("send sync user event: %w", err)
	}

	return &frontendv1.DeleteUserResponse{}, nil
}

func (s *Store) sendSyncUserEvent(ctx context.Context, qUser queries.User) error {
	qProjectWebhookSettings, err := s.q.GetProjectWebhookSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		// We want to ignore this error if the project does not have webhook settings
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("get project by id: %w", err)
	}

	message, err := s.svixClient.Message.Create(ctx, qProjectWebhookSettings.AppID, models.MessageIn{
		EventType: "sync.user",
		Payload: map[string]interface{}{
			"type":   "sync.user",
			"userId": idformat.User.Format(qUser.ID),
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	slog.InfoContext(ctx, "svix_message_created", "message_id", message.Id, "event_type", message.EventType, "user_id", idformat.User.Format(qUser.ID))

	return nil
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
		GithubUserId:        derefOrEmpty(qUser.GithubUserID),
		HasAuthenticatorApp: qUser.AuthenticatorAppSecretCiphertext != nil,
		DisplayName:         qUser.DisplayName,
		ProfilePictureUrl:   qUser.ProfilePictureUrl,
	}
}
