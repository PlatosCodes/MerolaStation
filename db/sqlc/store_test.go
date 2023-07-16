package db

import (
	"context"
	"testing"
	"time"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func TestRegisterTx(t *testing.T) {
	store := NewStore(testDB)

	hashedPassword, err := util.HashPassword("secret")
	require.NoError(t, err)

	// run n concurrent user registrations to ensure transaction works well
	n := 5

	errs := make(chan error)
	results := make(chan RegisterTxResult)

	for i := 0; i < 5; i++ {
		go func() {
			result, err := store.RegisterTx(context.Background(), CreateUserParams{
				Username:       "Alex",
				HashedPassword: hashedPassword,
				Email:          util.RandomEmail(),
			})

			errs <- err
			results <- result
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check createdUser. because of unique constraint on emails, need to use
		//random emails and so cannot check them.
		user := result.User
		require.NotEmpty(t, user)

		require.NotZero(t, user.ID)
		require.Equal(t, "Alex", user.Username)
		require.Equal(t, hashedPassword, user.HashedPassword)
		require.NotEmpty(t, user.Email)

		require.WithinDuration(t, user.CreatedAt, time.Now(), time.Second)

		_, err = store.GetUser(context.Background(), user.ID)
		require.NoError(t, err)
	}
}

func TestTradeTx(t *testing.T) {
	store := NewStore(testDB)

	OfferedTrain := createRandomCollectionTrain(t)
	RequestedTrain := createRandomCollectionTrain(t)

	// run n concurrent trade transactions to ensure transaction works well
	n := 5

	errs := make(chan error)
	results := make(chan TradeTxResult)

	for i := 0; i < n; i++ {

		if i%2 == 1 {
			OfferedTrain.UserID = RequestedTrain.UserID
			RequestedTrain.UserID = OfferedTrain.UserID
		}

		go func() {
			result, err := store.TradeTx(context.Background(), TradeTxParams{
				OfferedTrain:   OfferedTrain,
				RequestedTrain: RequestedTrain,
			})

			errs <- err
			results <- result

		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transaction record

		transactionRecord := result.TradeTransaction
		require.NotEmpty(t, transactionRecord)
		require.Equal(t, OfferedTrain.TrainID, transactionRecord.OfferedTrain)
		require.Equal(t, RequestedTrain.TrainID, transactionRecord.RequestedTrain)

		_, err = store.GetTradeTransaction(context.Background(), transactionRecord.ID)
		require.NoError(t, err)

	}

	//check final collection_trains
	updatedCollectionTrain1, err := store.GetCollectionTrain(context.Background(),
		GetCollectionTrainParams{
			UserID:  OfferedTrain.UserID,
			TrainID: OfferedTrain.TrainID,
		})
	require.NoError(t, err)

	updatedCollectionTrain2, err := store.GetCollectionTrain(context.Background(),
		GetCollectionTrainParams{
			UserID:  RequestedTrain.UserID,
			TrainID: RequestedTrain.TrainID,
		})
	require.NoError(t, err)

	require.Equal(t, OfferedTrain.UserID, updatedCollectionTrain1.UserID)
	require.Equal(t, RequestedTrain.UserID, updatedCollectionTrain2.UserID)
	require.Equal(t, OfferedTrain.TrainID, updatedCollectionTrain1.TrainID)
	require.Equal(t, RequestedTrain.TrainID, updatedCollectionTrain2.TrainID)
}

func TestTradeTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	OfferedTrain := createRandomCollectionTrain(t)
	RequestedTrain := createRandomCollectionTrain(t)

	// run n concurrent trade transactions to ensure transaction works well
	n := 10

	errs := make(chan error)
	results := make(chan TradeTxResult)

	for i := 0; i < n; i++ {

		if i%2 == 1 {
			OfferedTrain.UserID = RequestedTrain.UserID
			RequestedTrain.UserID = OfferedTrain.UserID
		}

		go func() {
			result, err := store.TradeTx(context.Background(), TradeTxParams{
				OfferedTrain:   OfferedTrain,
				RequestedTrain: RequestedTrain,
			})

			errs <- err
			results <- result

		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
	}

	//check final collection_trains
	updatedCollectionTrain1, err := store.GetCollectionTrain(context.Background(),
		GetCollectionTrainParams{
			UserID:  OfferedTrain.UserID,
			TrainID: OfferedTrain.TrainID,
		})
	require.NoError(t, err)

	updatedCollectionTrain2, err := store.GetCollectionTrain(context.Background(),
		GetCollectionTrainParams{
			UserID:  RequestedTrain.UserID,
			TrainID: RequestedTrain.TrainID,
		})
	require.NoError(t, err)

	require.Equal(t, OfferedTrain.UserID, updatedCollectionTrain1.UserID)
	require.Equal(t, RequestedTrain.UserID, updatedCollectionTrain2.UserID)
	require.Equal(t, OfferedTrain.TrainID, updatedCollectionTrain1.TrainID)
	require.Equal(t, RequestedTrain.TrainID, updatedCollectionTrain2.TrainID)
}
