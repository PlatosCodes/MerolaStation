 -- name: CreateTradeTransaction :one
INSERT INTO trade_transactions (
    offered_train, offered_train_owner, 
    requested_train, requested_train_owner
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTradeTransaction :one
SELECT * FROM trade_transactions
WHERE id = $1 LIMIT 1;

-- name: ListUserTradeTransactions :many
SELECT * FROM trade_transactions
WHERE offered_train_owner = $1 
OR requested_train_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListTrainTradeTransactions :many
SELECT * FROM trade_transactions
WHERE offered_train = $1 
OR requested_train = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListTradeTransactions :many
SELECT * FROM trade_transactions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteTradeTransaction :exec
DELETE from trade_transactions WHERE id = $1;

