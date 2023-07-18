package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDatabaseRepository is a mock of DatabaseRepository interface.
type MockDatabaseRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseRepositoryMockRecorder
}

// MockDatabaseRepositoryMockRecorder is the mock recorder for MockDatabaseRepository.
type MockDatabaseRepositoryMockRecorder struct {
	mock *MockDatabaseRepository
}

// NewMockDatabaseRepository creates a new mock instance.
func NewMockDatabaseRepository(ctrl *gomock.Controller) *MockDatabaseRepository {
	mock := &MockDatabaseRepository{ctrl: ctrl}
	mock.recorder = &MockDatabaseRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseRepository) EXPECT() *MockDatabaseRepositoryMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockDatabaseRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockDatabaseRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDatabaseRepository)(nil).Ping))
}
