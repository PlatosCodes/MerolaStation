// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package db

import (
	"context"
)

type Querier interface {
	CreateCollectionTrain(ctx context.Context, arg CreateCollectionTrainParams) (CollectionTrain, error)
	CreateTradeOffer(ctx context.Context, arg CreateTradeOfferParams) (TradeOffer, error)
	CreateTradeTransaction(ctx context.Context, arg CreateTradeTransactionParams) (TradeTransaction, error)
	CreateTrain(ctx context.Context, arg CreateTrainParams) (Train, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateWishlistTrain(ctx context.Context, arg CreateWishlistTrainParams) (WishlistTrain, error)
	DeleteCollectionTrain(ctx context.Context, arg DeleteCollectionTrainParams) error
	DeleteTradeOffer(ctx context.Context, id int64) error
	DeleteTradeTransaction(ctx context.Context, id int64) error
	DeleteTrain(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteWishlistTrain(ctx context.Context, arg DeleteWishlistTrainParams) error
	GetCollectionTrain(ctx context.Context, arg GetCollectionTrainParams) (CollectionTrain, error)
	GetCollectionTrainByID(ctx context.Context, id int64) (CollectionTrain, error)
	GetCollectionTrainforUpdate(ctx context.Context, arg GetCollectionTrainforUpdateParams) (CollectionTrain, error)
	GetCollectionTrainforUpdateByID(ctx context.Context, id int64) (CollectionTrain, error)
	GetTradeOfferByTradeID(ctx context.Context, id int64) (TradeOffer, error)
	GetTradeTransaction(ctx context.Context, id int64) (TradeTransaction, error)
	GetTrain(ctx context.Context, id int64) (Train, error)
	GetTrainByModel(ctx context.Context, modelNumber string) (Train, error)
	GetTrainByName(ctx context.Context, name string) (Train, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListAllUserTradeOffers(ctx context.Context, arg ListAllUserTradeOffersParams) ([]TradeOffer, error)
	ListCollectionTrainTradeOffers(ctx context.Context, arg ListCollectionTrainTradeOffersParams) (TradeOffer, error)
	ListCollectionTrains(ctx context.Context, arg ListCollectionTrainsParams) ([]CollectionTrain, error)
	ListTradeOffers(ctx context.Context, arg ListTradeOffersParams) (TradeOffer, error)
	ListTradeTransactions(ctx context.Context, arg ListTradeTransactionsParams) ([]TradeTransaction, error)
	ListTrainTradeTransactions(ctx context.Context, arg ListTrainTradeTransactionsParams) ([]TradeTransaction, error)
	ListTrains(ctx context.Context, arg ListTrainsParams) ([]Train, error)
	ListUserCollection(ctx context.Context, arg ListUserCollectionParams) ([]CollectionTrain, error)
	ListUserTradeOffers(ctx context.Context, arg ListUserTradeOffersParams) ([]TradeOffer, error)
	ListUserTradeRequests(ctx context.Context, arg ListUserTradeRequestsParams) ([]TradeOffer, error)
	ListUserTradeTransactions(ctx context.Context, arg ListUserTradeTransactionsParams) ([]TradeTransaction, error)
	ListUserWishlist(ctx context.Context, arg ListUserWishlistParams) ([]WishlistTrain, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListWishlists(ctx context.Context, arg ListWishlistsParams) ([]WishlistTrain, error)
	UpdateCollectionTrain(ctx context.Context, arg UpdateCollectionTrainParams) (CollectionTrain, error)
	UpdateTrainValue(ctx context.Context, arg UpdateTrainValueParams) error
}

var _ Querier = (*Queries)(nil)
