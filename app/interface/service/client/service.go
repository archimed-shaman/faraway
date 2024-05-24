package client

import (
	"context"
	"errors"
	"faraway/wow/app/infrastructure/server"
	"faraway/wow/pkg/pow"
	"faraway/wow/pkg/protocol"
	"io"
	"os"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	unknownRate   = -1
	challengeLen  = 32
	maxDifficulty = 32
)

type DDoSGuard interface {
	IncRate(ctx context.Context) (int64, error)
	DecRate(ctx context.Context) (int64, error)
}

type Encoder interface {
	Marshal(v any) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, v interface{}) error
}

type Codec interface {
	Encoder
	Decoder
}

type Service struct {
	buff  []byte
	codec Codec
	ddos  DDoSGuard
}

func NewService(buffSize int, codec Codec, ddos DDoSGuard) *Service {
	return &Service{
		codec: codec,
		buff:  make([]byte, buffSize),
		ddos:  ddos,
	}
}

func (s *Service) OnConnect(ctx context.Context, w server.ResponseWriter) error {
	rate, err := s.ddos.IncRate(ctx)
	if err != nil {
		zap.L().Error("Failed to increase current rate", zap.Error(err))

		rate = unknownRate
	}

	byteReq, err := s.mkChallengeReq(rate)
	if err != nil { // TODO: send error response to client
		zap.L().Error("Failed make challenge request", zap.Error(err))
		return io.EOF
	}

	_, err = w.Write(ctx, byteReq)

	switch {
	case err == nil:
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client")
		return io.EOF
	default:
		return pkgerr.Wrap(err, "client service failed send challenge request to client")
	}

	return nil
}

func (s *Service) OnData(ctx context.Context, r server.ResponseReader, w server.ResponseWriter) error {
	n, err := r.Data().Read(s.buff)

	switch {
	case err == nil: // Everything is ok, just read data
	case errors.Is(err, os.ErrDeadlineExceeded) || n == 0: // No data received
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client")
	default:
		zap.L().Debug("Error reading data from connection", zap.Error(err))
		return pkgerr.Wrap(err, "client service failed read data from the connection")
	}

	if _, err := w.Write(ctx, s.buff[:n]); err != nil {
		return pkgerr.Wrap(err, "client service failed write data to the connection")
	}

	return io.EOF
}

func (s *Service) OnDisconnect(ctx context.Context) {
	if _, err := s.ddos.DecRate(ctx); err != nil {
		zap.L().Error("Failed to decrease current rate", zap.Error(err))
	}
}

func (s *Service) mkChallengeReq(rate int64) ([]byte, error) {
	difficulty := int(rate) // In case of overflow, we will fix it at the next line
	if rate > maxDifficulty || rate <= 0 {
		difficulty = maxDifficulty
	}

	challenge, err := pow.GenChallenge(challengeLen, difficulty)
	if err != nil {
		return nil, pkgerr.Wrap(err, "failed generate challenge")
	}

	challengeReq := protocol.ChallengeReq{
		Challenge:  challenge,
		Difficulty: difficulty,
	}

	byteReq, err := s.codec.Marshal(&challengeReq)
	if err != nil {
		return nil, pkgerr.Wrap(err, "failed marshal challenge request")
	}

	return byteReq, nil
}
