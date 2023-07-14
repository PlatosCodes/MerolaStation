package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func createRandomTrain(t *testing.T) Train {
	model_number, name := util.RandomTrainRequest()

	arg := CreateTrainParams{
		ModelNumber: model_number,
		Name:        name,
	}

	train, err := testQueries.CreateTrain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, train)

	require.Equal(t, arg.ModelNumber, train.ModelNumber)
	require.Equal(t, arg.Name, train.Name)

	require.NotZero(t, train.ID)
	require.NotZero(t, train.CreatedAt)
	require.NotZero(t, train.Version)
	require.Zero(t, train.Value)

	return train
}

func TestCreateTrain(t *testing.T) {
	createRandomTrain(t)
}

func TestGetTrain(t *testing.T) {
	train1 := createRandomTrain(t)
	train2, err := testQueries.GetTrain(context.Background(), train1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, train2)

	require.Equal(t, train1.ID, train2.ID)
	require.Equal(t, train1.ModelNumber, train2.ModelNumber)
	require.Equal(t, train1.Name, train2.Name)
	require.Equal(t, train1.Value, train2.Value)
	require.Equal(t, train1.Version, train2.Version)

	require.WithinDuration(t, train1.CreatedAt, train2.CreatedAt, time.Second)
}

func TestGetTrainByModel(t *testing.T) {
	train1 := createRandomTrain(t)
	train2, err := testQueries.GetTrainByModel(context.Background(), train1.ModelNumber)
	require.NoError(t, err)
	require.NotEmpty(t, train2)

	require.Equal(t, train1.ID, train2.ID)
	require.Equal(t, train1.ModelNumber, train2.ModelNumber)
	require.Equal(t, train1.Name, train2.Name)
	require.Equal(t, train1.Value, train2.Value)
	require.Equal(t, train1.Version, train2.Version)

	require.WithinDuration(t, train1.CreatedAt, train2.CreatedAt, time.Second)
}

func TestGetTrainByName(t *testing.T) {
	train1 := createRandomTrain(t)
	train2, err := testQueries.GetTrainByName(context.Background(), train1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, train2)

	require.Equal(t, train1.ID, train2.ID)
	require.Equal(t, train1.ModelNumber, train2.ModelNumber)
	require.Equal(t, train1.Name, train2.Name)
	require.Equal(t, train1.Value, train2.Value)
	require.Equal(t, train1.Version, train2.Version)

	require.WithinDuration(t, train1.CreatedAt, train2.CreatedAt, time.Second)
}

func TestListTrains(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTrain(t)
	}

	arg := ListTrainsParams{
		Limit:  5,
		Offset: 5,
	}

	trains, err := testQueries.ListTrains(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, trains, 5)

	for _, train := range trains {
		require.NotEmpty(t, train)
	}
}

func TestUpdateTrainValue(t *testing.T) {
	train1 := createRandomTrain(t)
	arg := UpdateTrainValueParams{
		ID:    train1.ID,
		Value: util.RandomInt(100, 200),
	}

	err := testQueries.UpdateTrainValue(context.Background(), arg)
	require.NoError(t, err)

	train2, err := testQueries.GetTrain(context.Background(), train1.ID)
	require.NoError(t, err)

	require.Equal(t, arg.Value, train2.Value)
	require.Equal(t, int64(2), train2.Version)
}

func TestDeleteTrain(t *testing.T) {
	train1 := createRandomTrain(t)
	err := testQueries.DeleteTrain(context.Background(), train1.ID)
	require.NoError(t, err)

	train2, err := testQueries.GetTrain(context.Background(), train1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, train2)
}
