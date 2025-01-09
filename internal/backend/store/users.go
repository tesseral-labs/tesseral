package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
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
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	// authz
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
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
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	qUser, err := q.GetUser(ctx, queries.GetUserParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &backendv1.GetUserResponse{User: parseUser(qUser)}, nil
}

func parseUser(qUser queries.User) *backendv1.User {
	return &backendv1.User{
		Id:              idformat.User.Format(qUser.ID),
		OrganizationId:  idformat.Organization.Format(qUser.OrganizationID),
		Email:           qUser.Email,
		CreateTime:      timestamppb.New(*qUser.CreateTime),
		UpdateTime:      timestamppb.New(*qUser.UpdateTime),
		Owner:           &qUser.IsOwner,
		GoogleUserId:    derefOrEmpty(qUser.GoogleUserID),
		MicrosoftUserId: derefOrEmpty(qUser.MicrosoftUserID),
	}
}
