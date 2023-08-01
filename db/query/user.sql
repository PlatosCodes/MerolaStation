 -- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  email, first_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE from users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE(sqlc.narg(email), email),
  first_name = COALESCE(sqlc.narg(first_name), first_name)
WHERE
  username = sqlc.arg(username)
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET
  hashed_password = $1,
  password_changed_at = $2
WHERE
  username = sqlc.arg(username)
RETURNING *;