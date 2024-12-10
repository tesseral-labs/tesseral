-- name: GetSCIMAPIKeyByTokenSHA256 :one
SELECT
    scim_api_keys.*
FROM
    scim_api_keys
    JOIN organizations ON scim_api_keys.organization_id = organizations.id
WHERE
    token_sha256 = $1
    AND organizations.project_id = $2;

