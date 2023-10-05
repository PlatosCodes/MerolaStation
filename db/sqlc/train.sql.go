// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: train.sql

package db

import (
	"context"
	"database/sql"
)

const createTrain = `-- name: CreateTrain :one
INSERT INTO trains (
  model_number,
  name
) VALUES (
  $1, $2
) RETURNING id, model_number, name, value, created_at, version, last_edited_at
`

type CreateTrainParams struct {
	ModelNumber string `json:"model_number"`
	Name        string `json:"name"`
}

func (q *Queries) CreateTrain(ctx context.Context, arg CreateTrainParams) (Train, error) {
	row := q.db.QueryRowContext(ctx, createTrain, arg.ModelNumber, arg.Name)
	var i Train
	err := row.Scan(
		&i.ID,
		&i.ModelNumber,
		&i.Name,
		&i.Value,
		&i.CreatedAt,
		&i.Version,
		&i.LastEditedAt,
	)
	return i, err
}

const deleteTrain = `-- name: DeleteTrain :exec
DELETE from trains WHERE id = $1
`

func (q *Queries) DeleteTrain(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTrain, id)
	return err
}

const getTrain = `-- name: GetTrain :one
SELECT id, model_number, name, value, created_at, version, last_edited_at FROM trains
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTrain(ctx context.Context, id int64) (Train, error) {
	row := q.db.QueryRowContext(ctx, getTrain, id)
	var i Train
	err := row.Scan(
		&i.ID,
		&i.ModelNumber,
		&i.Name,
		&i.Value,
		&i.CreatedAt,
		&i.Version,
		&i.LastEditedAt,
	)
	return i, err
}

const getTrainByModel = `-- name: GetTrainByModel :one
SELECT id, model_number, name, value, created_at, version, last_edited_at FROM trains
WHERE model_number = $1 LIMIT 1
`

func (q *Queries) GetTrainByModel(ctx context.Context, modelNumber string) (Train, error) {
	row := q.db.QueryRowContext(ctx, getTrainByModel, modelNumber)
	var i Train
	err := row.Scan(
		&i.ID,
		&i.ModelNumber,
		&i.Name,
		&i.Value,
		&i.CreatedAt,
		&i.Version,
		&i.LastEditedAt,
	)
	return i, err
}

const getTrainByName = `-- name: GetTrainByName :one
SELECT id, model_number, name, value, created_at, version, last_edited_at FROM trains
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetTrainByName(ctx context.Context, name string) (Train, error) {
	row := q.db.QueryRowContext(ctx, getTrainByName, name)
	var i Train
	err := row.Scan(
		&i.ID,
		&i.ModelNumber,
		&i.Name,
		&i.Value,
		&i.CreatedAt,
		&i.Version,
		&i.LastEditedAt,
	)
	return i, err
}

const listTrains = `-- name: ListTrains :many
SELECT id, model_number, name, value, created_at, version, last_edited_at FROM trains
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListTrainsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTrains(ctx context.Context, arg ListTrainsParams) ([]Train, error) {
	rows, err := q.db.QueryContext(ctx, listTrains, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Train{}
	for rows.Next() {
		var i Train
		if err := rows.Scan(
			&i.ID,
			&i.ModelNumber,
			&i.Name,
			&i.Value,
			&i.CreatedAt,
			&i.Version,
			&i.LastEditedAt,
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

const searchTrainsByModelNumberSuggestions = `-- name: SearchTrainsByModelNumberSuggestions :many
SELECT DISTINCT id, model_number, name
FROM trains
WHERE model_number ILIKE $1 || '%'
ORDER BY model_number
LIMIT $2
OFFSET $3
`

type SearchTrainsByModelNumberSuggestionsParams struct {
	Column1 sql.NullString `json:"column_1"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
}

type SearchTrainsByModelNumberSuggestionsRow struct {
	ID          int64  `json:"id"`
	ModelNumber string `json:"model_number"`
	Name        string `json:"name"`
}

func (q *Queries) SearchTrainsByModelNumberSuggestions(ctx context.Context, arg SearchTrainsByModelNumberSuggestionsParams) ([]SearchTrainsByModelNumberSuggestionsRow, error) {
	rows, err := q.db.QueryContext(ctx, searchTrainsByModelNumberSuggestions, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchTrainsByModelNumberSuggestionsRow{}
	for rows.Next() {
		var i SearchTrainsByModelNumberSuggestionsRow
		if err := rows.Scan(&i.ID, &i.ModelNumber, &i.Name); err != nil {
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

const updateTrainValue = `-- name: UpdateTrainValue :exec
UPDATE trains SET value = $2, version = version + 1
WHERE id = $1
`

type UpdateTrainValueParams struct {
	ID    int64 `json:"id"`
	Value int64 `json:"value"`
}

func (q *Queries) UpdateTrainValue(ctx context.Context, arg UpdateTrainValueParams) error {
	_, err := q.db.ExecContext(ctx, updateTrainValue, arg.ID, arg.Value)
	return err
}
