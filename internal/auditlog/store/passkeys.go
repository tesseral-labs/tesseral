package store

import (
	"context"
	"encoding/pem"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetPasskey(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.Passkey, error) {
	qPasskey, err := queries.New(db).GetPasskey(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get passkey: %w", err)
	}

	return &auditlogv1.Passkey{
		Id:           idformat.Passkey.Format(qPasskey.ID),
		UserId:       idformat.User.Format(qPasskey.UserID),
		CreateTime:   timestamppb.New(*qPasskey.CreateTime),
		UpdateTime:   timestamppb.New(*qPasskey.UpdateTime),
		Disabled:     &qPasskey.Disabled,
		CredentialId: qPasskey.CredentialID,
		PublicKeyPkix: string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: qPasskey.PublicKey,
		})),
		Aaguid: qPasskey.Aaguid,
		RpId:   qPasskey.RpID,
	}, nil
}
