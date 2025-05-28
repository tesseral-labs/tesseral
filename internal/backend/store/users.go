package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/common/auditlog"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListUsers(ctx context.Context, req *backendv1.ListUsersRequest) (*backendv1.ListUsersResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// authz
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qUsers, err := q.ListUsers(ctx, queries.ListUsersParams{
		OrganizationID: orgID,
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	var users []*backendv1.User
	for _, qUser := range qUsers {
		users = append(users, parseUser(qUser))
	}

	var nextPageToken string
	if len(users) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUsers[limit].ID)
		users = users[:limit]
	}

	return &backendv1.ListUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetUser(ctx context.Context, req *backendv1.GetUserRequest) (*backendv1.GetUserResponse, error) {
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
		ProjectID: authn.ProjectID(ctx),
		ID:        userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get user: %w", err)
	}

	return &backendv1.GetUserResponse{User: parseUser(qUser)}, nil
}

func (s *Store) CreateUser(ctx context.Context, req *backendv1.CreateUserRequest) (*backendv1.CreateUserResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.User.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	projectID := authn.ProjectID(ctx)
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: projectID,
		ID:        orgID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization: %w", err))
		}
		return nil, fmt.Errorf("get organization: %w", err)
	}

	qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		Email:           req.User.Email,
		IsOwner:         req.User.GetOwner(),
		GoogleUserID:    req.User.GoogleUserId,
		MicrosoftUserID: req.User.MicrosoftUserId,
		GithubUserID:    req.User.GithubUserId,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// send create user event
	event, err := auditlog.NewCreateUserEvent(auditlog.CreateUserEventData{
		ProjectID:      projectID,
		OrganizationID: qUser.OrganizationID,
		User: auditlog.UserData{
			ID:    qUser.ID,
			Email: qUser.Email,
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "error", err)
	} else if _, err := s.common.CreateAuditLogEvent(ctx, event); err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "event", event, "error", err)
	}

	// send sync user event
	if err := s.sendSyncUserEvent(ctx, qUser); err != nil {
		return nil, fmt.Errorf("send sync user event: %w", err)
	}

	return &backendv1.CreateUserResponse{User: parseUser(qUser)}, nil
}

func (s *Store) UpdateUser(ctx context.Context, req *backendv1.UpdateUserRequest) (*backendv1.UpdateUserResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	projectID := authn.ProjectID(ctx)
	qUser, err := q.GetUser(ctx, queries.GetUserParams{
		ProjectID: projectID,
		ID:        userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user: %w", err))
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	var updates queries.UpdateUserParams
	var updatesEventData auditlog.UserUpdate
	updates.ID = userID

	updates.Email = qUser.Email
	if req.User.Email != "" {
		updates.Email = req.User.Email
		updatesEventData.Email = req.User.Email
	}

	updates.IsOwner = qUser.IsOwner
	if req.User.Owner != nil {
		updates.IsOwner = *req.User.Owner
		updatesEventData.IsOwner = req.User.Owner
	}

	updates.GoogleUserID = qUser.GoogleUserID
	if req.User.GoogleUserId != nil {
		// if value was actively set to empty string, then reset value in
		// database to null
		updates.GoogleUserID = refOrNil(*req.User.GoogleUserId)
		updatesEventData.GoogleUserID = *req.User.GoogleUserId
	}

	updates.MicrosoftUserID = qUser.MicrosoftUserID
	if req.User.MicrosoftUserId != nil {
		updates.MicrosoftUserID = refOrNil(*req.User.MicrosoftUserId)
		updatesEventData.MicrosoftUserID = *req.User.MicrosoftUserId
	}

	updates.GithubUserID = qUser.GithubUserID
	if req.User.GithubUserId != nil {
		updates.GithubUserID = refOrNil(*req.User.GithubUserId)
		updatesEventData.GithubUserID = *req.User.GithubUserId
	}

	updates.DisplayName = qUser.DisplayName
	if req.User.DisplayName != nil {
		updates.DisplayName = refOrNil(*req.User.DisplayName)
		updatesEventData.DisplayName = *req.User.DisplayName
	}

	updates.ProfilePictureUrl = qUser.ProfilePictureUrl
	if req.User.ProfilePictureUrl != nil {
		updates.ProfilePictureUrl = refOrNil(*req.User.ProfilePictureUrl)
		updatesEventData.ProfilePictureURL = *req.User.ProfilePictureUrl
	}

	qUpdatedUser, err := q.UpdateUser(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// send update user event
	event, err := auditlog.NewUpdateUserEvent(auditlog.UpdateUserEventData{
		ProjectID:      projectID,
		OrganizationID: qUser.OrganizationID,
		User: auditlog.UserData{
			ID:    qUser.ID,
			Email: qUser.Email,
		},
		Update: updatesEventData,
	})
	if err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "error", err)
	} else if _, err := s.common.CreateAuditLogEvent(ctx, event); err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "event", event, "error", err)
	}

	// send sync user event
	if err := s.sendSyncUserEvent(ctx, qUpdatedUser); err != nil {
		return nil, fmt.Errorf("send sync user event: %w", err)
	}

	return &backendv1.UpdateUserResponse{User: parseUser(qUpdatedUser)}, nil
}

func (s *Store) DeleteUser(ctx context.Context, req *backendv1.DeleteUserRequest) (*backendv1.DeleteUserResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	qUser, err := q.GetUser(ctx, queries.GetUserParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user: %w", err))
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	err = q.DeleteUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// send sync user event
	if err := s.sendSyncUserEvent(ctx, qUser); err != nil {
		return nil, fmt.Errorf("send sync user event: %w", err)
	}

	return &backendv1.DeleteUserResponse{}, nil
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

	if _, err := s.svixClient.Message.Create(ctx, qProjectWebhookSettings.AppID, models.MessageIn{
		EventType: "sync.user",
		Payload: map[string]interface{}{
			"type":   "sync.user",
			"userId": idformat.User.Format(qUser.ID),
		},
	}, nil); err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}

func parseUser(qUser queries.User) *backendv1.User {
	return &backendv1.User{
		Id:                  idformat.User.Format(qUser.ID),
		OrganizationId:      idformat.Organization.Format(qUser.OrganizationID),
		Email:               qUser.Email,
		CreateTime:          timestamppb.New(*qUser.CreateTime),
		UpdateTime:          timestamppb.New(*qUser.UpdateTime),
		Owner:               &qUser.IsOwner,
		GoogleUserId:        qUser.GoogleUserID,
		MicrosoftUserId:     qUser.MicrosoftUserID,
		GithubUserId:        qUser.GithubUserID,
		HasAuthenticatorApp: qUser.AuthenticatorAppSecretCiphertext != nil,
		DisplayName:         qUser.DisplayName,
		ProfilePictureUrl:   qUser.ProfilePictureUrl,
	}
}
