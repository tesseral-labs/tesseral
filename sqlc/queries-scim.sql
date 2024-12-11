-- name: GetSCIMAPIKeyByTokenSHA256 :one
SELECT
    scim_api_keys.*
FROM
    scim_api_keys
    JOIN organizations ON scim_api_keys.organization_id = organizations.id
WHERE
    token_sha256 = $1
    AND organizations.project_id = $2;

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

