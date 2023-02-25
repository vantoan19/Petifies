-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
  id, email, password, first_name, last_name, is_activated
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $2, password = $3, first_name = $4, last_name = $5, is_activated = $6, updated_at = $7
WHERE id = $1
RETURNING *;

-- name: DeleteUserByID :one
DELETE FROM users
WHERE id = $1
RETURNING *;

-- name: DeleteUserByEmail :one
DELETE FROM users
WHERE email = $1
RETURNING *;