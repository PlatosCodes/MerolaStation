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

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockdb "github.com/PlatosCodes/MerolaStation/db/mock"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
)

func TestCreateTrainAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username

	modelNumber, name := util.RandomTrainRequest()

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
				"model_number": modelNumber,
				"train_name":   name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTrainParams{
					ModelNumber: modelNumber,
					Name:        name,
				}
				store.EXPECT().
					CreateTrain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Train{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidModelNumber",
			requestBody: map[string]interface{}{
				"model_number": "",
				"train_name":   name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrain(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			requestBody: map[string]interface{}{
				"model_number": modelNumber,
				"train_name":   name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTrainParams{
					ModelNumber: modelNumber,
					Name:        name,
				}
				store.EXPECT().CreateTrain(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Train{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := "/trains"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestUpdateTrainValueAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	newValue := int64(1000)
	train := RandomTrain()

	testCases := []struct {
		name          string
		trainID       int64
		newValue      int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			trainID:  train.ID,
			newValue: newValue,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTrainValueParams{
					ID:    train.ID,
					Value: newValue,
				}
				store.EXPECT().
					UpdateTrainValue(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "NotFound",
			trainID:  train.ID,
			newValue: 1000,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTrainValueParams{
					ID:    train.ID,
					Value: 1000,
				}
				store.EXPECT().
					UpdateTrainValue(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			trainID:  train.ID,
			newValue: 1000,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTrainValueParams{
					ID:    train.ID,
					Value: 1000,
				}
				store.EXPECT().
					UpdateTrainValue(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidID",
			trainID:  0,
			newValue: 1000,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTrainValueParams{
					ID:    0,
					Value: 1000,
				}
				store.EXPECT().
					UpdateTrainValue(gomock.Any(), gomock.Eq(arg)).
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/trains/value"
			requestBody, err := json.Marshal(map[string]int64{
				"id":    tc.trainID,
				"value": tc.newValue,
			})
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(requestBody))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetTrainAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	train := RandomTrain()

	testCases := []struct {
		name          string
		trainID       int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrain(gomock.Any(), gomock.Eq(train.ID)).
					Times(1).
					Return(train, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrain(t, recorder.Body, train)
			},
		},
		{
			name:    "NotFound",
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrain(gomock.Any(), gomock.Eq(train.ID)).
					Times(1).
					Return(db.Train{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			trainID: train.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrain(gomock.Any(), gomock.Eq(train.ID)).
					Times(1).
					Return(db.Train{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "InvalidID",
			trainID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrain(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/trains/%d", tc.trainID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestGetTrainByModelAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username
	train := RandomTrain()

	testCases := []struct {
		name          string
		ModelNumber   string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			ModelNumber: train.ModelNumber,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrainByModel(gomock.Any(), gomock.Eq(train.ModelNumber)).
					Times(1).
					Return(train, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrain(t, recorder.Body, train)
			},
		},
		{
			name:        "NotFound",
			ModelNumber: train.ModelNumber,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrainByModel(gomock.Any(), gomock.Eq(train.ModelNumber)).
					Times(1).
					Return(db.Train{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "InternalError",
			ModelNumber: train.ModelNumber,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrainByModel(gomock.Any(), gomock.Eq(train.ModelNumber)).
					Times(1).
					Return(db.Train{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/trains/model/%s", tc.ModelNumber)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestListTrainAPI(t *testing.T) {
	user, _ := randomUser(t)
	userID := user.ID
	username := user.Username

	testCases := []struct {
		name             string
		listTrainRequest listTrainRequest
		setupAuth        func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			listTrainRequest: listTrainRequest{
				PageID:   1,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrains(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Train{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotFound",
			listTrainRequest: listTrainRequest{
				PageID:   1,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrains(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Train{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			listTrainRequest: listTrainRequest{
				PageID:   1,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrains(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Train{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			listTrainRequest: listTrainRequest{
				PageID:   0,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrains(gomock.Any(), gomock.Any()).
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

			// requestBody, err := json.Marshal(map[int32]int32{
			// 	"page_id": tc.listTrainRequest.PageID,
			// })
			// require.NoError(t, err)

			// url := "/Trains"

			url := fmt.Sprintf("/trains/all?page_id=%d&page_size=%d", tc.listTrainRequest.PageID, tc.listTrainRequest.PageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func requireBodyMatchTrain(t *testing.T, body *bytes.Buffer, expectedTrain db.Train) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTrain db.Train
	err = json.Unmarshal(data, &gotTrain)
	require.NoError(t, err)
	require.Equal(t, expectedTrain, gotTrain)
}

func RandomTrain() db.Train {
	return db.Train{
		ID:          util.RandomInt(1, 1000),
		ModelNumber: util.RandomModelNumber(),
		Name:        util.RandomTrainName(),
		Value:       0,
		Version:     1,
	}
}

//Have to EDIT THIS FIRST -- it is copied from ListTrainsAPI and not set up for User Trains
// func TestListUserTrainsAPI(t *testing.T) {
// 	user, _ := randomUser(t)
// 	userID := user.ID
// 	username := user.Username
// 	limit := int32(5)
// 	offset := int32(1)

// 	testCases := []struct {
// 		name          string
// 		userID        int64
// 		limit         int32
// 		offset        int32
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.ListUserTrainsParams{
// 					UserID: userID,
// 					Limit:  limit,
// 					Offset: (offset - 1) * limit,
// 				}
// 				store.EXPECT().
// 					ListUserTrains(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return([]db.ListUserTrainsRow{}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		// {
// 		// 	name: "NotFound",
// 		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		arg := db.ListUserTrainsParams{
// 		// 			UserID: userID,
// 		// 			Limit:  limit,
// 		// 			Offset: (offset - 1) * limit,
// 		// 		}
// 		// 		store.EXPECT().
// 		// 			ListUserTrains(gomock.Any(), gomock.Eq(arg)).
// 		// 			Times(1).
// 		// 			Return([]db.ListUserTrainsRow{}, sql.ErrNoRows)
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "InternalServerError",
// 		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		arg := db.ListUserTrainsParams{
// 		// 			UserID: userID,
// 		// 			Limit:  limit,
// 		// 			Offset: (offset - 1) * limit,
// 		// 		}
// 		// 		store.EXPECT().
// 		// 			ListUserTrains(gomock.Any(), gomock.Eq(arg)).
// 		// 			Times(1).
// 		// 			Return([]db.ListUserTrainsRow{}, sql.ErrConnDone)
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "InvalidPageID",
// 		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, username, time.Minute)
// 		// 	},
// 		// 	buildStubs: func(store *mockdb.MockStore) {
// 		// 		arg := db.ListUserTrainsParams{
// 		// 			UserID: userID,
// 		// 			Limit:  limit,
// 		// 			Offset: (offset - 1) * limit,
// 		// 		}
// 		// 		store.EXPECT().
// 		// 			ListUserTrains(gomock.Any(), gomock.Eq(arg)).
// 		// 			Times(0)
// 		// 	},
// 		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		// 	},
// 		// },
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			// start test server and send request
// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			// requestBody, err := json.Marshal(map[int32]int32{
// 			// 	"page_id": tc.listTrainRequest.PageID,
// 			// })
// 			// require.NoError(t, err)

// 			// url := "/Trains"

// 			url := fmt.Sprintf("/trains?page_id=%d&page_size=%d", tc.offset, tc.limit)

// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.tokenMaker)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})

// 	}
// }
