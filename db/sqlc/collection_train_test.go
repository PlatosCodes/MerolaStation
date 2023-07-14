package db

import (
	"context"
	"testing"
	"time"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func createRandomCollectionTrain(t *testing.T) (User, Train) {
	user := createRandomUser(t)
	train := createRandomTrain(t)

	ct, err := testQueries.CreateCollectionTrain(context.Background(), CreateCollectionTrainParams{
		UserID:  user.ID,
		TrainID: train.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, ct.UserID)

	require.Equal(t, user.ID, ct.UserID)
	require.Equal(t, train.ID, ct.TrainID)

	require.WithinDuration(t, user.CreatedAt, ct.CreatedAt, time.Second)

	return user, train
}

func TestCreateCollectionTrain(t *testing.T) {
	createRandomCollectionTrain(t)
}

func TestGetUserCollectionTrains(t *testing.T) {
	user, train := createRandomCollectionTrain(t)

	arg := GetUserCollectionTrainsParams{
		UserID: user.ID,
		Limit:  1,
		Offset: 0,
	}

	wishlistTrains, err := testQueries.GetUserCollectionTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, wishlistTrains, 1)

	gotCollectionTrain := wishlistTrains[0]
	require.Equal(t, user.ID, gotCollectionTrain.UserID)
	require.Equal(t, train.ID, gotCollectionTrain.TrainID)
}

func TestDeleteCollectionTrain(t *testing.T) {
	user, train := createRandomCollectionTrain(t)
	page_id := int32(util.RandomInt(1, 100))
	page_size := int32(util.RandomInt(1, 1000))
	arg := GetUserCollectionTrainsParams{
		UserID: user.ID,
		Limit:  page_id,
		Offset: page_size,
	}
	err := testQueries.DeleteCollectionTrain(context.Background(), DeleteCollectionTrainParams{
		UserID:  user.ID,
		TrainID: train.ID,
	})
	require.NoError(t, err)

	collectionTrains, err := testQueries.GetUserCollectionTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, collectionTrains)
}

func TestListCollectionTrains(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCollectionTrain(t)
	}

	arg := ListCollectionTrainsParams{
		Limit:  5,
		Offset: 5,
	}

	collectionTrains, err := testQueries.ListCollectionTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, collectionTrains, 5)

	for _, collectionTrain := range collectionTrains {
		require.NotEmpty(t, collectionTrain)
	}
}
