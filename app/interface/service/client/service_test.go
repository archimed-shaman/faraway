//nolint:errchkjson // To avoid test code in common space
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"faraway/wow/pkg/protocol"
	"faraway/wow/pkg/test"
	"testing"

	mock "faraway/wow/app/interface/service/client/mock"

	"go.uber.org/mock/gomock"
)

func TestService_onNonceReq_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockCodec(ctrl)
	mockDDoSGuard := mock.NewMockDDoSGuard(ctrl)
	mockUserLogic := mock.NewMockUserLogic(ctrl)

	rate := 3
	factor := 2.0
	maxDifficulty := 10

	svc := NewService(1024, maxDifficulty, factor, mockCodec, mockUserLogic, mockDDoSGuard)

	ctx := context.Background()
	writer := &bytes.Buffer{}

	mockDDoSGuard.EXPECT().IncRate(gomock.Any()).Return(int64(rate), nil)

	nonceResp := &protocol.NonceResp{
		Nonce:      []byte("test_nonce"),
		Difficulty: rate * int(factor),
	}

	payloadBytes, _ := json.Marshal(nonceResp)

	packageNonceResp := &protocol.Package{
		Tag:     "NonceResp",
		Payload: payloadBytes,
	}

	packageBytes, _ := json.Marshal(packageNonceResp)

	mockCodec.EXPECT().Marshal(gomock.AssignableToTypeOf(nonceResp)).Return(payloadBytes, nil)
	mockCodec.EXPECT().Marshal(gomock.AssignableToTypeOf(packageNonceResp)).Return(packageBytes, nil)

	err := svc.onNonceReq(ctx, &protocol.NonceReq{}, writer)
	test.Nil(t, "onNonceReq error", err)

	var pkg protocol.Package
	err = json.Unmarshal(writer.Bytes(), &pkg)
	test.Nil(t, "Unmarshal package error", err)

	var returnedNonceResp protocol.NonceResp
	err = json.Unmarshal(pkg.Payload, &returnedNonceResp)

	test.Nil(t, "Unmarshal NonceResp error", err)
	test.Check(t, "NonceResp", *nonceResp, returnedNonceResp)
}

func TestService_onDataReq_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockCodec(ctrl)
	mockDDoSGuard := mock.NewMockDDoSGuard(ctrl)
	mockUserLogic := mock.NewMockUserLogic(ctrl)

	svc := NewService(1024, 10, 1.0, mockCodec, mockUserLogic, mockDDoSGuard)

	ctx := context.Background()
	writer := &bytes.Buffer{}

	svc.challenge = []byte("challenge")
	svc.difficulty = 1

	mockUserLogic.EXPECT().GetQuote(gomock.Any()).Return("test quote", nil)

	dataResp := &protocol.DataResp{
		Payload: []byte("test quote"),
	}

	payloadBytes, _ := json.Marshal(dataResp)

	packageDataResp := &protocol.Package{
		Tag:     "DataResp",
		Payload: payloadBytes,
	}

	packageBytes, _ := json.Marshal(packageDataResp)

	mockCodec.EXPECT().Marshal(gomock.AssignableToTypeOf(dataResp)).Return(payloadBytes, nil)
	mockCodec.EXPECT().Marshal(gomock.AssignableToTypeOf(packageDataResp)).Return(packageBytes, nil)

	err := svc.onDataReq(ctx, &protocol.DataReq{CNonce: []byte("solution")}, writer)
	test.Nil(t, "onDataReq error", err)

	var pkg protocol.Package
	err = json.Unmarshal(writer.Bytes(), &pkg)
	test.Nil(t, "Unmarshal package error", err)

	var returnedDataResp protocol.DataResp
	err = json.Unmarshal(pkg.Payload, &returnedDataResp)

	test.Nil(t, "Unmarshal DataResp error", err)
	test.Check(t, "DataResp", *dataResp, returnedDataResp)
}

func TestService_onDataReq_BadSolution(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockCodec(ctrl)
	mockDDoSGuard := mock.NewMockDDoSGuard(ctrl)
	mockUserLogic := mock.NewMockUserLogic(ctrl)

	svc := NewService(1024, 10, 1.0, mockCodec, mockUserLogic, mockDDoSGuard)

	ctx := context.Background()
	writer := &bytes.Buffer{}

	svc.challenge = []byte("challenge")
	svc.difficulty = 1

	errorResp := &protocol.ErrorResp{
		Reason: "Bad challenge solution",
	}

	payloadBytes, _ := json.Marshal(errorResp)

	packageErrorResp := &protocol.Package{
		Tag:     "ErrorResp",
		Payload: payloadBytes,
	}

	packageBytes, _ := json.Marshal(packageErrorResp)

	mockCodec.EXPECT().Marshal(gomock.AssignableToTypeOf(errorResp)).Return(payloadBytes, nil)
	mockCodec.EXPECT().Marshal(gomock.AssignableToTypeOf(packageErrorResp)).Return(packageBytes, nil)

	err := svc.onDataReq(ctx, &protocol.DataReq{CNonce: []byte("bad_solution")}, writer)
	test.Nil(t, "onDataReq error", err)

	var pkg protocol.Package
	err = json.Unmarshal(writer.Bytes(), &pkg)
	test.Nil(t, "Unmarshal package error", err)

	var returnedErrorResp protocol.ErrorResp
	err = json.Unmarshal(pkg.Payload, &returnedErrorResp)

	test.Nil(t, "Unmarshal ErrorResp error", err)
	test.Check(t, "ErrorResp", *errorResp, returnedErrorResp)
}
