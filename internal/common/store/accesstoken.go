package store

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	commonv1 "github.com/openauth/openauth/internal/common/gen/openauth/common/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/ujwt"
	"google.golang.org/protobuf/encoding/protojson"
)

const accessTokenDuration = time.Minute * 5

func (s *Store) IssueAccessToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenUUID, err := idformat.SessionRefreshToken.Parse(refreshToken)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %w", err)
	}

	refreshTokenSHA := sha256.Sum256(refreshTokenUUID[:])
	qDetails, err := s.q.GetSessionDetailsByRefreshTokenSHA256(ctx, refreshTokenSHA[:])
	if err != nil {
		return "", fmt.Errorf("get session details by refresh token sha256: %w", err)
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

	accessToken := ujwt.Sign(idformat.SessionSigningKey.Format(qSessionSigningKey.ID), priv, json.RawMessage(encodedClaims))
	return accessToken, nil
}
