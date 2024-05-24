package main

import (
	"context"
	jsonC "faraway/wow/app/interface/service/codec/json"
	"faraway/wow/pkg/pow"
	"faraway/wow/pkg/protocol"
	"net"
	"os"
	"time"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logLevel := os.Getenv("LOG_LEVEL")
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		// TODO: warning
		level = zap.InfoLevel
	}

	logger := zap.Must(zap.NewProduction(zap.IncreaseLevel(level)))

	zap.ReplaceGlobals(logger)

	defer logger.Sync()

	codec := jsonC.NewCodec()

	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		logger.Fatal("Failed to connect to server", zap.Error(err))
	}

	defer conn.Close()

	logger.Info("Connected to server")

	// Read initial challenge request
	challengeReq, err := readData[protocol.ChallengeReq](ctx, conn, codec)
	if err != nil {
		logger.Fatal("Failed to connect to server", zap.Error(err))
	}

	solution, err := pow.Resolve(challengeReq.Challenge, challengeReq.Difficulty)
	if err != nil {
		logger.Fatal("Failed to find solution for challenge", zap.Error(err))
	}

	if err := writeData(ctx, conn, codec, &protocol.ChallengeResp{
		Challenge:  challengeReq.Challenge,
		Difficulty: challengeReq.Difficulty,
		Solution:   solution,
	}); err != nil {
		logger.Fatal("Failed to send challenge response", zap.Error(err))
	}

	// Read initial challenge request
	data, err := readData[protocol.Data](ctx, conn, codec)
	if err != nil {
		logger.Fatal("Failed to get data from server", zap.Error(err))
	}

	logger.Info("Quote received", zap.String("quote", string(data.Payload)))

	logger.Info("Client shutting down")
}

func readData[A any](ctx context.Context, conn net.Conn, codec Codec) (*A, error) {
	buffer := make([]byte, 1024)

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		return nil, pkgerr.Wrap(err, "failed to set read timeout")
	}

	n, err := conn.Read(buffer)
	if err != nil {
		return nil, pkgerr.Wrap(err, "failed to read data from server")
	}

	var inst A

	if err := codec.Unmarshal(buffer[:n], &inst); err != nil {
		return nil, pkgerr.Wrapf(err, "failed to unmarshal data from server")
	}

	return &inst, nil
}

func writeData[A any](ctx context.Context, conn net.Conn, codec Codec, obj *A) error {
	byteObj, err := codec.Marshal(obj)
	if err != nil {
		return pkgerr.Wrap(err, "failed to marshal data")
	}

	if err := conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return pkgerr.Wrap(err, "failed to set write timeout")
	}

	if _, err := conn.Write(byteObj); err != nil {
		return pkgerr.Wrap(err, "failed to write data to server")
	}

	return nil
}
