 -- name: CreateTrain :one
INSERT INTO trains (
  model_number,
  name
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetTrain :one
SELECT * FROM trains
WHERE id = $1 LIMIT 1;

-- name: GetTrainByModel :one
SELECT * FROM trains
WHERE model_number = $1 LIMIT 1;

-- name: GetTrainByName :one
SELECT * FROM trains
WHERE name = $1 LIMIT 1;

-- name: ListTrains :many
SELECT * FROM trains
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTrainValue :exec
UPDATE trains SET value = $2, version = version + 1
WHERE id = $1;

-- name: DeleteTrain :exec
DELETE from trains WHERE id = $1;

-- name: SearchTrainsByModelNumberSuggestions :many
SELECT DISTINCT id, model_number, name
FROM trains
WHERE model_number ILIKE $1 || '%'
ORDER BY model_number
LIMIT $2
OFFSET $3;

