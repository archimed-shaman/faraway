//nolint:errchkjson,goerr113 // To avoid test code in common space
package dispatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"faraway/wow/pkg/protocol"
	"faraway/wow/pkg/test"
	"io"
	"testing"

	mock "faraway/wow/app/interface/service/dispatcher/mock"

	"go.uber.org/mock/gomock"
)

func Test_Handle_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	testPayload := protocol.NonceReq{}
	payloadBytes, _ := json.Marshal(testPayload)
	mockCodec.EXPECT().Unmarshal(payloadBytes, gomock.Any()).SetArg(1, testPayload).Return(nil)

	handler := func(ctx context.Context, pkg *protocol.NonceReq, w io.Writer) error {
		test.Check(t, "Handler received package", testPayload, *pkg)
		return nil
	}

	processor := NewProcessor(handler)

	err := processor.Handle(context.Background(), mockCodec, payloadBytes, &bytes.Buffer{})

	test.Nil(t, "Error", err)
}

func Test_Handle_UnmarshalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	payloadBytes := []byte("invalid payload")
	err := errors.New("unmarshal error")

	mockCodec.EXPECT().Unmarshal(payloadBytes, gomock.Any()).Return(err)

	handler := func(ctx context.Context, pkg *protocol.NonceReq, w io.Writer) error {
		return nil
	}

	processor := NewProcessor(handler)

	got := processor.Handle(context.Background(), mockCodec, payloadBytes, &bytes.Buffer{})

	test.Err(t, "Unmarshal error", err, got)
}

func Test_Handle_HandlerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockStreamDecoder(ctrl)

	testPayload := protocol.NonceReq{}
	payloadBytes, _ := json.Marshal(testPayload)

	mockCodec.EXPECT().Unmarshal(payloadBytes, gomock.Any()).SetArg(1, testPayload).Return(nil)

	err := errors.New("handler error")

	handler := func(ctx context.Context, pkg *protocol.NonceReq, w io.Writer) error {
		return err
	}

	processor := NewProcessor(handler)

	got := processor.Handle(context.Background(), mockCodec, payloadBytes, &bytes.Buffer{})

	test.Err(t, "Handler error", err, got)
}
