// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/server/model/gacha_probability.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	model "game-gacha/pkg/server/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGachaProbabilityRepositoryInterface is a mock of GachaProbabilityRepositoryInterface interface.
type MockGachaProbabilityRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGachaProbabilityRepositoryInterfaceMockRecorder
}

// MockGachaProbabilityRepositoryInterfaceMockRecorder is the mock recorder for MockGachaProbabilityRepositoryInterface.
type MockGachaProbabilityRepositoryInterfaceMockRecorder struct {
	mock *MockGachaProbabilityRepositoryInterface
}

// NewMockGachaProbabilityRepositoryInterface creates a new mock instance.
func NewMockGachaProbabilityRepositoryInterface(ctrl *gomock.Controller) *MockGachaProbabilityRepositoryInterface {
	mock := &MockGachaProbabilityRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockGachaProbabilityRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGachaProbabilityRepositoryInterface) EXPECT() *MockGachaProbabilityRepositoryInterfaceMockRecorder {
	return m.recorder
}

// SelectGachaProbabilities mocks base method.
func (m *MockGachaProbabilityRepositoryInterface) SelectGachaProbabilities() ([]*model.GachaProbability, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectGachaProbabilities")
	ret0, _ := ret[0].([]*model.GachaProbability)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectGachaProbabilities indicates an expected call of SelectGachaProbabilities.
func (mr *MockGachaProbabilityRepositoryInterfaceMockRecorder) SelectGachaProbabilities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectGachaProbabilities", reflect.TypeOf((*MockGachaProbabilityRepositoryInterface)(nil).SelectGachaProbabilities))
}
