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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateCollectionTrainAPI(t *testing.T) {
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
				arg := db.CreateCollectionTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateCollectionTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.CollectionTrain{}, nil)
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
				arg := db.CreateCollectionTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateCollectionTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.CollectionTrain{}, sql.ErrConnDone)
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
				arg := db.CreateCollectionTrainParams{
					UserID:  userID + 1,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateCollectionTrain(gomock.Any(), gomock.Eq(arg)).
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
				arg := db.CreateCollectionTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					CreateCollectionTrain(gomock.Any(), gomock.Eq(arg)).
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
				url = "/users/invalid/collection/123" // Provide invalid URL for this scenario
			} else {
				url = fmt.Sprintf("/users/%d/collection/%d", tc.userID, tc.trainID) // Normal URL for other scenarios
			}
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestDeleteCollectionTrainAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	train := RandomTrain()
	collectionTrain := db.CollectionTrain{
		UserID:  userID,
		TrainID: train.ID,
	}

	testCases := []struct {
		name            string
		collectionTrain db.CollectionTrain
		setupAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:            "OK",
			collectionTrain: collectionTrain,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteCollectionTrainParams{
					UserID:  collectionTrain.UserID,
					TrainID: collectionTrain.TrainID,
				}
				store.EXPECT().
					DeleteCollectionTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:            "InternalError",
			collectionTrain: collectionTrain,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteCollectionTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					DeleteCollectionTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotAuthorized",
			collectionTrain: db.CollectionTrain{
				UserID:  userID + 1,
				TrainID: collectionTrain.TrainID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteCollectionTrainParams{
					UserID:  userID + 1,
					TrainID: train.ID,
				}
				store.EXPECT().
					DeleteCollectionTrain(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:            "InvalidRequestURI",
			collectionTrain: collectionTrain,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteCollectionTrainParams{
					UserID:  userID,
					TrainID: train.ID,
				}
				store.EXPECT().
					DeleteCollectionTrain(gomock.Any(), gomock.Eq(arg)).
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
				url = "/users/invalid/collection/123" // Provide invalid URL for this scenario
			} else {
				url = fmt.Sprintf("/users/%d/collection/%d", tc.collectionTrain.UserID, tc.collectionTrain.TrainID) // Normal URL for other scenarios
			}
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}
