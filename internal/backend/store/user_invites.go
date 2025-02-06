package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListUserInvites(ctx context.Context, req *backendv1.ListUserInvitesRequest) (*backendv1.ListUserInvitesResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qUserInvites, err := q.ListUserInvites(ctx, queries.ListUserInvitesParams{
		OrganizationID: orgID,
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list user invites: %w", err)
	}

	var userInvites []*backendv1.UserInvite
	for _, qUserInvite := range qUserInvites {
		userInvites = append(userInvites, parseUserInvite(qUserInvite))
	}

	var nextPageToken string
	if len(userInvites) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUserInvites[limit].ID)
		userInvites = userInvites[:limit]
	}

	return &backendv1.ListUserInvitesResponse{
		UserInvites:   userInvites,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetUserInvite(ctx context.Context, req *backendv1.GetUserInviteRequest) (*backendv1.GetUserInviteResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userInviteID, err := idformat.UserInvite.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse user invite id: %w", err)
	}

	qInvite, err := q.GetUserInvite(ctx, queries.GetUserInviteParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        userInviteID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user invite not found", fmt.Errorf("get user invite: %w", err))
		}

		return nil, fmt.Errorf("get user invite: %w", err)
	}

	return &backendv1.GetUserInviteResponse{UserInvite: parseUserInvite(qInvite)}, nil
}

func (s *Store) CreateUserInvite(ctx context.Context, req *backendv1.CreateUserInviteRequest) (*backendv1.CreateUserInviteResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.UserInvite.OrganizationId)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	// See note in CreateUserInvite in frontend/store/user_invites.go
	emailTaken, err := q.ExistsUserWithEmailInOrganization(ctx, queries.ExistsUserWithEmailInOrganizationParams{
		OrganizationID: orgID,
		Email:          req.UserInvite.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("exists user with email: %w", err)
	}

	if emailTaken {
		return nil, apierror.NewFailedPreconditionError("a user with that email already exists", nil)
	}

	qUserInvite, err := q.CreateUserInvite(ctx, queries.CreateUserInviteParams{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          req.UserInvite.Email,
		IsOwner:        req.UserInvite.IsOwner,
	})
	if err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreateUserInviteResponse{UserInvite: parseUserInvite(qUserInvite)}, nil
}

func (s *Store) DeleteUserInvite(ctx context.Context, req *backendv1.DeleteUserInviteRequest) (*backendv1.DeleteUserInviteResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userInviteID, err := idformat.UserInvite.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse user invite id: %w", err)
	}

	if _, err := q.GetUserInvite(ctx, queries.GetUserInviteParams{
		ID:        userInviteID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user invite not found", fmt.Errorf("get user invite: %w", err))
		}

		return nil, fmt.Errorf("get user invite: %w", err)
	}

	if err := q.DeleteUserInvite(ctx, userInviteID); err != nil {
		return nil, fmt.Errorf("delete user invite: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteUserInviteResponse{}, nil
}

func parseUserInvite(qUserInvite queries.UserInvite) *backendv1.UserInvite {
	return &backendv1.UserInvite{
		Id:             idformat.UserInvite.Format(qUserInvite.ID),
		OrganizationId: idformat.Organization.Format(qUserInvite.OrganizationID),
		CreateTime:     timestamppb.New(*qUserInvite.CreateTime),
		UpdateTime:     timestamppb.New(*qUserInvite.UpdateTime),
		Email:          qUserInvite.Email,
		IsOwner:        qUserInvite.IsOwner,
	}
}
