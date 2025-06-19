package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetUserInvite(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.UserInvite, error) {
	qUserInvite, err := queries.New(db).GetUserInvite(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user invite: %w", err)
	}

	return &auditlogv1.UserInvite{
		Id:         idformat.UserInvite.Format(qUserInvite.ID),
		CreateTime: timestamppb.New(*qUserInvite.CreateTime),
		UpdateTime: timestamppb.New(*qUserInvite.UpdateTime),
		Email:      qUserInvite.Email,
		Owner:      qUserInvite.IsOwner,
	}, nil
}
