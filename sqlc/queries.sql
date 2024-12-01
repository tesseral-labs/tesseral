-- name: CreateEmailVerificationChallenge :one
INSERT INTO email_verification_challenges (id, intermediate_session_id, project_id, email, challenge_sha256, expire_time, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateIntermediateSession :one
INSERT INTO intermediate_sessions (id, project_id, unverified_email, verified_email, expire_time, token, token_sha256)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: CreateIntermediateSessionSigningKey :one
INSERT INTO intermediate_session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, google_hosted_domain, microsoft_tenant_id, override_log_in_with_google_enabled, override_log_in_with_microsoft_enabled, override_log_in_with_password_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateProject :one
INSERT INTO projects (id, organization_id, log_in_with_password_enabled, log_in_with_google_enabled, log_in_with_microsoft_enabled, google_oauth_client_id, google_oauth_client_secret, microsoft_oauth_client_id, microsoft_oauth_client_secret)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    *;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, create_time, expire_time, revoked)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateSessionSigningKey :one
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, unverified_email, verified_email, password_bcrypt, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: CreateGoogleUser :one
INSERT INTO users (id, organization_id, google_user_id, verified_email)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateMicrosoftUser :one
INSERT INTO users (id, organization_id, microsoft_user_id, verified_email)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateUnverifiedUser :one
INSERT INTO users (id, organization_id, unverified_email)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: CreateVerifiedEmail :one
INSERT INTO verified_emails (id, project_id, email, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetEmailVerificationChallenge :one
SELECT
    *
FROM
    email_verification_challenges
WHERE
    project_id = $1
    AND intermediate_session_id = $2
    AND challenge_sha256 = $3
    AND expire_time > $4
    AND (email = $5
        OR google_user_id = $6
        OR microsoft_user_id = $7)
LIMIT 1;

-- name: GetIntermediateSessionByID :one
SELECT
    *
FROM
    intermediate_sessions
WHERE
    id = $1;

-- name: GetIntermediateSessionSigningKeyByID :one
SELECT
    *
FROM
    intermediate_session_signing_keys
WHERE
    id = $1;

-- name: GetIntermediateSessionSigningKeyByProjectID :one
SELECT
    *
FROM
    intermediate_session_signing_keys
WHERE
    project_id = $1
ORDER BY
    create_time DESC
LIMIT 1;

-- name: GetOrganizationByID :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1;

-- name: GetOrganizationByProjectIDAndID :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1
    AND project_id = $2;

-- name: GetProjectByID :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetSessionByID :one
SELECT
    *
FROM
    sessions
WHERE
    id = $1;

-- name: GetSessionSigningKeyByID :one
SELECT
    *
FROM
    session_signing_keys
WHERE
    id = $1;

-- name: GetSessionSigningKeyByProjectID :one
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1
ORDER BY
    create_time DESC
LIMIT 1;

-- name: GetOrganizationByGoogleHostedDomain :one
SELECT
    *
FROM
    organizations
WHERE
    google_hosted_domain = $1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    verified_email = $1
    OR unverified_email = $1;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: GetUserByGoogleUserID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND google_user_id = $2;

-- name: GetUserByMicrosoftUserID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND microsoft_user_id = $2;

-- name: GetUserBySessionID :one
SELECT
    *
FROM
    users
WHERE
    users.id = (
        SELECT
            user_id
        FROM
            sessions
        WHERE
            sessions.id = $1);

-- name: GetUserByUnverifiedEmail :one
SELECT
    *
FROM
    users
WHERE
    unverified_email = $1;

-- name: GetUserByVerifiedEmail :one
SELECT
    *
FROM
    users
WHERE
    verified_email = $1;

-- name: ListOrganizations :many
SELECT
    *
FROM
    organizations;

-- name: ListOrganizationsByProjectId :many
SELECT
    *
FROM
    organizations
WHERE
    project_id = $1
ORDER BY
    id
LIMIT $2;

-- name: ListOrganizationsByProjectIdAndEmail :many
SELECT
    o.*
FROM
    organizations AS o
    JOIN users AS u ON o.id = users.organization_id
WHERE
    project_id = $1
    AND u.verified_email = $2
    OR u.unverified_email = $2
ORDER BY
    o.display_name
LIMIT $3;

-- name: ListProjects :many
SELECT
    *
FROM
    projects
ORDER BY
    id
LIMIT $1;

-- name: ListUsersByEmail :many
SELECT
    *
FROM
    users
WHERE
    unverified_email = $1
    OR verified_email = $1;

-- name: ListUsersByOrganization :many
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
ORDER BY
    id
LIMIT $2;

-- name: ListVerifiedEmails :many
SELECT
    *
FROM
    verified_emails
WHERE
    project_id = $1
    AND email = $2
    AND (google_user_id = $3
        OR microsoft_user_id = $4)
ORDER BY
    id;

-- name: RevokeIntermediateSession :one
UPDATE
    intermediate_sessions
SET
    revoked = TRUE
WHERE
    id = $1
RETURNING
    *;

-- name: RevokeSession :one
UPDATE
    sessions
SET
    revoked = TRUE
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateOrganization :one
UPDATE
    organizations
SET
    display_name = $2,
    google_hosted_domain = $3,
    microsoft_tenant_id = $4,
    override_log_in_with_password_enabled = $5,
    override_log_in_with_google_enabled = $6,
    override_log_in_with_microsoft_enabled = $7
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateOrganizationDisplayName :one
UPDATE
    organizations
SET
    display_name = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateOrganizationGoogleHostedDomain :one
UPDATE
    organizations
SET
    google_hosted_domain = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateOrganizationMicrosoftTenantID :one
UPDATE
    organizations
SET
    microsoft_tenant_id = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateOrganizationOverrides :one
UPDATE
    organizations
SET
    override_log_in_with_password_enabled = $2,
    override_log_in_with_google_enabled = $3,
    override_log_in_with_microsoft_enabled = $4
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProject :one
UPDATE
    projects
SET
    log_in_with_password_enabled = $2,
    log_in_with_google_enabled = $3,
    log_in_with_microsoft_enabled = $4,
    google_oauth_client_id = $5,
    google_oauth_client_secret = $6,
    microsoft_oauth_client_id = $7,
    microsoft_oauth_client_secret = $8
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProjectOrganizationID :one
UPDATE
    projects
SET
    organization_id = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProjectGoogleOAuthClient :one
UPDATE
    projects
SET
    google_oauth_client_id = $2,
    google_oauth_client_secret = $3
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProjectMicrosoftOAuthClient :one
UPDATE
    projects
SET
    microsoft_oauth_client_id = $2,
    microsoft_oauth_client_secret = $3
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProjectLoginMethods :one
UPDATE
    projects
SET
    log_in_with_password_enabled = $2,
    log_in_with_google_enabled = $3,
    log_in_with_microsoft_enabled = $4
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUser :one
UPDATE
    users
SET
    organization_id = $2,
    unverified_email = $3,
    verified_email = $4,
    password_bcrypt = $5,
    google_user_id = $6,
    microsoft_user_id = $7
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUserPassword :one
UPDATE
    users
SET
    password_bcrypt = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUserGoogleUserID :one
UPDATE
    users
SET
    google_user_id = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUserMicrosoftUserID :one
UPDATE
    users
SET
    microsoft_user_id = $2
WHERE
    id = $1
RETURNING
    *;

-- name: VerifyIntermediateSessionEmail :one
UPDATE
    intermediate_sessions
SET
    unverified_email = NULL,
    verified_email = $2
WHERE
    id = $1
RETURNING
    *;

-- name: VerifyUserEmail :one
UPDATE
    users
SET
    unverified_email = NULL,
    verified_email = $2
WHERE
    id = $1
RETURNING
    *;

-- name: CountAllProjects :one
SELECT
    count(*)
FROM
    projects;

-- name: CreateProjectAPIKey :one
INSERT INTO project_api_keys (id, project_id, create_time, revoked, secret_token_sha256)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetProjectAPIKeyBySecretTokenSHA256 :one
SELECT
    *
FROM
    project_api_keys
WHERE
    secret_token_sha256 = $1;

