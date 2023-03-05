// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user_event.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const createUserEvent = `-- name: CreateUserEvent :one
INSERT INTO user_events (
  id, payload, outbox_state, created_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, payload, outbox_state, locked_by, locked_at, error, completed_at, created_at
`

type CreateUserEventParams struct {
	ID          uuid.UUID       `json:"id"`
	Payload     json.RawMessage `json:"payload"`
	OutboxState OutboxState     `json:"outbox_state"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (q *Queries) CreateUserEvent(ctx context.Context, arg CreateUserEventParams) (UserEvent, error) {
	row := q.db.QueryRowContext(ctx, createUserEvent,
		arg.ID,
		arg.Payload,
		arg.OutboxState,
		arg.CreatedAt,
	)
	var i UserEvent
	err := row.Scan(
		&i.ID,
		&i.Payload,
		&i.OutboxState,
		&i.LockedBy,
		&i.LockedAt,
		&i.Error,
		&i.CompletedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEventsBeforeDatetime = `-- name: DeleteEventsBeforeDatetime :exec
DELETE FROM user_events
WHERE created_at < $1
`

func (q *Queries) DeleteEventsBeforeDatetime(ctx context.Context, createdAt time.Time) error {
	_, err := q.db.ExecContext(ctx, deleteEventsBeforeDatetime, createdAt)
	return err
}

const getUserEventByLockerID = `-- name: GetUserEventByLockerID :many
SELECT id, payload, outbox_state, locked_by, locked_at, error, completed_at, created_at FROM user_events
WHERE locked_by = $1
`

func (q *Queries) GetUserEventByLockerID(ctx context.Context, lockedBy uuid.NullUUID) ([]UserEvent, error) {
	rows, err := q.db.QueryContext(ctx, getUserEventByLockerID, lockedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserEvent
	for rows.Next() {
		var i UserEvent
		if err := rows.Scan(
			&i.ID,
			&i.Payload,
			&i.OutboxState,
			&i.LockedBy,
			&i.LockedAt,
			&i.Error,
			&i.CompletedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const lockStartedEvents = `-- name: LockStartedEvents :exec
UPDATE user_events SET locked_by = $1, locked_at = $2
WHERE outbox_state = 'STARTED'
`

type LockStartedEventsParams struct {
	LockedBy uuid.NullUUID `json:"locked_by"`
	LockedAt sql.NullTime  `json:"locked_at"`
}

func (q *Queries) LockStartedEvents(ctx context.Context, arg LockStartedEventsParams) error {
	_, err := q.db.ExecContext(ctx, lockStartedEvents, arg.LockedBy, arg.LockedAt)
	return err
}

const unlockEventsBeforeDatetime = `-- name: UnlockEventsBeforeDatetime :exec
UPDATE user_events SET locked_by = NULL, locked_at = NULL
WHERE locked_at < $1
`

func (q *Queries) UnlockEventsBeforeDatetime(ctx context.Context, lockedAt sql.NullTime) error {
	_, err := q.db.ExecContext(ctx, unlockEventsBeforeDatetime, lockedAt)
	return err
}

const unlockEventsByLockerID = `-- name: UnlockEventsByLockerID :exec
UPDATE user_events SET locked_by = NULL, locked_at = NULL
WHERE locked_by = $1
`

func (q *Queries) UnlockEventsByLockerID(ctx context.Context, lockedBy uuid.NullUUID) error {
	_, err := q.db.ExecContext(ctx, unlockEventsByLockerID, lockedBy)
	return err
}

const updateEvent = `-- name: UpdateEvent :exec
UPDATE user_events SET outbox_state = $2, locked_by = $3, locked_at = $4, error = $5, completed_at = $6
WHERE id = $1
`

type UpdateEventParams struct {
	ID          uuid.UUID      `json:"id"`
	OutboxState OutboxState    `json:"outbox_state"`
	LockedBy    uuid.NullUUID  `json:"locked_by"`
	LockedAt    sql.NullTime   `json:"locked_at"`
	Error       sql.NullString `json:"error"`
	CompletedAt sql.NullTime   `json:"completed_at"`
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) error {
	_, err := q.db.ExecContext(ctx, updateEvent,
		arg.ID,
		arg.OutboxState,
		arg.LockedBy,
		arg.LockedAt,
		arg.Error,
		arg.CompletedAt,
	)
	return err
}
