// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockICrawlerStore is a mock of ICrawlerStore interface.
type MockICrawlerStore struct {
	ctrl     *gomock.Controller
	recorder *MockICrawlerStoreMockRecorder
}

// MockICrawlerStoreMockRecorder is the mock recorder for MockICrawlerStore.
type MockICrawlerStoreMockRecorder struct {
	mock *MockICrawlerStore
}

// NewMockICrawlerStore creates a new mock instance.
func NewMockICrawlerStore(ctrl *gomock.Controller) *MockICrawlerStore {
	mock := &MockICrawlerStore{ctrl: ctrl}
	mock.recorder = &MockICrawlerStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICrawlerStore) EXPECT() *MockICrawlerStoreMockRecorder {
	return m.recorder
}

// WasAlreadyVisited mocks base method.
func (m *MockICrawlerStore) WasAlreadyVisited(site string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WasAlreadyVisited", site)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WasAlreadyVisited indicates an expected call of WasAlreadyVisited.
func (mr *MockICrawlerStoreMockRecorder) WasAlreadyVisited(site interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WasAlreadyVisited", reflect.TypeOf((*MockICrawlerStore)(nil).WasAlreadyVisited), site)
}
