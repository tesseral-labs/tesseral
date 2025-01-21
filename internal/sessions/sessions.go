package sessions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/ujwt"
)

type sessionClaims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Nbf int64  `json:"nbf"`
	Iat int64  `json:"iat"`

	Session      json.RawMessage `json:"session"`
	User         json.RawMessage `json:"user"`
	Organization json.RawMessage `json:"organization"`
	Project      json.RawMessage `json:"project"`
}

type Organization struct {
	ID                        string    `json:"id"`
	ProjectID                 string    `json:"projectId"`
	CreateTime                time.Time `json:"createTime"`
	UpdateTime                time.Time `json:"updateTime"`
	DisplayName               string    `json:"displayName"`
	LogInWithGoogleEnabled    bool      `json:"logInWithGoogleEnabled"`
	LogInWithMicrosoftEnabled bool      `json:"logInWithMicrosoftEnabled"`
	LogInWithPasswordEnabled  bool      `json:"logInWithPasswordEnabled"`
	SamlEnabled               bool      `json:"samlEnabled"`
}

type Project struct {
	ID                        string `json:"id"`
	CreateTime                time.Time
	UpdateTime                time.Time
	AuthDomain                string `json:"authDomain"`
	DisplayName               string `json:"displayName"`
	LogInWithGoogleEnabled    bool   `json:"logInWithGoogleEnabled"`
	LogInWithMicrosoftEnabled bool   `json:"logInWithMicrosoftEnabled"`
	LogInWithPasswordEnabled  bool   `json:"logInWithPasswordEnabled"`
}

type Session struct {
	ID         string `json:"id"`
	CreateTime time.Time
	ExpireTime time.Time
	UserID     string
	Revoked    bool
}

type User struct {
	ID              string `json:"id"`
	CreateTime      time.Time
	UpdateTime      time.Time
	Email           string
	GoogleUserID    string
	MicrosoftUserID string
}

func GetAccessToken(ctx context.Context, organization *Organization, project *Project, session *Session, user *User, privateKeyID uuid.UUID, privateKey *ecdsa.PrivateKey) (string, error) {
	now := time.Now()
	exp := now.Add(5 * time.Minute) // TODO(ucarion) parameterize

	organizationClaim, err := json.Marshal(organization)
	if err != nil {
		return "", fmt.Errorf("marshal organization claim: %w", err)
	}

	projectClaim, err := json.Marshal(project)
	if err != nil {
		return "", fmt.Errorf("marshal project claim: %w", err)
	}

	sessionClaim, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("marshal session claim: %w", err)
	}

	userClaim, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("marshal user claim: %w", err)
	}

	claims := sessionClaims{
		Iss: "TODO",
		Sub: user.ID,
		Aud: "TODO",
		Exp: exp.Unix(),
		Nbf: now.Unix(),
		Iat: now.Unix(),

		Session:      sessionClaim,
		User:         userClaim,
		Organization: organizationClaim,
		Project:      projectClaim,
	}

	accessToken := ujwt.Sign(idformat.SessionSigningKey.Format(privateKeyID), privateKey, claims)

	return accessToken, nil
}
