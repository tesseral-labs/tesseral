package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListUserInvites(ctx context.Context, req *frontendv1.ListUserInvitesRequest) (*frontendv1.ListUserInvitesResponse, error) {
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
	qUserInvites, err := q.ListUserInvites(ctx, queries.ListUserInvitesParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list user invites: %w", err)
	}

	var userInvites []*frontendv1.UserInvite
	for _, qUserInvite := range qUserInvites {
		userInvites = append(userInvites, parseUserInvite(qUserInvite))
	}

	var nextPageToken string
	if len(userInvites) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUserInvites[limit].ID)
		userInvites = userInvites[:limit]
	}

	return &frontendv1.ListUserInvitesResponse{
		UserInvites:   userInvites,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetUserInvite(ctx context.Context, req *frontendv1.GetUserInviteRequest) (*frontendv1.GetUserInviteResponse, error) {
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
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userInviteID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user invite not found", fmt.Errorf("get user invite: %w", err))
		}

		return nil, fmt.Errorf("get user invite: %w", err)
	}

	return &frontendv1.GetUserInviteResponse{UserInvite: parseUserInvite(qInvite)}, nil
}

func (s *Store) CreateUserInvite(ctx context.Context, req *frontendv1.CreateUserInviteRequest) (*frontendv1.CreateUserInviteResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// As a sanity measure, prevent invites to existing users with the same
	// email. Otherwise, those invites will remain present even after the user
	// is deleted, and so someone with that email address could redeem the
	// invite afterword.
	emailTaken, err := q.ExistsUserWithEmail(ctx, queries.ExistsUserWithEmailParams{
		OrganizationID: authn.OrganizationID(ctx),
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
		OrganizationID: authn.OrganizationID(ctx),
		Email:          req.UserInvite.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.CreateUserInviteResponse{UserInvite: parseUserInvite(qUserInvite)}, nil
}

func (s *Store) DeleteUserInvite(ctx context.Context, req *frontendv1.DeleteUserInviteRequest) (*frontendv1.DeleteUserInviteResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

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
		ID:             userInviteID,
		OrganizationID: authn.OrganizationID(ctx),
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

	return &frontendv1.DeleteUserInviteResponse{}, nil
}

func parseUserInvite(qUserInvite queries.UserInvite) *frontendv1.UserInvite {
	return &frontendv1.UserInvite{
		Id:         idformat.UserInvite.Format(qUserInvite.ID),
		CreateTime: timestamppb.New(*qUserInvite.CreateTime),
		UpdateTime: timestamppb.New(*qUserInvite.UpdateTime),
		Email:      qUserInvite.Email,
	}
}
