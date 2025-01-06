-- name: GetSessionSigningKeyPublicKey :one
SELECT
    public_key
FROM
    session_signing_keys
WHERE
    project_id = $1
    AND id = $2
    AND expire_time > @now;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, password_bcrypt, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

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

-- name: GetOrganizationByID :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1;

-- name: UpdateOrganization :one
UPDATE
    organizations
SET
    display_name = $2,
    google_hosted_domain = $3,
    microsoft_tenant_id = $4,
    override_log_in_methods = $5,
    override_log_in_with_password_enabled = $6,
    override_log_in_with_google_enabled = $7,
    override_log_in_with_microsoft_enabled = $8
WHERE
    id = $1
RETURNING
    *;

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

-- name: GetSessionDetailsByRefreshTokenSHA256 :one
SELECT
    sessions.id AS session_id,
    users.id AS user_id,
    organizations.id AS organization_id,
    projects.id AS project_id
FROM
    sessions
    JOIN users ON sessions.user_id = users.id
    JOIN organizations ON users.organization_id = organizations.id
    JOIN projects ON organizations.id = projects.organization_id
WHERE
    revoked = FALSE
    AND refresh_token_sha256 = $1;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: ListSCIMAPIKeys :many
SELECT
    *
FROM
    scim_api_keys
WHERE
    organization_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetSCIMAPIKey :one
SELECT
    *
FROM
    scim_api_keys
WHERE
    id = $1
    AND organization_id = $2;

-- name: CreateSCIMAPIKey :one
INSERT INTO scim_api_keys (id, organization_id, display_name, token_sha256)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateSCIMAPIKey :one
UPDATE
    scim_api_keys
SET
    display_name = $1
WHERE
    id = $2
RETURNING
    *;

-- name: DeleteSCIMAPIKey :exec
DELETE FROM scim_api_keys
WHERE id = $1;

-- name: RevokeSCIMAPIKey :one
UPDATE
    scim_api_keys
SET
    token_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: ListUsers :many
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = $1
    AND organization_id = $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    is_owner = $1
WHERE
    id = $2
RETURNING
    *;
