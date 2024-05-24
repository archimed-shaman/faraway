// Code generated by MockGen. DO NOT EDIT.
// Source: server.go
//
// Generated by this command:
//
//	mockgen -source=server.go -destination=mock/server.go
//

// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	types "faraway/wow/app/infrastructure/server/types"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// OnConnect mocks base method.
func (m *MockService) OnConnect(ctx context.Context, w types.ResponseWriter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnConnect", ctx, w)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnConnect indicates an expected call of OnConnect.
func (mr *MockServiceMockRecorder) OnConnect(ctx, w any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnConnect", reflect.TypeOf((*MockService)(nil).OnConnect), ctx, w)
}

// OnData mocks base method.
func (m *MockService) OnData(ctx context.Context, r types.ResponseReader, w types.ResponseWriter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnData", ctx, r, w)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnData indicates an expected call of OnData.
func (mr *MockServiceMockRecorder) OnData(ctx, r, w any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnData", reflect.TypeOf((*MockService)(nil).OnData), ctx, r, w)
}

// OnDisconnect mocks base method.
func (m *MockService) OnDisconnect(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDisconnect", ctx)
}

// OnDisconnect indicates an expected call of OnDisconnect.
func (mr *MockServiceMockRecorder) OnDisconnect(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDisconnect", reflect.TypeOf((*MockService)(nil).OnDisconnect), ctx)
}
