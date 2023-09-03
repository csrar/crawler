// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"

	models "github.com/csrar/tick-tock-bong/pkg/models"
	gomock "github.com/golang/mock/gomock"
)

// MockIConfig is a mock of IConfig interface.
type MockIConfig struct {
	ctrl     *gomock.Controller
	recorder *MockIConfigMockRecorder
}

// MockIConfigMockRecorder is the mock recorder for MockIConfig.
type MockIConfigMockRecorder struct {
	mock *MockIConfig
}

// NewMockIConfig creates a new mock instance.
func NewMockIConfig(ctrl *gomock.Controller) *MockIConfig {
	mock := &MockIConfig{ctrl: ctrl}
	mock.recorder = &MockIConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConfig) EXPECT() *MockIConfigMockRecorder {
	return m.recorder
}

// GetConfig mocks base method.
func (m *MockIConfig) GetConfig() models.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(models.Config)
	return ret0
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockIConfigMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockIConfig)(nil).GetConfig))
}
