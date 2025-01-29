package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

	if qIntermediateSession.EmailVerificationChallengeSha256 != nil && qIntermediateSession.EmailVerificationChallengeCompleted {
		return true, nil
	}

	return false, nil
}

func enforceOrganizationLoginEnabled(qOrganization queries.Organization) error {
	if qOrganization.LoginsDisabled {
		return apierror.NewPermissionDeniedError("login disabled", fmt.Errorf("organization login disabled"))
	}
	return nil
}

func enforceProjectLoginEnabled(qProject queries.Project) error {
	if qProject.LoginsDisabled {
		return apierror.NewPermissionDeniedError("login disabled", fmt.Errorf("project login disabled"))
	}
	return nil
}

func parseIntermediateSession(qIntermediateSession queries.IntermediateSession, emailVerified bool) *intermediatev1.IntermediateSession {
	var organizationID string
	if qIntermediateSession.OrganizationID != nil {
		organizationID = idformat.Organization.Format(*qIntermediateSession.OrganizationID)
	}

	return &intermediatev1.IntermediateSession{
		Id:                                   idformat.IntermediateSession.Format(qIntermediateSession.ID),
		ProjectId:                            idformat.Project.Format(qIntermediateSession.ProjectID),
		Email:                                derefOrEmpty(qIntermediateSession.Email),
		EmailVerified:                        emailVerified,
		EmailVerificationChallengeRegistered: qIntermediateSession.EmailVerificationChallengeSha256 != nil,
		GoogleUserId:                         derefOrEmpty(qIntermediateSession.GoogleUserID),
		GoogleHostedDomain:                   derefOrEmpty(qIntermediateSession.GoogleHostedDomain),
		MicrosoftUserId:                      derefOrEmpty(qIntermediateSession.MicrosoftUserID),
		MicrosoftTenantId:                    derefOrEmpty(qIntermediateSession.MicrosoftTenantID),
		PasswordVerified:                     qIntermediateSession.PasswordVerified,
		NewUserPasswordRegistered:            qIntermediateSession.NewUserPasswordBcrypt != nil,
		OrganizationId:                       organizationID,
	}
}
