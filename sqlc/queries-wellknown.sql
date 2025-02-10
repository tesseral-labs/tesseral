-- name: GetProjectTrustedDomains :many
SELECT
    *
FROM
    project_trusted_domains
WHERE
    project_id = $1;

