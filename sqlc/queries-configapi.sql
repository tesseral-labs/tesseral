-- name: GetPublishableKeyConfiguration :one
SELECT
    projects.id AS project_id,
    projects.vault_domain
FROM
    publishable_keys
    JOIN projects ON publishable_keys.project_id = projects.id
WHERE
    publishable_keys.id = $1;

-- name: GetPublishableKeySessionSigningPublicKeys :many
SELECT
    session_signing_keys.id,
    session_signing_keys.public_key
FROM
    publishable_keys
    JOIN projects ON publishable_keys.project_id = projects.id
    JOIN session_signing_keys ON projects.id = session_signing_keys.project_id;

