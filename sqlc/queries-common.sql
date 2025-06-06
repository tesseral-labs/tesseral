-- name: GetProjectIDByVaultDomain :one
SELECT
    id
FROM
    projects
WHERE
    vault_domain = $1;

-- name: GetProjectTrustedDomains :many
SELECT
    project_trusted_domains.domain
FROM
    project_trusted_domains
WHERE
    project_id = $1;

-- name: GetSessionDetailsByRelayedSessionRefreshTokenSHA256 :one
SELECT
    sessions.id AS session_id,
    users.id AS user_id,
    users.is_owner AS user_is_owner,
    users.email AS user_email,
    users.display_name AS user_display_name,
    users.profile_picture_url AS user_profile_picture_url,
    organizations.id AS organization_id,
    organizations.display_name AS organization_display_name,
    organizations.project_id AS project_id,
    sessions.impersonator_user_id
FROM
    relayed_sessions
    JOIN sessions ON relayed_sessions.session_id = sessions.id
    JOIN users ON sessions.user_id = users.id
    JOIN organizations ON users.organization_id = organizations.id
WHERE
    relayed_sessions.relayed_refresh_token_sha256 = $1;

-- name: GetSessionDetailsByRefreshTokenSHA256 :one
SELECT
    sessions.id AS session_id,
    users.is_owner AS user_is_owner,
    users.id AS user_id,
    users.email AS user_email,
    users.display_name AS user_display_name,
    users.profile_picture_url AS user_profile_picture_url,
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

-- name: GetProjectActions :many
SELECT
    name
FROM
    actions
WHERE
    project_id = $1;

-- name: GetUserActions :many
SELECT DISTINCT
    (actions.name)
FROM
    users
    JOIN user_role_assignments ON users.id = user_role_assignments.user_id
    JOIN roles ON user_role_assignments.role_id = roles.id
    JOIN role_actions ON roles.id = role_actions.role_id
    JOIN actions ON role_actions.action_id = actions.id
WHERE
    user_id = $1;

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

-- name: GetProjectCookieDomainByProjectID :one
SELECT
    cookie_domain
FROM
    projects
WHERE
    id = $1;

