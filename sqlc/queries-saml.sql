-- name: GetProject :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetSAMLConnection :one
SELECT
    saml_connections.*
FROM
    saml_connections
    JOIN organizations ON saml_connections.organization_id = organizations.id
WHERE
    organizations.project_id = $1
    AND organizations.log_in_with_saml
    AND saml_connections.id = $2;

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
INSERT INTO sessions (id, user_id, expire_time, refresh_token_sha256)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;
