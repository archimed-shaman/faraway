package main

import (
	"faraway/wow/app/infrastructure/config"
	jsonC "faraway/wow/app/interface/service/codec/json"
	"faraway/wow/pkg/pow"
	"faraway/wow/pkg/protocol"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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

func init() {
	const serviceName = "client"

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "/etc/faraway/wow/conf/client.yaml"
	}

	// Create temporary logger for initial loging
	logger := zap.Must(zap.NewProduction(zap.Fields(zap.String("service", serviceName))))

	cfg := config.NewConfig(configPath, logger)

	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		// TODO: wraning
		level = zap.InfoLevel
	}

	logCfg := zap.NewProductionConfig()
	logCfg.Level.SetLevel(level)
	zap.ReplaceGlobals(zap.Must(logCfg.Build(zap.Fields(zap.String("service", serviceName)))))

	_ = logger.Sync()
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)

	cfg := config.Get()

	for {
		select {
		case sig := <-signals:
			zap.L().Info("Client interrupted", zap.String("signal", sig.String()))
			return

		default:
			run(cfg)
		}
	}
}

func run(cfg *config.Config) {
	defer func() { _ = zap.L().Sync() }()

	codec := jsonC.NewCodec()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Net.Host, cfg.Net.Port))
	if err != nil {
		zap.L().Fatal("Failed to connect to server", zap.Error(err))
	}

	defer conn.Close()

	zap.L().Info("Connected to server")

	// Read initial challenge request
	challengeReq, err := readData[protocol.ChallengeReq](conn, codec,
		time.Duration(cfg.Net.Timeout), cfg.Net.BuffSize)
	if err != nil {
		zap.L().Error("Failed to read data from server", zap.Error(err))
		return
	}

	start := time.Now()

	solution, err := pow.Resolve(challengeReq.Challenge, challengeReq.Difficulty)
	if err != nil {
		zap.L().Fatal("Failed to find solution for challenge", zap.Error(err))
	}

	zap.L().Info("Challenge solution found",
		zap.Int("difficulty", challengeReq.Difficulty), zap.Duration("elapsed", time.Since(start)))

	if err = writeData(conn, codec, time.Duration(cfg.Net.Timeout), &protocol.ChallengeResp{
		Challenge:  challengeReq.Challenge,
		Difficulty: challengeReq.Difficulty,
		Solution:   solution,
	}); err != nil {
		zap.L().Error("Failed to send challenge response", zap.Error(err))
		return
	}

	// Read initial challenge request
	data, err := readData[protocol.Data](conn, codec, time.Duration(cfg.Net.Timeout), cfg.Net.BuffSize)
	if err != nil {
		zap.L().Error("Failed to get data from server", zap.Error(err))
		return
	}

	zap.L().Info("Quote received", zap.String("quote", string(data.Payload)))

	zap.L().Info("Client shutting down")
}

func readData[A any](conn net.Conn, codec Codec, timeout time.Duration, buffSize int) (*A, error) {
	buffer := make([]byte, buffSize)

	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
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

func writeData[A any](conn net.Conn, codec Codec, timeout time.Duration, obj *A) error {
	byteObj, err := codec.Marshal(obj)
	if err != nil {
		return pkgerr.Wrap(err, "failed to marshal data")
	}

	if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return pkgerr.Wrap(err, "failed to set write timeout")
	}

	if _, err := conn.Write(byteObj); err != nil {
		return pkgerr.Wrap(err, "failed to write data to server")
	}

	return nil
}
