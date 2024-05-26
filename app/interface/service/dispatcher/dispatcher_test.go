//nolint:errchkjson,goerr113 // To avoid test code in common space
package dispatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mock "faraway/wow/app/interface/service/dispatcher/mock"
	"faraway/wow/pkg/protocol"
	"faraway/wow/pkg/test"
	"io"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_Dispatch_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	testPayload := protocol.NonceReq{}
	payloadBytes, _ := json.Marshal(testPayload)
	pkg := &protocol.Package{
		Tag:     "NonceReq",
		Payload: payloadBytes,
	}
	buff := make([]byte, 1024)

	mockCodec.EXPECT().GetRaw(gomock.Any(), buff).Return(pkg, nil)
	mockCodec.EXPECT().Unmarshal(payloadBytes, gomock.Any()).SetArg(1, testPayload).Return(nil)

	handler := func(ctx context.Context, pkg *protocol.NonceReq, w io.Writer) error {
		test.Check(t, "Handler received package", testPayload, *pkg)
		return nil
	}

	processor := NewProcessor(handler)
	dispatcher := New(mockCodec)
	dispatcher.Register(processor)

	err := dispatcher.Dispatch(context.Background(), &bytes.Buffer{}, &bytes.Buffer{}, buff)

	test.Nil(t, "Error", err)
}

func Test_Dispatch_GetRawError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	buff := make([]byte, 1024)

	err := errors.New("get raw error")
	mockCodec.EXPECT().GetRaw(gomock.Any(), buff).Return(nil, err)

	dispatcher := New(mockCodec)

	got := dispatcher.Dispatch(context.Background(), &bytes.Buffer{}, &bytes.Buffer{}, buff)

	test.Err(t, "GetRaw error", err, got)
}

func Test_Dispatch_UnknownType(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	pkg := &protocol.Package{
		Tag:     "UnknownTag",
		Payload: []byte{},
	}
	buff := make([]byte, 1024)

	mockCodec.EXPECT().GetRaw(gomock.Any(), buff).Return(pkg, nil)

	dispatcher := New(mockCodec)

	err := dispatcher.Dispatch(context.Background(), &bytes.Buffer{}, &bytes.Buffer{}, buff)

	test.Err(t, "Unknown type", ErrBadType, err)
}

func Test_Dispatch_HandlerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	testPayload := protocol.NonceReq{}
	payloadBytes, _ := json.Marshal(testPayload)
	pkg := &protocol.Package{
		Tag:     "NonceReq",
		Payload: payloadBytes,
	}
	buff := make([]byte, 1024)

	mockCodec.EXPECT().GetRaw(gomock.Any(), buff).Return(pkg, nil)
	mockCodec.EXPECT().Unmarshal(payloadBytes, gomock.Any()).SetArg(1, testPayload).Return(nil)

	err := errors.New("handler error")
	handler := func(ctx context.Context, pkg *protocol.NonceReq, w io.Writer) error {
		return err
	}

	processor := NewProcessor(handler)
	dispatcher := New(mockCodec)
	dispatcher.Register(processor)

	got := dispatcher.Dispatch(context.Background(), &bytes.Buffer{}, &bytes.Buffer{}, buff)

	test.Err(t, "Handler error", err, got)
}
