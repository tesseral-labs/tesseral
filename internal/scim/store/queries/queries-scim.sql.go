// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries-scim.sql

package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countUsers = `-- name: CountUsers :one
SELECT
    count(*)
FROM
    users
WHERE
    organization_id = $1
`

func (q *Queries) CountUsers(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countUsers, organizationID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAuditLogEvent = `-- name: CreateAuditLogEvent :one
INSERT INTO audit_log_events (id, project_id, organization_id, actor_scim_api_key_id, resource_type, resource_id, event_name, event_time, event_details)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, coalesce($9, '{}'::jsonb))
RETURNING
    id, project_id, organization_id, actor_user_id, actor_session_id, actor_api_key_id, actor_console_user_id, actor_console_session_id, actor_backend_api_key_id, actor_intermediate_session_id, resource_type, resource_id, event_name, event_time, event_details, actor_scim_api_key_id
`

type CreateAuditLogEventParams struct {
	ID                uuid.UUID
	ProjectID         uuid.UUID
	OrganizationID    *uuid.UUID
	ActorScimApiKeyID *uuid.UUID
	ResourceType      *AuditLogEventResourceType
	ResourceID        *uuid.UUID
	EventName         string
	EventTime         *time.Time
	EventDetails      interface{}
}

func (q *Queries) CreateAuditLogEvent(ctx context.Context, arg CreateAuditLogEventParams) (AuditLogEvent, error) {
	row := q.db.QueryRow(ctx, createAuditLogEvent,
		arg.ID,
		arg.ProjectID,
		arg.OrganizationID,
		arg.ActorScimApiKeyID,
		arg.ResourceType,
		arg.ResourceID,
		arg.EventName,
		arg.EventTime,
		arg.EventDetails,
	)
	var i AuditLogEvent
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.OrganizationID,
		&i.ActorUserID,
		&i.ActorSessionID,
		&i.ActorApiKeyID,
		&i.ActorConsoleUserID,
		&i.ActorConsoleSessionID,
		&i.ActorBackendApiKeyID,
		&i.ActorIntermediateSessionID,
		&i.ResourceType,
		&i.ResourceID,
		&i.EventName,
		&i.EventTime,
		&i.EventDetails,
		&i.ActorScimApiKeyID,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, is_owner)
    VALUES ($1, $2, $3, $4)
RETURNING
    id, organization_id, password_bcrypt, google_user_id, microsoft_user_id, email, create_time, update_time, is_owner, failed_password_attempts, password_lockout_expire_time, authenticator_app_secret_ciphertext, failed_authenticator_app_attempts, authenticator_app_lockout_expire_time, authenticator_app_recovery_code_sha256s, display_name, profile_picture_url, github_user_id
`

type CreateUserParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	IsOwner        bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.OrganizationID,
		arg.Email,
		arg.IsOwner,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.PasswordBcrypt,
		&i.GoogleUserID,
		&i.MicrosoftUserID,
		&i.Email,
		&i.CreateTime,
		&i.UpdateTime,
		&i.IsOwner,
		&i.FailedPasswordAttempts,
		&i.PasswordLockoutExpireTime,
		&i.AuthenticatorAppSecretCiphertext,
		&i.FailedAuthenticatorAppAttempts,
		&i.AuthenticatorAppLockoutExpireTime,
		&i.AuthenticatorAppRecoveryCodeSha256s,
		&i.DisplayName,
		&i.ProfilePictureUrl,
		&i.GithubUserID,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
    AND organization_id = $2
RETURNING
    id, organization_id, password_bcrypt, google_user_id, microsoft_user_id, email, create_time, update_time, is_owner, failed_password_attempts, password_lockout_expire_time, authenticator_app_secret_ciphertext, failed_authenticator_app_attempts, authenticator_app_lockout_expire_time, authenticator_app_recovery_code_sha256s, display_name, profile_picture_url, github_user_id
`

type DeleteUserParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

func (q *Queries) DeleteUser(ctx context.Context, arg DeleteUserParams) (User, error) {
	row := q.db.QueryRow(ctx, deleteUser, arg.ID, arg.OrganizationID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.PasswordBcrypt,
		&i.GoogleUserID,
		&i.MicrosoftUserID,
		&i.Email,
		&i.CreateTime,
		&i.UpdateTime,
		&i.IsOwner,
		&i.FailedPasswordAttempts,
		&i.PasswordLockoutExpireTime,
		&i.AuthenticatorAppSecretCiphertext,
		&i.FailedAuthenticatorAppAttempts,
		&i.AuthenticatorAppLockoutExpireTime,
		&i.AuthenticatorAppRecoveryCodeSha256s,
		&i.DisplayName,
		&i.ProfilePictureUrl,
		&i.GithubUserID,
	)
	return i, err
}

const getOrganizationDomains = `-- name: GetOrganizationDomains :many
SELECT
    DOMAIN
FROM
    organization_domains
WHERE
    organization_id = $1
`

func (q *Queries) GetOrganizationDomains(ctx context.Context, organizationID uuid.UUID) ([]string, error) {
	rows, err := q.db.Query(ctx, getOrganizationDomains, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return nil, err
		}
		items = append(items, domain)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSCIMAPIKeyByTokenSHA256 = `-- name: GetSCIMAPIKeyByTokenSHA256 :one
SELECT
    scim_api_keys.id, scim_api_keys.organization_id, scim_api_keys.secret_token_sha256, scim_api_keys.display_name, scim_api_keys.create_time, scim_api_keys.update_time
FROM
    scim_api_keys
    JOIN organizations ON scim_api_keys.organization_id = organizations.id
WHERE
    secret_token_sha256 = $1
    AND organizations.project_id = $2
`

type GetSCIMAPIKeyByTokenSHA256Params struct {
	SecretTokenSha256 []byte
	ProjectID         uuid.UUID
}

func (q *Queries) GetSCIMAPIKeyByTokenSHA256(ctx context.Context, arg GetSCIMAPIKeyByTokenSHA256Params) (ScimApiKey, error) {
	row := q.db.QueryRow(ctx, getSCIMAPIKeyByTokenSHA256, arg.SecretTokenSha256, arg.ProjectID)
	var i ScimApiKey
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.SecretTokenSha256,
		&i.DisplayName,
		&i.CreateTime,
		&i.UpdateTime,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id, organization_id, password_bcrypt, google_user_id, microsoft_user_id, email, create_time, update_time, is_owner, failed_password_attempts, password_lockout_expire_time, authenticator_app_secret_ciphertext, failed_authenticator_app_attempts, authenticator_app_lockout_expire_time, authenticator_app_recovery_code_sha256s, display_name, profile_picture_url, github_user_id
FROM
    users
WHERE
    organization_id = $1
    AND email = $2
`

type GetUserByEmailParams struct {
	OrganizationID uuid.UUID
	Email          string
}

func (q *Queries) GetUserByEmail(ctx context.Context, arg GetUserByEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, arg.OrganizationID, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.PasswordBcrypt,
		&i.GoogleUserID,
		&i.MicrosoftUserID,
		&i.Email,
		&i.CreateTime,
		&i.UpdateTime,
		&i.IsOwner,
		&i.FailedPasswordAttempts,
		&i.PasswordLockoutExpireTime,
		&i.AuthenticatorAppSecretCiphertext,
		&i.FailedAuthenticatorAppAttempts,
		&i.AuthenticatorAppLockoutExpireTime,
		&i.AuthenticatorAppRecoveryCodeSha256s,
		&i.DisplayName,
		&i.ProfilePictureUrl,
		&i.GithubUserID,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT
    id, organization_id, password_bcrypt, google_user_id, microsoft_user_id, email, create_time, update_time, is_owner, failed_password_attempts, password_lockout_expire_time, authenticator_app_secret_ciphertext, failed_authenticator_app_attempts, authenticator_app_lockout_expire_time, authenticator_app_recovery_code_sha256s, display_name, profile_picture_url, github_user_id
FROM
    users
WHERE
    organization_id = $1
    AND id = $2
`

type GetUserByIDParams struct {
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (q *Queries) GetUserByID(ctx context.Context, arg GetUserByIDParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, arg.OrganizationID, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.PasswordBcrypt,
		&i.GoogleUserID,
		&i.MicrosoftUserID,
		&i.Email,
		&i.CreateTime,
		&i.UpdateTime,
		&i.IsOwner,
		&i.FailedPasswordAttempts,
		&i.PasswordLockoutExpireTime,
		&i.AuthenticatorAppSecretCiphertext,
		&i.FailedAuthenticatorAppAttempts,
		&i.AuthenticatorAppLockoutExpireTime,
		&i.AuthenticatorAppRecoveryCodeSha256s,
		&i.DisplayName,
		&i.ProfilePictureUrl,
		&i.GithubUserID,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT
    id, organization_id, password_bcrypt, google_user_id, microsoft_user_id, email, create_time, update_time, is_owner, failed_password_attempts, password_lockout_expire_time, authenticator_app_secret_ciphertext, failed_authenticator_app_attempts, authenticator_app_lockout_expire_time, authenticator_app_recovery_code_sha256s, display_name, profile_picture_url, github_user_id
FROM
    users
WHERE
    organization_id = $1
ORDER BY
    id
LIMIT $2 OFFSET $3
`

type ListUsersParams struct {
	OrganizationID uuid.UUID
	Limit          int32
	Offset         int32
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.OrganizationID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.PasswordBcrypt,
			&i.GoogleUserID,
			&i.MicrosoftUserID,
			&i.Email,
			&i.CreateTime,
			&i.UpdateTime,
			&i.IsOwner,
			&i.FailedPasswordAttempts,
			&i.PasswordLockoutExpireTime,
			&i.AuthenticatorAppSecretCiphertext,
			&i.FailedAuthenticatorAppAttempts,
			&i.AuthenticatorAppLockoutExpireTime,
			&i.AuthenticatorAppRecoveryCodeSha256s,
			&i.DisplayName,
			&i.ProfilePictureUrl,
			&i.GithubUserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE
    users
SET
    email = $1
WHERE
    id = $2
    AND organization_id = $3
RETURNING
    id, organization_id, password_bcrypt, google_user_id, microsoft_user_id, email, create_time, update_time, is_owner, failed_password_attempts, password_lockout_expire_time, authenticator_app_secret_ciphertext, failed_authenticator_app_attempts, authenticator_app_lockout_expire_time, authenticator_app_recovery_code_sha256s, display_name, profile_picture_url, github_user_id
`

type UpdateUserParams struct {
	Email          string
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.Email, arg.ID, arg.OrganizationID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.PasswordBcrypt,
		&i.GoogleUserID,
		&i.MicrosoftUserID,
		&i.Email,
		&i.CreateTime,
		&i.UpdateTime,
		&i.IsOwner,
		&i.FailedPasswordAttempts,
		&i.PasswordLockoutExpireTime,
		&i.AuthenticatorAppSecretCiphertext,
		&i.FailedAuthenticatorAppAttempts,
		&i.AuthenticatorAppLockoutExpireTime,
		&i.AuthenticatorAppRecoveryCodeSha256s,
		&i.DisplayName,
		&i.ProfilePictureUrl,
		&i.GithubUserID,
	)
	return i, err
}
