package client

import (
	"context"
	"errors"
	serverTypes "faraway/wow/app/infrastructure/server/types"
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

var StopSignal = io.EOF

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE
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

type UserLogic interface {
	GetQuote(ctx context.Context) (string, error)
}

type Service struct {
	buff       []byte
	challenge  []byte
	difficulty int

	codec Codec
	ddos  DDoSGuard
	logic UserLogic
}

func NewService(buffSize int, logic UserLogic, codec Codec, ddos DDoSGuard) *Service {
	return &Service{
		buff:       make([]byte, buffSize),
		challenge:  nil,
		difficulty: maxDifficulty,

		codec: codec,
		ddos:  ddos,
		logic: logic,
	}
}

func (s *Service) OnConnect(ctx context.Context, w serverTypes.ResponseWriter) error {
	rate, err := s.ddos.IncRate(ctx)
	if err != nil {
		zap.L().Error("Failed to increase current rate", zap.Error(err))

		rate = unknownRate
	}

	challengeReq, err := s.mkChallengeReq(rate)
	if err != nil {
		zap.L().Error("Failed to make challenge request", zap.Error(err))
		s.sendErrorResponse(ctx, w, "Failed to make challenge request")

		return err
	}

	s.challenge, s.difficulty = challengeReq.Challenge, challengeReq.Difficulty

	byteReq, err := s.codec.Marshal(challengeReq)
	if err != nil {
		return pkgerr.Wrap(err, "client service failed to marshal challenge request")
	}

	_, err = w.Write(ctx, byteReq)
	if hErr := handleWriteError(err); hErr != nil {
		return hErr
	}

	return nil
}

func (s *Service) OnData(ctx context.Context, r serverTypes.ResponseReader, w serverTypes.ResponseWriter) error {
	n, err := r.Data().Read(s.buff)
	if hErr := handleReadError(err); hErr != nil {
		return hErr
	}

	// Assuming the buffer is enough to receive the response at once.
	// For the production use here must be a check if any other data is expected.
	var challengeResp protocol.ChallengeResp
	if err = s.codec.Unmarshal(s.buff[:n], &challengeResp); err != nil {
		zap.L().Debug("Failed to unmarshal quote", zap.Error(err))
		s.sendErrorResponse(ctx, w, "Bad response")

		return pkgerr.Wrap(err, "client service failed to unmarshal challenge response")
	}

	ok, err := pow.CheckSolution(s.challenge, challengeResp.Solution, s.difficulty)
	if err != nil {
		s.sendErrorResponse(ctx, w, "Bad challenge resolution")
		return pkgerr.Wrap(err, "client service failed to check challenge solution")
	}

	if !ok {
		zap.L().Debug("Bad challenge, disconnection response",
			zap.ByteString("challenge", s.challenge),
			zap.ByteString("solution", challengeResp.Solution),
			zap.Int("difficulty", s.difficulty),
		)

		s.sendErrorResponse(ctx, w, "Bad challenge resolution")

		return StopSignal
	}

	quote, err := s.logic.GetQuote(ctx)
	if err != nil {
		zap.L().Error("Failed to get quote from user logic", zap.Error(err))
		s.sendErrorResponse(ctx, w, "Internal error")

		return pkgerr.Wrap(err, "client service failed to get quote")
	}

	data, err := s.codec.Marshal(&protocol.Data{
		Payload: []byte(quote),
	})
	if err != nil {
		return pkgerr.Wrap(err, "client service failed to marshal data")
	}

	_, err = w.Write(ctx, data)
	if hErr := handleWriteError(err); hErr != nil {
		return hErr
	}

	return StopSignal
}

func (s *Service) OnDisconnect(ctx context.Context) {
	s.challenge = nil
	s.difficulty = 0

	if _, err := s.ddos.DecRate(ctx); err != nil {
		zap.L().Error("Failed to decrease current rate", zap.Error(err))
	}
}

func (s *Service) mkChallengeReq(rate int64) (*protocol.ChallengeReq, error) {
	difficulty := int(rate)
	if rate > maxDifficulty || rate <= 0 {
		difficulty = maxDifficulty
	}

	challenge, err := pow.GenChallenge(challengeLen, difficulty)
	if err != nil {
		return nil, pkgerr.Wrap(err, "client service failed to generate challenge")
	}

	return &protocol.ChallengeReq{
		Challenge:  challenge,
		Difficulty: difficulty,
	}, nil
}

func (s *Service) sendErrorResponse(ctx context.Context, w serverTypes.ResponseWriter, msg string) {
	// Fixed set of error codes would be better
	errResp := protocol.ErrorResp{Reason: msg}

	data, err := s.codec.Marshal(&errResp)
	if err != nil {
		zap.L().Error("Failed to marshal error response", zap.Error(err))
		return
	}

	if _, err := w.Write(ctx, data); err != nil {
		zap.L().Error("Failed to send error response", zap.Error(err))
	}
}

func handleReadError(err error) error {
	switch {
	case err == nil: // Everything is ok, just read data
	case errors.Is(err, os.ErrDeadlineExceeded): // No data received
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client")
		return StopSignal
	default:
		zap.L().Debug("Error reading data from connection", zap.Error(err))
		return pkgerr.Wrap(err, "client service failed to read data from the connection")
	}

	return nil
}

func handleWriteError(err error) error {
	switch {
	case err == nil:
	case errors.Is(err, io.EOF):
		zap.L().Info("Connection closed by client")
		return StopSignal
	default:
		return pkgerr.Wrap(err, "client service failed to send challenge request to client")
	}

	return nil
}
