package main

import (
	"context"
	"ebus/internal/broker"
	"ebus/internal/worker"
	"ebus/pkg/config"
	"ebus/pkg/logging"
	"ebus/pkg/redis"

	"go.uber.org/zap"
)

func main() {
	logger, err := logging.InitLogger()
	if err != nil {
		panic("Erorr initializing logger")
	}

	logger.Info("Successfully initialized logger")

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Recovered from panic", zap.Error(err.(error)))
		}
		cancel()
		logger.Info("Stopping gateway")
	}()

	// Read config
	config, err := config.ReadConfig[worker.Config]("internal/worker/config.yaml")
	if err != nil {
		panic("Error reading config")
	}

	// Create Redis Client
	redisClient := redis.NewClient(config.Redis)

	// Create Broker
	broker := broker.NewRedisStreamsBroker(logger, redisClient)

	worker := worker.NewWorker(logger, broker)

	logger.Info("Starting worker")
	worker.Start(ctx)
}
