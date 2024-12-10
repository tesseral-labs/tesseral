-- name: CompleteEmailVerificationChallenge :one
UPDATE
    email_verification_challenges
SET
    complete_time = $1
WHERE
    id = $2
RETURNING
    *;

-- name: GetIntermediateSessionByTokenSHA256 :one
SELECT
    *
FROM
    intermediate_sessions
WHERE
    token_sha256 = $1;

-- name: CreateEmailVerificationChallenge :one
INSERT INTO email_verification_challenges (id, project_id, intermediate_session_id, challenge_sha256, expire_time)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateIntermediateSession :one
INSERT INTO intermediate_sessions (id, project_id, expire_time, email, google_user_id, microsoft_user_id, token_sha256)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, google_hosted_domain, microsoft_tenant_id, override_log_in_with_google_enabled, override_log_in_with_microsoft_enabled, override_log_in_with_password_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expire_time, revoked)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateVerifiedEmail :one
INSERT INTO verified_emails (id, project_id, email, google_user_id, google_hosted_domain, microsoft_user_id, microsoft_tenant_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: GetEmailVerificationChallengeByID :one
SELECT
    *
FROM
    email_verification_challenges
WHERE
    id = $1;

-- name: GetEmailVerificationChallengeForCompletion :one
SELECT
    *
FROM
    email_verification_challenges
WHERE
    project_id = $1
    AND intermediate_session_id = $2
    AND expire_time > $3
    AND revoked = FALSE
    AND complete_time IS NULL
ORDER BY
    create_time DESC
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

-- name: GetCurrentSessionKeyByProjectID :one
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1
ORDER BY
    create_time DESC
LIMIT 1;

-- name: GetOrganizationUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND email = $2;

-- name: GetOrganizationUserByGoogleUserID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND google_user_id = $2;

-- name: GetOrganizationUserByMicrosoftUserID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND microsoft_user_id = $2;

-- name: GetProjectOrganizationByID :one
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

-- name: GetSessionSigningKeysByProjectID :many
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1;

-- name: ListOrganizationsByProjectIdAndEmail :many
SELECT
    o.*
FROM
    organizations AS o
    JOIN users AS u ON o.id = users.organization_id
WHERE
    project_id = $1
    AND u.email = $2
ORDER BY
    o.display_name
LIMIT $3;

-- name: ListUsersByEmail :many
SELECT
    *
FROM
    users
WHERE
    email = $1;

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

-- name: RevokeEmailVerificationChallenge :one
UPDATE
    email_verification_challenges
SET
    revoked = TRUE
WHERE
    id = $1
RETURNING
    *;

-- name: RevokeIntermediateSession :one
UPDATE
    intermediate_sessions
SET
    revoked = TRUE
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateIntermediateSessionGoogleOAuthStateSHA256 :one
UPDATE
    intermediate_sessions
SET
    google_oauth_state_sha256 = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionGoogleDetails :one
UPDATE
    intermediate_sessions
SET
    email = $1,
    google_user_id = $2,
    google_hosted_domain = $3
WHERE
    id = $4
RETURNING
    *;

-- name: UpdateIntermediateSessionMicrosoftOAuthStateSHA256 :one
UPDATE
    intermediate_sessions
SET
    microsoft_oauth_state_sha256 = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionMicrosoftDetails :one
UPDATE
    intermediate_sessions
SET
    email = $1,
    microsoft_user_id = $2,
    microsoft_tenant_id = $3
WHERE
    id = $4
RETURNING
    *;

