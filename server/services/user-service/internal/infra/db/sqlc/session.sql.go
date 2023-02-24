// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const bulkCreateSession = `-- name: BulkCreateSession :many
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  exprires_at,
  client_ip,
  is_disabled
) 
SELECT
UNNEST($1::uuid[]) AS id,
UNNEST($2::uuid[]) as user_id,
UNNEST($3::VARCHAR[]) as refresh_token,
UNNEST($4::timestamptz[]) as exprires_at,
UNNEST($5::VARCHAR[]) as client_ip,
UNNEST($6::BOOLEAN[]) as is_disabled
RETURNING id, user_id, refresh_token, exprires_at, client_ip, is_disabled, created_at
`

type BulkCreateSessionParams struct {
	ID           []uuid.UUID `json:"id"`
	UserID       []uuid.UUID `json:"user_id"`
	RefreshToken []string    `json:"refresh_token"`
	ExpriresAt   []time.Time `json:"exprires_at"`
	ClientIp     []string    `json:"client_ip"`
	IsDisabled   []bool      `json:"is_disabled"`
}

func (q *Queries) BulkCreateSession(ctx context.Context, arg BulkCreateSessionParams) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, bulkCreateSession,
		pq.Array(arg.ID),
		pq.Array(arg.UserID),
		pq.Array(arg.RefreshToken),
		pq.Array(arg.ExpriresAt),
		pq.Array(arg.ClientIp),
		pq.Array(arg.IsDisabled),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RefreshToken,
			&i.ExpriresAt,
			&i.ClientIp,
			&i.IsDisabled,
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

const bulkUpsertSessions = `-- name: BulkUpsertSessions :many
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  exprires_at,
  client_ip,
  is_disabled
) 
SELECT
UNNEST($1::uuid[]) AS id,
UNNEST($2::uuid[]) as user_id,
UNNEST($3::VARCHAR[]) as refresh_token,
UNNEST($4::timestamptz[]) as exprires_at,
UNNEST($5::VARCHAR[]) as client_ip,
UNNEST($6::BOOLEAN[]) as is_disabled
ON CONFLICT (id) DO UPDATE 
  SET id = excluded.id, 
      user_id = excluded.user_id,
      refresh_token = excluded.refresh_token,
      exprires_at = excluded.exprires_at,
      client_ip = excluded.client_ip,
      is_disabled = excluded.is_disabled
RETURNING id, user_id, refresh_token, exprires_at, client_ip, is_disabled, created_at
`

type BulkUpsertSessionsParams struct {
	ID           []uuid.UUID `json:"id"`
	UserID       []uuid.UUID `json:"user_id"`
	RefreshToken []string    `json:"refresh_token"`
	ExpriresAt   []time.Time `json:"exprires_at"`
	ClientIp     []string    `json:"client_ip"`
	IsDisabled   []bool      `json:"is_disabled"`
}

func (q *Queries) BulkUpsertSessions(ctx context.Context, arg BulkUpsertSessionsParams) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, bulkUpsertSessions,
		pq.Array(arg.ID),
		pq.Array(arg.UserID),
		pq.Array(arg.RefreshToken),
		pq.Array(arg.ExpriresAt),
		pq.Array(arg.ClientIp),
		pq.Array(arg.IsDisabled),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RefreshToken,
			&i.ExpriresAt,
			&i.ClientIp,
			&i.IsDisabled,
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

const disableSession = `-- name: DisableSession :one
UPDATE sessions SET is_disabled = true
WHERE id = $1 
RETURNING id, user_id, refresh_token, exprires_at, client_ip, is_disabled, created_at
`

func (q *Queries) DisableSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, disableSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.ExpriresAt,
		&i.ClientIp,
		&i.IsDisabled,
		&i.CreatedAt,
	)
	return i, err
}

const getSessionById = `-- name: GetSessionById :one
SELECT id, user_id, refresh_token, exprires_at, client_ip, is_disabled, created_at FROM sessions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSessionById(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionById, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.ExpriresAt,
		&i.ClientIp,
		&i.IsDisabled,
		&i.CreatedAt,
	)
	return i, err
}

const getSessionsForUser = `-- name: GetSessionsForUser :many
SELECT id, user_id, refresh_token, exprires_at, client_ip, is_disabled, created_at FROM sessions
WHERE user_id = $1
ORDER BY id
`

func (q *Queries) GetSessionsForUser(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, getSessionsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RefreshToken,
			&i.ExpriresAt,
			&i.ClientIp,
			&i.IsDisabled,
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