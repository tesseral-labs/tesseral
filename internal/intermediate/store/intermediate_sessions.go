package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) getIntermediateSessionEmailVerified(ctx context.Context, q *queries.Queries, id uuid.UUID) (bool, error) {
	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("get intermediate session by id: %w", err)
	}

	// An email may be verified on account of a previously authenticated
	// google or microsoft user ID-to-email pair. The google or microsoft
	// user ID on an intermediate session is always authentic.
	if qIntermediateSession.GoogleUserID != nil {
		qVerifiedGoogleUserID, err := q.GetEmailVerifiedByGoogleUserID(ctx, queries.GetEmailVerifiedByGoogleUserIDParams{
			ProjectID:    authn.ProjectID(ctx),
			Email:        *qIntermediateSession.Email,
			GoogleUserID: qIntermediateSession.GoogleUserID,
		})
		if err != nil {
			return false, fmt.Errorf("get email verified by google user id: %w", err)
		}

		if qVerifiedGoogleUserID {
			return true, nil
		}
	}
	if qIntermediateSession.MicrosoftUserID != nil {
		qVerifiedMicrosoftUserID, err := q.GetEmailVerifiedByMicrosoftUserID(ctx, queries.GetEmailVerifiedByMicrosoftUserIDParams{
			ProjectID:       authn.ProjectID(ctx),
			Email:           *qIntermediateSession.Email,
			MicrosoftUserID: qIntermediateSession.MicrosoftUserID,
		})
		if err != nil {
			return false, fmt.Errorf("get email verified by microsoft user id: %w", err)
		}

		if qVerifiedMicrosoftUserID {
			return true, nil
		}
	}

	// If there's a successful email verification challenge associated with
	// this intermediate session, then the email is verified.
	qVerifiedEmailVerificationChallenge, err := q.GetEmailVerifiedByEmailVerificationChallenge(ctx, qIntermediateSession.ID)
	if err != nil {
		return false, fmt.Errorf("get email verified by email verification challenge: %w", err)
	}

	if qVerifiedEmailVerificationChallenge {
		return true, nil
	}
	return false, nil
}

func (s *Store) verifyOrganization(ctx context.Context, q queries.Queries) (*queries.Organization, error) {
	// intermediateSession := authn.IntermediateSession(ctx)
	// orgID, err := idformat.Organization.Parse(intermediateSession.OrganizationID)
	// if err != nil {
	// 	return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	// }

	// qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
	// 	ProjectID: authn.ProjectID(ctx),
	// 	ID:        orgID,
	// })
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
	// 	}

	// 	return nil, fmt.Errorf("get organization by id: %w", err)
	// }

	// if qOrg.LoginDisabled != nil && *qOrg.LoginDisabled {
	// 	return nil, apierror.NewPermissionDeniedError("login disabled", fmt.Errorf("organization login disabled"))
	// }

	// return &qOrg, nil
}

func (s *Store) verifyProject(ctx context.Context, q queries.Queries) (*queries.Project, error) {
	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if qProject.LoginDisabled != nil && *qProject.LoginDisabled {
		return nil, apierror.NewPermissionDeniedError("login disabled", fmt.Errorf("project login disabled"))
	}

	return &qProject, nil
}

func parseIntermediateSession(qIntermediateSession queries.IntermediateSession, emailVerified bool) *intermediatev1.IntermediateSession {
	return &intermediatev1.IntermediateSession{
		Id:                 idformat.IntermediateSession.Format(qIntermediateSession.ID),
		ProjectId:          idformat.Project.Format(qIntermediateSession.ProjectID),
		Email:              derefOrEmpty(qIntermediateSession.Email),
		EmailVerified:      emailVerified,
		GoogleUserId:       derefOrEmpty(qIntermediateSession.GoogleUserID),
		GoogleHostedDomain: derefOrEmpty(qIntermediateSession.GoogleHostedDomain),
		MicrosoftUserId:    derefOrEmpty(qIntermediateSession.MicrosoftUserID),
		MicrosoftTenantId:  derefOrEmpty(qIntermediateSession.MicrosoftTenantID),
		PasswordVerified:   derefOrEmpty(qIntermediateSession.PasswordVerified),
	}
}
