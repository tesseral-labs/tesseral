-- name: GetProject :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

