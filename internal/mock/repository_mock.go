// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/paulbahush/projects/yp/GophKeeper/internal/storage/repository.go

// Package mock_storage is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/PaBah/GophKeeper/internal/models"
	"go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AuthorizeUser mocks base method.
func (m *MockRepository) AuthorizeUser(ctx context.Context, email string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizeUser", ctx, email)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthorizeUser indicates an expected call of AuthorizeUser.
func (mr *MockRepositoryMockRecorder) AuthorizeUser(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizeUser", reflect.TypeOf((*MockRepository)(nil).AuthorizeUser), ctx, email)
}

// CreateCard mocks base method.
func (m *MockRepository) CreateCard(ctx context.Context, card models.Card) (models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCard", ctx, card)
	ret0, _ := ret[0].(models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCard indicates an expected call of CreateCard.
func (mr *MockRepositoryMockRecorder) CreateCard(ctx, card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCard", reflect.TypeOf((*MockRepository)(nil).CreateCard), ctx, card)
}

// CreateCredentials mocks base method.
func (m *MockRepository) CreateCredentials(ctx context.Context, credentials models.Credentials) (models.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCredentials", ctx, credentials)
	ret0, _ := ret[0].(models.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCredentials indicates an expected call of CreateCredentials.
func (mr *MockRepositoryMockRecorder) CreateCredentials(ctx, credentials interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCredentials", reflect.TypeOf((*MockRepository)(nil).CreateCredentials), ctx, credentials)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, user)
}

// DeleteCard mocks base method.
func (m *MockRepository) DeleteCard(ctx context.Context, cardID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", ctx, cardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard.
func (mr *MockRepositoryMockRecorder) DeleteCard(ctx, cardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockRepository)(nil).DeleteCard), ctx, cardID)
}

// DeleteCredentials mocks base method.
func (m *MockRepository) DeleteCredentials(ctx context.Context, credentialsID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCredentials", ctx, credentialsID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCredentials indicates an expected call of DeleteCredentials.
func (mr *MockRepositoryMockRecorder) DeleteCredentials(ctx, credentialsID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCredentials", reflect.TypeOf((*MockRepository)(nil).DeleteCredentials), ctx, credentialsID)
}

// GetCards mocks base method.
func (m *MockRepository) GetCards(ctx context.Context) ([]models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCards", ctx)
	ret0, _ := ret[0].([]models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCards indicates an expected call of GetCards.
func (mr *MockRepositoryMockRecorder) GetCards(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCards", reflect.TypeOf((*MockRepository)(nil).GetCards), ctx)
}

// GetCredentials mocks base method.
func (m *MockRepository) GetCredentials(ctx context.Context) ([]models.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentials", ctx)
	ret0, _ := ret[0].([]models.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredentials indicates an expected call of GetCredentials.
func (mr *MockRepositoryMockRecorder) GetCredentials(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentials", reflect.TypeOf((*MockRepository)(nil).GetCredentials), ctx)
}

// UpdateCard mocks base method.
func (m *MockRepository) UpdateCard(ctx context.Context, card models.Card) (models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", ctx, card)
	ret0, _ := ret[0].(models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockRepositoryMockRecorder) UpdateCard(ctx, card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockRepository)(nil).UpdateCard), ctx, card)
}

// UpdateCredentials mocks base method.
func (m *MockRepository) UpdateCredentials(ctx context.Context, credentials models.Credentials) (models.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCredentials", ctx, credentials)
	ret0, _ := ret[0].(models.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCredentials indicates an expected call of UpdateCredentials.
func (mr *MockRepositoryMockRecorder) UpdateCredentials(ctx, credentials interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCredentials", reflect.TypeOf((*MockRepository)(nil).UpdateCredentials), ctx, credentials)
}
