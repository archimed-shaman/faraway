package main

import (
	"context"
	"errors"
	"faraway/wow/app/infrastructure/config"
	jsonC "faraway/wow/app/interface/service/codec/json"
	"faraway/wow/app/interface/service/dispatcher"
	"faraway/wow/pkg/pow"
	"faraway/wow/pkg/protocol"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pkgerr "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrServerError = errors.New("server error")

type Codec interface {
	GetRaw(r io.Reader, buff []byte) (*protocol.Package, error)
	Unmarshal(data []byte, v any) error
	Marshal(v any) ([]byte, error)
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
		// TODO: warning
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

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				runConnect(ctx, cfg)
			}
		}
	}()

	sig := <-signals

	zap.L().Info("Client interrupted", zap.String("signal", sig.String()))
	cancel()

	wg.Wait()
}

func runConnect(ctx context.Context, cfg *config.Config) {
	defer func() { _ = zap.L().Sync() }()

	codec := jsonC.NewCodec()
	disp := dispatcher.New(codec)

	disp.Register(
		dispatcher.NewProcessor(func(ctx context.Context, pkg *protocol.NonceResp, w io.Writer) error {
			return onNonceResp(ctx, time.Duration(cfg.Net.Timeout), codec, pkg, w)
		}),
		dispatcher.NewProcessor(func(_ context.Context, pkg *protocol.DataResp, w io.Writer) error {
			return onDataResp(codec, pkg, w)
		}),
		dispatcher.NewProcessor(onErrorResp),
	)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Net.Host, cfg.Net.Port))
	if err != nil {
		zap.L().Fatal("Failed to connect to server", zap.Error(err))
	}

	defer conn.Close()

	zap.L().Info("Connected to server")

	if err := loop(ctx, cfg, conn, codec, disp); err != nil {
		zap.L().Error("Disconnecting on error...", zap.Error(err))
	}
}

func loop(ctx context.Context,
	cfg *config.Config, conn net.Conn, codec Codec, disp *dispatcher.Dispatcher,
) error {
	if err := sendNonceReq(codec, conn); err != nil {
		return err
	}

	buff := make([]byte, cfg.Net.BuffSize)

	for {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(cfg.Net.Timeout))); err != nil {
			return pkgerr.Wrap(err, "failed to set connection timeout")
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			if err := disp.Dispatch(ctx, conn, conn, buff); err != nil {
				return pkgerr.Wrap(err, "failed to dispatch")
			}
		}
	}
}

func onNonceResp(ctx context.Context,
	timeout time.Duration, codec Codec, pkg *protocol.NonceResp, w io.Writer,
) error {
	start := time.Now()

	dlCtx, cancel := context.WithDeadline(ctx, start.Add(timeout))
	defer cancel()

	cNonce, err := pow.Resolve(dlCtx, pkg.Nonce, pkg.Difficulty)
	if err != nil {
		zap.L().Error("Failed to find solution for challenge",
			zap.Int("difficulty", pkg.Difficulty), zap.Error(err))
	}

	zap.L().Info("Challenge solution found",
		zap.Int("difficulty", pkg.Difficulty), zap.Duration("elapsed", time.Since(start)))

	data, err := dispatcher.EncodePackage(&protocol.DataReq{
		Nonce:      pkg.Nonce,
		Difficulty: pkg.Difficulty,
		CNonce:     cNonce,
	}, codec)
	if err != nil {
		return pkgerr.Wrap(err, "failed to encode package with DataReq")
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgerr.Wrap(err, "failed to send DataReq")
	}

	return nil
}

func onDataResp(codec Codec, pkg *protocol.DataResp, w io.Writer) error {
	zap.L().Info("Client received data", zap.ByteString("quote", pkg.Payload))

	if err := sendNonceReq(codec, w); err != nil {
		return err
	}

	return nil
}

func onErrorResp(ctx context.Context, pkg *protocol.ErrorResp, w io.Writer) error {
	zap.L().Error("Server returned error", zap.ByteString("error", []byte(pkg.Reason)))
	return pkgerr.Wrap(ErrServerError, pkg.Reason)
}

func sendNonceReq(codec Codec, w io.Writer) error {
	data, err := dispatcher.EncodePackage(&protocol.NonceReq{}, codec)
	if err != nil {
		return pkgerr.Wrap(err, "failed to encode package with NonceReq")
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgerr.Wrap(err, "failed to send NonceReq")
	}

	return nil
}
