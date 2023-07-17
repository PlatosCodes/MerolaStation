package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/PlatosCodes/MerolaStation/db/mock"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTradeOffer(t *testing.T) {
	user1, _ := randomUser(t)
	ct1 := randomCollection(user1.ID)
	fmt.Print(ct1[0])

	user2, _ := randomUser(t)
	ct2 := randomCollection(user2.ID)

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
				"offered_train":         ct1[0].TrainID,
				"offered_train_owner":   ct1[0].UserID,
				"requested_train":       ct2[0].TrainID,
				"requested_train_owner": ct2[0].UserID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTradeOfferParams{
					OfferedTrain:        ct1[0].TrainID,
					OfferedTrainOwner:   ct1[0].UserID,
					RequestedTrain:      ct2[0].TrainID,
					RequestedTrainOwner: ct2[0].UserID,
				}
				store.EXPECT().
					CreateTradeOffer(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TradeOffer{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			requestBody: map[string]interface{}{
				"offered_train":         ct1[0].TrainID,
				"offered_train_owner":   ct1[0].UserID,
				"requested_train":       ct2[0].TrainID,
				"requested_train_owner": ct2[0].UserID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTradeOfferParams{
					OfferedTrain:        ct1[0].TrainID,
					OfferedTrainOwner:   ct1[0].UserID,
					RequestedTrain:      ct2[0].TrainID,
					RequestedTrainOwner: ct2[0].UserID,
				}
				store.EXPECT().
					CreateTradeOffer(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TradeOffer{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotAuthorized",
			requestBody: map[string]interface{}{
				"offered_train":         ct1[0].TrainID,
				"offered_train_owner":   ct1[0].UserID + 1,
				"requested_train":       ct2[0].TrainID,
				"requested_train_owner": ct2[0].UserID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTradeOfferParams{
					OfferedTrain:        ct1[0].TrainID,
					OfferedTrainOwner:   ct1[0].UserID + 1,
					RequestedTrain:      ct2[0].TrainID,
					RequestedTrainOwner: ct2[0].UserID,
				}
				store.EXPECT().
					CreateTradeOffer(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidRequest",
			requestBody: map[string]interface{}{
				"offered_train":         0,
				"offered_train_owner":   ct1[0].UserID,
				"requested_train":       ct2[0].TrainID,
				"requested_train_owner": ct2[0].UserID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTradeOfferParams{
					OfferedTrain:        0,
					OfferedTrainOwner:   ct1[0].UserID,
					RequestedTrain:      ct2[0].TrainID,
					RequestedTrainOwner: ct2[0].UserID,
				}
				store.EXPECT().
					CreateTradeOffer(gomock.Any(), gomock.Eq(arg)).
					Times(0)
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

			url := "/trade_offer"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestGetTradeOfferByIDAPI(t *testing.T) {
	to, username := randomTradeOffer(t)

	testCases := []struct {
		name          string
		trade_offer   int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			trade_offer: to.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(to, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTradeOffer(t, recorder.Body, to)
			},
		},
		{
			name:        "NotFound",
			trade_offer: to.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
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
			name:        "InternalError",
			trade_offer: to.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
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
			name:        "InvalidID",
			trade_offer: to.ID + 1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq(to.ID+1)).
					Times(1).
					Return(db.TradeOffer{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "InvalidRequestURI",
			trade_offer: to.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTradeOfferByTradeID(gomock.Any(), gomock.Eq("string")).
					Times(0)
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

			var url string
			if tc.name == "InvalidRequestURI" {
				url = "/trade_offer/string" // Provide invalid URL for this scenario
			} else {
				url = fmt.Sprintf("/trade_offer/%d", tc.trade_offer)
			}
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestListUserTradeOfferAPI(t *testing.T) {
	to, username := randomTradeOffer(t)
	page_id := 1
	page_size := 1

	testCases := []struct {
		name                string
		offered_train_owner int64
		page_id             int
		page_size           int
		setupAuth           func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs          func(store *mockdb.MockStore)
		checkResponse       func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:                "OK",
			offered_train_owner: to.OfferedTrainOwner,
			page_id:             page_id,
			page_size:           page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := UserTradeOffersPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.ListUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.TradeOffer{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
				// requireBodyMatchTradeOffer(t, recorder.Body, to)
			},
		},
		{
			name:                "EmptyTradeOffer",
			offered_train_owner: to.OfferedTrainOwner,
			page_id:             page_id,
			page_size:           page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := UserTradeOffersPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.ListUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.TradeOffer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:                "NotAuthorized",
			offered_train_owner: to.OfferedTrainOwner + 1,
			page_id:             page_id,
			page_size:           page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := UserTradeOffersPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.ListUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner + 1,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:                "InvalidUserID",
			offered_train_owner: 0,
			page_id:             page_id,
			page_size:           page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: 0,
				}
				arg2 := UserTradeOffersPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.ListUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:                "InvalidPageSize",
			offered_train_owner: to.OfferedTrainOwner,
			page_id:             page_id,
			page_size:           0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := UserTradeOffersPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.ListUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:                "InvalidPageID",
			offered_train_owner: to.OfferedTrainOwner,
			page_id:             0,
			page_size:           page_id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := UserTradeOffersPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.ListUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:                "MissingQueryParameters",
			offered_train_owner: to.OfferedTrainOwner,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				// arg2 := UserTradeOffersPaginationQuery{
				// 	PageID:   page_id,
				// 	PageSize: page_size,
				// }
				// arg := db.ListUserTradeOffersParams{
				// 	OfferedTrainOwner: arg1.OfferedTrainOwner,
				// 	Limit:             int32(arg2.PageSize),
				// 	Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				// }
				store.EXPECT().
					ListUserTradeOffers(gomock.Any(), gomock.Eq(arg1)).
					Times(0)
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

			var url string
			if tc.name == "MissingQueryParameters" {
				url = fmt.Sprintf("/users/trade_offers/offered/%d", tc.offered_train_owner)

			} else {
				url = fmt.Sprintf("/users/trade_offers/offered/%d?page_id=%d&page_size=%d", tc.offered_train_owner, tc.page_id, tc.page_size)
			}
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestListUserTradeRequestsAPI(t *testing.T) {
	to, username := randomTradeOfferforRequests(t)
	page_id := 1
	page_size := 5
	testCases := []struct {
		name                  string
		requested_train_owner int64
		page_size             int32
		page_id               int32
		setupAuth             func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs            func(store *mockdb.MockStore)
		checkResponse         func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:                  "OK",
			requested_train_owner: to.RequestedTrainOwner,
			page_size:             int32(page_size),
			page_id:               int32(page_id),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeRequestsPath{
					RequestedTrainOwner: to.RequestedTrainOwner,
				}
				arg2 := UserTradeRequestsPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListUserTradeRequestsParams{
					RequestedTrainOwner: arg1.RequestedTrainOwner,
					Limit:               int32(arg2.PageSize),
					Offset:              int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeRequests(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.TradeOffer{to}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
				// requireBodyMatchTradeOffer(t, recorder.Body, to)
			},
		},
		{
			name:                  "EmptyTradeOffer",
			requested_train_owner: to.RequestedTrainOwner,
			page_id:               int32(page_id),
			page_size:             int32(page_size),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeRequestsPath{
					RequestedTrainOwner: to.RequestedTrainOwner,
				}
				arg2 := UserTradeRequestsPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListUserTradeRequestsParams{
					RequestedTrainOwner: arg1.RequestedTrainOwner,
					Limit:               int32(arg2.PageSize),
					Offset:              int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeRequests(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.TradeOffer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:                  "NotAuthorized",
			requested_train_owner: to.RequestedTrainOwner + 1,
			page_id:               int32(page_id),
			page_size:             int32(page_size),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeRequestsPath{
					RequestedTrainOwner: to.RequestedTrainOwner,
				}
				arg2 := UserTradeRequestsPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListUserTradeRequestsParams{
					RequestedTrainOwner: arg1.RequestedTrainOwner,
					Limit:               int32(arg2.PageSize),
					Offset:              int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeRequests(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:                  "InvalidUserID",
			requested_train_owner: 0,
			page_id:               int32(page_id),
			page_size:             int32(page_size),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserTradeRequestsPath{
					RequestedTrainOwner: to.RequestedTrainOwner,
				}
				arg2 := UserTradeRequestsPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListUserTradeRequestsParams{
					RequestedTrainOwner: arg1.RequestedTrainOwner,
					Limit:               int32(arg2.PageSize),
					Offset:              int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserTradeRequests(gomock.Any(), gomock.Eq(arg)).
					Times(0)
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

			url := fmt.Sprintf("/users/trade_offers/requests/%d?page_id=%d&page_size=%d", tc.requested_train_owner, tc.page_id, tc.page_size)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestListAllUserTradeOffersAPI(t *testing.T) {
	to, username := randomTradeOffer(t)
	page_id := 1
	page_size := 5

	testCases := []struct {
		name              string
		OfferedTrainOwner int64
		page_size         int
		page_id           int
		setupAuth         func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs        func(store *mockdb.MockStore)
		checkResponse     func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:              "OK",
			OfferedTrainOwner: to.OfferedTrainOwner,
			page_size:         page_size,
			page_id:           page_id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := AllUserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := AllUserTradeOffersPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListAllUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListAllUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.TradeOffer{to}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
				// requireBodyMatchTradeOffer(t, recorder.Body, to)
			},
		},
		{
			name:              "EmptyTradeOffer",
			OfferedTrainOwner: to.OfferedTrainOwner,
			page_id:           page_id,
			page_size:         page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := AllUserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := AllUserTradeOffersPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListAllUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListAllUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.TradeOffer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:              "NotAuthorized",
			OfferedTrainOwner: to.OfferedTrainOwner + 1,
			page_id:           page_id,
			page_size:         page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := AllUserTradeOffersPath{
					OfferedTrainOwner: to.OfferedTrainOwner,
				}
				arg2 := AllUserTradeOffersPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListAllUserTradeOffersParams{
					OfferedTrainOwner: arg1.OfferedTrainOwner,
					Limit:             int32(arg2.PageSize),
					Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListAllUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:              "InvalidOfferedTrainOwner",
			OfferedTrainOwner: 0,
			page_id:           page_id,
			page_size:         page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllUserTradeOffersParams{
					OfferedTrainOwner: 0,
					Limit:             int32(page_size),
					Offset:            int32((page_id - 1) * page_size),
				}
				store.EXPECT().
					ListUserTradeRequests(gomock.Any(), gomock.Eq(arg)).
					Times(0)
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

			url := fmt.Sprintf("/users/trade_offers/all/%d?page_id=%d&page_size=%d", tc.OfferedTrainOwner, tc.page_id, tc.page_size)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestListCollectionTrainTradeOffersAPI(t *testing.T) {
	to, username := randomTradeOfferforRequests(t)
	page_id := 1
	page_size := 5

	testCases := []struct {
		name            string
		requested_train int64
		page_size       int
		page_id         int
		setupAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:            "OK",
			requested_train: to.RequestedTrain,
			page_size:       page_size,
			page_id:         page_id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrain, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := listCollectionTrainTradeOffersPath{
					TrainID: to.RequestedTrain,
				}
				arg2 := AllUserTradeOffersPaginationQuery{
					PageID:   int32(page_id),
					PageSize: int32(page_size),
				}
				arg := db.ListCollectionTrainTradeOffersParams{
					RequestedTrain: arg1.TrainID,
					Limit:          int32(arg2.PageSize),
					Offset:         int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListCollectionTrainTradeOffers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(to, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
				// requireBodyMatchTradeOffer(t, recorder.Body, to)
			},
		},
		// 	{
		// 		name:              "EmptyTradeOffer",
		// 		OfferedTrainOwner: to.OfferedTrainOwner,
		// 		page_id:           page_id,
		// 		page_size:         page_size,
		// 		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 			addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
		// 		},
		// 		buildStubs: func(store *mockdb.MockStore) {
		// 			arg1 := AllUserTradeOffersPath{
		// 				OfferedTrainOwner: to.OfferedTrainOwner,
		// 			}
		// 			arg2 := AllUserTradeOffersPaginationQuery{
		// 				PageID:   int32(page_id),
		// 				PageSize: int32(page_size),
		// 			}
		// 			arg := db.ListAllUserTradeOffersParams{
		// 				OfferedTrainOwner: arg1.OfferedTrainOwner,
		// 				Limit:             int32(arg2.PageSize),
		// 				Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
		// 			}
		// 			store.EXPECT().
		// 				ListAllUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
		// 				Times(1).
		// 				Return([]db.TradeOffer{}, sql.ErrNoRows)
		// 		},
		// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 			require.Equal(t, http.StatusNotFound, recorder.Code)
		// 		},
		// 	},
		// 	{
		// 		name:              "NotAuthorized",
		// 		OfferedTrainOwner: to.OfferedTrainOwner + 1,
		// 		page_id:           page_id,
		// 		page_size:         page_size,
		// 		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 			addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
		// 		},
		// 		buildStubs: func(store *mockdb.MockStore) {
		// 			arg1 := AllUserTradeOffersPath{
		// 				OfferedTrainOwner: to.OfferedTrainOwner,
		// 			}
		// 			arg2 := AllUserTradeOffersPaginationQuery{
		// 				PageID:   int32(page_id),
		// 				PageSize: int32(page_size),
		// 			}
		// 			arg := db.ListAllUserTradeOffersParams{
		// 				OfferedTrainOwner: arg1.OfferedTrainOwner,
		// 				Limit:             int32(arg2.PageSize),
		// 				Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
		// 			}
		// 			store.EXPECT().
		// 				ListAllUserTradeOffers(gomock.Any(), gomock.Eq(arg)).
		// 				Times(0)
		// 		},
		// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 		},
		// 	},
		// 	{
		// 		name:              "InvalidOfferedTrainOwner",
		// 		OfferedTrainOwner: 0,
		// 		page_id:           page_id,
		// 		page_size:         page_size,
		// 		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 			addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.RequestedTrainOwner, username, time.Minute)
		// 		},
		// 		buildStubs: func(store *mockdb.MockStore) {
		// 			arg := db.ListAllUserTradeOffersParams{
		// 				OfferedTrainOwner: 0,
		// 				Limit:             int32(page_size),
		// 				Offset:            int32((page_id - 1) * page_size),
		// 			}
		// 			store.EXPECT().
		// 				ListUserTradeRequests(gomock.Any(), gomock.Eq(arg)).
		// 				Times(0)
		// 		},
		// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 			require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 		},
		// 	},
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

			url := fmt.Sprintf("/collection_trains/trade_offers/%d?page_id=%d&page_size=%d", tc.requested_train, tc.page_id, tc.page_size)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestDeleteTradeOfferAPI(t *testing.T) {
	to, username := randomTradeOffer(t)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"id":                  to.ID,
				"offered_train_owner": to.OfferedTrainOwner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTradeOffer(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id":                  to.ID,
				"offered_train_owner": to.OfferedTrainOwner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTradeOffer(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotAuthorized",
			body: gin.H{
				"id":                  to.ID,
				"offered_train_owner": to.OfferedTrainOwner + 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTradeOffer(gomock.Any(), gomock.Eq(to.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidRequest",
			body: gin.H{
				"id":                  0,
				"offered_train_owner": to.OfferedTrainOwner,
			}, setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTradeOffer(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"id":                  to.ID,
				"offered_train_owner": to.OfferedTrainOwner,
			}, setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, to.OfferedTrainOwner, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteTradeOffer(gomock.Any(), gomock.Eq(to.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			requestBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/trade_offer"

			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(requestBody))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

// List train trade offers tests

// List all trade offers tests

// {
// 	name:      "GetUserTradeOfferTrains_InternalError",
// offered_train_owner: to.OfferedTrainOwner,
// 	page_id:   page_id,
// 	page_size: page_size,
// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
// 	},
// 	buildStubs: func(store *mockdb.MockStore) {
// 		arg1 := UserTradeOffersPath{
// 			OfferedTrainOwner: user.ID,
// 		}
// 		arg2 := UserTradeOffersPaginationQuery{
// 			PageID:   page_id,
// 			PageSize: page_size,
// 		}
// 		arg := db.ListUserTradeOffersParams{
// 			OfferedTrainOwner: arg1.OfferedTrainOwner,
// 			Limit:             int32(arg2.PageSize),
// 			Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
// 		}
// 		store.EXPECT().
// 			ListUserTradeOffer(gomock.Any(), gomock.Eq(arg)).
// 			Times(1).
// 			Return(trade_offer, nil)
// 		for _, train := range trade_offer {
// 			store.EXPECT().
// 				GetTrain(gomock.Any(), train.TrainID).
// 				Times(1).
// 				Return(db.Train{}, sql.ErrConnDone) // or return some mock train data
// 		}
// 	},
// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 	},
// },
// {
// 	name:      "TrainInTradeOfferDoesNotExist",
// offered_train_owner: to.OfferedTrainOwner,
// 	page_id:   page_id,
// 	page_size: page_size,
// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
// 	},
// 	buildStubs: func(store *mockdb.MockStore) {
// 		arg1 := UserTradeOffersPath{
// 			OfferedTrainOwner: user.ID,
// 		}
// 		arg2 := UserTradeOffersPaginationQuery{
// 			PageID:   page_id,
// 			PageSize: page_size,
// 		}
// 		arg := db.ListUserTradeOffersParams{
// 			OfferedTrainOwner: arg1.OfferedTrainOwner,
// 			Limit:             int32(arg2.PageSize),
// 			Offset:            int32((arg2.PageID - 1) * arg2.PageSize),
// 		}
// 		store.EXPECT().
// 			ListUserTradeOffer(gomock.Any(), gomock.Eq(arg)).
// 			Times(1).
// 			Return(trade_offer, nil)
// 		for _, train := range trade_offer {
// 			store.EXPECT().
// 				GetTrain(gomock.Any(), train.TrainID).
// 				Times(1).
// 				Return(db.Train{}, sql.ErrNoRows) // or return some mock train data
// 		}
// 	},
// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		require.Equal(t, http.StatusNotFound, recorder.Code)
// 	},

func randomTradeOffer(t *testing.T) (tradeOffer db.TradeOffer, username string) {
	user1, _ := randomUser(t)
	offered_trade_username := user1.Username
	ct1 := randomCollection(user1.ID)[0]

	user2, _ := randomUser(t)
	ct2 := randomCollection(user2.ID)[0]

	return db.TradeOffer{
		ID:                  util.RandomInt(1, 1000),
		OfferedTrain:        ct1.TrainID,
		OfferedTrainOwner:   ct1.UserID,
		RequestedTrain:      ct2.TrainID,
		RequestedTrainOwner: ct2.UserID,
		CreatedAt:           time.Now(),
	}, offered_trade_username
}

func randomTradeOfferforRequests(t *testing.T) (tradeOffer db.TradeOffer, username2 string) {
	user1, _ := randomUser(t)
	ct1 := randomCollection(user1.ID)[0]

	user2, _ := randomUser(t)
	ct2 := randomCollection(user2.ID)[0]

	return db.TradeOffer{
		ID:                  util.RandomInt(1, 1000),
		OfferedTrain:        ct1.TrainID,
		OfferedTrainOwner:   ct1.UserID,
		RequestedTrain:      ct2.TrainID,
		RequestedTrainOwner: ct2.UserID,
		CreatedAt:           time.Now(),
	}, user2.Username
}

func requireBodyMatchTradeOffer(t *testing.T, body *bytes.Buffer, tradeOffer db.TradeOffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTradeOffer db.TradeOffer
	err = json.Unmarshal(data, &gotTradeOffer)
	require.NoError(t, err)
	require.Equal(t, tradeOffer.ID, gotTradeOffer.ID)
	require.Equal(t, tradeOffer.OfferedTrain, gotTradeOffer.OfferedTrain)
	require.Equal(t, tradeOffer.OfferedTrainOwner, gotTradeOffer.OfferedTrainOwner)
	require.Equal(t, tradeOffer.RequestedTrain, gotTradeOffer.RequestedTrain)
	require.Equal(t, tradeOffer.RequestedTrainOwner, gotTradeOffer.RequestedTrainOwner)
	require.WithinDuration(t, tradeOffer.CreatedAt, gotTradeOffer.CreatedAt, time.Second)

}
