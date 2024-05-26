package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mockServer "faraway/wow/app/infrastructure/server/types/mock"
	mock "faraway/wow/app/interface/service/client/mock"
	"faraway/wow/pkg/protocol"
	"faraway/wow/pkg/test"
	"io"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

const (
	buffSize             = 1024
	maxDifficulty        = 5
	rateDifficultyFactor = 1
)

var errTestAssertionFailed = errors.New("type assertion failed")

func TestService_OnConnect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedChallenge := []byte("mocked")

	mockCodec := mock.NewMockCodec(ctrl)
	mockDDoSGuard := mock.NewMockDDoSGuard(ctrl)
	mockLogic := mock.NewMockUserLogic(ctrl)
	mockWriter := mockServer.NewMockResponseWriter(ctrl)

	var challenge []byte

	rate := 4

	mockDDoSGuard.EXPECT().IncRate(gomock.Any()).Return(int64(rate), nil)

	mockCodec.EXPECT().Marshal(gomock.Any()).
		DoAndReturn(func(v any) ([]byte, error) {
			req, ok := v.(*protocol.NonceResp)
			if !ok {
				t.Errorf("Bad type: %s (%+v)", reflect.TypeOf(v), v)
				return nil, errTestAssertionFailed
			}

			challenge = req.Nonce

			return mockedChallenge, nil
		})
	mockWriter.EXPECT().Write(gomock.Any(), mockedChallenge).Return(len(mockedChallenge), nil)

	s := NewService(maxDifficulty, rateDifficultyFactor, buffSize, mockLogic, mockCodec, mockDDoSGuard)

	err := s.OnConnect(context.Background(), mockWriter)

	test.Nil(t, "OnConnect error", err)
	test.Check(t, "OnConnect challenge", challenge, s.challenge)
	test.Check(t, "OnConnect difficulty", rate*rateDifficultyFactor, s.difficulty)
}

// TODO: add more tests to check errors from different interfaces

func TestService_OnData(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockCodec(ctrl)
	mockDDoSGuard := mock.NewMockDDoSGuard(ctrl)
	mockLogic := mock.NewMockUserLogic(ctrl)
	mockReader := mockServer.NewMockResponseReader(ctrl)
	mockWriter := mockServer.NewMockResponseWriter(ctrl)

	// Precalculated challenge and solution
	rate := 4
	challenge := []byte{
		142, 235, 122, 84, 235, 172, 46, 185, 5, 54, 158, 113, 220, 139, 151, 91,
		200, 37, 143, 77, 64, 125, 13, 129, 124, 100, 58, 7, 97, 180, 245, 3,
	}
	solution := []byte{20}

	challengeResp := protocol.GetDataReq{
		Nonce:      challenge,
		CNonce:     solution,
		Difficulty: rate,
	}

	challengeRespBytes, err := json.Marshal(&challengeResp)
	test.Nil(t, "Marshal error", err)

	mockReader.EXPECT().Data().Return(io.NopCloser(bytes.NewReader(challengeRespBytes)))
	mockCodec.EXPECT().Unmarshal(challengeRespBytes, gomock.Any()).
		DoAndReturn(func(data []byte, v any) error {
			resp, ok := v.(*protocol.GetDataReq)
			if !ok {
				t.Errorf("Bad type: %s (%+v)", reflect.TypeOf(v), v)
				return errTestAssertionFailed
			}

			*resp = challengeResp

			return nil
		})

	quote := "Word of Wisdom quote"

	mockLogic.EXPECT().GetQuote(gomock.Any()).Return(quote, nil)

	data := protocol.DataResp{
		Payload: []byte(quote),
	}
	mockedData := []byte("mocked data")

	mockCodec.EXPECT().Marshal(gomock.Any()).
		DoAndReturn(func(v any) ([]byte, error) {
			req, ok := v.(*protocol.DataResp)
			if !ok {
				t.Errorf("Bad type: %s (%+v)", reflect.TypeOf(v), v)
				return nil, errTestAssertionFailed
			}

			test.Check(t, "returned data", data, *req)

			return mockedData, nil
		})

	mockWriter.EXPECT().Write(gomock.Any(), mockedData).Return(len(mockedData), nil)

	s := NewService(maxDifficulty, rateDifficultyFactor, buffSize, mockLogic, mockCodec, mockDDoSGuard)
	s.challenge = challenge
	s.difficulty = rate

	err = s.OnData(context.Background(), mockReader, mockWriter)
	test.Err(t, "OnData error", io.EOF, err)
}

func TestService_OnDisconnect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCodec := mock.NewMockCodec(ctrl)
	mockDDoSGuard := mock.NewMockDDoSGuard(ctrl)
	mockLogic := mock.NewMockUserLogic(ctrl)

	mockDDoSGuard.EXPECT().DecRate(gomock.Any()).Return(int64(5), nil).Times(1)

	s := NewService(maxDifficulty, rateDifficultyFactor, buffSize, mockLogic, mockCodec, mockDDoSGuard)
	s.challenge = []byte("challenge")
	s.difficulty = 5

	s.OnDisconnect(context.Background())

	test.Check(t, "OnDisconnect challenge", []uint8(nil), s.challenge)
	test.Check(t, "OnDisconnect difficulty", 0, s.difficulty)
}
