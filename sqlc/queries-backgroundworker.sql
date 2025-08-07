-- name: GetProjectWebhookSettings :one
SELECT
    *
FROM
    project_webhook_settings
WHERE
    project_id = $1;

-- name: GetProject :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetUserInvite :one
SELECT
    *
FROM
    user_invites
WHERE
    id = $1;

-- name: GetOrganization :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1;

