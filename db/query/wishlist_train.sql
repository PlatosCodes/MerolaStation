 -- name: CreateWishlistTrain :one
INSERT INTO wishlist_trains (
  user_id,
  train_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListUserWishlist :many
SELECT * FROM wishlist_trains
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: ListWishlists :many
SELECT * FROM wishlist_trains
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: DeleteWishlistTrain :exec
DELETE from wishlist_trains WHERE user_id = $1 AND train_id = $2;

-- name: GetWishlistTrain :one
SELECT * FROM wishlist_trains
WHERE user_id = $1 AND train_id = $2
LIMIT 1;

-- name: GetUserWishlistWithCollectionStatus :many
SELECT 
    w.*,
    CASE WHEN c.train_id IS NULL THEN false ELSE true END AS is_in_collection
FROM wishlist_trains w
LEFT JOIN collection_trains c ON w.train_id = c.train_id AND w.user_id = c.user_id
WHERE w.user_id = $1;
