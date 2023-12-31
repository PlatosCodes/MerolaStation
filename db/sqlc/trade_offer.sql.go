// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: trade_offer.sql

package db

import (
	"context"
)

const createTradeOffer = `-- name: CreateTradeOffer :one
INSERT INTO trade_offers (
    offered_train, offered_train_owner, 
    requested_train, requested_train_owner
) VALUES (
  $1, $2, $3, $4
) RETURNING id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at
`

type CreateTradeOfferParams struct {
	OfferedTrain        int64 `json:"offered_train"`
	OfferedTrainOwner   int64 `json:"offered_train_owner"`
	RequestedTrain      int64 `json:"requested_train"`
	RequestedTrainOwner int64 `json:"requested_train_owner"`
}

func (q *Queries) CreateTradeOffer(ctx context.Context, arg CreateTradeOfferParams) (TradeOffer, error) {
	row := q.db.QueryRowContext(ctx, createTradeOffer,
		arg.OfferedTrain,
		arg.OfferedTrainOwner,
		arg.RequestedTrain,
		arg.RequestedTrainOwner,
	)
	var i TradeOffer
	err := row.Scan(
		&i.ID,
		&i.OfferedTrain,
		&i.OfferedTrainOwner,
		&i.RequestedTrain,
		&i.RequestedTrainOwner,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTradeOffer = `-- name: DeleteTradeOffer :exec
DELETE from trade_offers WHERE id = $1
`

func (q *Queries) DeleteTradeOffer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTradeOffer, id)
	return err
}

const getTradeOfferByTradeID = `-- name: GetTradeOfferByTradeID :one
SELECT id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at FROM trade_offers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTradeOfferByTradeID(ctx context.Context, id int64) (TradeOffer, error) {
	row := q.db.QueryRowContext(ctx, getTradeOfferByTradeID, id)
	var i TradeOffer
	err := row.Scan(
		&i.ID,
		&i.OfferedTrain,
		&i.OfferedTrainOwner,
		&i.RequestedTrain,
		&i.RequestedTrainOwner,
		&i.CreatedAt,
	)
	return i, err
}

const listAllUserTradeOffers = `-- name: ListAllUserTradeOffers :many
SELECT id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at FROM trade_offers
WHERE offered_train_owner = $1 
OR requested_train_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListAllUserTradeOffersParams struct {
	OfferedTrainOwner int64 `json:"offered_train_owner"`
	Limit             int32 `json:"limit"`
	Offset            int32 `json:"offset"`
}

func (q *Queries) ListAllUserTradeOffers(ctx context.Context, arg ListAllUserTradeOffersParams) ([]TradeOffer, error) {
	rows, err := q.db.QueryContext(ctx, listAllUserTradeOffers, arg.OfferedTrainOwner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TradeOffer{}
	for rows.Next() {
		var i TradeOffer
		if err := rows.Scan(
			&i.ID,
			&i.OfferedTrain,
			&i.OfferedTrainOwner,
			&i.RequestedTrain,
			&i.RequestedTrainOwner,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCollectionTrainTradeOffers = `-- name: ListCollectionTrainTradeOffers :one
SELECT id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at FROM trade_offers
WHERE requested_train = $1 
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListCollectionTrainTradeOffersParams struct {
	RequestedTrain int64 `json:"requested_train"`
	Limit          int32 `json:"limit"`
	Offset         int32 `json:"offset"`
}

func (q *Queries) ListCollectionTrainTradeOffers(ctx context.Context, arg ListCollectionTrainTradeOffersParams) (TradeOffer, error) {
	row := q.db.QueryRowContext(ctx, listCollectionTrainTradeOffers, arg.RequestedTrain, arg.Limit, arg.Offset)
	var i TradeOffer
	err := row.Scan(
		&i.ID,
		&i.OfferedTrain,
		&i.OfferedTrainOwner,
		&i.RequestedTrain,
		&i.RequestedTrainOwner,
		&i.CreatedAt,
	)
	return i, err
}

const listTradeOffers = `-- name: ListTradeOffers :one
SELECT id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at FROM trade_offers
WHERE offered_train = $1 
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTradeOffersParams struct {
	OfferedTrain int64 `json:"offered_train"`
	Limit        int32 `json:"limit"`
	Offset       int32 `json:"offset"`
}

func (q *Queries) ListTradeOffers(ctx context.Context, arg ListTradeOffersParams) (TradeOffer, error) {
	row := q.db.QueryRowContext(ctx, listTradeOffers, arg.OfferedTrain, arg.Limit, arg.Offset)
	var i TradeOffer
	err := row.Scan(
		&i.ID,
		&i.OfferedTrain,
		&i.OfferedTrainOwner,
		&i.RequestedTrain,
		&i.RequestedTrainOwner,
		&i.CreatedAt,
	)
	return i, err
}

const listUserTradeOffers = `-- name: ListUserTradeOffers :many
SELECT id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at FROM trade_offers
WHERE offered_train_owner = $1 
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListUserTradeOffersParams struct {
	OfferedTrainOwner int64 `json:"offered_train_owner"`
	Limit             int32 `json:"limit"`
	Offset            int32 `json:"offset"`
}

func (q *Queries) ListUserTradeOffers(ctx context.Context, arg ListUserTradeOffersParams) ([]TradeOffer, error) {
	rows, err := q.db.QueryContext(ctx, listUserTradeOffers, arg.OfferedTrainOwner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TradeOffer{}
	for rows.Next() {
		var i TradeOffer
		if err := rows.Scan(
			&i.ID,
			&i.OfferedTrain,
			&i.OfferedTrainOwner,
			&i.RequestedTrain,
			&i.RequestedTrainOwner,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUserTradeRequests = `-- name: ListUserTradeRequests :many
SELECT id, offered_train, offered_train_owner, requested_train, requested_train_owner, created_at FROM trade_offers
WHERE requested_train_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListUserTradeRequestsParams struct {
	RequestedTrainOwner int64 `json:"requested_train_owner"`
	Limit               int32 `json:"limit"`
	Offset              int32 `json:"offset"`
}

func (q *Queries) ListUserTradeRequests(ctx context.Context, arg ListUserTradeRequestsParams) ([]TradeOffer, error) {
	rows, err := q.db.QueryContext(ctx, listUserTradeRequests, arg.RequestedTrainOwner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TradeOffer{}
	for rows.Next() {
		var i TradeOffer
		if err := rows.Scan(
			&i.ID,
			&i.OfferedTrain,
			&i.OfferedTrainOwner,
			&i.RequestedTrain,
			&i.RequestedTrainOwner,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
