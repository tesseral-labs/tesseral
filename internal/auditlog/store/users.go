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

func (s *Store) GetUser(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.User, error) {
	qUser, err := queries.New(db).GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &auditlogv1.User{
		Id:                  idformat.User.Format(qUser.ID),
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
	}, nil
}
