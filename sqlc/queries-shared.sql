-- name: GetProjectIDByCustomAuthDomain :one
SELECT
    id
FROM
    projects
WHERE
    custom_auth_domain = $1;

