package db

import (
	"context"
	"testing"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func createRandomWishlistTrain(t *testing.T) (User, Train) {
	train := createRandomTrain(t)
	user := createRandomUser(t)

	_, err := testQueries.CreateWishlistTrain(context.Background(), CreateWishlistTrainParams{
		UserID:  user.ID,
		TrainID: train.ID,
	})
	require.NoError(t, err)

	return user, train
}

func TestCreateWishlistTrain(t *testing.T) {
	createRandomWishlistTrain(t)
}

func TestDeleteWishlistTrain(t *testing.T) {
	user, train := createRandomWishlistTrain(t)
	page_id := int32(util.RandomInt(1, 100))
	page_size := int32(util.RandomInt(1, 1000))
	arg := GetUserWishlistTrainsParams{
		UserID: user.ID,
		Limit:  page_id,
		Offset: page_size,
	}
	err := testQueries.DeleteWishlistTrain(context.Background(), DeleteWishlistTrainParams{
		UserID:  user.ID,
		TrainID: train.ID,
	})
	require.NoError(t, err)

	wishlistTrains, err := testQueries.GetUserWishlistTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, wishlistTrains)
}

func TestGetUserWishlistTrains(t *testing.T) {
	user, train := createRandomWishlistTrain(t)

	arg := GetUserWishlistTrainsParams{
		UserID: user.ID,
		Limit:  1,
		Offset: 0,
	}

	wishlistTrains, err := testQueries.GetUserWishlistTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, wishlistTrains, 1)

	gotWishlistTrain := wishlistTrains[0]
	require.Equal(t, user.ID, gotWishlistTrain.UserID)
	require.Equal(t, train.ID, gotWishlistTrain.TrainID)
}

func TestListWishlistTrains(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomWishlistTrain(t)
	}

	arg := ListWishlistTrainsParams{
		Limit:  5,
		Offset: 5,
	}

	wishlistTrains, err := testQueries.ListWishlistTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, wishlistTrains, 5)

	for _, wishlistTrain := range wishlistTrains {
		require.NotEmpty(t, wishlistTrain)
	}
}
