package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Storage provides all functions to execure db queries and transactions using
// composition and extending the functionality of queries for single db operations
type Storage struct {
	*Queries
	db *sql.DB
}

// NewStorage creates a new Storage
func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (storage *Storage) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// In future will add more params to extend functionality of RegisterTx
// type RegisterTxParams struct {
// 	Username       string `json:"username"`
// 	HashedPassword string `json:"hashed_password"`
// 	Email          string `json:"email"`
// **Other fields here**
// }

type RegisterTxResult struct {
	User User `json:"user"`
}

// RegisterTx performs a new user registration.
// It creates a new user only, so there is no reason to actually use
// this besides getting practice for now, and adding new
// multi-operation database transaction features later
// **RegisterTxResult is also rather useless for now, but will be useful when
// we have actual transcations occuring.
func (storage *Storage) RegisterTx(ctx context.Context, arg CreateUserParams) (RegisterTxResult, error) {
	var result RegisterTxResult

	err := storage.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg)

		if err != nil {
			return err
		}

		// future additonal operations here
		// ex. result.StoreProfilePhoto, err = q.StoreProfilePhoto(...)..

		return nil
	})

	return result, err

}
