package client

import (
	"context"
	"faraway/wow/app/interface/service/dispatcher"
	"faraway/wow/pkg/pow"
	"faraway/wow/pkg/protocol"
	"fmt"
	"io"
	"math"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	unknownRate  = -1
	challengeLen = 32
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE

type Codec interface {
	GetRaw(r io.Reader, buff []byte) (*protocol.Package, error)
	Unmarshal(data []byte, v any) error
	Marshal(v any) ([]byte, error)
}

type DDoSGuard interface {
	IncRate(ctx context.Context) (int64, error)
	DecRate(ctx context.Context) (int64, error)
}

type UserLogic interface {
	GetQuote(ctx context.Context) (string, error)
}

type Service struct {
	dispatcher           *dispatcher.Dispatcher
	buff                 []byte
	maxDifficulty        int
	rateDifficultyFactor float64
	codec                Codec
	ddos                 DDoSGuard
	logic                UserLogic

	// Client state
	// FIXME: add outstanding state holder or make it stateless (like JWT)
	challenge  []byte
	difficulty int
}

func NewService(
	buffSize int, maxDifficulty int, rateDifficultyFactor float64,
	codec Codec, logic UserLogic, ddos DDoSGuard,
) *Service {
	service := &Service{
		dispatcher:           dispatcher.New(codec),
		buff:                 make([]byte, buffSize),
		maxDifficulty:        maxDifficulty,
		rateDifficultyFactor: rateDifficultyFactor,
		codec:                codec,
		ddos:                 ddos,
		logic:                logic,
		challenge:            nil,
		difficulty:           0,
	}

	service.dispatcher.Register(
		dispatcher.NewProcessor(service.onNonceReq),
		dispatcher.NewProcessor(service.onDataReq),
	)

	return service
}

func (s *Service) Handle(ctx context.Context, r io.Reader, w io.Writer) error {
	if err := s.dispatcher.Dispatch(ctx, r, w, s.buff); err != nil {
		return pkgerr.Wrap(err, "client service failed to route data")
	}

	return nil
}

func (s *Service) onNonceReq(ctx context.Context, pkg *protocol.NonceReq, w io.Writer) error {
	rate, err := s.ddos.IncRate(ctx)
	if err != nil {
		zap.L().Error("Failed to increase current rate", zap.Error(err))

		rate = unknownRate
	}

	nonceResp, err := s.mkNonceResp(rate)
	if err != nil {
		zap.L().Error("Failed to make challenge request", zap.Error(err))

		if sendErr := s.sendErrorResponse(ctx, w, "Failed to make challenge request"); sendErr != nil {
			zap.L().Error("Failed to send error", zap.Error(err))
		}

		return err
	}

	s.challenge = nonceResp.Nonce
	s.difficulty = nonceResp.Difficulty

	bytes, err := dispatcher.EncodePackage(nonceResp, s.codec)
	if err != nil {
		return pkgerr.Wrap(err, "client service failed to encode nonce response")
	}

	if _, err := w.Write(bytes); err != nil {
		return pkgerr.Wrap(err, "client service failed respond nonce")
	}

	return nil
}

func (s *Service) onDataReq(ctx context.Context, pkg *protocol.DataReq, w io.Writer) error {
	ok, err := pow.CheckSolution(s.challenge, pkg.CNonce, s.difficulty)
	if err != nil {
		return pkgerr.Wrap(err, "client service failed to check challenge solution")
	}

	if !ok {
		zap.L().Debug("Bad challenge solution",
			zap.ByteString("challenge", s.challenge),
			zap.ByteString("solution", pkg.CNonce),
			zap.Int("difficulty", s.difficulty),
		)

		if sendErr := s.sendErrorResponse(ctx, w, "Bad challenge solution"); sendErr != nil {
			zap.L().Error("Failed to send error", zap.Error(err))
		}

		return nil
	}

	quote, err := s.logic.GetQuote(ctx)
	if err != nil {
		zap.L().Error("Failed to get quote from user logic", zap.Error(err))
		return pkgerr.Wrap(err, "client service failed to get quote")
	}

	bytes, err := dispatcher.EncodePackage(&protocol.DataResp{Payload: []byte(quote)}, s.codec)
	if err != nil {
		return pkgerr.Wrap(err, "client service failed to encode data response")
	}

	if _, err := w.Write(bytes); err != nil {
		return pkgerr.Wrap(err, "client service failed respond data")
	}

	return nil
}

func (s *Service) mkNonceResp(rate int64) (*protocol.NonceResp, error) {
	difficulty := int(math.Floor(float64(rate) * s.rateDifficultyFactor))
	if difficulty > s.maxDifficulty {
		difficulty = s.maxDifficulty
	}

	fmt.Println(rate, s.maxDifficulty, difficulty)

	challenge, err := pow.GenChallenge(challengeLen, difficulty)
	if err != nil {
		return nil, pkgerr.Wrap(err, "client service failed to generate challenge")
	}

	return &protocol.NonceResp{
		Nonce:      challenge,
		Difficulty: difficulty,
	}, nil
}

func (s *Service) sendErrorResponse(ctx context.Context, w io.Writer, msg string) error {
	// Fixed set of error codes would be better
	errResp := protocol.ErrorResp{Reason: msg}

	data, err := dispatcher.EncodePackage(&errResp, s.codec)
	if err != nil {
		return pkgerr.Wrap(err, "client service failed to encode ErrorResp package")
	}

	if _, err := w.Write(data); err != nil {
		return pkgerr.Wrap(err, "client service failed to send ErrorResp package")
	}

	return nil
}
