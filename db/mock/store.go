// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PlatosCodes/MerolaStation/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// ActivateUser mocks base method.
func (m *MockStore) ActivateUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateUser indicates an expected call of ActivateUser.
func (mr *MockStoreMockRecorder) ActivateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateUser", reflect.TypeOf((*MockStore)(nil).ActivateUser), arg0, arg1)
}

// BlockSession mocks base method.
func (m *MockStore) BlockSession(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// BlockSession indicates an expected call of BlockSession.
func (mr *MockStoreMockRecorder) BlockSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockSession", reflect.TypeOf((*MockStore)(nil).BlockSession), arg0, arg1)
}

// CreateCollectionTrain mocks base method.
func (m *MockStore) CreateCollectionTrain(arg0 context.Context, arg1 db.CreateCollectionTrainParams) (db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCollectionTrain", arg0, arg1)
	ret0, _ := ret[0].(db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCollectionTrain indicates an expected call of CreateCollectionTrain.
func (mr *MockStoreMockRecorder) CreateCollectionTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCollectionTrain", reflect.TypeOf((*MockStore)(nil).CreateCollectionTrain), arg0, arg1)
}

// CreateImageTrain mocks base method.
func (m *MockStore) CreateImageTrain(arg0 context.Context, arg1 db.CreateImageTrainParams) (db.Train, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImageTrain", arg0, arg1)
	ret0, _ := ret[0].(db.Train)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateImageTrain indicates an expected call of CreateImageTrain.
func (mr *MockStoreMockRecorder) CreateImageTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImageTrain", reflect.TypeOf((*MockStore)(nil).CreateImageTrain), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateTradeOffer mocks base method.
func (m *MockStore) CreateTradeOffer(arg0 context.Context, arg1 db.CreateTradeOfferParams) (db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTradeOffer", arg0, arg1)
	ret0, _ := ret[0].(db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTradeOffer indicates an expected call of CreateTradeOffer.
func (mr *MockStoreMockRecorder) CreateTradeOffer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTradeOffer", reflect.TypeOf((*MockStore)(nil).CreateTradeOffer), arg0, arg1)
}

// CreateTradeTransaction mocks base method.
func (m *MockStore) CreateTradeTransaction(arg0 context.Context, arg1 db.CreateTradeTransactionParams) (db.TradeTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTradeTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.TradeTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTradeTransaction indicates an expected call of CreateTradeTransaction.
func (mr *MockStoreMockRecorder) CreateTradeTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTradeTransaction", reflect.TypeOf((*MockStore)(nil).CreateTradeTransaction), arg0, arg1)
}

// CreateTrain mocks base method.
func (m *MockStore) CreateTrain(arg0 context.Context, arg1 db.CreateTrainParams) (db.Train, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrain", arg0, arg1)
	ret0, _ := ret[0].(db.Train)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTrain indicates an expected call of CreateTrain.
func (mr *MockStoreMockRecorder) CreateTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrain", reflect.TypeOf((*MockStore)(nil).CreateTrain), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateWishlistTrain mocks base method.
func (m *MockStore) CreateWishlistTrain(arg0 context.Context, arg1 db.CreateWishlistTrainParams) (db.WishlistTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWishlistTrain", arg0, arg1)
	ret0, _ := ret[0].(db.WishlistTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWishlistTrain indicates an expected call of CreateWishlistTrain.
func (mr *MockStoreMockRecorder) CreateWishlistTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWishlistTrain", reflect.TypeOf((*MockStore)(nil).CreateWishlistTrain), arg0, arg1)
}

// DeleteActivationToken mocks base method.
func (m *MockStore) DeleteActivationToken(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActivationToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActivationToken indicates an expected call of DeleteActivationToken.
func (mr *MockStoreMockRecorder) DeleteActivationToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActivationToken", reflect.TypeOf((*MockStore)(nil).DeleteActivationToken), arg0, arg1)
}

// DeleteCollectionTrain mocks base method.
func (m *MockStore) DeleteCollectionTrain(arg0 context.Context, arg1 db.DeleteCollectionTrainParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCollectionTrain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCollectionTrain indicates an expected call of DeleteCollectionTrain.
func (mr *MockStoreMockRecorder) DeleteCollectionTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCollectionTrain", reflect.TypeOf((*MockStore)(nil).DeleteCollectionTrain), arg0, arg1)
}

// DeleteTradeOffer mocks base method.
func (m *MockStore) DeleteTradeOffer(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTradeOffer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTradeOffer indicates an expected call of DeleteTradeOffer.
func (mr *MockStoreMockRecorder) DeleteTradeOffer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTradeOffer", reflect.TypeOf((*MockStore)(nil).DeleteTradeOffer), arg0, arg1)
}

// DeleteTradeTransaction mocks base method.
func (m *MockStore) DeleteTradeTransaction(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTradeTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTradeTransaction indicates an expected call of DeleteTradeTransaction.
func (mr *MockStoreMockRecorder) DeleteTradeTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTradeTransaction", reflect.TypeOf((*MockStore)(nil).DeleteTradeTransaction), arg0, arg1)
}

// DeleteTrain mocks base method.
func (m *MockStore) DeleteTrain(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrain indicates an expected call of DeleteTrain.
func (mr *MockStoreMockRecorder) DeleteTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrain", reflect.TypeOf((*MockStore)(nil).DeleteTrain), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// DeleteWishlistTrain mocks base method.
func (m *MockStore) DeleteWishlistTrain(arg0 context.Context, arg1 db.DeleteWishlistTrainParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWishlistTrain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWishlistTrain indicates an expected call of DeleteWishlistTrain.
func (mr *MockStoreMockRecorder) DeleteWishlistTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWishlistTrain", reflect.TypeOf((*MockStore)(nil).DeleteWishlistTrain), arg0, arg1)
}

// GetActivationToken mocks base method.
func (m *MockStore) GetActivationToken(arg0 context.Context, arg1 string) (db.ActivationToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivationToken", arg0, arg1)
	ret0, _ := ret[0].(db.ActivationToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActivationToken indicates an expected call of GetActivationToken.
func (mr *MockStoreMockRecorder) GetActivationToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivationToken", reflect.TypeOf((*MockStore)(nil).GetActivationToken), arg0, arg1)
}

// GetCollectionTrain mocks base method.
func (m *MockStore) GetCollectionTrain(arg0 context.Context, arg1 db.GetCollectionTrainParams) (db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionTrain", arg0, arg1)
	ret0, _ := ret[0].(db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionTrain indicates an expected call of GetCollectionTrain.
func (mr *MockStoreMockRecorder) GetCollectionTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionTrain", reflect.TypeOf((*MockStore)(nil).GetCollectionTrain), arg0, arg1)
}

// GetCollectionTrainByID mocks base method.
func (m *MockStore) GetCollectionTrainByID(arg0 context.Context, arg1 int64) (db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionTrainByID", arg0, arg1)
	ret0, _ := ret[0].(db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionTrainByID indicates an expected call of GetCollectionTrainByID.
func (mr *MockStoreMockRecorder) GetCollectionTrainByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionTrainByID", reflect.TypeOf((*MockStore)(nil).GetCollectionTrainByID), arg0, arg1)
}

// GetCollectionTrainforUpdate mocks base method.
func (m *MockStore) GetCollectionTrainforUpdate(arg0 context.Context, arg1 db.GetCollectionTrainforUpdateParams) (db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionTrainforUpdate", arg0, arg1)
	ret0, _ := ret[0].(db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionTrainforUpdate indicates an expected call of GetCollectionTrainforUpdate.
func (mr *MockStoreMockRecorder) GetCollectionTrainforUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionTrainforUpdate", reflect.TypeOf((*MockStore)(nil).GetCollectionTrainforUpdate), arg0, arg1)
}

// GetCollectionTrainforUpdateByID mocks base method.
func (m *MockStore) GetCollectionTrainforUpdateByID(arg0 context.Context, arg1 int64) (db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionTrainforUpdateByID", arg0, arg1)
	ret0, _ := ret[0].(db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionTrainforUpdateByID indicates an expected call of GetCollectionTrainforUpdateByID.
func (mr *MockStoreMockRecorder) GetCollectionTrainforUpdateByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionTrainforUpdateByID", reflect.TypeOf((*MockStore)(nil).GetCollectionTrainforUpdateByID), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetTotalCollectionValue mocks base method.
func (m *MockStore) GetTotalCollectionValue(arg0 context.Context, arg1 int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalCollectionValue", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalCollectionValue indicates an expected call of GetTotalCollectionValue.
func (mr *MockStoreMockRecorder) GetTotalCollectionValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalCollectionValue", reflect.TypeOf((*MockStore)(nil).GetTotalCollectionValue), arg0, arg1)
}

// GetTotalTrainCount mocks base method.
func (m *MockStore) GetTotalTrainCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalTrainCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalTrainCount indicates an expected call of GetTotalTrainCount.
func (mr *MockStoreMockRecorder) GetTotalTrainCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalTrainCount", reflect.TypeOf((*MockStore)(nil).GetTotalTrainCount), arg0)
}

// GetTradeOfferByTradeID mocks base method.
func (m *MockStore) GetTradeOfferByTradeID(arg0 context.Context, arg1 int64) (db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTradeOfferByTradeID", arg0, arg1)
	ret0, _ := ret[0].(db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTradeOfferByTradeID indicates an expected call of GetTradeOfferByTradeID.
func (mr *MockStoreMockRecorder) GetTradeOfferByTradeID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTradeOfferByTradeID", reflect.TypeOf((*MockStore)(nil).GetTradeOfferByTradeID), arg0, arg1)
}

// GetTradeTransaction mocks base method.
func (m *MockStore) GetTradeTransaction(arg0 context.Context, arg1 int64) (db.TradeTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTradeTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.TradeTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTradeTransaction indicates an expected call of GetTradeTransaction.
func (mr *MockStoreMockRecorder) GetTradeTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTradeTransaction", reflect.TypeOf((*MockStore)(nil).GetTradeTransaction), arg0, arg1)
}

// GetTrain mocks base method.
func (m *MockStore) GetTrain(arg0 context.Context, arg1 int64) (db.Train, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrain", arg0, arg1)
	ret0, _ := ret[0].(db.Train)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrain indicates an expected call of GetTrain.
func (mr *MockStoreMockRecorder) GetTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrain", reflect.TypeOf((*MockStore)(nil).GetTrain), arg0, arg1)
}

// GetTrainByModel mocks base method.
func (m *MockStore) GetTrainByModel(arg0 context.Context, arg1 string) (db.Train, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrainByModel", arg0, arg1)
	ret0, _ := ret[0].(db.Train)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrainByModel indicates an expected call of GetTrainByModel.
func (mr *MockStoreMockRecorder) GetTrainByModel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrainByModel", reflect.TypeOf((*MockStore)(nil).GetTrainByModel), arg0, arg1)
}

// GetTrainByName mocks base method.
func (m *MockStore) GetTrainByName(arg0 context.Context, arg1 string) (db.Train, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrainByName", arg0, arg1)
	ret0, _ := ret[0].(db.Train)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrainByName indicates an expected call of GetTrainByName.
func (mr *MockStoreMockRecorder) GetTrainByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrainByName", reflect.TypeOf((*MockStore)(nil).GetTrainByName), arg0, arg1)
}

// GetTrainDetail mocks base method.
func (m *MockStore) GetTrainDetail(arg0 context.Context, arg1 db.GetTrainDetailParams) (db.GetTrainDetailRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrainDetail", arg0, arg1)
	ret0, _ := ret[0].(db.GetTrainDetailRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrainDetail indicates an expected call of GetTrainDetail.
func (mr *MockStoreMockRecorder) GetTrainDetail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrainDetail", reflect.TypeOf((*MockStore)(nil).GetTrainDetail), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 int64) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserByUsername mocks base method.
func (m *MockStore) GetUserByUsername(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockStoreMockRecorder) GetUserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockStore)(nil).GetUserByUsername), arg0, arg1)
}

// GetUserCollectionWithWishlistStatus mocks base method.
func (m *MockStore) GetUserCollectionWithWishlistStatus(arg0 context.Context, arg1 int64) ([]db.GetUserCollectionWithWishlistStatusRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCollectionWithWishlistStatus", arg0, arg1)
	ret0, _ := ret[0].([]db.GetUserCollectionWithWishlistStatusRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCollectionWithWishlistStatus indicates an expected call of GetUserCollectionWithWishlistStatus.
func (mr *MockStoreMockRecorder) GetUserCollectionWithWishlistStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCollectionWithWishlistStatus", reflect.TypeOf((*MockStore)(nil).GetUserCollectionWithWishlistStatus), arg0, arg1)
}

// GetUserWishlistWithCollectionStatus mocks base method.
func (m *MockStore) GetUserWishlistWithCollectionStatus(arg0 context.Context, arg1 int64) ([]db.GetUserWishlistWithCollectionStatusRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWishlistWithCollectionStatus", arg0, arg1)
	ret0, _ := ret[0].([]db.GetUserWishlistWithCollectionStatusRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWishlistWithCollectionStatus indicates an expected call of GetUserWishlistWithCollectionStatus.
func (mr *MockStoreMockRecorder) GetUserWishlistWithCollectionStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWishlistWithCollectionStatus", reflect.TypeOf((*MockStore)(nil).GetUserWishlistWithCollectionStatus), arg0, arg1)
}

// GetWishlistTrain mocks base method.
func (m *MockStore) GetWishlistTrain(arg0 context.Context, arg1 db.GetWishlistTrainParams) (db.WishlistTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWishlistTrain", arg0, arg1)
	ret0, _ := ret[0].(db.WishlistTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWishlistTrain indicates an expected call of GetWishlistTrain.
func (mr *MockStoreMockRecorder) GetWishlistTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWishlistTrain", reflect.TypeOf((*MockStore)(nil).GetWishlistTrain), arg0, arg1)
}

// InsertActivationToken mocks base method.
func (m *MockStore) InsertActivationToken(arg0 context.Context, arg1 db.InsertActivationTokenParams) (db.ActivationToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertActivationToken", arg0, arg1)
	ret0, _ := ret[0].(db.ActivationToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertActivationToken indicates an expected call of InsertActivationToken.
func (mr *MockStoreMockRecorder) InsertActivationToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertActivationToken", reflect.TypeOf((*MockStore)(nil).InsertActivationToken), arg0, arg1)
}

// ListAllUserTradeOffers mocks base method.
func (m *MockStore) ListAllUserTradeOffers(arg0 context.Context, arg1 db.ListAllUserTradeOffersParams) ([]db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllUserTradeOffers", arg0, arg1)
	ret0, _ := ret[0].([]db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllUserTradeOffers indicates an expected call of ListAllUserTradeOffers.
func (mr *MockStoreMockRecorder) ListAllUserTradeOffers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllUserTradeOffers", reflect.TypeOf((*MockStore)(nil).ListAllUserTradeOffers), arg0, arg1)
}

// ListCollectionTrainTradeOffers mocks base method.
func (m *MockStore) ListCollectionTrainTradeOffers(arg0 context.Context, arg1 db.ListCollectionTrainTradeOffersParams) (db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCollectionTrainTradeOffers", arg0, arg1)
	ret0, _ := ret[0].(db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCollectionTrainTradeOffers indicates an expected call of ListCollectionTrainTradeOffers.
func (mr *MockStoreMockRecorder) ListCollectionTrainTradeOffers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCollectionTrainTradeOffers", reflect.TypeOf((*MockStore)(nil).ListCollectionTrainTradeOffers), arg0, arg1)
}

// ListCollectionTrains mocks base method.
func (m *MockStore) ListCollectionTrains(arg0 context.Context, arg1 db.ListCollectionTrainsParams) ([]db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCollectionTrains", arg0, arg1)
	ret0, _ := ret[0].([]db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCollectionTrains indicates an expected call of ListCollectionTrains.
func (mr *MockStoreMockRecorder) ListCollectionTrains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCollectionTrains", reflect.TypeOf((*MockStore)(nil).ListCollectionTrains), arg0, arg1)
}

// ListTradeOffers mocks base method.
func (m *MockStore) ListTradeOffers(arg0 context.Context, arg1 db.ListTradeOffersParams) (db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTradeOffers", arg0, arg1)
	ret0, _ := ret[0].(db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTradeOffers indicates an expected call of ListTradeOffers.
func (mr *MockStoreMockRecorder) ListTradeOffers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTradeOffers", reflect.TypeOf((*MockStore)(nil).ListTradeOffers), arg0, arg1)
}

// ListTradeTransactions mocks base method.
func (m *MockStore) ListTradeTransactions(arg0 context.Context, arg1 db.ListTradeTransactionsParams) ([]db.TradeTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTradeTransactions", arg0, arg1)
	ret0, _ := ret[0].([]db.TradeTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTradeTransactions indicates an expected call of ListTradeTransactions.
func (mr *MockStoreMockRecorder) ListTradeTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTradeTransactions", reflect.TypeOf((*MockStore)(nil).ListTradeTransactions), arg0, arg1)
}

// ListTrainTradeTransactions mocks base method.
func (m *MockStore) ListTrainTradeTransactions(arg0 context.Context, arg1 db.ListTrainTradeTransactionsParams) ([]db.TradeTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTrainTradeTransactions", arg0, arg1)
	ret0, _ := ret[0].([]db.TradeTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTrainTradeTransactions indicates an expected call of ListTrainTradeTransactions.
func (mr *MockStoreMockRecorder) ListTrainTradeTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTrainTradeTransactions", reflect.TypeOf((*MockStore)(nil).ListTrainTradeTransactions), arg0, arg1)
}

// ListTrains mocks base method.
func (m *MockStore) ListTrains(arg0 context.Context, arg1 db.ListTrainsParams) ([]db.Train, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTrains", arg0, arg1)
	ret0, _ := ret[0].([]db.Train)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTrains indicates an expected call of ListTrains.
func (mr *MockStoreMockRecorder) ListTrains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTrains", reflect.TypeOf((*MockStore)(nil).ListTrains), arg0, arg1)
}

// ListUserCollection mocks base method.
func (m *MockStore) ListUserCollection(arg0 context.Context, arg1 db.ListUserCollectionParams) ([]db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserCollection", arg0, arg1)
	ret0, _ := ret[0].([]db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserCollection indicates an expected call of ListUserCollection.
func (mr *MockStoreMockRecorder) ListUserCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserCollection", reflect.TypeOf((*MockStore)(nil).ListUserCollection), arg0, arg1)
}

// ListUserTradeOffers mocks base method.
func (m *MockStore) ListUserTradeOffers(arg0 context.Context, arg1 db.ListUserTradeOffersParams) ([]db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserTradeOffers", arg0, arg1)
	ret0, _ := ret[0].([]db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserTradeOffers indicates an expected call of ListUserTradeOffers.
func (mr *MockStoreMockRecorder) ListUserTradeOffers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserTradeOffers", reflect.TypeOf((*MockStore)(nil).ListUserTradeOffers), arg0, arg1)
}

// ListUserTradeRequests mocks base method.
func (m *MockStore) ListUserTradeRequests(arg0 context.Context, arg1 db.ListUserTradeRequestsParams) ([]db.TradeOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserTradeRequests", arg0, arg1)
	ret0, _ := ret[0].([]db.TradeOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserTradeRequests indicates an expected call of ListUserTradeRequests.
func (mr *MockStoreMockRecorder) ListUserTradeRequests(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserTradeRequests", reflect.TypeOf((*MockStore)(nil).ListUserTradeRequests), arg0, arg1)
}

// ListUserTradeTransactions mocks base method.
func (m *MockStore) ListUserTradeTransactions(arg0 context.Context, arg1 db.ListUserTradeTransactionsParams) ([]db.TradeTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserTradeTransactions", arg0, arg1)
	ret0, _ := ret[0].([]db.TradeTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserTradeTransactions indicates an expected call of ListUserTradeTransactions.
func (mr *MockStoreMockRecorder) ListUserTradeTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserTradeTransactions", reflect.TypeOf((*MockStore)(nil).ListUserTradeTransactions), arg0, arg1)
}

// ListUserTrains mocks base method.
func (m *MockStore) ListUserTrains(arg0 context.Context, arg1 db.ListUserTrainsParams) ([]db.ListUserTrainsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserTrains", arg0, arg1)
	ret0, _ := ret[0].([]db.ListUserTrainsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserTrains indicates an expected call of ListUserTrains.
func (mr *MockStoreMockRecorder) ListUserTrains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserTrains", reflect.TypeOf((*MockStore)(nil).ListUserTrains), arg0, arg1)
}

// ListUserWishlist mocks base method.
func (m *MockStore) ListUserWishlist(arg0 context.Context, arg1 db.ListUserWishlistParams) ([]db.WishlistTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserWishlist", arg0, arg1)
	ret0, _ := ret[0].([]db.WishlistTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserWishlist indicates an expected call of ListUserWishlist.
func (mr *MockStoreMockRecorder) ListUserWishlist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserWishlist", reflect.TypeOf((*MockStore)(nil).ListUserWishlist), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(arg0 context.Context, arg1 db.ListUsersParams) ([]db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), arg0, arg1)
}

// ListWishlists mocks base method.
func (m *MockStore) ListWishlists(arg0 context.Context, arg1 db.ListWishlistsParams) ([]db.WishlistTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWishlists", arg0, arg1)
	ret0, _ := ret[0].([]db.WishlistTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWishlists indicates an expected call of ListWishlists.
func (mr *MockStoreMockRecorder) ListWishlists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWishlists", reflect.TypeOf((*MockStore)(nil).ListWishlists), arg0, arg1)
}

// RegisterTx mocks base method.
func (m *MockStore) RegisterTx(arg0 context.Context, arg1 db.CreateUserParams) (db.RegisterTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTx", arg0, arg1)
	ret0, _ := ret[0].(db.RegisterTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterTx indicates an expected call of RegisterTx.
func (mr *MockStoreMockRecorder) RegisterTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTx", reflect.TypeOf((*MockStore)(nil).RegisterTx), arg0, arg1)
}

// SearchTrainsByModelNumberSuggestions mocks base method.
func (m *MockStore) SearchTrainsByModelNumberSuggestions(arg0 context.Context, arg1 db.SearchTrainsByModelNumberSuggestionsParams) ([]db.SearchTrainsByModelNumberSuggestionsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchTrainsByModelNumberSuggestions", arg0, arg1)
	ret0, _ := ret[0].([]db.SearchTrainsByModelNumberSuggestionsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchTrainsByModelNumberSuggestions indicates an expected call of SearchTrainsByModelNumberSuggestions.
func (mr *MockStoreMockRecorder) SearchTrainsByModelNumberSuggestions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchTrainsByModelNumberSuggestions", reflect.TypeOf((*MockStore)(nil).SearchTrainsByModelNumberSuggestions), arg0, arg1)
}

// TradeTx mocks base method.
func (m *MockStore) TradeTx(arg0 context.Context, arg1 db.TradeTxParams) (db.TradeTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TradeTx", arg0, arg1)
	ret0, _ := ret[0].(db.TradeTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TradeTx indicates an expected call of TradeTx.
func (mr *MockStoreMockRecorder) TradeTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TradeTx", reflect.TypeOf((*MockStore)(nil).TradeTx), arg0, arg1)
}

// UpdateCollectionTrain mocks base method.
func (m *MockStore) UpdateCollectionTrain(arg0 context.Context, arg1 db.UpdateCollectionTrainParams) (db.CollectionTrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCollectionTrain", arg0, arg1)
	ret0, _ := ret[0].(db.CollectionTrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCollectionTrain indicates an expected call of UpdateCollectionTrain.
func (mr *MockStoreMockRecorder) UpdateCollectionTrain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCollectionTrain", reflect.TypeOf((*MockStore)(nil).UpdateCollectionTrain), arg0, arg1)
}

// UpdatePassword mocks base method.
func (m *MockStore) UpdatePassword(arg0 context.Context, arg1 db.UpdatePasswordParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockStoreMockRecorder) UpdatePassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockStore)(nil).UpdatePassword), arg0, arg1)
}

// UpdateTrainImageUrl mocks base method.
func (m *MockStore) UpdateTrainImageUrl(arg0 context.Context, arg1 db.UpdateTrainImageUrlParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrainImageUrl", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTrainImageUrl indicates an expected call of UpdateTrainImageUrl.
func (mr *MockStoreMockRecorder) UpdateTrainImageUrl(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrainImageUrl", reflect.TypeOf((*MockStore)(nil).UpdateTrainImageUrl), arg0, arg1)
}

// UpdateTrainValue mocks base method.
func (m *MockStore) UpdateTrainValue(arg0 context.Context, arg1 db.UpdateTrainValueParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrainValue", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTrainValue indicates an expected call of UpdateTrainValue.
func (mr *MockStoreMockRecorder) UpdateTrainValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrainValue", reflect.TypeOf((*MockStore)(nil).UpdateTrainValue), arg0, arg1)
}

// UpdateTrainsValuesBatch mocks base method.
func (m *MockStore) UpdateTrainsValuesBatch(arg0 context.Context, arg1, arg2 []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrainsValuesBatch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTrainsValuesBatch indicates an expected call of UpdateTrainsValuesBatch.
func (mr *MockStoreMockRecorder) UpdateTrainsValuesBatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrainsValuesBatch", reflect.TypeOf((*MockStore)(nil).UpdateTrainsValuesBatch), arg0, arg1, arg2)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 db.UpdateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}
