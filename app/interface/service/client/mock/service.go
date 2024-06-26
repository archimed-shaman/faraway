// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mock/service.go
//

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	protocol "faraway/wow/pkg/protocol"
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCodec is a mock of Codec interface.
type MockCodec struct {
	ctrl     *gomock.Controller
	recorder *MockCodecMockRecorder
}

// MockCodecMockRecorder is the mock recorder for MockCodec.
type MockCodecMockRecorder struct {
	mock *MockCodec
}

// NewMockCodec creates a new mock instance.
func NewMockCodec(ctrl *gomock.Controller) *MockCodec {
	mock := &MockCodec{ctrl: ctrl}
	mock.recorder = &MockCodecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodec) EXPECT() *MockCodecMockRecorder {
	return m.recorder
}

// GetRaw mocks base method.
func (m *MockCodec) GetRaw(r io.Reader, buff []byte) (*protocol.Package, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRaw", r, buff)
	ret0, _ := ret[0].(*protocol.Package)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRaw indicates an expected call of GetRaw.
func (mr *MockCodecMockRecorder) GetRaw(r, buff any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRaw", reflect.TypeOf((*MockCodec)(nil).GetRaw), r, buff)
}

// Marshal mocks base method.
func (m *MockCodec) Marshal(v any) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal", v)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Marshal indicates an expected call of Marshal.
func (mr *MockCodecMockRecorder) Marshal(v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*MockCodec)(nil).Marshal), v)
}

// Unmarshal mocks base method.
func (m *MockCodec) Unmarshal(data []byte, v any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", data, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockCodecMockRecorder) Unmarshal(data, v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockCodec)(nil).Unmarshal), data, v)
}

// MockDDoSGuard is a mock of DDoSGuard interface.
type MockDDoSGuard struct {
	ctrl     *gomock.Controller
	recorder *MockDDoSGuardMockRecorder
}

// MockDDoSGuardMockRecorder is the mock recorder for MockDDoSGuard.
type MockDDoSGuardMockRecorder struct {
	mock *MockDDoSGuard
}

// NewMockDDoSGuard creates a new mock instance.
func NewMockDDoSGuard(ctrl *gomock.Controller) *MockDDoSGuard {
	mock := &MockDDoSGuard{ctrl: ctrl}
	mock.recorder = &MockDDoSGuardMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDDoSGuard) EXPECT() *MockDDoSGuardMockRecorder {
	return m.recorder
}

// IncRate mocks base method.
func (m *MockDDoSGuard) IncRate(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncRate", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncRate indicates an expected call of IncRate.
func (mr *MockDDoSGuardMockRecorder) IncRate(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncRate", reflect.TypeOf((*MockDDoSGuard)(nil).IncRate), ctx)
}

// MockUserLogic is a mock of UserLogic interface.
type MockUserLogic struct {
	ctrl     *gomock.Controller
	recorder *MockUserLogicMockRecorder
}

// MockUserLogicMockRecorder is the mock recorder for MockUserLogic.
type MockUserLogicMockRecorder struct {
	mock *MockUserLogic
}

// NewMockUserLogic creates a new mock instance.
func NewMockUserLogic(ctrl *gomock.Controller) *MockUserLogic {
	mock := &MockUserLogic{ctrl: ctrl}
	mock.recorder = &MockUserLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserLogic) EXPECT() *MockUserLogicMockRecorder {
	return m.recorder
}

// GetQuote mocks base method.
func (m *MockUserLogic) GetQuote(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuote", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuote indicates an expected call of GetQuote.
func (mr *MockUserLogicMockRecorder) GetQuote(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuote", reflect.TypeOf((*MockUserLogic)(nil).GetQuote), ctx)
}
