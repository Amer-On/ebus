package main

import (
	"ebus/internal/broker"
	"ebus/internal/gateway/api"
	"ebus/internal/gateway/config"
	"ebus/internal/gateway/server"
	"ebus/internal/gateway/service"
	configPkg "ebus/pkg/config"
	"ebus/pkg/logging"
	"ebus/pkg/redis"

	"go.uber.org/zap"
)

func main() {
	logger, err := logging.InitLogger()
	if err != nil {
		panic("Erorr initializing logger")
	}

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Recovered from panic", zap.Error(err.(error)))
		}
		logger.Info("Stopping gateway")
	}()

	logger.Info("Successfully initialized logger")

	config, err := configPkg.ReadConfig[config.Config]("internal/gateway/config/config.yaml")
	if err != nil {
		panic(err)
	}

	// Create Redis Client
	redisClient := redis.NewClient(config.Redis)

	// Create Broker
	broker := broker.NewRedisStreamsBroker(logger, redisClient)

	// Create Services
	subscriptionService := service.NewSubscriptionService(logger)
	publishService := service.NewPublishService(logger, subscriptionService, broker)

	api := api.NewAPI(logger, subscriptionService, publishService)

	// start http server
	server, err := server.NewServer(logger, config.Server, api)
	if err != nil {
		panic(err)
	}

	logger.Info("Running HTTP Server")
	err = server.Run()
	if err != nil {
		panic(err)
	}
}
