package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	frontendv1 "github.com/openauth/openauth/internal/gen/frontend/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
	"github.com/openauth/openauth/internal/ujwt"
	"google.golang.org/protobuf/encoding/protojson"
)

type sessionClaims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Nbf int64  `json:"nbf"`
	Iat int64  `json:"iat"`
	Jti string `json:"jti"`

	Session      json.RawMessage `json:"session"`
	User         json.RawMessage `json:"user"`
	Organization json.RawMessage `json:"organization"`
}

func (s *Store) GetAccessToken(ctx context.Context, req *frontendv1.GetAccessTokenRequest) (*frontendv1.GetAccessTokenResponse, error) {
	// TODO(ucarion): this endpoint will also look at + update state related to
	// latest activity; calling GetAccessToken is precisely what we define
	// "activity" to be
	qSession, qUser, qOrganization, qSessionSigningKey, err := s.getAccessTokenSessionDetails(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("get access token session details: %w", err)
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      qSessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt private key cipher texts: %w", err)
	}

	priv, err := openauthecdsa.PrivateKeyFromBytes(decryptRes.Value)
	if err != nil {
		return nil, fmt.Errorf("unmarshal private key: %w", err)
	}

	now := time.Now()
	exp := now.Add(5 * time.Minute) // TODO(ucarion) parameterize

	sessionClaim, err := protojson.Marshal(parseSession(qSession))
	if err != nil {
		return nil, fmt.Errorf("marshal session claim: %w", err)
	}

	userClaim, err := protojson.Marshal(parseUser(qUser))
	if err != nil {
		return nil, fmt.Errorf("marshal user claim: %w", err)
	}

	organizationClaim, err := protojson.Marshal(parseOrganization(*qOrganization))
	if err != nil {
		return nil, fmt.Errorf("marshal organization claim: %w", err)
	}

	claims := sessionClaims{
		Iss: "TODO",
		Sub: idformat.User.Format(qUser.ID),
		Aud: "TODO",
		Exp: exp.Unix(),
		Nbf: now.Unix(),
		Iat: now.Unix(),
		Jti: "TODO",

		Session:      sessionClaim,
		User:         userClaim,
		Organization: organizationClaim,
	}

	accessToken := ujwt.Sign(idformat.SessionSigningKey.Format(qSessionSigningKey.ID), priv, claims)
	return &frontendv1.GetAccessTokenResponse{AccessToken: accessToken}, nil
}

// getAccessTokenSessionDetails gets details on the session and its encrypted
// signing private key given a refreshToken.
//
// Conceptually, this exists to do database operations for GetAccessToken that
// come before calling out to AWS KMS.
func (s *Store) getAccessTokenSessionDetails(ctx context.Context, refreshToken string) (*queries.Session, *queries.User, *queries.Organization, *queries.SessionSigningKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rollback()

	refreshTokenBytes, err := idformat.SessionRefreshToken.Parse(refreshToken)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parse refresh token: %w", err)
	}

	refreshTokenSHA := sha256.Sum256(refreshTokenBytes[:])
	qSessionDetails, err := q.GetSessionDetailsByRefreshTokenSHA256(ctx, refreshTokenSHA[:])
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get session by refresh token sha256: %w", err)
	}

	qSessionSigningKey, err := q.GetCurrentSessionKeyByProjectID(ctx, qSessionDetails.ProjectID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get current session key by project id: %w", err)
	}

	qSession, err := q.GetSessionByID(ctx, qSessionDetails.SessionID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get session by id: %w", err)
	}

	qUser, err := q.GetUserByID(ctx, qSessionDetails.UserID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get user by id: %w", err)
	}

	qOrganization, err := q.GetOrganizationByID(ctx, qSessionDetails.OrganizationID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get organization by id: %w", err)
	}

	return &qSession, &qUser, &qOrganization, &qSessionSigningKey, nil
}
