-- name: GetSessionSigningKeysByProjectID :many
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1;

