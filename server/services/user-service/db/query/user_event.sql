-- name: CreateUserEvent :one
INSERT INTO user_events (
  id, payload, outbox_state, created_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserEventByLockerID :many
SELECT * FROM user_events
WHERE locked_by = $1;

-- name: LockStartedEvents :exec
UPDATE user_events SET locked_by = $1, locked_at = $2
WHERE outbox_state = "STARTED";

-- name: UpdateEvent :exec
UPDATE user_events SET outbox_state = $2, locked_by = $3, locked_at = $4, error = $5, completed_at = $6
WHERE id = $1;

-- name: UnlockEventsByLockerID :exec
UPDATE user_events SET locked_by = NULL, locked_at = NULL
WHERE locked_by = $1;

-- name: UnlockEventsBeforeDatetime :exec
UPDATE user_events SET locked_by = NULL, locked_at = NULL
WHERE locked_at < $1;

-- name: DeleteEventsBeforeDatetime :exec
DELETE FROM user_events
WHERE created_at < $1;
