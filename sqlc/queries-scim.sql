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
INSERT INTO users (id, organization_id, email, is_owner)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = $1
WHERE
    id = $2
    AND organization_id = $3
RETURNING
    *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
    AND organization_id = $2
RETURNING
    *;

-- name: CreateAuditLogEvent :one
INSERT INTO audit_log_events (id, project_id, organization_id, actor_scim_api_key_id, resource_type, resource_id, event_name, event_time, event_details)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, coalesce(@event_details, '{}'::jsonb))
RETURNING
    *;

