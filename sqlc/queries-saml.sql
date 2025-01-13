-- name: GetSAMLConnection :one
SELECT
    saml_connections.*
FROM
    saml_connections
    JOIN organizations ON saml_connections.organization_id = organizations.id
WHERE
    organizations.project_id = $1
    AND saml_connections.id = $2;

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
    $1 = ANY (custom_domains);

