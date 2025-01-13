-- name: GetSCIMAPIKeyByTokenSHA256 :one
SELECT
    scim_api_keys.*
FROM
    scim_api_keys
    JOIN organizations ON scim_api_keys.organization_id = organizations.id
WHERE
    secret_token_sha256 = $1
    AND organizations.project_id = $2;

-- name: GetOrganizationDomains :many
SELECT
    DOMAIN
FROM
    organization_domains
WHERE
    organization_id = $1;

-- name: GetProjectIDByCustomDomain :one
SELECT
    id
FROM
    projects
WHERE
    custom_auth_domain = $1;

-- name: CountUsers :one
SELECT
    count(*)
FROM
    users
WHERE
    organization_id = $1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
ORDER BY
    id
LIMIT $2 OFFSET $3;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND id = $2;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND email = $2;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = $1,
    deactivate_time = $2
WHERE
    id = $3
    AND organization_id = $4
RETURNING
    *;

-- name: DeactivateUser :one
UPDATE
    users
SET
    deactivate_time = $1
WHERE
    id = $2
    AND organization_id = $3
RETURNING
    *;

