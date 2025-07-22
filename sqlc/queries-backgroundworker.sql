-- name: GetProjectWebhookSettings :one
SELECT
    *
FROM
    project_webhook_settings
WHERE
    project_id = $1;

