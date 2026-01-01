-- name: GetVersion :one
SELECT
  *
FROM
  version
LIMIT
  1;
