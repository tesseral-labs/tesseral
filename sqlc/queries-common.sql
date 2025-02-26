-- name: GetProjectIDByVaultDomain :one
SELECT
    id
FROM
    projects
WHERE
    vault_domain = $1;

-- name: GetSessionDetailsByRefreshTokenSHA256 :one
SELECT
    sessions.id AS session_id,
    users.id AS user_id,
    users.email AS user_email,
    organizations.id AS organization_id,
    organizations.display_name AS organization_display_name,
    organizations.project_id AS project_id,
    sessions.impersonator_user_id
FROM
    sessions
    JOIN users ON sessions.user_id = users.id
    JOIN organizations ON users.organization_id = organizations.id
WHERE
    sessions.refresh_token_sha256 = $1;

-- name: GetImpersonatorUserByID :one
SELECT
    id,
    email
FROM
    users
WHERE
    id = $1;

-- name: GetCurrentSessionSigningKeyByProjectID :one
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1
ORDER BY
    create_time DESC
LIMIT 1;

-- name: BumpSessionLastActiveTime :exec
UPDATE
    sessions
SET
    last_active_time = now()
WHERE
    id = $1;

