 -- name: CreateCollectionTrain :one
INSERT INTO collection_trains (
  user_id,
  train_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetCollectionTrainByID :one
SELECT * FROM collection_trains
WHERE id = $1
LIMIT 1;

-- name: GetCollectionTrainforUpdateByID :one
SELECT * FROM collection_trains
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetCollectionTrain :one
SELECT * FROM collection_trains
WHERE user_id = $1 AND train_id = $2
LIMIT 1;

-- name: GetCollectionTrainforUpdate :one
SELECT * FROM collection_trains
WHERE user_id = $1 AND train_id = $2
LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUserCollection :many
SELECT * FROM collection_trains
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListCollectionTrains :many
SELECT * FROM collection_trains
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: UpdateCollectionTrain :one
UPDATE collection_trains 
SET user_id = $2, times_traded = times_traded + 1
WHERE id = $1
RETURNING *;

-- name: DeleteCollectionTrain :exec
DELETE from collection_trains WHERE user_id = $1 AND train_id = $2;

-- name: ListUserTrains :many
SELECT 
    trains.*, 
    CASE WHEN collection_trains.train_id IS NULL THEN FALSE ELSE TRUE END AS is_in_collection,
    CASE WHEN wishlist_trains.train_id IS NULL THEN FALSE ELSE TRUE END AS is_in_wishlist
FROM 
    trains 
LEFT JOIN 
    collection_trains ON trains.id = collection_trains.train_id AND collection_trains.user_id = $1
LEFT JOIN
    wishlist_trains ON trains.id = wishlist_trains.train_id AND wishlist_trains.user_id = $1
ORDER BY 
    trains.id
LIMIT $2 
OFFSET $3;

-- name: GetUserCollectionWithWishlistStatus :many
SELECT 
    c.*,
    CASE WHEN w.train_id IS NULL THEN false ELSE true END AS is_in_wishlist
FROM collection_trains c
LEFT JOIN wishlist_trains w ON c.train_id = w.train_id AND w.user_id = c.user_id
WHERE c.user_id = $1;
