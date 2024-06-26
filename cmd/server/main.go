package main

import (
	"context"
	"faraway/wow/app/infrastructure/config"
	"faraway/wow/app/infrastructure/server"
	"faraway/wow/app/infrastructure/version"
	"faraway/wow/app/interface/service/client"
	jsonC "faraway/wow/app/interface/service/codec/json"
	"faraway/wow/app/interface/service/ddos"
	"faraway/wow/app/usecase/logic"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	const serviceName = "server"

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "/etc/faraway/wow/conf/server.yaml"
	}

	// Create temporary logger for initial loging
	logger := zap.Must(zap.NewProduction(zap.Fields(zap.String("service", serviceName))))

	cfg := config.NewConfig(configPath, logger)

	// In this case, we are using the global zap logger for simplicity.
	// In production scenarios with distributed microservices, logging and tracing
	// becomes much more complex and requires more comprehensive solutions.

	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zap.InfoLevel
		logger.Warn("Failed to parse logging level, falling back to Info level", zap.Error(err))
	}

	logCfg := zap.NewProductionConfig()
	logCfg.Level.SetLevel(level)
	zap.ReplaceGlobals(zap.Must(logCfg.Build(zap.Fields(zap.String("service", serviceName)))))

	_ = logger.Sync()
}

func main() {
	zap.L().Info("Running Word Of Wisdom server",
		zap.String("name", version.GetProductName()),
		zap.String("version", version.GetVersion()),
		zap.String("build_date", version.GetDate()),
		zap.String("git", version.GetGit()),
	)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)

	cfg := config.Get()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	ddosGuard := ddos.NewGuard(time.Duration(cfg.App.Window))

	wg.Add(1)

	go func() {
		defer wg.Done()
		ddosGuard.Run(ctx)
	}()

	codec := jsonC.NewCodec()
	userLogic := logic.New()

	logicAllocator := func() server.Service {
		return client.NewService(cfg.Net.BuffSize, cfg.App.MaxDifficulty, cfg.App.RateDifficultyFactor,
			codec, userLogic, ddosGuard)
	}

	srv := server.New(cfg.Net.MaxConnection, time.Duration(cfg.Net.Timeout), logicAllocator)

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
