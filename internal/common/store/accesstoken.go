package store

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/ujwt"
	"google.golang.org/protobuf/encoding/protojson"
)

const accessTokenDuration = time.Minute * 5

func (s *Store) IssueAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// this type exists to unify the datatypes we get from refresh tokens that
	// belong to sessions vs relayed sessions
	var qDetails struct {
		SessionID               uuid.UUID
		UserID                  uuid.UUID
		OrganizationID          uuid.UUID
		UserEmail               string
		OrganizationDisplayName string
		ImpersonatorUserID      *uuid.UUID
		ProjectID               uuid.UUID
	}

	switch {
	case strings.HasPrefix(refreshToken, "tesseral_secret_session_refresh_token_"):
		refreshTokenUUID, err := idformat.SessionRefreshToken.Parse(refreshToken)
		if err != nil {
			return "", fmt.Errorf("parse refresh token: %w", err)
		}

		refreshTokenSHA := sha256.Sum256(refreshTokenUUID[:])
		qSessionDetails, err := s.q.GetSessionDetailsByRefreshTokenSHA256(ctx, refreshTokenSHA[:])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
			}

			return "", fmt.Errorf("get session details by refresh token sha256: %w", err)
		}

		qDetails.SessionID = qSessionDetails.SessionID
		qDetails.UserID = qSessionDetails.UserID
		qDetails.OrganizationID = qSessionDetails.OrganizationID
		qDetails.UserEmail = qSessionDetails.UserEmail
		qDetails.OrganizationDisplayName = qSessionDetails.OrganizationDisplayName
		qDetails.ImpersonatorUserID = qSessionDetails.ImpersonatorUserID
		qDetails.ProjectID = qSessionDetails.ProjectID
	case strings.HasPrefix(refreshToken, "tesseral_secret_relayed_session_refresh_token_"):
		relayedRefreshTokenUUID, err := idformat.RelayedSessionRefreshToken.Parse(refreshToken)
		if err != nil {
			return "", fmt.Errorf("parse refresh token: %w", err)
		}

		relayedRefreshTokenSHA := sha256.Sum256(relayedRefreshTokenUUID[:])
		qSessionDetails, err := s.q.GetSessionDetailsByRelayedSessionRefreshTokenSHA256(ctx, relayedRefreshTokenSHA[:])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
			}

			return "", fmt.Errorf("get session details by refresh token sha256: %w", err)
		}

		qDetails.SessionID = qSessionDetails.SessionID
		qDetails.UserID = qSessionDetails.UserID
		qDetails.OrganizationID = qSessionDetails.OrganizationID
		qDetails.UserEmail = qSessionDetails.UserEmail
		qDetails.OrganizationDisplayName = qSessionDetails.OrganizationDisplayName
		qDetails.ImpersonatorUserID = qSessionDetails.ImpersonatorUserID
		qDetails.ProjectID = qSessionDetails.ProjectID
	}

	issAndAud := fmt.Sprintf("https://%s.tesseral.app", strings.ReplaceAll(idformat.Project.Format(qDetails.ProjectID), "_", "-"))
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
			Id:             idformat.User.Format(qDetails.UserID),
			OrganizationId: idformat.Organization.Format(qDetails.OrganizationID),
			Email:          qDetails.UserEmail,
		},
		Organization: &commonv1.AccessTokenOrganization{
			Id:          idformat.Organization.Format(qDetails.OrganizationID),
			DisplayName: qDetails.OrganizationDisplayName,
		},
		Impersonator: impersonator,
	}

	// claims is a proto message, so we have to use protojson to encode it first
	encodedClaims, err := protojson.Marshal(claims)
	if err != nil {
		panic(fmt.Errorf("marshal claims: %w", err))
	}

	qSessionSigningKey, err := s.q.GetCurrentSessionSigningKeyByProjectID(ctx, qDetails.ProjectID)
	if err != nil {
		return "", fmt.Errorf("get current session signing key by project id: %w", err)
	}

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

	accessToken := ujwt.Sign(idformat.SessionSigningKey.Format(qSessionSigningKey.ID), priv, json.RawMessage(encodedClaims))
	return accessToken, nil
}
