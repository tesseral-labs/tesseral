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

func (s *Store) GetSession(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.Session, error) {
	qSession, err := queries.New(db).GetSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user invite: %w", err)
	}

	var primaryAuthFactor auditlogv1.PrimaryAuthFactor
	switch qSession.PrimaryAuthFactor {
	case queries.PrimaryAuthFactorEmail:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_EMAIL
	case queries.PrimaryAuthFactorGoogle:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GOOGLE
	case queries.PrimaryAuthFactorMicrosoft:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_MICROSOFT
	case queries.PrimaryAuthFactorGithub:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GITHUB
	case queries.PrimaryAuthFactorSaml:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_SAML
	case queries.PrimaryAuthFactorOidc:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_OIDC
	default:
		primaryAuthFactor = auditlogv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_UNSPECIFIED
	}

	var impersonatorEmail string
	if qSession.ImpersonatorUserID != nil {
		qImpersonator, err := queries.New(db).GetUser(ctx, *qSession.ImpersonatorUserID)
		if err != nil {
			return nil, fmt.Errorf("get user: %w", err)
		}
		impersonatorEmail = qImpersonator.Email
	}

	return &auditlogv1.Session{
		Id:                idformat.Session.Format(qSession.ID),
		UserId:            idformat.User.Format(qSession.UserID),
		ExpireTime:        timestamppb.New(derefOrEmpty(qSession.ExpireTime)),
		LastActiveTime:    timestamppb.New(derefOrEmpty(qSession.LastActiveTime)),
		PrimaryAuthFactor: primaryAuthFactor,
		ImpersonatorEmail: impersonatorEmail,
	}, nil
}
