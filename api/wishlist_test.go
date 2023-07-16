package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockdb "github.com/PlatosCodes/MerolaStation/db/mock"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
)

func TestGetUserWishlistAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	wishlist := RandomWishlist(userID)
	page_id := 1
	page_size := 1

	testCases := []struct {
		name          string
		userID        int64
		page_id       int
		page_size     int
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			userID:    userID,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(wishlist, nil)
				for _, train := range wishlist {
					store.EXPECT().
						GetTrain(gomock.Any(), train.TrainID).
						Times(1).
						Return(db.Train{}, nil) // or return some mock train data
				}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code, "Response body: %s", recorder.Body.String())
			},
		},
		{
			name:      "EmptyWishlist",
			userID:    userID,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "WishlistInternalError",
			userID:    userID,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidUserID",
			userID:    0,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: 0,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotAuthorized",
			userID:    userID + 1,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID + 1,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: 0,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "InvalidPageSize",
			userID:    userID,
			page_id:   page_id,
			page_size: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: 0,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidPageID",
			userID:    userID,
			page_id:   0,
			page_size: page_id,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   0,
					PageSize: page_id,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "GetUserWishlistTrains_InternalError",
			userID:    userID,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(wishlist, nil)
				for _, train := range wishlist {
					store.EXPECT().
						GetTrain(gomock.Any(), train.TrainID).
						Times(1).
						Return(db.Train{}, sql.ErrConnDone) // or return some mock train data
				}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "TrainInWishlistDoesNotExist",
			userID:    userID,
			page_id:   page_id,
			page_size: page_size,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg1 := UserWishlistPath{
					UserID: userID,
				}
				arg2 := WishlistPaginationQuery{
					PageID:   page_id,
					PageSize: page_size,
				}
				arg := db.GetUserWishlistTrainsParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					GetUserWishlistTrains(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(wishlist, nil)
				for _, train := range wishlist {
					store.EXPECT().
						GetTrain(gomock.Any(), train.TrainID).
						Times(1).
						Return(db.Train{}, sql.ErrNoRows) // or return some mock train data
				}
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

			url := fmt.Sprintf("/users/%d/wishlist?page_id=%d&page_size=%d", tc.userID, tc.page_id, tc.page_size)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

// func TestListTrainAPI(t *testing.T) {

// 	testCases := []struct {
// 		name             string
// 		listTrainRequest listTrainRequest
// 		buildStubs       func(store *mockdb.MockStore)
// 		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			listTrainRequest: listTrainRequest{
// 				PageID:   1,
// 				PageSize: 5,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListTrains(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return([]db.Train{}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NotFound",
// 			listTrainRequest: listTrainRequest{
// 				PageID:   1,
// 				PageSize: 5,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListTrains(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return([]db.Train{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalServerError",
// 			listTrainRequest: listTrainRequest{
// 				PageID:   1,
// 				PageSize: 5,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListTrains(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return([]db.Train{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageID",
// 			listTrainRequest: listTrainRequest{
// 				PageID:   0,
// 				PageSize: 5,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListTrains(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			// start test server and send request
// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			// requestBody, err := json.Marshal(map[int32]int32{
// 			// 	"page_id": tc.listTrainRequest.PageID,
// 			// })
// 			// require.NoError(t, err)

// 			// url := "/Trains"

// 			url := fmt.Sprintf("/trains?page_id=%d&page_size=%d", tc.listTrainRequest.PageID, tc.listTrainRequest.PageSize)

// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})

// 	}
// }

// func requireBodyMatchTrain(t *testing.T, body *bytes.Buffer, expectedTrain db.Train) {
// 	data, err := ioutil.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotTrain db.Train
// 	err = json.Unmarshal(data, &gotTrain)
// 	require.NoError(t, err)
// 	require.Equal(t, expectedTrain, gotTrain)
// }

// func RandomTrain() db.Train {
// 	return db.Train{
// 		ID:          util.RandomInt(1, 1000),
// 		ModelNumber: util.RandomModelNumber(),
// 		Year:        util.RandomYear(),
// 		Gauge:       util.RandomGauge(),
// 		Value:       0,
// 		Version:     1,
// 	}
// }

func RandomWishlist(userID int64) []db.WishlistTrain {
	return []db.WishlistTrain{{
		UserID:    userID,
		TrainID:   util.RandomInt(1, 1000),
		CreatedAt: time.Now(),
	}}
}
