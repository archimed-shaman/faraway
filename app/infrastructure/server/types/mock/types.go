// Code generated by MockGen. DO NOT EDIT.
// Source: types.go
//
// Generated by this command:
//
//	mockgen -source=types.go -destination=mock/types.go
//

// Package mock_types is a generated GoMock package.
package mock_types

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockResponseWriter is a mock of ResponseWriter interface.
type MockResponseWriter struct {
	ctrl     *gomock.Controller
	recorder *MockResponseWriterMockRecorder
}

// MockResponseWriterMockRecorder is the mock recorder for MockResponseWriter.
type MockResponseWriterMockRecorder struct {
	mock *MockResponseWriter
}

// NewMockResponseWriter creates a new mock instance.
func NewMockResponseWriter(ctrl *gomock.Controller) *MockResponseWriter {
	mock := &MockResponseWriter{ctrl: ctrl}
	mock.recorder = &MockResponseWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResponseWriter) EXPECT() *MockResponseWriterMockRecorder {
	return m.recorder
}

// Write mocks base method.
func (m *MockResponseWriter) Write(ctx context.Context, data []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", ctx, data)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockResponseWriterMockRecorder) Write(ctx, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockResponseWriter)(nil).Write), ctx, data)
}

// MockResponseReader is a mock of ResponseReader interface.
type MockResponseReader struct {
	ctrl     *gomock.Controller
	recorder *MockResponseReaderMockRecorder
}

// MockResponseReaderMockRecorder is the mock recorder for MockResponseReader.
type MockResponseReaderMockRecorder struct {
	mock *MockResponseReader
}

// NewMockResponseReader creates a new mock instance.
func NewMockResponseReader(ctrl *gomock.Controller) *MockResponseReader {
	mock := &MockResponseReader{ctrl: ctrl}
	mock.recorder = &MockResponseReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResponseReader) EXPECT() *MockResponseReaderMockRecorder {
	return m.recorder
}

// Data mocks base method.
func (m *MockResponseReader) Data() io.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Data")
	ret0, _ := ret[0].(io.Reader)
	return ret0
}

// Data indicates an expected call of Data.
func (mr *MockResponseReaderMockRecorder) Data() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockResponseReader)(nil).Data))
}
