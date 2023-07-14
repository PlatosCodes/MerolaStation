package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTradeTransaction(t *testing.T) TradeTransaction {
	user1, train1 := createRandomUser(t), createRandomTrain(t)
	user2, train2 := createRandomUser(t), createRandomTrain(t)

	arg := CreateTradeTransactionParams{
		OfferedTrain:        train1.ID,
		OfferedTrainOwner:   user1.ID,
		RequestedTrain:      train2.ID,
		RequestedTrainOwner: user2.ID,
	}

	tt, err := testQueries.CreateTradeTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tt)

	require.Equal(t, arg.OfferedTrain, tt.OfferedTrain)
	require.Equal(t, arg.OfferedTrainOwner, tt.OfferedTrainOwner)
	require.Equal(t, arg.RequestedTrain, tt.RequestedTrain)
	require.Equal(t, arg.RequestedTrainOwner, tt.RequestedTrainOwner)

	require.WithinDuration(t, time.Now(), tt.CreatedAt, time.Second)

	return tt
}

func TestCreateTradeTransaction(t *testing.T) {
	createRandomTradeTransaction(t)
}

func TestDeleteTradeTransaction(t *testing.T) {
	tt := createRandomTradeTransaction(t)
	err := testQueries.DeleteTradeTransaction(context.Background(), tt.ID)
	require.NoError(t, err)
}

func TestGetTradeTransaction(t *testing.T) {
	tt1 := createRandomTradeTransaction(t)
	tt2, err := testQueries.GetTradeTransaction(context.Background(), tt1.ID)

	require.NoError(t, err)
	require.Equal(t, tt1.ID, tt2.ID)
	require.Equal(t, tt1.OfferedTrain, tt2.OfferedTrain)
	require.Equal(t, tt1.OfferedTrainOwner, tt2.OfferedTrainOwner)
	require.Equal(t, tt1.RequestedTrain, tt2.RequestedTrain)
	require.Equal(t, tt1.RequestedTrainOwner, tt2.RequestedTrainOwner)
}

func TestListTrainTradeTransactions(t *testing.T) {
	tt := createRandomTradeTransaction(t)

	arg := ListTrainTradeTransactionsParams{
		OfferedTrain:   tt.OfferedTrain,
		RequestedTrain: tt.RequestedTrain,
		Limit:          1,
		Offset:         0,
	}

	gotTT, err := testQueries.ListTrainTradeTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, tt.ID, gotTT[0].ID)
}

func TestListUserTradeTransactions(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		tt := CreateTradeTransactionParams{
			OfferedTrain:        createRandomTrain(t).ID,
			OfferedTrainOwner:   user.ID,
			RequestedTrain:      createRandomTrain(t).ID,
			RequestedTrainOwner: createRandomUser(t).ID,
		}
		_, err := testQueries.CreateTradeTransaction(context.Background(), tt)
		require.NoError(t, err)
	}

	arg := ListUserTradeTransactionsParams{
		OfferedTrainOwner: user.ID,
		Limit:             5,
		Offset:            0,
	}

	tradeTransactions, err := testQueries.ListUserTradeTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tradeTransactions, 5)
}

func TestListTradeTransactions(t *testing.T) {
	// Create a number of trade transactions
	for i := 0; i < 10; i++ {
		createRandomTradeTransaction(t)
	}

	// Test function
	arg := ListTradeTransactionsParams{
		Limit:  5,
		Offset: 0,
	}

	transactions, err := testQueries.ListTradeTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	// Check the transactions are ordered by ID
	for i := 0; i < len(transactions)-1; i++ {
		require.True(t, transactions[i].ID < transactions[i+1].ID)
	}
}
