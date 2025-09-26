-- name: CreateSession :one
INSERT INTO session (
  id,
  account_id,
  refresh_token,
  account_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: DeleteSession :exec
DELETE FROM session
WHERE id = $1;

-- name: GetSession :one
SELECT * FROM session
WHERE id = $1 LIMIT 1;
