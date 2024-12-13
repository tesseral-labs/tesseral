package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/googleoauth"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetGoogleOAuthRedirectURL(ctx context.Context, req *intermediatev1.GetGoogleOAuthRedirectURLRequest) (*intermediatev1.GetGoogleOAuthRedirectURLResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, projectid.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %v", err)
	}

	if qProject.GoogleOauthClientID == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("google oauth not configured"))
	}

	token := uuid.New()
	tokenSha256 := sha256.Sum256(token[:])
	expiresAt := time.Now().Add(155 * time.Minute)

	// Since this is the entrypoint for the google oauth flow, we create the intermediate session here
	intermediateSession, err := q.CreateIntermediateSession(ctx, queries.CreateIntermediateSessionParams{
		ID:          uuid.New(),
		ProjectID:   projectid.ProjectID(ctx),
		ExpireTime:  &expiresAt,
		TokenSha256: tokenSha256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create intermediate session: %v", err)
	}

	state := uuid.NewString()
	stateSHA := sha256.Sum256([]byte(state))
	if _, err := q.UpdateIntermediateSessionGoogleOAuthStateSHA256(ctx, queries.UpdateIntermediateSessionGoogleOAuthStateSHA256Params{
		ID:                     intermediateSession.ID,
		GoogleOauthStateSha256: stateSHA[:],
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session google oauth state: %v", err)
	}

	url := googleoauth.GetAuthorizeURL(&googleoauth.GetAuthorizeURLRequest{
		GoogleOAuthClientID: *qProject.GoogleOauthClientID,
		RedirectURI:         fmt.Sprintf("%s/google-oauth-callback", s.uiUrl),
		State:               state,
	})

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %v", err)
	}

	return &intermediatev1.GetGoogleOAuthRedirectURLResponse{
		IntermediateSessionToken: idformat.IntermediateSessionToken.Format(token),
		Url:                      url,
	}, nil
}

func (s *Store) RedeemGoogleOAuthCode(ctx context.Context, req *intermediatev1.RedeemGoogleOAuthCodeRequest) (*intermediatev1.RedeemGoogleOAuthCodeResponse, error) {
	qProject, qIntermediateSession, err := s.getProjectAndIntermediateSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project and intermediate session: %v", err)
	}

	if qProject.GoogleOauthClientID == nil || qProject.GoogleOauthClientSecretCiphertext == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("google oauth not configured"))
	}

	stateSHA := sha256.Sum256([]byte(req.State))
	if !bytes.Equal(qIntermediateSession.GoogleOauthStateSha256, stateSHA[:]) {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad state parameter"))
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      qProject.GoogleOauthClientSecretCiphertext,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.googleOAuthClientSecretsKMSKeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt google oauth client secret: %v", err)
	}

	// start new tx now that all kms and oauth http i/o is done
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	redeemRes, err := s.googleOAuthClient.RedeemCode(ctx, &googleoauth.RedeemCodeRequest{
		GoogleOAuthClientID:     *qProject.GoogleOauthClientID,
		GoogleOAuthClientSecret: string(decryptRes.Value),
		RedirectURI:             fmt.Sprintf("%s/google-oauth-callback", s.uiUrl),
		Code:                    req.Code,
	})
	if err != nil {
		return nil, fmt.Errorf("redeem google oauth code: %v", err)
	}

	// todo what if the intermediate session already has an email
	// todo send an email verification challenge?
	// todo what if redeem doesn't come back with a google hosted domain?
	if _, err := q.UpdateIntermediateSessionGoogleDetails(ctx, queries.UpdateIntermediateSessionGoogleDetailsParams{
		ID:                 authn.IntermediateSessionID(ctx),
		Email:              &redeemRes.Email,
		GoogleUserID:       &redeemRes.GoogleUserID,
		GoogleHostedDomain: &redeemRes.GoogleHostedDomain,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session google details: %v", err)
	}

	shouldVerifyEmail, err := s.shouldVerifyGoogleEmail(ctx, redeemRes.Email, redeemRes.GoogleUserID)
	if err != nil {
		return nil, fmt.Errorf("should verify google email: %v", err)
	}

	var evcID uuid.UUID

	if shouldVerifyEmail {
		// Create a new secret token for the challenge
		secretToken, err := generateSecretToken()
		if err != nil {
			return nil, err
		}

		secretTokenSha256 := sha256.Sum256([]byte(secretToken))
		expiresAt := time.Now().Add(15 * time.Minute)

		evc, err := q.CreateEmailVerificationChallenge(ctx, queries.CreateEmailVerificationChallengeParams{
			ID:                    uuid.New(),
			ChallengeSha256:       secretTokenSha256[:],
			ExpireTime:            &expiresAt,
			IntermediateSessionID: authn.IntermediateSessionID(ctx),
			ProjectID:             authn.ProjectID(ctx),
		})
		if err != nil {
			return nil, fmt.Errorf("create email verification challenge: %v", err)
		}

		evcID = evc.ID

		// TODO: Send the secret token to the user's email address
		slog.InfoContext(ctx, "RedeemGoogleOAuthCode", "challenge", secretToken)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	response := &intermediatev1.RedeemGoogleOAuthCodeResponse{
		ShouldVerifyEmail: shouldVerifyEmail,
	}

	if evcID != uuid.Nil {
		response.EmailVerificationChallengeId = idformat.EmailVerificationChallenge.Format(evcID)
	}

	return response, nil
}

func (s *Store) getProjectAndIntermediateSession(ctx context.Context) (*queries.Project, *queries.IntermediateSession, error) {
	// this function exists to avoid doing non-database i/o operations mid-tx
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf("get project by id: %v", err)
	}

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf("get intermediate session by id: %v", err)
	}

	return &qProject, &qIntermediateSession, nil
}

func (s *Store) shouldVerifyGoogleEmail(ctx context.Context, email string, googleUserID string) (bool, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return false, err
	}

	// Check if the email is already verified
	verified, err := q.IsGoogleEmailVerified(ctx, queries.IsGoogleEmailVerifiedParams{
		Email:        email,
		GoogleUserID: &googleUserID,
		ProjectID:    authn.ProjectID(ctx),
	})
	if err != nil {
		return false, err
	}

	return !verified, nil
}
