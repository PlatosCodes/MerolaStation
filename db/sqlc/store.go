package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execure db queries and transactions
// Uses composition and extending the functionality of queries for single db operations
type Store interface {
	Querier
	RegisterTx(ctx context.Context, arg CreateUserParams) (RegisterTxResult, error)
	TradeTx(ctx context.Context, arg TradeTxParams) (TradeTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
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
func (store *SQLStore) RegisterTx(ctx context.Context, arg CreateUserParams) (RegisterTxResult, error) {
	var result RegisterTxResult

	err := store.execTx(ctx, func(q *Queries) error {
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

type TradeTxParams struct {
	TradeOfferID   int64           `json:"trade_offer_id"`
	OfferedTrain   CollectionTrain `json:"offered_train"`
	RequestedTrain CollectionTrain `json:"requested_train"`
}

// Only returns the requested train as the offered train is now
// a collection train belonging to someone else, so up to them to
// keep it in their collection as private vs public
type TradeTxResult struct {
	TradeTransaction TradeTransaction `json:"trade_transaction"`
	OfferedTrain     CollectionTrain  `json:"offered_train"`
	RequestedTrain   CollectionTrain  `json:"requested_train"`
}

// TradeTx performs a train trade between users
// It creates a trade record, verifies the collection trains of the two users,
// and updates the train owner for each collection train.
func (store *SQLStore) TradeTx(ctx context.Context, arg TradeTxParams) (TradeTxResult, error) {
	var result TradeTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create trade record
		result.TradeTransaction, err = q.CreateTradeTransaction(ctx, CreateTradeTransactionParams{
			OfferedTrain:        arg.OfferedTrain.TrainID,
			OfferedTrainOwner:   arg.OfferedTrain.UserID,
			RequestedTrain:      arg.RequestedTrain.TrainID,
			RequestedTrainOwner: arg.RequestedTrain.UserID,
		})
		if err != nil {
			return err
		}

		// Need to check to make sure both collection trains exist (equivalent to getAccount)
		ct1, err := q.GetCollectionTrainforUpdateByID(ctx, arg.OfferedTrain.ID)
		if err != nil {
			return err
		}

		ct2, err := q.GetCollectionTrainforUpdateByID(ctx, arg.RequestedTrain.ID)
		if err != nil {
			return err
		}

		// Ensure train collections are always updated in a consistent order to avoid deadlock scenario
		// It is pretty artificial as it requires the immediate reversal of a trade to ever be possible,
		// but it was good to understand and have a reference in place for future scenarios.

		if arg.OfferedTrain.ID < arg.RequestedTrain.ID {
			result.OfferedTrain, err = q.UpdateCollectionTrain(ctx, UpdateCollectionTrainParams{
				ID:     ct1.ID,
				UserID: arg.RequestedTrain.UserID,
			})
			if err != nil {
				return err
			}

			result.RequestedTrain, err = q.UpdateCollectionTrain(ctx, UpdateCollectionTrainParams{
				ID:     ct2.ID,
				UserID: arg.OfferedTrain.UserID,
			})
			if err != nil {

				return err
			}
		} else {

			result.RequestedTrain, err = q.UpdateCollectionTrain(ctx, UpdateCollectionTrainParams{
				ID:     ct2.ID,
				UserID: arg.OfferedTrain.UserID,
			})
			if err != nil {

				return err
			}

			result.OfferedTrain, err = q.UpdateCollectionTrain(ctx, UpdateCollectionTrainParams{
				ID:     ct1.ID,
				UserID: arg.RequestedTrain.UserID,
			})
			if err != nil {
				return err
			}

		}

		// Deletes Trade Offer after trade completion to avoid ability to trigger trade again if
		//  one of OG owners adds the previously traded train back into their collection
		// (buying or trading for a new one) and forgets to delete this old trade themselves
		err = q.DeleteTradeOffer(ctx, arg.TradeOfferID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err

}
