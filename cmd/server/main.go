package main

import (
	"context"
	"faraway/wow/app/infrastructure/config"
	"faraway/wow/app/infrastructure/server"
	"faraway/wow/app/infrastructure/server/acceptor"
	"faraway/wow/app/usecase/server/logic"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "/etc/faraway/wow/conf/server.yaml"
	}

	// Create temporary logger for initial loging
	logger := zap.Must(zap.NewProduction())

	cfg := config.NewConfig(configPath, logger)

	// In this case, we are using the global zap logger for simplicity.
	// In production scenarios with distributed microservices, logging and tracing
	// becomes much more complex and requires more comprehensive solutions.

	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zap.InfoLevel
		logger.Warn("Failed to parse logging level, falling back to Info level", zap.Error(err))
	}

	zap.ReplaceGlobals(zap.Must(zap.NewProduction(zap.IncreaseLevel(level))))

	_ = logger.Sync()
}

func main() {
	zap.L().Info("Running Word Of Wisdom server")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)

	cfg := config.Get()

	ctx, cancel := context.WithCancel(context.Background())
	srv := server.New(func() server.ReusableAcceptor {
		return acceptor.New(
			logic.New(),
			cfg.Net.BuffSize,
			time.Duration(cfg.Net.Timeout),
		)
	})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		// Server runs until context cancel is called
		if err := srv.Run(ctx, cfg.Net.Host, cfg.Net.Port); err != nil {
			zap.L().Fatal("Server failed to run", zap.Error(err))
		}
	}()

	sig := <-signals

	zap.L().Info("Word Of Wisdom server shutdowned", zap.String("signal", sig.String()))

	cancel()

	wg.Wait()
}
