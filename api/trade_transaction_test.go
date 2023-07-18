package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/PlatosCodes/MerolaStation/db/mock"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTradeTransaction(t *testing.T) {

	to, ct1, ct2, username, _ := createTestTradeOffer(t)

	testCases := []struct {
		name          string
		requestBody   map[string]interface{}
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ct1_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.OfferedTrainOwner,
					TrainID: to.OfferedTrain,
				}
				ct2_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.RequestedTrainOwner,
					TrainID: to.RequestedTrain,
				}

				tradeParams := db.TradeTxParams{
					TradeOfferID:   to.ID,
					OfferedTrain:   ct1,
					RequestedTrain: ct2,
				}
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
					Times(1).
					Return(ct1, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct2_req)).
					Times(1).
					Return(ct2, nil)
				store.EXPECT().
					TradeTx(gomock.Any(), gomock.Eq(tradeParams)).
					Times(1).
					Return(db.TradeTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
			},
		},
		{
			name: "InternalServerError",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ct1_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.OfferedTrainOwner,
					TrainID: to.OfferedTrain,
				}
				ct2_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.RequestedTrainOwner,
					TrainID: to.RequestedTrain,
				}

				tradeParams := db.TradeTxParams{
					TradeOfferID:   to.ID,
					OfferedTrain:   ct1,
					RequestedTrain: ct2,
				}
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
					Times(1).
					Return(ct1, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct2_req)).
					Times(1).
					Return(ct2, nil)
				store.EXPECT().
					TradeTx(gomock.Any(), gomock.Eq(tradeParams)).
					Times(1).
					Return(db.TradeTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "GetTradeOfferByIDInternalError",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// ct1_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.OfferedTrainOwner,
				// 	TrainID: to.OfferedTrain,
				// }
				// ct2_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.RequestedTrainOwner,
				// 	TrainID: to.RequestedTrain,
				// }

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(db.TradeOffer{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// ct1_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.OfferedTrainOwner,
				// 	TrainID: to.OfferedTrain,
				// }
				// ct2_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.RequestedTrainOwner,
				// 	TrainID: to.RequestedTrain,
				// }

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(db.TradeOffer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "CollectionTrain1InternalError",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ct1_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.OfferedTrainOwner,
					TrainID: to.OfferedTrain,
				}
				// ct2_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.RequestedTrainOwner,
				// 	TrainID: to.RequestedTrain,
				// }

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
					Times(1).
					Return(db.CollectionTrain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "CollectionTrain1NotFound",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ct1_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.OfferedTrainOwner,
					TrainID: to.OfferedTrain,
				}
				// ct2_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.RequestedTrainOwner,
				// 	TrainID: to.RequestedTrain,
				// }

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
					Times(1).
					Return(db.CollectionTrain{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "CollectionTrain2InternalError",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ct1_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.OfferedTrainOwner,
					TrainID: to.OfferedTrain,
				}
				ct2_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.RequestedTrainOwner,
					TrainID: to.RequestedTrain,
				}

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
					Times(1).
					Return(ct1, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct2_req)).
					Times(1).
					Return(db.CollectionTrain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "CollectionTrain2NotFound",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ct1_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.OfferedTrainOwner,
					TrainID: to.OfferedTrain,
				}
				ct2_req := db.GetCollectionTrainforUpdateParams{
					UserID:  to.RequestedTrainOwner,
					TrainID: to.RequestedTrain,
				}

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
					Times(1).
					Return(ct1, nil)
				store.EXPECT().
					GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct2_req)).
					Times(1).
					Return(db.CollectionTrain{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NotAuthorized",
			requestBody: map[string]interface{}{
				"id": to.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner+1, "TestName", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// ct1_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.OfferedTrainOwner,
				// 	TrainID: to.OfferedTrain,
				// }
				// ct2_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.RequestedTrainOwner,
				// 	TrainID: to.RequestedTrain,
				// }

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(db.TradeOffer{}, nil)
				// store.EXPECT().
				// 	GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
				// 	Times(1).
				// 	Return(ct1, nil)
				// store.EXPECT().
				// 	GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct2_req)).
				// 	Times(1).
				// 	Return(ct2, nil)
				// store.EXPECT().
				// 	TradeTx(gomock.Any(), gomock.Eq(tradeParams)).
				// 	Times(1).
				// 	Return(db.TradeTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidRequest",
			requestBody: map[string]interface{}{
				"id": 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// ct1_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.OfferedTrainOwner,
				// 	TrainID: 0,
				// }
				// ct2_req := db.GetCollectionTrainforUpdateParams{
				// 	UserID:  to.RequestedTrainOwner,
				// 	TrainID: to.RequestedTrain,
				// }

				// tradeParams := db.TradeTxParams{
				// 	TradeOfferID:   to.ID,
				// 	OfferedTrain:   ct1,
				// 	RequestedTrain: ct2,
				// }
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(0)).
					Times(0)
				// store.EXPECT().
				// 	GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct1_req)).
				// 	Times(1).
				// 	Return(ct1, nil)
				// store.EXPECT().
				// 	GetCollectionTrainforUpdate(gomock.Any(), gomock.Eq(ct2_req)).
				// 	Times(1).
				// 	Return(ct2, nil)
				// store.EXPECT().
				// 	TradeTx(gomock.Any(), gomock.Eq(tradeParams)).
				// 	Times(1).
				// 	Return(db.TradeTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			requestBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			url := "/trade"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func createTestCollectionTrain(t *testing.T) (db.CollectionTrain, db.User) {
	user1, _ := randomUser(t)

	return randomCollectionTrain(user1.ID), user1
}

func createTestTradeOffer(t *testing.T) (db.TradeOffer, db.CollectionTrain, db.CollectionTrain, string, string) {
	ct1, user1 := createTestCollectionTrain(t)
	ct2, user2 := createTestCollectionTrain(t)

	return db.TradeOffer{
		ID:                  util.RandomInt(1, 1000),
		OfferedTrain:        ct1.TrainID,
		OfferedTrainOwner:   ct1.UserID,
		RequestedTrain:      ct2.TrainID,
		RequestedTrainOwner: ct2.UserID,
		CreatedAt:           time.Now(),
	}, ct1, ct2, user1.Username, user2.Username
}
