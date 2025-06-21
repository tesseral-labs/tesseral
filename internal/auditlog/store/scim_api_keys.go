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

func (s *Store) GetSCIMAPIKey(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.SCIMAPIKey, error) {
	qSCIMAPIKey, err := queries.New(db).GetSCIMAPIKey(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	return &auditlogv1.SCIMAPIKey{
		Id:          idformat.SCIMAPIKey.Format(qSCIMAPIKey.ID),
		CreateTime:  timestamppb.New(*qSCIMAPIKey.CreateTime),
		UpdateTime:  timestamppb.New(*qSCIMAPIKey.UpdateTime),
		DisplayName: qSCIMAPIKey.DisplayName,
		Revoked:     qSCIMAPIKey.SecretTokenSha256 == nil,
	}, nil
}
