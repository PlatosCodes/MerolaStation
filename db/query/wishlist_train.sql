 -- name: CreateWishlistTrain :one
INSERT INTO wishlist_trains (
  user_id,
  train_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserWishlistTrains :many
SELECT * FROM wishlist_trains
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: ListWishlistTrains :many
SELECT * FROM wishlist_trains
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: DeleteWishlistTrain :exec
DELETE from wishlist_trains WHERE user_id = $1 AND train_id = $2;