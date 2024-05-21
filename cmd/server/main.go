package main

import "go.uber.org/zap"

func main() {
	logger := zap.Must(zap.NewProduction(zap.IncreaseLevel(zap.InfoLevel)))
	logger.Info("Running Word Of Wisdom server")
}
