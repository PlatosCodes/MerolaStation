 -- name: CreateTradeOffer :one
INSERT INTO trade_offers (
    offered_train, offered_train_owner, 
    requested_train, requested_train_owner
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTradeOfferByTradeID :one
SELECT * FROM trade_offers
WHERE id = $1 LIMIT 1;

-- name: ListUserTradeOffers :many
SELECT * FROM trade_offers
WHERE offered_train_owner = $1 
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListUserTradeRequests :many
SELECT * FROM trade_offers
WHERE requested_train_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListAllUserTradeOffers :many
SELECT * FROM trade_offers
WHERE offered_train_owner = $1 
OR requested_train_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListCollectionTrainTradeOffers :one
SELECT * FROM trade_offers
WHERE requested_train = $1 
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListTradeOffers :one
SELECT * FROM trade_offers
WHERE offered_train = $1 
ORDER BY id
LIMIT $2
OFFSET $3;


-- name: DeleteTradeOffer :exec
DELETE from trade_offers WHERE id = $1;

