-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: GetSessionsForUser :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY id;

-- name: BulkCreateSession :many
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  exprires_at,
  client_ip,
  is_disabled
) 
SELECT
UNNEST(@id::uuid[]) AS id,
UNNEST(@user_id::uuid[]) as user_id,
UNNEST(@refresh_token::VARCHAR[]) as refresh_token,
UNNEST(@exprires_at::timestamptz[]) as exprires_at,
UNNEST(@client_ip::VARCHAR[]) as client_ip,
UNNEST(@is_disabled::BOOLEAN[]) as is_disabled
RETURNING *;

-- name: DisableSession :one
UPDATE sessions SET is_disabled = true
WHERE id = $1 
RETURNING *;

-- name: BulkUpsertSessions :many
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  exprires_at,
  client_ip,
  is_disabled
) 
SELECT
UNNEST(@id::uuid[]) AS id,
UNNEST(@user_id::uuid[]) as user_id,
UNNEST(@refresh_token::VARCHAR[]) as refresh_token,
UNNEST(@exprires_at::timestamptz[]) as exprires_at,
UNNEST(@client_ip::VARCHAR[]) as client_ip,
UNNEST(@is_disabled::BOOLEAN[]) as is_disabled
ON CONFLICT (id) DO UPDATE 
  SET id = excluded.id, 
      user_id = excluded.user_id,
      refresh_token = excluded.refresh_token,
      exprires_at = excluded.exprires_at,
      client_ip = excluded.client_ip,
      is_disabled = excluded.is_disabled
RETURNING *;