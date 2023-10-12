 -- name: CreateTrain :one
INSERT INTO trains (
  model_number,
  name
) VALUES (
  $1, $2
) RETURNING *;

 -- name: CreateImageTrain :one
INSERT INTO trains (
  model_number,
  name,
  img_url
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTrain :one
SELECT * FROM trains
WHERE id = $1 LIMIT 1;

-- name: GetTrainDetail :one
SELECT 
    trains.*,
    CASE WHEN collection_trains.train_id IS NULL THEN FALSE ELSE TRUE END AS is_in_collection,
    CASE WHEN wishlist_trains.train_id IS NULL THEN FALSE ELSE TRUE END AS is_in_wishlist
FROM 
    trains 
LEFT JOIN 
    collection_trains ON trains.id = collection_trains.train_id AND collection_trains.user_id = $2
LEFT JOIN 
    wishlist_trains ON trains.id = wishlist_trains.train_id AND wishlist_trains.user_id = $2
WHERE 
    trains.id = $1 
LIMIT 1;

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

-- name: GetTotalTrainCount :one
SELECT COUNT(*) FROM trains;

-- name: UpdateTrainValue :exec
UPDATE trains SET value = $2, version = version + 1
WHERE id = $1;

-- IN FUTURE WHEN UPGRADE TO PGX
-- -- name: UpdateTrainsBatch :batchexec
-- UPDATE trains SET value = $2 WHERE id = $1; 

-- name: UpdateTrainImageUrl :exec
UPDATE trains SET img_url = $2, version = version + 1
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

-- name: SearchTrainSuggestionsByModelNumberWithListStatus :many
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
WHERE model_number ILIKE $2 || '%'
ORDER BY model_number
LIMIT $3 
OFFSET $4;

-- name: SearchTrainSuggestionsByNameWithListStatus :many
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
WHERE name ILIKE $2 || '%'
ORDER BY name
LIMIT $3 
OFFSET $4;

-- name: SearchTrainSuggestionsWithListStatus :many
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
WHERE 
    CASE WHEN $4 = 'model' THEN model_number ELSE name END ILIKE $2 || '%'
ORDER BY 
    CASE WHEN $4 = 'model' THEN model_number ELSE name END
LIMIT $3 
OFFSET $5;

-- name: GetTotalSearchSuggestionsTrainCount :one
SELECT COUNT(*)
FROM trains
WHERE 
CASE WHEN $1 = 'model' THEN model_number ELSE name END ILIKE $2 || '%';

-- name: GetTotalSearchSuggestionsByModelNumberTrainCount :one
SELECT COUNT(*)
FROM trains
WHERE model_number ILIKE $1 || '%';

-- name: GetTotalSearchSuggestionsByNameTrainCount :one
SELECT COUNT(*)
FROM trains
WHERE name ILIKE $1 || '%';

-- name: SearchTrainsByNameSuggestions :many
SELECT DISTINCT id, model_number, name
FROM trains
WHERE name ILIKE $1 || '%'
ORDER BY name
LIMIT $2
OFFSET $3;
