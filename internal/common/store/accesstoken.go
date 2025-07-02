package store

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/ujwt"
	"google.golang.org/protobuf/encoding/protojson"
)

const accessTokenDuration = time.Minute * 5

func (s *Store) IssueAccessToken(ctx context.Context, projectID uuid.UUID, refreshToken string) (string, error) {
	// this type exists to unify the datatypes we get from refresh tokens that
	// belong to sessions vs relayed sessions
	var qDetails struct {
		SessionID               uuid.UUID
		UserID                  uuid.UUID
		OrganizationID          uuid.UUID
		UserIsOwner             bool
		UserEmail               string
		UserDisplayName         *string
		UserProfilePictureUrl   *string
		OrganizationDisplayName string
		ImpersonatorUserID      *uuid.UUID
	}

	switch {
	case strings.HasPrefix(refreshToken, "tesseral_secret_session_refresh_token_"):
		slog.InfoContext(ctx, "refresh_session_token")

		refreshTokenUUID, err := idformat.SessionRefreshToken.Parse(refreshToken)
		if err != nil {
			return "", fmt.Errorf("parse refresh token: %w", err)
		}

		refreshTokenSHA := sha256.Sum256(refreshTokenUUID[:])
		qSessionDetails, err := s.q.GetSessionDetailsByRefreshTokenSHA256(ctx, queries.GetSessionDetailsByRefreshTokenSHA256Params{
			ProjectID:          projectID,
			RefreshTokenSha256: refreshTokenSHA[:],
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
			}

			return "", fmt.Errorf("get session details by refresh token sha256: %w", err)
		}

		qDetails.SessionID = qSessionDetails.SessionID
		qDetails.UserID = qSessionDetails.UserID
		qDetails.OrganizationID = qSessionDetails.OrganizationID
		qDetails.UserIsOwner = qSessionDetails.UserIsOwner
		qDetails.UserEmail = qSessionDetails.UserEmail
		qDetails.UserDisplayName = qSessionDetails.UserDisplayName
		qDetails.UserProfilePictureUrl = qSessionDetails.UserProfilePictureUrl
		qDetails.OrganizationDisplayName = qSessionDetails.OrganizationDisplayName
		qDetails.ImpersonatorUserID = qSessionDetails.ImpersonatorUserID
	case strings.HasPrefix(refreshToken, "tesseral_secret_relayed_session_refresh_token_"):
		slog.InfoContext(ctx, "refresh_relayed_session_token")

		relayedRefreshTokenUUID, err := idformat.RelayedSessionRefreshToken.Parse(refreshToken)
		if err != nil {
			return "", fmt.Errorf("parse refresh token: %w", err)
		}

		relayedRefreshTokenSHA := sha256.Sum256(relayedRefreshTokenUUID[:])
		qSessionDetails, err := s.q.GetSessionDetailsByRelayedSessionRefreshTokenSHA256(ctx, queries.GetSessionDetailsByRelayedSessionRefreshTokenSHA256Params{
			ProjectID:                 projectID,
			RelayedRefreshTokenSha256: relayedRefreshTokenSHA[:],
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
			}

			return "", fmt.Errorf("get session details by refresh token sha256: %w", err)
		}

		qDetails.SessionID = qSessionDetails.SessionID
		qDetails.UserID = qSessionDetails.UserID
		qDetails.OrganizationID = qSessionDetails.OrganizationID
		qDetails.UserIsOwner = qSessionDetails.UserIsOwner
		qDetails.UserEmail = qSessionDetails.UserEmail
		qDetails.UserDisplayName = qSessionDetails.UserDisplayName
		qDetails.UserProfilePictureUrl = qSessionDetails.UserProfilePictureUrl
		qDetails.OrganizationDisplayName = qSessionDetails.OrganizationDisplayName
		qDetails.ImpersonatorUserID = qSessionDetails.ImpersonatorUserID
	}

	issAndAud := fmt.Sprintf("https://%s.tesseral.app", strings.ReplaceAll(idformat.Project.Format(projectID), "_", "-"))
	now := time.Now()

	// Add details about the creator of the impersonation token to the session.
	//
	// We could in principle add this data using a LEFT JOIN, but the vast
	// majority of sessions are not impersonated, and so this branch is rarely
	// exercised.
	//
	// Plus sqlc does not at the time of writing correctly handle the types
	// associated with having two joins (one INNER, one LEFT) on users in the
	// GetSessionDetailsByRefreshTokenSHA256 query.
	var impersonator *commonv1.AccessTokenImpersonator
	if qDetails.ImpersonatorUserID != nil {
		qImpersonator, err := s.q.GetImpersonatorUserByID(ctx, *qDetails.ImpersonatorUserID)
		if err != nil {
			return "", fmt.Errorf("get impersonator user by id: %w", err)
		}

		impersonator = &commonv1.AccessTokenImpersonator{
			Email: qImpersonator.Email,
		}
	}

	var actions []string
	if qDetails.UserIsOwner {
		projectActions, err := s.q.GetProjectActions(ctx, projectID)
		if err != nil {
			return "", fmt.Errorf("get project actions: %w", err)
		}

		actions = projectActions
	} else {
		userActions, err := s.q.GetUserActions(ctx, qDetails.UserID)
		if err != nil {
			return "", fmt.Errorf("get user actions: %w", err)
		}

		actions = userActions
	}

	slices.Sort(actions)

	claims := &commonv1.AccessTokenData{
		Iss: issAndAud,
		Sub: idformat.User.Format(qDetails.UserID),
		Aud: issAndAud,
		Exp: float64(now.Add(accessTokenDuration).Unix()),
		Nbf: float64(now.Unix()),
		Iat: float64(now.Unix()),
		Session: &commonv1.AccessTokenSession{
			Id: idformat.Session.Format(qDetails.SessionID),
		},
		User: &commonv1.AccessTokenUser{
			Id:                idformat.User.Format(qDetails.UserID),
			Email:             qDetails.UserEmail,
			DisplayName:       derefOrEmpty(qDetails.UserDisplayName),
			ProfilePictureUrl: derefOrEmpty(qDetails.UserProfilePictureUrl),
		},
		Organization: &commonv1.AccessTokenOrganization{
			Id:          idformat.Organization.Format(qDetails.OrganizationID),
			DisplayName: qDetails.OrganizationDisplayName,
		},
		Actions:      actions,
		Impersonator: impersonator,
	}

	slog.InfoContext(ctx, "issue_access_token",
		"project_id", idformat.Project.Format(projectID),
		"claims", claims)

	// claims is a proto message, so we have to use protojson to encode it first
	encodedClaims, err := protojson.Marshal(claims)
	if err != nil {
		panic(fmt.Errorf("marshal claims: %w", err))
	}

	qSessionSigningKey, err := s.q.GetCurrentSessionSigningKeyByProjectID(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("get current session signing key by project id: %w", err)
	}

	sessionSigningKeyID := idformat.SessionSigningKey.Format(qSessionSigningKey.ID)
	slog.InfoContext(ctx, "sign_with_session_key", "session_signing_key_id", sessionSigningKeyID)

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      qSessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKMSKeyID,
	})
	if err != nil {
		return "", fmt.Errorf("decrypt session signing key ciphertext: %w", err)
	}

	priv, err := x509.ParseECPrivateKey(decryptRes.Plaintext)
	if err != nil {
		panic(fmt.Errorf("private key from bytes: %w", err))
	}

	if err := s.q.BumpSessionLastActiveTime(ctx, qDetails.SessionID); err != nil {
		return "", fmt.Errorf("bump session last active time: %w", err)
	}

	accessToken := ujwt.Sign(sessionSigningKeyID, priv, json.RawMessage(encodedClaims))
	return accessToken, nil
}
