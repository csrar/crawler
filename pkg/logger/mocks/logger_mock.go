// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go

// Package mock_logger is a generated GoMock package.
package mock_logger

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIlogger is a mock of Ilogger interface.
type MockIlogger struct {
	ctrl     *gomock.Controller
	recorder *MockIloggerMockRecorder
}

// MockIloggerMockRecorder is the mock recorder for MockIlogger.
type MockIloggerMockRecorder struct {
	mock *MockIlogger
}

// NewMockIlogger creates a new mock instance.
func NewMockIlogger(ctrl *gomock.Controller) *MockIlogger {
	mock := &MockIlogger{ctrl: ctrl}
	mock.recorder = &MockIloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIlogger) EXPECT() *MockIloggerMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockIlogger) Error(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", err)
}

// Error indicates an expected call of Error.
func (mr *MockIloggerMockRecorder) Error(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockIlogger)(nil).Error), err)
}

// Info mocks base method.
func (m *MockIlogger) Info(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", message)
}

// Info indicates an expected call of Info.
func (mr *MockIloggerMockRecorder) Info(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockIlogger)(nil).Info), message)
}

// Warn mocks base method.
func (m *MockIlogger) Warn(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warn", message)
}

// Warn indicates an expected call of Warn.
func (mr *MockIloggerMockRecorder) Warn(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockIlogger)(nil).Warn), message)
}
