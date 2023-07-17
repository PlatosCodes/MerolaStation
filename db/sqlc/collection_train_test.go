package db

import (
	"context"
	"testing"
	"time"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func createRandomCollectionTrain(t *testing.T) CollectionTrain {
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

	return ct
}

func TestCreateCollectionTrain(t *testing.T) {
	createRandomCollectionTrain(t)
}

func TestGetUserCollectionTrains(t *testing.T) {
	ct := createRandomCollectionTrain(t)

	arg := ListUserCollectionParams{
		UserID: ct.UserID,
		Limit:  1,
		Offset: 0,
	}

	collectionTrains, err := testQueries.ListUserCollection(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, collectionTrains, 1)

	gotCollectionTrain := collectionTrains[0]
	require.Equal(t, ct.UserID, gotCollectionTrain.UserID)
	require.Equal(t, ct.TrainID, gotCollectionTrain.TrainID)
}

func TestGetCollectionTrainByID(t *testing.T) {
	ct := createRandomCollectionTrain(t)

	collectionTrain, err := testQueries.GetCollectionTrainByID(context.Background(), ct.ID)
	require.NoError(t, err)

	require.Equal(t, ct.UserID, collectionTrain.UserID)
	require.Equal(t, ct.TrainID, collectionTrain.TrainID)
}

func TestGetCollectionTrainForUpdate(t *testing.T) {
	collectionTrain := createRandomCollectionTrain(t)

	arg := GetCollectionTrainforUpdateParams{
		UserID:  collectionTrain.UserID,
		TrainID: collectionTrain.TrainID,
	}

	ct, err := testQueries.GetCollectionTrainforUpdate(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, collectionTrain.UserID, ct.UserID)
	require.Equal(t, collectionTrain.TrainID, ct.TrainID)
}

func TestDeleteCollectionTrain(t *testing.T) {
	ct := createRandomCollectionTrain(t)
	page_id := int32(util.RandomInt(1, 100))
	page_size := int32(util.RandomInt(1, 1000))
	arg := ListUserCollectionParams{
		UserID: ct.UserID,
		Limit:  page_id,
		Offset: page_size,
	}
	err := testQueries.DeleteCollectionTrain(context.Background(), DeleteCollectionTrainParams{
		UserID:  ct.UserID,
		TrainID: ct.TrainID,
	})
	require.NoError(t, err)

	collectionTrains, err := testQueries.ListUserCollection(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, collectionTrains)
}

func TestListUserCollection(t *testing.T) {
	//create a single user and 10 random trains
	testUser := createRandomUser(t)
	for i := 0; i < 10; i++ {
		testTrain := createRandomTrain(t)
		testCT := CreateCollectionTrainParams{
			UserID:  testUser.ID,
			TrainID: testTrain.ID,
		}
		testQueries.CreateCollectionTrain(context.Background(), testCT)
	}

	arg := ListUserCollectionParams{
		UserID: testUser.ID,
		Limit:  5,
		Offset: 5,
	}

	collectionTrains, err := testQueries.ListUserCollection(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, collectionTrains, 5)

	for _, collectionTrain := range collectionTrains {
		require.NotEmpty(t, collectionTrain)
	}
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

func TestUpdateCollectionTrain(t *testing.T) {
	ct := createRandomCollectionTrain(t)
	user2 := createRandomUser(t)
	arg := UpdateCollectionTrainParams{
		ID:     ct.ID,
		UserID: user2.ID,
	}

	ct, err := testQueries.UpdateCollectionTrain(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, user2.ID, ct.UserID)
	require.Equal(t, ct.TrainID, ct.TrainID)
	require.WithinDuration(t, ct.CreatedAt, ct.CreatedAt, time.Second)

	train2, err := testQueries.GetCollectionTrain(context.Background(), GetCollectionTrainParams{
		UserID:  user2.ID,
		TrainID: ct.TrainID,
	})
	require.NoError(t, err)

	require.Equal(t, user2.ID, train2.UserID)
	require.Equal(t, ct.TrainID, train2.TrainID)
	require.WithinDuration(t, ct.CreatedAt, train2.CreatedAt, time.Second)
}
