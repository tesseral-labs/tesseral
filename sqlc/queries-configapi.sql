-- name: GetPublishableKeyConfiguration :one
SELECT
    projects.id AS project_id,
    projects.vault_domain,
    projects.cookie_domain,
    publishable_keys.dev_mode
FROM
    publishable_keys
    JOIN projects ON publishable_keys.project_id = projects.id
WHERE
    publishable_keys.id = $1;

-- name: GetPublishableKeyTrustedDomains :many
SELECT
    project_trusted_domains.domain
FROM
    publishable_keys
    JOIN projects ON publishable_keys.project_id = projects.id
    JOIN project_trusted_domains ON projects.id = project_trusted_domains.project_id
WHERE
    publishable_keys.id = $1;

-- name: GetPublishableKeySessionSigningPublicKeys :many
SELECT
    session_signing_keys.id,
    session_signing_keys.public_key
FROM
    publishable_keys
    JOIN projects ON publishable_keys.project_id = projects.id
    JOIN session_signing_keys ON projects.id = session_signing_keys.project_id
WHERE
    publishable_keys.id = $1;

