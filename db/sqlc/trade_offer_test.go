package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTradeOffer(t *testing.T) TradeOffer {
	user1, train1 := createRandomUser(t), createRandomTrain(t)
	user2, train2 := createRandomUser(t), createRandomTrain(t)

	arg := CreateTradeOfferParams{
		OfferedTrain:        train1.ID,
		OfferedTrainOwner:   user1.ID,
		RequestedTrain:      train2.ID,
		RequestedTrainOwner: user2.ID,
	}

	to, err := testQueries.CreateTradeOffer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, to)

	require.Equal(t, arg.OfferedTrain, to.OfferedTrain)
	require.Equal(t, arg.OfferedTrainOwner, to.OfferedTrainOwner)
	require.Equal(t, arg.RequestedTrain, to.RequestedTrain)
	require.Equal(t, arg.RequestedTrainOwner, to.RequestedTrainOwner)

	require.WithinDuration(t, time.Now(), to.CreatedAt, time.Second)

	return to
}

func TestCreateTradeOffer(t *testing.T) {
	createRandomTradeOffer(t)
}

func TestDeleteTradeOffer(t *testing.T) {
	to := createRandomTradeOffer(t)
	err := testQueries.DeleteTradeOffer(context.Background(), to.ID)
	require.NoError(t, err)
}

func TestGetTradeOfferByTradeID(t *testing.T) {
	to1 := createRandomTradeOffer(t)
	to2, err := testQueries.GetTradeOfferByTradeID(context.Background(), to1.ID)

	require.NoError(t, err)
	require.Equal(t, to1.ID, to2.ID)
	require.Equal(t, to1.OfferedTrain, to2.OfferedTrain)
	require.Equal(t, to1.OfferedTrainOwner, to2.OfferedTrainOwner)
	require.Equal(t, to1.RequestedTrain, to2.RequestedTrain)
	require.Equal(t, to1.RequestedTrainOwner, to2.RequestedTrainOwner)
}

func TestListCollectionTrainTradeOffers(t *testing.T) {
	to := createRandomTradeOffer(t)

	arg := ListCollectionTrainTradeOffersParams{
		RequestedTrain: to.RequestedTrain,
		Limit:          1,
		Offset:         0,
	}

	gotTo, err := testQueries.ListCollectionTrainTradeOffers(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, to.ID, gotTo.ID)
}

func TestListUserTradeOffers(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		to := CreateTradeOfferParams{
			OfferedTrain:        createRandomTrain(t).ID,
			OfferedTrainOwner:   user.ID,
			RequestedTrain:      createRandomTrain(t).ID,
			RequestedTrainOwner: createRandomUser(t).ID,
		}
		_, err := testQueries.CreateTradeOffer(context.Background(), to)
		require.NoError(t, err)
	}

	arg := ListUserTradeOffersParams{
		OfferedTrainOwner: user.ID,
		Limit:             5,
		Offset:            0,
	}

	tradeOffers, err := testQueries.ListUserTradeOffers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tradeOffers, 5)
}

func TestListAllUserTradeOffers(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	// Create a number of trade offers
	for i := 0; i < 10; i++ {
		to := CreateTradeOfferParams{
			OfferedTrain:        createRandomTrain(t).ID,
			OfferedTrainOwner:   user1.ID,
			RequestedTrain:      createRandomTrain(t).ID,
			RequestedTrainOwner: user2.ID,
		}
		_, err := testQueries.CreateTradeOffer(context.Background(), to)
		require.NoError(t, err)
	}

	// Test function
	arg := ListAllUserTradeOffersParams{
		OfferedTrainOwner:   user1.ID,
		RequestedTrainOwner: user2.ID,
		Limit:               5,
		Offset:              0,
	}

	offers, err := testQueries.ListAllUserTradeOffers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, offers, 5)

	// Check the offers are ordered by ID
	for i := 0; i < len(offers)-1; i++ {
		require.True(t, offers[i].ID < offers[i+1].ID)
	}
}
