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

