 -- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  email
) VALUES (
  $1, $2, $3
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
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  activated = COALESCE(sqlc.narg(activated), activated)
WHERE
  username = sqlc.arg(username)
RETURNING *;