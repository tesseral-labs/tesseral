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

func (s *Store) GetAPIKey(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.APIKey, error) {
	qAPIKey, err := queries.New(db).GetAPIKey(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get api key: %w", err)
	}

	return &auditlogv1.APIKey{
		Id:                idformat.APIKey.Format(qAPIKey.ID),
		CreateTime:        timestamppb.New(*qAPIKey.CreateTime),
		UpdateTime:        timestamppb.New(*qAPIKey.UpdateTime),
		ExpireTime:        timestampOrNil(qAPIKey.ExpireTime),
		DisplayName:       qAPIKey.DisplayName,
		SecretTokenSuffix: derefOrEmpty(qAPIKey.SecretTokenSuffix),
		Revoked:           qAPIKey.SecretTokenSha256 == nil,
	}, nil
}
