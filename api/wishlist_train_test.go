package api

import (
	"database/sql"
	"fmt"
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

func TestCreateWishlistTrainAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	train := RandomTrain()

	testCases := []struct {
		name          string
		userID        int64
		trainID       int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			userID:  userID,
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateWishlistTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateWishlistTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.WishlistTrain{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			userID:  userID,
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateWishlistTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateWishlistTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.WishlistTrain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "NotAuthorized",
			userID:  userID + 1,
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateWishlistTrainParams{
					UserID:  userID + 1,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateWishlistTrain(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:    "InvalidRequestURI",
			userID:  userID,
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateWishlistTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateWishlistTrain(gomock.Any(), gomock.Eq(arg)).
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

			// requestBody, err := json.Marshal(tc.body)
			// require.NoError(t, err)

			var url string
			if tc.name == "InvalidRequestURI" {
				url = "/users/invalid/wishlist/123" // Provide invalid URL for this scenario
			} else {
				url = fmt.Sprintf("/users/%d/wishlist/%d", tc.userID, tc.trainID) // Normal URL for other scenarios
			}
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestGetUserWishlistAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	wishlist := randomWishlist(userID)
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.ListUserWishlistParams{
					UserID: arg1.UserID,
					Limit:  int32(arg2.PageSize),
					Offset: int32((arg2.PageID - 1) * arg2.PageSize),
				}
				store.EXPECT().
					ListUserWishlist(gomock.Any(), gomock.Eq(arg)).
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

func TestDeleteWishlistTrainAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	train := RandomTrain()
	wishlistTrain := db.WishlistTrain{
		UserID:  userID,
		TrainID: train.ID,
	}

	testCases := []struct {
		name          string
		wishlistTrain db.WishlistTrain
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			wishlistTrain: wishlistTrain,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteWishlistTrainParams{
					UserID:  wishlistTrain.UserID,
					TrainID: wishlistTrain.TrainID,
				}
				store.EXPECT().
					DeleteWishlistTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "InternalError",
			wishlistTrain: wishlistTrain,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteWishlistTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					DeleteWishlistTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotAuthorized",
			wishlistTrain: db.WishlistTrain{
				UserID:  userID + 1,
				TrainID: wishlistTrain.TrainID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteWishlistTrainParams{
					UserID:  userID + 1,
					TrainID: train.ID,
				}
				store.EXPECT().
					DeleteWishlistTrain(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:          "InvalidRequestURI",
			wishlistTrain: wishlistTrain,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteWishlistTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					DeleteWishlistTrain(gomock.Any(), gomock.Eq(arg)).
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

			// requestBody, err := json.Marshal(tc.body)
			// require.NoError(t, err)

			var url string
			if tc.name == "InvalidRequestURI" {
				url = "/users/invalid/wishlist/123" // Provide invalid URL for this scenario
			} else {
				url = fmt.Sprintf("/users/%d/wishlist/%d", tc.wishlistTrain.UserID, tc.wishlistTrain.TrainID) // Normal URL for other scenarios
			}
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func randomWishlist(userID int64) []db.WishlistTrain {
	return []db.WishlistTrain{{
		UserID:    userID,
		TrainID:   util.RandomInt(1, 1000),
		CreatedAt: time.Now(),
	}}
}
