-- name: GetProject :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetAPIKeyRoleAssignment :one
SELECT
    *
FROM
    api_key_role_assignments
WHERE
    id = $1;

-- name: GetAPIKey :one
SELECT
    *
FROM
    api_keys
WHERE
    id = $1;

-- name: GetOrganization :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1;

-- name: GetPasskey :one
SELECT
    *
FROM
    passkeys
WHERE
    id = $1;

-- name: GetRole :one
SELECT
    *
FROM
    roles
WHERE
    id = $1;

-- name: GetActions :many
SELECT
    *
FROM
    actions
WHERE
    project_id = $1;

-- name: BatchGetRoleActionsByRoleID :many
SELECT
    *
FROM
    role_actions
WHERE
    role_id = ANY ($1::uuid[]);

-- name: GetSAMLConnection :one
SELECT
    *
FROM
    saml_connections
WHERE
    id = $1;

-- name: GetSCIMAPIKey :one
SELECT
    *
FROM
    scim_api_keys
WHERE
    id = $1;

-- name: GetUserRoleAssignment :one
SELECT
    *
FROM
    user_role_assignments
WHERE
    id = $1;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: GetUserInvite :one
SELECT
    *
FROM
    user_invites
WHERE
    id = $1;

-- name: GetSession :one
SELECT
    *
FROM
    sessions
WHERE
    id = $1;
