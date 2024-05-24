package server

import (
	"context"
	"errors"
	mock "faraway/wow/app/infrastructure/server/mock"
	"io"
	"net"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=$GOFILE -destination=mock/accept_test_mock.go

// Just hack to generate appropriate mock
type conn interface {
	net.Conn
}

var errConnectionErr = errors.New("connect error")

func TestHandleClient_OnConnectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	mockConn := mock.NewMockconn(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeout := time.Second

	mockService.EXPECT().OnConnect(ctx, gomock.Any()).Return(errConnectionErr)

	handleClient(ctx, mockConn, timeout, mockService)

	// Check no other logic is called
}

func TestHandleClient_OnDataEOF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	mockConn := mock.NewMockconn(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeout := time.Second

	mockService.EXPECT().OnConnect(ctx, gomock.Any()).Return(nil)
	mockConn.EXPECT().SetDeadline(gomock.Any()).Return(nil)
	mockService.EXPECT().OnData(ctx, gomock.Any(), gomock.Any()).Return(io.EOF)
	mockService.EXPECT().OnDisconnect(ctx)

	handleClient(ctx, mockConn, timeout, mockService)

	// Check onData returned io.EOF and no other logic is called
}

func TestHandleClient_OnDataError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	mockConn := mock.NewMockconn(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeout := time.Second

	mockService.EXPECT().OnConnect(ctx, gomock.Any()).Return(nil)
	mockConn.EXPECT().SetDeadline(gomock.Any()).Return(nil)
	mockService.EXPECT().OnData(ctx, gomock.Any(), gomock.Any()).Return(io.EOF)
	mockService.EXPECT().OnDisconnect(ctx)

	handleClient(ctx, mockConn, timeout, mockService)

	// Check onData returned error and no other logic is called
}

func TestHandleClient_ContextDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockService(ctrl)
	mockConn := mock.NewMockconn(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	timeout := time.Second

	mockService.EXPECT().OnConnect(ctx, gomock.Any()).Return(nil)
	mockConn.EXPECT().SetDeadline(gomock.Any()).Return(nil).AnyTimes()

	cancel()

	mockService.EXPECT().OnDisconnect(ctx)

	handleClient(ctx, mockConn, timeout, mockService)
}
