 -- name: CreateCollectionTrain :one
INSERT INTO collection_trains (
  user_id,
  train_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserCollectionTrains :many
SELECT * FROM collection_trains
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: ListCollectionTrains :many
SELECT * FROM collection_trains
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: DeleteCollectionTrain :exec
DELETE from collection_trains WHERE user_id = $1 AND train_id = $2;