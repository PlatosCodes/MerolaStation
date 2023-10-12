 -- name: InsertActivationToken :one
INSERT INTO activation_tokens (
  user_id,
  activation_token,
  expires_at
) VALUES (
  $1, $2, $3
) RETURNING *;

 -- name: GetActivationToken :one
SELECT * FROM activation_tokens
WHERE activation_token = $1 LIMIT 1;

-- name: DeleteActivationToken :exec
DELETE from activation_tokens WHERE user_id = $1;