-- name: GetProject :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetOIDCConnection :one
SELECT
    oidc_connections.*
FROM
    oidc_connections
    JOIN organizations ON oidc_connections.organization_id = organizations.id
WHERE
    organizations.project_id = $1
    AND organizations.log_in_with_oidc
    AND oidc_connections.id = $2;

-- name: CreateOIDCIntermediateSession :one
INSERT INTO oidc_intermediate_sessions (oidc_intermediate_session_id, oidc_connection_id, code_verifier)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: DeleteOIDCIntermediateSession :one
DELETE FROM oidc_intermediate_sessions
WHERE oidc_intermediate_session_id = $1
RETURNING
    *;

-- name: GetOrganizationDomains :many
SELECT
    DOMAIN
FROM
    organization_domains
WHERE
    organization_id = $1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND email = $2;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, is_owner)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expire_time, refresh_token_sha256, primary_auth_factor)
    VALUES ($1, $2, $3, $4, 'oidc')
RETURNING
    *;

